package sbanken

import "encoding/json"

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
	UpdatedAmount       string  `json:"updatedAmount"`
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

func (conn *APIConnection) GetNewEfakturas() []EFaktura {
	r := newAPIRequest()
	r.target = newEfakturas
	var a eFakturaListResponse
	json.Unmarshal(conn.makeAPIRequest(r), &a)
	return a.Items
}

func (conn *APIConnection) GetAllEfakturas() []EFaktura {
	r := newAPIRequest()
	r.target = efakturas
	var a eFakturaListResponse
	json.Unmarshal(conn.makeAPIRequest(r), &a)
	return a.Items
}

func (conn *APIConnection) GetEfaktura(efakturaId string) EFaktura {
	r := newAPIRequest()
	r.target = efakturas + "/" + efakturaId
	var a eFakturaItemResponse
	json.Unmarshal(conn.makeAPIRequest(r), &a)
	return a.Item
}
