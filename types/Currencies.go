package types

import "fmt"

type Currency string

const (
	USD Currency = "USD"
	JPY Currency = "JPY"
	CZK Currency = "CZK"
	DKK Currency = "DKK"
	GBP Currency = "GBP"
	HUF Currency = "HUF"
	PLN Currency = "PLN"
	RON Currency = "RON"
	SEK Currency = "SEK"
	CHF Currency = "CHF"
	ISK Currency = "ISK"
	NOK Currency = "NOK"
	TRY Currency = "TRY"
	AUD Currency = "AUD"
	BRL Currency = "BRL"
	CAD Currency = "CAD"
	CNY Currency = "CNY"
	HKD Currency = "HKD"
	IDR Currency = "IDR"
	ILS Currency = "ILS"
	INR Currency = "INR"
	KRW Currency = "KRW"
	MXN Currency = "MXN"
	MYR Currency = "MYR"
	NZD Currency = "NZD"
	PHP Currency = "PHP"
	SGD Currency = "SGD"
	THB Currency = "THB"
	ZAR Currency = "ZAR"
	EUR Currency = "EUR"
)

func (c Currency) IsValid() bool {
	switch c {
	case EUR, USD, JPY, CZK, DKK, GBP, HUF, PLN, RON, SEK, CHF,
		ISK, NOK, TRY, AUD, BRL, CAD, CNY, HKD, IDR, ILS,
		INR, KRW, MXN, MYR, NZD, PHP, SGD, THB, ZAR:
		return true
	default:
		return false
	}
}

func ParseCurrency(s string) (Currency, error) {
	c := Currency(s)

	switch c {
	case EUR, USD, JPY, CZK, DKK, GBP, HUF, PLN, RON, SEK, CHF,
		ISK, NOK, TRY, AUD, BRL, CAD, CNY, HKD, IDR, ILS,
		INR, KRW, MXN, MYR, NZD, PHP, SGD, THB, ZAR:
		return c, nil
	default:
		return "", fmt.Errorf("invalid currency: %s", s)
	}
}

func (c Currency) String() string {
	return string(c)
}
