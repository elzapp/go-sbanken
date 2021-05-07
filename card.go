package sbanken

import (
	"encoding/json"
)

const cards = "https://publicapi.sbanken.no/apibeta/api/v1/Cards"

// Card ...
type Card struct {
	CardID            string `json:"cardId"`
	CardNumber        string `json:"cardNumber"`
	CardVersionNumber string `json:"cardVersionNumber"`
	AccountNumber     string `json:"accountNumber"`
	CustomerID        string `json:"customerId"`
	ExpiryDate        string `json:"expiryDate"`
	AccountOwner      string `json:"accountOwner"`
	Status            string `json:"status"`
	CardType          string `json:"cardType"`
	ProductCode       string `json:"productCode"`
}

type cardListResponse struct {
	AvailableItems int64  `json:"availableItems"`
	Items          []Card `json:"items"`
	errorInformation
}

type cardItemResponse struct {
	Item Card `json:"item"`
	errorInformation
}

// GetCards ...
func (conn *APIConnection) GetCards() ([]Card, error) {
	r := newAPIRequest()
	r.target = cards
	var a cardListResponse
	resp, err := conn.makeAPIRequest(r)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(resp, &a)
	return a.Items, nil
}
