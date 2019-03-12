package sbanken

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
