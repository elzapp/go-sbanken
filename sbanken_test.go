package sbanken

import (
	"fmt"
	"testing"
	"time"
)

func TestGetAccounts(t *testing.T) {
	var cred Credentials
	conn := NewAPIConnection(cred)
	conn.makeAPIRequest = func(r apirequest) ([]byte, error) {
		return []byte(`{
			"availableItems": 1,
			"items": [
			  {
				"accountId": "972219XXXXX",
				"accountNumber": "972219YYYYY",
				"ownerCustomerId": "01017012345",
				"name": "checking account",
				"accountType": "string",
				"available": 12345,
				"balance": 100,
				"creditLimit": 12245
			  }
			],
			"errorType": "System",
			"isError": false,
			"errorMessage": "string",
			"traceId": "string"
		  }`), nil
	}
	accounts, _ := conn.GetAccounts()
	if len(accounts) != 1 {
		t.Errorf("Expected number of returned accounts to be 1, got %d", len(accounts))
	}
	if accounts[0].AccountID != "972219XXXXX" {
		t.Errorf("Expected number of returned accounts to be 1, got %s", accounts[0].AccountID)
	}
}

func TestGetTransactions(t *testing.T) {
	var cred Credentials
	conn := NewAPIConnection(cred)
	conn.makeAPIRequest = func(r apirequest) ([]byte, error) {
		return []byte(`{
			"availableItems": 2,
			"items": [
			  {
				"accountingDate": "2019-10-14T00:00:00",
				"interestDate": "2019-10-15T00:00:00",
				"otherAccountNumberSpecified": false,
				"amount": -58.000,
				"text": "*3100 13.10 NOK 58.00 PAYPAL INC Kurs: 1.0000",
				"transactionType": "VISA VARE",
				"transactionTypeCode": 714,
				"transactionTypeText": "VISA VARE",
				"isReservation": false,
				"reservationType": null,
				"source": "Archive",
				"cardDetailsSpecified": true,
				"cardDetails": {
					"cardNumber": "*3100",
					"currencyAmount": 58.000,
					"currencyRate": 1.00000,
					"merchantCategoryCode": "5941",
					"merchantCategoryDescription": "Sportsutstyr",
					"merchantCity": "4029357733",
					"merchantName": "PAYPAL *STRAVA INC",
					"originalCurrencyCode": "NOK",
					"purchaseDate": "2019-10-13T00:00:00",
					"transactionId": "5892862195085990"
				},
				"transactionDetailSpecified": false
			},
			{
				"accountingDate": "2019-10-23T00:00:00",
				"interestDate": "2019-10-23T00:00:00",
				"otherAccountNumberSpecified": false,
				"amount": -16.410,
				"text": "23.10 REMA SPECTRUM FOLKE BERNAD FYLLINGSDALEN",
				"transactionType": "VARER",
				"transactionTypeCode": 710,
				"transactionTypeText": "VARER",
				"isReservation": false,
				"reservationType": null,
				"source": "Archive",
				"cardDetailsSpecified": false,
				"transactionDetailSpecified": false
			}
			],
			"errorType": "System",
			"isError": true,
			"errorMessage": "string",
			"traceId": "string"
		  }`), nil
	}
	transactions, _ := conn.GetTransactions("972219XXXXX")
	if len(transactions) != 2 {
		t.Errorf("Expected number of returned transactions to be 2, got %d", len(transactions))
	}
	if transactions[0].AccountingDate != "2019-10-14T00:00:00" {
		t.Errorf("Expected accounting date to be 2019-10-14T00:00:00, got %s", transactions[0].AccountingDate)
	}
	if transactions[0].GetAccountingDate().Unix() != 1571011200 {
		t.Errorf("Expected accounting date to be %d, got %d", 1571011200, transactions[0].GetAccountingDate().Unix())
	}
	if transactions[0].GetInterestDate().Unix() != 1571097600 {
		t.Errorf("Expected interest date to be %d, got %d", 1571097600, transactions[0].GetInterestDate().Unix())
	}
	if transactions[0].CardDetails.CardNumber != "*3100" {
		t.Errorf("Expected cardNumber to be %s, but got %s", "*3100", transactions[0].CardDetails.CardNumber)
	}
	if transactions[1].Source != "Archive" {
		t.Errorf("Expected Source on second tx to be %s, but got %s", "Archive", transactions[0].Source)
	}
}

const efakturaList = `{
	"availableItems": 0,
	"items": [
	  {
		"eFakturaId": "XYZXYZ",
		"issuerId": "string",
		"eFakturaReference": "string",
		"documentType": "string",
		"status": "string",
		"kid": "string",
		"originalDueDate": "2019-03-12T20:30:21.730Z",
		"originalAmount": 0,
		"minimumAmount": 0,
		"notificationDate": "2019-03-12T20:30:21.730Z",
		"issuerName": "string"
	  }
	],
	"errorType": "System",
	"isError": true,
	"errorMessage": "string",
	"traceId": "string"
  }`
const singleEFaktura = `{
	"item": {
	  "eFakturaId": "XXXYZXYZ",
	  "issuerId": "a",
	  "eFakturaReference": "b",
	  "documentType": "c",
	  "status": "d",
	  "kid": "e",
	  "originalDueDate": "2019-03-12T20:44:24.879Z",
	  "originalAmount": 100,
	  "minimumAmount": 10,
	  "updatedDueDate": "2019-03-12T20:44:24.879Z",
	  "updatedAmount": 30,
	  "notificationDate": "2019-03-12T20:44:24.879Z",
	  "creditAccountNumber": "x",
	  "issuerName": "Telenor"
	},
	"errorType": "System",
	"isError": true,
	"errorMessage": "string",
	"traceId": "string"
  }`

func TestGetNewEfakturas(t *testing.T) {
	var cred Credentials
	conn := NewAPIConnection(cred)
	conn.makeAPIRequest = func(r apirequest) ([]byte, error) {
		return []byte(efakturaList), nil
	}
	efakturas, _ := conn.GetNewEFakturas()
	if len(efakturas) != 1 {
		t.Errorf("Expected number of returned transactions to be 1, got %d", len(efakturas))
	}

	if efakturas[0].EFakturaID != "XYZXYZ" {
		t.Errorf("Expected efaktura id to be XYZXYZ, got %s", efakturas[0].EFakturaID)
	}
}
func TestGetAllEfakturas(t *testing.T) {
	var cred Credentials
	conn := NewAPIConnection(cred)
	conn.makeAPIRequest = func(r apirequest) ([]byte, error) {
		if r.target != "https://publicapi.sbanken.no/apibeta/api/v1/EFakturas" {
			t.Errorf("GetAllEFakturas is calling wrong endpoint: %s", r.target)
		}
		return []byte(efakturaList), nil
	}
	efakturas, _ := conn.GetAllEFakturas()
	if len(efakturas) != 1 {
		t.Errorf("Expected number of returned transactions to be 1, got %d", len(efakturas))
	}

	if efakturas[0].EFakturaID != "XYZXYZ" {
		t.Errorf("Expected efaktura id to be XYZXYZ, got %s", efakturas[0].EFakturaID)
	}
}

func TestGetSingleEFaktura(t *testing.T) {
	var cred Credentials
	conn := NewAPIConnection(cred)
	conn.makeAPIRequest = func(r apirequest) ([]byte, error) {
		if r.target != "https://publicapi.sbanken.no/apibeta/api/v1/EFakturas/XYZXYZ" {
			t.Errorf("GetEfaktura is calling wrong endpoint: %s", r.target)
		}
		return []byte(singleEFaktura), nil
	}
	efaktura := conn.GetEFaktura("XYZXYZ")

	if efaktura.EFakturaID != "XXXYZXYZ" {
		t.Errorf("Expected efaktura id to be XXXYZXYZ, got %s", efaktura.EFakturaID)
	}
	if efaktura.IssuerName != "Telenor" {
		t.Errorf("Expected issuer name to be Telenor, got %s", efaktura.IssuerName)
	}
}

func TestGetPurchaseDateWithFallback(t *testing.T) {
	tx := Transaction{AccountingDate: "2020-03-01T00:00:00"}
	expect := "2020-03-01"
	got := tx.GetTransactionDate().Format(time.RFC3339)[0:10]
	if got != expect {
		fmt.Printf("Using fallback, got %s, expected %s\n", got, expect)
		t.Fail()
	}
}
func TestGetPurchaseDateFromCardDetails(t *testing.T) {
	tx := Transaction{AccountingDate: "2020-03-01T00:00:00"}

	tx.CardDetails.PurchaseDate = "2020-03-05T00:00:00"
	expect := "2020-03-05"
	got := tx.GetTransactionDate().Format(time.RFC3339)[0:10]
	if got != expect {
		fmt.Printf("Using CardDetails, got %s, expected %s\n", got, expect)
		t.Fail()
	}
}

func TestGetPurchaseDateFromText(t *testing.T) {
	tx := Transaction{
		AccountingDate: "2020-03-01T00:00:00",
		Text:           "28.02 REMA KALMARHUSE JON SMØRSGT  BERGEN",
	}
	expect := "2020-02-28"
	got := tx.GetTransactionDate().Format(time.RFC3339)[0:10]
	if got != expect {
		fmt.Printf("Using Text, got %s, expected %s\n", got, expect)
		t.Fail()
	}
}

func TestGetPurchaseDateFromTextNewyear(t *testing.T) {
	tx := Transaction{
		AccountingDate: "2020-01-01T00:00:00",
		Text:           "31.12 REMA KALMARHUSE JON SMØRSGT  BERGEN",
	}
	expect := "2019-12-31"
	got := tx.GetTransactionDate().Format(time.RFC3339)[0:10]
	if got != expect {
		fmt.Printf("Using Text prev year, got %s, expected %s\n", got, expect)
		t.Fail()
	}
}
func TestGetPurchaseDateFromTextSameday(t *testing.T) {
	tx := Transaction{
		AccountingDate: "2020-01-01T00:00:00",
		Text:           "01.01 REMA KALMARHUSE JON SMØRSGT  BERGEN",
	}
	expect := "2020-01-01"
	got := tx.GetTransactionDate().Format(time.RFC3339)[0:10]
	if got != expect {
		fmt.Printf("Using Text on same day, got %s, expected %s\n", got, expect)
		t.Fail()
	}
}

func TestGetPurchaseDateFromTextCreditCardNoDetails(t *testing.T) {
	tx := Transaction{
		AccountingDate: "2021-03-23T00:00:00",
		Text:           "*1234 22.03 NOK 49.30 EXTRA NESTTUN 837625 KURS: 1.0000",
	}
	expect := "2021-03-22"
	got := tx.GetTransactionDate().Format(time.RFC3339)[0:10]
	if got != expect {
		fmt.Printf("Using Text on same day, got %s, expected %s\n", got, expect)
		t.Fail()
	}
}

func TestGetText(t *testing.T) {
	tx := Transaction{
		Text: "01.01 REMA KALMARHUSE JON SMØRSGT  BERGEN",
	}
	expect := "REMA KALMARHUSE JON SMØRSGT  BERGEN"
	got := tx.GetText()
	if got != expect {
		fmt.Printf("Got %s, expected %s\n", got, expect)
		t.Fail()
	}
}

func TestGetTextPayment(t *testing.T) {
	tx := Transaction{
		Text: "Til: BONNIER PUBLICA Betalt: 17.03.21",
	}
	expect := "BONNIER PUBLICA"
	got := tx.GetText()
	if got != expect {
		fmt.Printf("Got %s, expected %s\n", got, expect)
		t.Fail()
	}
}

func TestGetTextNettgiro(t *testing.T) {
	tx := Transaction{
		Text: "Nettgiro til: BONNIER PUBLICA Betalt: 12.03.21",
	}
	expect := "BONNIER PUBLICA"
	got := tx.GetText()
	if got != expect {
		fmt.Printf("Got %s, expected %s\n", got, expect)
		t.Fail()
	}
}

func TestGetTextCreditCard(t *testing.T) {
	tx := Transaction{
		Text:        "*1234 15.03 NOK 67.50 BUNNPRIS SLETTE Kurs: 1.0000",
		CardDetails: cardDetails{MerchantName: "BUNNPRIS SLETTE", MerchantCity: "BERGEN"},
	}
	expect := "BUNNPRIS SLETTE, BERGEN"
	got := tx.GetText()
	if got != expect {
		fmt.Printf("Got %s, expected %s\n", got, expect)
		t.Fail()
	}
}

func TestGetTextCreditCardWithoutCardDetails(t *testing.T) {
	tx := Transaction{
		Text: "*1234 22.03 NOK 49.30 EXTRA NESTTUN 837625 KURS: 1.0000",
	}
	expect := "EXTRA NESTTUN 837625"
	got := tx.GetText()
	if got != expect {
		fmt.Printf("Got %s, expected %s\n", got, expect)
		t.Fail()
	}
}
