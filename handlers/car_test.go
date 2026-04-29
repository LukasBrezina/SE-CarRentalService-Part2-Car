// handlers/car_handler_test.go
package handlers

import (
	"SE-CarRentalService/types"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupCarTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := NewCarHandler()

	r.GET("/cars", h.GetCars)
	r.GET("/car/:id", h.GetCar)
	r.POST("/car", h.CreateCar)
	r.PUT("/car/:id", h.UpdateCar)
	r.DELETE("/car/:id", h.DeleteCar)

	return r
}

func request(method string, path string, body any) *httptest.ResponseRecorder {
	r := setupCarTestRouter()

	var buffer bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buffer).Encode(body)
	}

	req := httptest.NewRequest(method, path, &buffer)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestGetCarInvalidID(t *testing.T) {
	w := request(http.MethodGet, "/car/abc", nil)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCarInvalidID(t *testing.T) {
	w := request(http.MethodPut, "/car/abc", types.Car{})

	if w.Code != http.StatusBadRequest && w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 400 or 401, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateCarWithoutToken(t *testing.T) {
	w := request(http.MethodPost, "/car", types.CreateCarRequest{})

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCarWithoutToken(t *testing.T) {
	w := request(http.MethodPut, "/car/1", types.Car{})

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteCarWithoutToken(t *testing.T) {
	w := request(http.MethodDelete, "/car/1", nil)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteCarInvalidIDShouldReturn400(t *testing.T) {
	w := request(http.MethodDelete, "/car/abc", nil)

	if w.Code != http.StatusUnauthorized && w.Code != http.StatusBadRequest {
		t.Fatalf("expected 401 or 400, got %d, body: %s", w.Code, w.Body.String())
	}
}
