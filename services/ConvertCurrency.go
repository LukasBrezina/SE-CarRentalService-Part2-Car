package services

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type SoapResponse struct {
	Body struct {
		Response struct {
			ConvertedAmount float64 `xml:"converted_amount"`
			TargetCurrency  string  `xml:"target_currency"`
		} `xml:"convert_currency_response"`
	} `xml:"Body"`
}

func ConvertCurrency(currentCurrency string, initialAmount float32, targetCurrency string) (float32, error) {
	url := os.Getenv("CURRENCY_CONVERTER_URL")
	soapBody := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<env:Envelope
    xmlns:env="http://schemas.xmlsoap.org/soap/envelope/"
    xmlns:tns="urn:CurrencyConverter">
  <env:Body>
    <tns:convert_currency>
      <initial_currency>%s</initial_currency>
      <initial_amount>%f</initial_amount>
      <target_currency>%s</target_currency>
    </tns:convert_currency>
  </env:Body>
</env:Envelope>`, currentCurrency, initialAmount*1.0, targetCurrency)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(soapBody))
	if err != nil {
		log.Printf("[ERROR] ConvertCurrency: Failed to create HTTP request: %v", err)
		return 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "\"convert_currency\"")
	req.SetBasicAuth(os.Getenv("BE_USER"), os.Getenv("BE_PASSWORD"))
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("[ERROR] ConvertCurrency: Failed to send HTTP request: %v", err)
		return 0, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("[ERROR] ConvertCurrency: Failed to close response body: %v", closeErr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] ConvertCurrency: Failed to read response body: %v", err)
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var result SoapResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		log.Printf("[ERROR] ConvertCurrency: Failed to parse XML response: %v, body: %s", err, string(body))
		return 0, fmt.Errorf("failed to parse XML response: %w", err)
	}

	convertedAmount := float32(result.Body.Response.ConvertedAmount)
	log.Printf("[INFO] ConvertCurrency: %s→%s %.2f = %.2f %s", currentCurrency, targetCurrency, initialAmount, convertedAmount, result.Body.Response.TargetCurrency)

	return convertedAmount, nil
}
