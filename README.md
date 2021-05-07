# go-sbanken
Sbanken API client library in Golang


[![Build Status](https://travis-ci.org/elzapp/go-sbanken.svg?branch=master)](https://travis-ci.org/elzapp/go-sbanken) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=elzapp_go-sbanken&metric=alert_status)](https://sonarcloud.io/dashboard?id=elzapp_go-sbanken) [![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=elzapp_go-sbanken&metric=sqale_index)](https://sonarcloud.io/dashboard?id=elzapp_go-sbanken)

--
```go
import sbanken "github.com/elzapp/go-sbanken"
```

## Usage


## Example

This small program will print your accounts and their balance
```go
package main
import {
	"fmt"
	sbanken "github.com/elzapp/go-sbanken"
}
func main() {
	creds := sbanken.Credentials{"MYAPIKEY","MYSECRET"}
	conn := sbanken.NewAPIConnection(creds)
	accounts := conn.GetAccounts()

	for _, account range accounts {
		fmt.Printf("%-14s %9.6f", account.AccountNumber, account.Balance)
	}
}
```



#### type APIConnection

```go
type APIConnection struct {
}
```

APIConnection is the Api client

#### func  NewAPIConnection

```go
func NewAPIConnection(cred Credentials) APIConnection
```
NewAPIConnection creates an API connection for you This is your starting point,
supply it with a Credentials struct, which you easily can read from a JSON file.

The returned APIConnection struct contains all the methods to communicate with
the public Sbanken API

#### func (*APIConnection) GetAccounts

```go
func (conn *APIConnection) GetAccounts() []Account
```
GetAccounts returns a list of all your bank accounts See the Account struct for
details

#### func (*APIConnection) GetAllEFakturas

```go
func (conn *APIConnection) GetAllEFakturas() []EFaktura
```
GetAllEFakturas returns all pending eFakturas

#### func (*APIConnection) GetEFaktura

```go
func (conn *APIConnection) GetEFaktura(eFakturaID string) EFaktura
```
GetEFaktura returns information on a single EFaktura specified by eFakturaID

#### func (*APIConnection) GetNewEFakturas

```go
func (conn *APIConnection) GetNewEFakturas() []EFaktura
```
GetNewEFakturas returns eFakturas that has not been accepted yet

#### func (*APIConnection) GetTransactions

```go
func (conn *APIConnection) GetTransactions(accountid string) []Transaction
```
GetTransactions returns the latest transactions on a given account using the
default limits set by Sbanken

#### type Account

```go
type Account struct {
	AccountID       string  `json:"accountId"`
	AccountNumber   string  `json:"accountNumber"`
	OwnerCustomerID string  `json:"ownerCustomerId"`
	Name            string  `json:"name"`
	AccountType     string  `json:"accountType"`
	Available       float64 `json:"available"`
	Balance         float64 `json:"balance"`
	CreditLimit     float64 `json:"creditLimit"`
}
```

Account information

#### type Credentials

```go
type Credentials struct {
	Apikey string `json:"apikey"`
	Secret string `json:"secret"`
}
```

Credentials holds login information

#### type EFaktura

```go
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
```

EFaktura as received from the Sbanken public API

#### type EFakturaPayRequest

```go
type EFakturaPayRequest struct {
	EFakturaID           string `json:"eFakturaId"`
	AccountID            string `json:"accountId"`
	PayOnlyMinimumAmount bool   `json:"payOnlyMinimumAmount"`
}
```

EFakturaPayRequest is used to accept an eFaktura, charging the AccountID on the
due date with the maximum or the minimun amount

#### type Transaction

```go
type Transaction struct {
	TransactionID      string  `json:"transactionId"`
	AccountingDate     string  `json:"accountingDate"`
	InterestDate       string  `json:"interestDate"`
	OtherAccountNumber string  `json:"otherAccountNumber"`
	Amount             float64 `json:"amount"`
	Text               string  `json:"text"`
	Source             string  `json:"source"`
}
```

Transaction information
