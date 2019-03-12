package sbanken

import "testing"

func TestGetAccounts(t *testing.T) {
	var cred Credentials
	conn := NewApiConnection(cred)
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
	conn := NewApiConnection(cred)
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
		t.Errorf("Expected number of returned accounts to be 2019-03-12T20:15:12.477Z, got %s", transactions[0].AccountingDate)
	}
}
