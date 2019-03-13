package sbanken

import (
	"testing"
)

func TestGetAccounts(t *testing.T) {
	var cred Credentials
	conn := NewAPIConnection(cred)
	conn.makeAPIRequest = func(r apirequest) []byte {
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
		  }`)
	}
	accounts := conn.GetAccounts()
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
	conn.makeAPIRequest = func(r apirequest) []byte {
		return []byte(`{
			"availableItems": 1,
			"items": [
			  {
				"accountingDate": "2019-03-12T20:15:12.477Z",
				"interestDate": "2019-03-12T20:15:12.477Z",
				"otherAccountNumber": "string",
				"otherAccountNumberSpecified": true,
				"amount": 0,
				"text": "string",
				"transactionType": "string",
				"transactionTypeCode": 0,
				"transactionTypeText": "string",
				"isReservation": true,
				"reservationType": "NotReservation",
				"source": "AccountStatement",
				"cardDetails": {
				  "cardNumber": "string",
				  "currencyAmount": 0,
				  "currencyRate": 0,
				  "merchantCategoryCode": "string",
				  "merchantCategoryDescription": "string",
				  "merchantCity": "string",
				  "merchantName": "string",
				  "originalCurrencyCode": "string",
				  "purchaseDate": "2019-03-12T20:15:12.477Z",
				  "transactionId": "string"
				},
				"cardDetailsSpecified": true,
				"hasTransactionDetail": true,
				"transactionDetail": {
				  "formattedAccountNumber": "string",
				  "transactionId": 0,
				  "cid": "string",
				  "amountDescription": "string",
				  "receiverName": "string",
				  "numericReference": 0,
				  "payerName": "string",
				  "registrationDate": "2019-03-12T20:15:12.477Z"
				}
			  }
			],
			"errorType": "System",
			"isError": true,
			"errorMessage": "string",
			"traceId": "string"
		  }`)
	}
	transactions := conn.GetTransactions("972219XXXXX")
	if len(transactions) != 1 {
		t.Errorf("Expected number of returned transactions to be 1, got %d", len(transactions))
	}

	if transactions[0].AccountingDate != "2019-03-12T20:15:12.477Z" {
		t.Errorf("Expected accounting date to be 2019-03-12T20:15:12.477Z, got %s", transactions[0].AccountingDate)
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
	conn.makeAPIRequest = func(r apirequest) []byte {
		return []byte(efakturaList)
	}
	efakturas := conn.GetNewEFakturas()
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
	conn.makeAPIRequest = func(r apirequest) []byte {
		if r.target != "https://api.sbanken.no/Bank/api/v1/EFakturas" {
			t.Errorf("GetAllEFakturas is calling wrong endpoint: %s", r.target)
		}
		return []byte(efakturaList)
	}
	efakturas := conn.GetAllEFakturas()
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
	conn.makeAPIRequest = func(r apirequest) []byte {
		if r.target != "https://api.sbanken.no/Bank/api/v1/EFakturas/XYZXYZ" {
			t.Errorf("GetEfaktura is calling wrong endpoint: %s", r.target)
		}
		return []byte(singleEFaktura)
	}
	efaktura := conn.GetEFaktura("XYZXYZ")

	if efaktura.EFakturaID != "XXXYZXYZ" {
		t.Errorf("Expected efaktura id to be XXXYZXYZ, got %s", efaktura.EFakturaID)
	}
	if efaktura.IssuerName != "Telenor" {
		t.Errorf("Expected issuer name to be Telenor, got %s", efaktura.IssuerName)
	}
}
