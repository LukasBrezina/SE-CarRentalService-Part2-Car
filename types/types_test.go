package types

import (
	"encoding/json"
	"testing"
)

func TestCarJSONMarshalUnmarshal(t *testing.T) {
	car := Car{
		ID:        1,
		Model:     "TestModel",
		Brand:     "TestBrand",
		CollectAt: "2023-01-01",
		AccountId: 0,
		Year:      2020,
		Price:     10000.5,
		PS:        150,
	}

	data, err := json.Marshal(car)
	if err != nil {
		t.Fatalf("failed to marshal Car: %v", err)
	}

	var unmarshaled Car
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal Car: %v", err)
	}

	if unmarshaled != car {
		t.Fatalf("expected %+v, got %+v", car, unmarshaled)
	}
}

func TestCreateCarRequestJSONMarshalUnmarshal(t *testing.T) {
	req := CreateCarRequest{
		Model:     "TestModel",
		Brand:     "TestBrand",
		CollectAt: "2023-01-01",
		AccountId: 0,
		Year:      2020,
		Price:     10000.5,
		PS:        150,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal CreateCarRequest: %v", err)
	}

	var unmarshaled CreateCarRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal CreateCarRequest: %v", err)
	}

	if unmarshaled != req {
		t.Fatalf("expected %+v, got %+v", req, unmarshaled)
	}
}

func TestAccountJSONMarshalUnmarshal(t *testing.T) {
	account := Account{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
		Currency: "USD",
		IsAdmin:  true,
	}

	data, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("failed to marshal Account: %v", err)
	}

	var unmarshaled Account
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal Account: %v", err)
	}

	if unmarshaled != account {
		t.Fatalf("expected %+v, got %+v", account, unmarshaled)
	}
}

func TestTokenRequestJSONMarshalUnmarshal(t *testing.T) {
	req := TokenRequest{
		Token: "testtoken",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal TokenRequest: %v", err)
	}

	var unmarshaled TokenRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal TokenRequest: %v", err)
	}

	if unmarshaled != req {
		t.Fatalf("expected %+v, got %+v", req, unmarshaled)
	}
}

func TestTokenResponseJSONMarshalUnmarshal(t *testing.T) {
	resp := TokenResponse{
		Account: Account{ID: 1, Username: "test", Currency: "USD", IsAdmin: true},
		Valid:   true,
		Error:   "",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal TokenResponse: %v", err)
	}

	var unmarshaled TokenResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal TokenResponse: %v", err)
	}

	if unmarshaled.Account != resp.Account || unmarshaled.Valid != resp.Valid || unmarshaled.Error != resp.Error {
		t.Fatalf("expected %+v, got %+v", resp, unmarshaled)
	}
}

func TestCurrencyRequestJSONMarshalUnmarshal(t *testing.T) {
	req := CurrencyRequest{
		Base: "USD",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal CurrencyRequest: %v", err)
	}

	var unmarshaled CurrencyRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal CurrencyRequest: %v", err)
	}

	if unmarshaled != req {
		t.Fatalf("expected %+v, got %+v", req, unmarshaled)
	}
}

func TestCurrencyResponseJSONMarshalUnmarshal(t *testing.T) {
	resp := CurrencyResponse{
		Currencies: map[string]float64{"EUR": 0.85, "GBP": 0.75},
		Error:      "",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal CurrencyResponse: %v", err)
	}

	var unmarshaled CurrencyResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal CurrencyResponse: %v", err)
	}

	if len(unmarshaled.Currencies) != len(resp.Currencies) || unmarshaled.Error != resp.Error {
		t.Fatalf("expected %+v, got %+v", resp, unmarshaled)
	}
	for k, v := range resp.Currencies {
		if unmarshaled.Currencies[k] != v {
			t.Fatalf("currency mismatch for %s", k)
		}
	}
}

func TestCarJSONMarshalUnmarshalWithSpecialChars(t *testing.T) {
	car := Car{
		ID:        1,
		Model:     "Model with spaces & symbols !@#",
		Brand:     "Brand\nwith\tnewlines",
		CollectAt: "2023-01-01",
		AccountId: 0,
		Year:      2020,
		Price:     10000.5,
		PS:        150,
	}

	data, err := json.Marshal(car)
	if err != nil {
		t.Fatalf("failed to marshal Car: %v", err)
	}

	var unmarshaled Car
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal Car: %v", err)
	}

	if unmarshaled != car {
		t.Fatalf("expected %+v, got %+v", car, unmarshaled)
	}
}

func TestCreateCarRequestJSONMarshalUnmarshalWithEmptyStrings(t *testing.T) {
	req := CreateCarRequest{
		Model:     "",
		Brand:     "",
		CollectAt: "",
		AccountId: 0,
		Year:      0,
		Price:     0,
		PS:        0,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal CreateCarRequest: %v", err)
	}

	var unmarshaled CreateCarRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal CreateCarRequest: %v", err)
	}

	if unmarshaled != req {
		t.Fatalf("expected %+v, got %+v", req, unmarshaled)
	}
}

func TestAccountJSONMarshalUnmarshalWithLongStrings(t *testing.T) {
	longString := string(make([]byte, 500))
	account := Account{
		ID:       1,
		Username: longString,
		Email:    longString + "@example.com",
		Password: longString,
		Currency: "USD",
		IsAdmin:  true,
	}

	data, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("failed to marshal Account: %v", err)
	}

	var unmarshaled Account
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal Account: %v", err)
	}

	if unmarshaled != account {
		t.Fatalf("expected %+v, got %+v", account, unmarshaled)
	}
}

func TestTokenResponseJSONMarshalUnmarshalWithError(t *testing.T) {
	resp := TokenResponse{
		Account: Account{},
		Valid:   false,
		Error:   "some error message with spaces and !@#",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal TokenResponse: %v", err)
	}

	var unmarshaled TokenResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal TokenResponse: %v", err)
	}

	if unmarshaled.Account != resp.Account || unmarshaled.Valid != resp.Valid || unmarshaled.Error != resp.Error {
		t.Fatalf("expected %+v, got %+v", resp, unmarshaled)
	}
}

func TestCurrencyRequestJSONMarshalUnmarshalWithInvalidBase(t *testing.T) {
	req := CurrencyRequest{
		Base: "INVALID",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal CurrencyRequest: %v", err)
	}

	var unmarshaled CurrencyRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal CurrencyRequest: %v", err)
	}

	if unmarshaled != req {
		t.Fatalf("expected %+v, got %+v", req, unmarshaled)
	}
}

func TestCurrencyResponseJSONMarshalUnmarshalWithEmptyMap(t *testing.T) {
	resp := CurrencyResponse{
		Currencies: map[string]float64{},
		Error:      "",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal CurrencyResponse: %v", err)
	}

	var unmarshaled CurrencyResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal CurrencyResponse: %v", err)
	}

	if len(unmarshaled.Currencies) != 0 || unmarshaled.Error != resp.Error {
		t.Fatalf("expected %+v, got %+v", resp, unmarshaled)
	}
}

func TestCurrencyResponseJSONMarshalUnmarshalWithNegativeRates(t *testing.T) {
	resp := CurrencyResponse{
		Currencies: map[string]float64{"EUR": -0.85, "GBP": -0.75},
		Error:      "negative rates error",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal CurrencyResponse: %v", err)
	}

	var unmarshaled CurrencyResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal CurrencyResponse: %v", err)
	}

	if len(unmarshaled.Currencies) != len(resp.Currencies) || unmarshaled.Error != resp.Error {
		t.Fatalf("expected %+v, got %+v", resp, unmarshaled)
	}
	for k, v := range resp.Currencies {
		if unmarshaled.Currencies[k] != v {
			t.Fatalf("currency mismatch for %s", k)
		}
	}
}

func TestParseCurrencyInvalid(t *testing.T) {
	_, err := ParseCurrency("INVALID")
	if err == nil {
		t.Fatalf("expected error for invalid currency")
	}
}

func TestCurrencyString(t *testing.T) {
	c := Currency("USD")
	if c.String() != "USD" {
		t.Fatalf("expected USD, got %s", c.String())
	}
}

func TestCurrencyIsValid(t *testing.T) {
	if !Currency("USD").IsValid() {
		t.Fatalf("USD should be valid")
	}
	if Currency("INVALID").IsValid() {
		t.Fatalf("INVALID should not be valid")
	}
}
