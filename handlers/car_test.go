// handlers/car_handler_test.go
package handlers

import (
	"SE-CarRentalService/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Mock functions for dependencies
var mockGetCarsFromDatabase = func() ([]types.Car, error) {
	return []types.Car{
		{ID: 1, Model: "Model1", Brand: "Brand1", CollectAt: "2023-01-01", AccountId: 0, Year: 2020, Price: 10000, PS: 150},
		{ID: 2, Model: "Model2", Brand: "Brand2", CollectAt: "2023-01-02", AccountId: 1, Year: 2021, Price: 20000, PS: 200},
	}, nil
}

var mockGetCarByID = func(id int) (types.Car, error) {
	if id == 1 {
		return types.Car{ID: 1, Model: "Model1", Brand: "Brand1", CollectAt: "2023-01-01", AccountId: 0, Year: 2020, Price: 10000, PS: 150}, nil
	}
	return types.Car{}, fmt.Errorf("car not found")
}

var mockCreateCar = func(c types.CreateCarRequest, currency string) (types.Car, error) {
	return types.Car{ID: 3, Model: c.Model, Brand: c.Brand, CollectAt: c.CollectAt, AccountId: 0, Year: c.Year, Price: c.Price, PS: c.PS}, nil
}

var mockUpdateCar = func(id int, c types.Car, account types.Account) error {
	if id == 1 {
		return nil
	}
	return fmt.Errorf("car not found")
}

var mockDeleteCar = func(id int) error {
	if id == 1 {
		return nil
	}
	return fmt.Errorf("car not found")
}

var mockVerifyToken = func(ch *amqp.Channel, token string) (types.TokenResponse, error) {
	if token == "valid" {
		return types.TokenResponse{Account: types.Account{ID: 1, Currency: "USD", IsAdmin: true}, Valid: true}, nil
	} else if token == "valid_non_admin" {
		return types.TokenResponse{Account: types.Account{ID: 2, Currency: "EUR", IsAdmin: false}, Valid: true}, nil
	}
	return types.TokenResponse{Valid: false, Error: "invalid token"}, nil
}

var mockConvertCurrency = func(initialCurrency string, initialAmount float64, targetCurrency string) (float64, string, error) {
	// Simple mock: assume 1 USD = 0.85 EUR, etc.
	if initialCurrency == "USD" && targetCurrency == "EUR" {
		return initialAmount * 0.85, "EUR", nil
	}
	if initialCurrency == "EUR" && targetCurrency == "USD" {
		return initialAmount / 0.85, "USD", nil
	}
	return initialAmount, targetCurrency, nil
}

func setupCarTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := NewCarHandler()

	// Assign mocks
	getCarsFromDatabase = mockGetCarsFromDatabase
	getCarByID = mockGetCarByID
	createCar = mockCreateCar
	updateCar = mockUpdateCar
	deleteCar = mockDeleteCar
	verifyToken = mockVerifyToken
	convertCurrency = mockConvertCurrency

	r.GET("/cars", h.GetCars)
	r.GET("/car/:id", h.GetCar)
	r.POST("/car", h.CreateCar)
	r.PUT("/car/:id", h.UpdateCar)
	r.DELETE("/car/:id", h.DeleteCar)

	return r
}

func requestWithAuth(method string, path string, body any, token string) *httptest.ResponseRecorder {
	r := setupCarTestRouter()

	var buffer bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buffer).Encode(body)
	}

	req := httptest.NewRequest(method, path, &buffer)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestGetCarInvalidID(t *testing.T) {
	w := requestWithAuth(http.MethodGet, "/car/abc", nil, "")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCarInvalidID(t *testing.T) {
	w := requestWithAuth(http.MethodPut, "/car/abc", types.Car{}, "")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteCarInvalidIDShouldReturn400(t *testing.T) {
	w := requestWithAuth(http.MethodDelete, "/car/abc", nil, "")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

// Additional tests

func TestGetCarsSuccess(t *testing.T) {
	w := requestWithAuth(http.MethodGet, "/cars", nil, "valid")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var cars []types.Car
	if err := json.Unmarshal(w.Body.Bytes(), &cars); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(cars) != 2 {
		t.Fatalf("expected 2 cars, got %d", len(cars))
	}
}

func TestGetCarSuccess(t *testing.T) {
	w := requestWithAuth(http.MethodGet, "/car/1", nil, "valid")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var car types.Car
	if err := json.Unmarshal(w.Body.Bytes(), &car); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if car.ID != 1 {
		t.Fatalf("expected car ID 1, got %d", car.ID)
	}
}

func TestGetCarNotFound(t *testing.T) {
	w := requestWithAuth(http.MethodGet, "/car/999", nil, "valid")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateCarSuccess(t *testing.T) {
	req := types.CreateCarRequest{
		Model:     "NewModel",
		Brand:     "NewBrand",
		CollectAt: "2023-01-03",
		Year:      2022,
		Price:     30000,
		PS:        250,
	}
	w := requestWithAuth(http.MethodPost, "/car", req, "valid")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var car types.Car
	if err := json.Unmarshal(w.Body.Bytes(), &car); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if car.Model != "NewModel" {
		t.Fatalf("expected model NewModel, got %s", car.Model)
	}
}

func TestCreateCarInvalidToken(t *testing.T) {
	req := types.CreateCarRequest{
		Model:     "NewModel",
		Brand:     "NewBrand",
		CollectAt: "2023-01-03",
		Year:      2022,
		Price:     30000,
		PS:        250,
	}
	w := requestWithAuth(http.MethodPost, "/car", req, "invalid")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateCarInvalidBody(t *testing.T) {
	w := requestWithAuth(http.MethodPost, "/car", map[string]string{"invalid": "body"}, "valid")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCarSuccess(t *testing.T) {
	car := types.Car{
		Model:     "UpdatedModel",
		Brand:     "UpdatedBrand",
		CollectAt: "2023-01-04",
		Year:      2023,
		Price:     40000,
		PS:        300,
	}
	w := requestWithAuth(http.MethodPut, "/car/1", car, "valid")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCarInvalidToken(t *testing.T) {
	car := types.Car{
		Model:     "UpdatedModel",
		Brand:     "UpdatedBrand",
		CollectAt: "2023-01-04",
		Year:      2023,
		Price:     40000,
		PS:        300,
	}
	w := requestWithAuth(http.MethodPut, "/car/1", car, "invalid")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCarNotFound(t *testing.T) {
	car := types.Car{
		Model:     "UpdatedModel",
		Brand:     "UpdatedBrand",
		CollectAt: "2023-01-04",
		Year:      2023,
		Price:     40000,
		PS:        300,
	}
	w := requestWithAuth(http.MethodPut, "/car/999", car, "valid")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteCarSuccess(t *testing.T) {
	w := requestWithAuth(http.MethodDelete, "/car/1", nil, "valid")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteCarInvalidToken(t *testing.T) {
	w := requestWithAuth(http.MethodDelete, "/car/1", nil, "invalid")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteCarNotFound(t *testing.T) {
	w := requestWithAuth(http.MethodDelete, "/car/999", nil, "valid")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetCarsWithInvalidToken(t *testing.T) {
	w := requestWithAuth(http.MethodGet, "/cars", nil, "invalid")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var cars []types.Car
	if err := json.Unmarshal(w.Body.Bytes(), &cars); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(cars) != 2 {
		t.Fatalf("expected 2 cars, got %d", len(cars))
	}
	// Assuming default currency EUR, prices should be converted
}

func TestGetCarWithInvalidToken(t *testing.T) {
	w := requestWithAuth(http.MethodGet, "/car/1", nil, "invalid")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var car types.Car
	if err := json.Unmarshal(w.Body.Bytes(), &car); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if car.ID != 1 {
		t.Fatalf("expected car ID 1, got %d", car.ID)
	}
}

func TestCreateCarWithNonAdminToken(t *testing.T) {
	req := types.CreateCarRequest{
		Model:     "NewModel",
		Brand:     "NewBrand",
		CollectAt: "2023-01-03",
		Year:      2022,
		Price:     30000,
		PS:        250,
	}
	w := requestWithAuth(http.MethodPost, "/car", req, "valid_non_admin")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCarWithNonAdminTokenRental(t *testing.T) {
	// Assuming car 1 is available (accountId=0)
	car := types.Car{
		AccountId: 2, // Non-admin ID
	}
	w := requestWithAuth(http.MethodPut, "/car/1", car, "valid_non_admin")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCarWithNonAdminTokenAlreadyRented(t *testing.T) {
	// Car 2 is rented by someone else (accountId=1), non-admin (2) tries to rent
	car := types.Car{
		AccountId: 2,
	}
	w := requestWithAuth(http.MethodPut, "/car/2", car, "valid_non_admin")

	// Should fail with 500 (car not found error, but actually not updatable)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteCarWithNonAdminToken(t *testing.T) {
	w := requestWithAuth(http.MethodDelete, "/car/1", nil, "valid_non_admin")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateCarWithEmptyBody(t *testing.T) {
	w := requestWithAuth(http.MethodPost, "/car", "", "valid")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateCarWithInvalidJSON(t *testing.T) {
	w := requestWithAuth(http.MethodPost, "/car", `{"model": "test", "invalid":}`, "valid")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateCarWithMissingRequiredFields(t *testing.T) {
	req := map[string]interface{}{
		"model": "Test",
		// Missing brand, year, price, ps
	}
	w := requestWithAuth(http.MethodPost, "/car", req, "valid")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCarWithEmptyBody(t *testing.T) {
	w := requestWithAuth(http.MethodPut, "/car/1", "", "valid")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetCarWithNegativeID(t *testing.T) {
	w := requestWithAuth(http.MethodGet, "/car/-1", nil, "valid")

	// Should fail as strconv.Atoi("-1") succeeds, but db query might fail
	// Assuming mock handles it as not found
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetCarWithVeryLargeID(t *testing.T) {
	w := requestWithAuth(http.MethodGet, "/car/999999999", nil, "valid")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateCarWithNegativeYear(t *testing.T) {
	req := types.CreateCarRequest{
		Model:     "NewModel",
		Brand:     "NewBrand",
		CollectAt: "2023-01-03",
		Year:      -1,
		Price:     30000,
		PS:        250,
	}
	w := requestWithAuth(http.MethodPost, "/car", req, "valid")

	// Gin's binding might allow it, but test if code handles
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 (if allowed), got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateCarWithVeryLongStrings(t *testing.T) {
	longString := string(make([]byte, 1000)) // 1000 chars
	req := types.CreateCarRequest{
		Model:     longString,
		Brand:     "NewBrand",
		CollectAt: "2023-01-03",
		Year:      2022,
		Price:     30000,
		PS:        250,
	}
	w := requestWithAuth(http.MethodPost, "/car", req, "valid")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetCarsWithNoAuthHeader(t *testing.T) {
	r := setupCarTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/cars", nil)
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateCarWithNoAuthHeader(t *testing.T) {
	req := types.CreateCarRequest{
		Model:     "NewModel",
		Brand:     "NewBrand",
		CollectAt: "2023-01-03",
		Year:      2022,
		Price:     30000,
		PS:        250,
	}
	r := setupCarTestRouter()
	data, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/car", bytes.NewReader(data))
	httpReq.Header.Set("Content-Type", "application/json")
	// No Authorization header
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httpReq)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}
