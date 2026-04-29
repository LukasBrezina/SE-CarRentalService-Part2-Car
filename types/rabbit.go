package types

type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Currency string `json:"currency"`
	IsAdmin  bool   `json:"is_admin"`
}

type TokenRequest struct {
	Token string `json:"token"`
}

type TokenResponse struct {
	Account Account `json:"account"`
	Valid   bool    `json:"valid"`
	Error   string  `json:"error,omitempty"`
}

type CurrencyRequest struct {
	Base string `json:"base"`
}

type CurrencyResponse struct {
	Currencies map[string]float64 `json:"currencies"`
	Error      string             `json:"error,omitempty"`
}
