package sbanken

import (
	"encoding/json"
)

const payments = `https://publicapi.sbanken.no/apibeta/api/v1/Payments/`

// Payment ...
type Payment struct {
	ID                     string   `json:"paymentId"`
	RecipientAccountNumber string   `json:"recipientAccountNumber"`
	Amount                 float64  `json:"amount"`
	DueDate                string   `json:"dueDate"`
	KID                    string   `json:"kid"`
	Text                   string   `json:"text"`
	IsActive               bool     `json:"isActive"`
	Status                 string   `json:"status"`
	AllowedNewStatusTypes  []string `json:"allowedNewStatusTypes"`
	StatusDetails          string   `json:"statusDetails"`
	ProductType            string   `json:"productType"`
	PaymentType            string   `json:"paymentType"`
	PaymentNumber          int64    `json:"paymentNumber"`
	BeneficiaryName        string   `json:"beneficiaryName"`
}

type paymentListResponse struct {
	AvailableItems int64     `json:"availableItems"`
	Items          []Payment `json:"items"`
	errorInformation
}

type parymentItemResponse struct {
	Item Payment `json:"item"`
	errorInformation
}

// GetPayments ...
func (conn *APIConnection) GetPayments(accountID string) ([]Payment, error) {
	r := newAPIRequest()
	r.target = payments + accountID
	var a paymentListResponse
	resp, err := conn.makeAPIRequest(r)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(resp, &a)
	return a.Items, nil
}
