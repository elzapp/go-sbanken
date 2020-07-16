package sbanken

import (
	"encoding/json"
)

const newEfakturas = "https://api.sbanken.no/Bank/api/v1/EFakturas/new"
const efakturas = "https://api.sbanken.no/Bank/api/v1/EFakturas"

// EFaktura as received from the Sbanken public API
type EFaktura struct {
	EFakturaID          string  `json:"eFakturaId"`
	IssuerID            string  `json:"issuerId"`
	EFakturaReference   string  `json:"eFakturaReference"`
	DocumentType        string  `json:"documentType"`
	Status              string  `json:"status"`
	KID                 string  `json:"kid"`
	OriginalDueDate     string  `json:"originalDueDate"`
	OriginalAmount      float64 `json:"originalAmount"`
	MinimumAmount       float64 `json:"minimumAmount"`
	NotificationDate    string  `json:"notificationDate"`
	IssuerName          string  `json:"issuerName"`
	UpdatedDueDate      string  `json:"updatedDueDate"`
	UpdatedAmount       float64 `json:"updatedAmount"`
	CreditAccountNumber string  `json:"creditAccountNumber"`
}

// EFakturaPayRequest is used to accept an eFaktura, charging the AccountID
// on the due date with the maximum or the minimun amount
type EFakturaPayRequest struct {
	EFakturaID           string `json:"eFakturaId"`
	AccountID            string `json:"accountId"`
	PayOnlyMinimumAmount bool   `json:"payOnlyMinimumAmount"`
}

type eFakturaListResponse struct {
	AvailableItems int64      `json:"availableItems"`
	Items          []EFaktura `json:"items"`
	errorInformation
}
type eFakturaItemResponse struct {
	Item EFaktura `json:"item"`
	errorInformation
}

// GetNewEFakturas returns eFakturas that has not been accepted yet
func (conn *APIConnection) GetNewEFakturas() ([]EFaktura, error) {
	r := newAPIRequest()
	r.target = newEfakturas
	var a eFakturaListResponse
	resp, err := conn.makeAPIRequest(r)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(resp, &a)
	return a.Items, nil
}

// GetAllEFakturas returns all pending eFakturas
func (conn *APIConnection) GetAllEFakturas() ([]EFaktura, error) {
	r := newAPIRequest()
	r.target = efakturas
	var a eFakturaListResponse
	resp, err := conn.makeAPIRequest(r)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(resp, &a)
	return a.Items, nil
}

// GetEFaktura returns information on a single EFaktura specified by eFakturaID
func (conn *APIConnection) GetEFaktura(eFakturaID string) EFaktura {
	r := newAPIRequest()
	r.target = efakturas + "/" + eFakturaID
	var a eFakturaItemResponse
	resp, _ := conn.makeAPIRequest(r)
	json.Unmarshal(resp, &a)
	return a.Item
}
