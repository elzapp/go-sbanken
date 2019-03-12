package sbanken

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const identityserver = "https://auth.sbanken.no/IdentityServer/connect/token"
const apiAccounts = "https://api.sbanken.no/Bank/api/v1/Accounts"
const apiTransactions = "https://api.sbanken.no/Bank/api/v1/Transactions/%s"

type accounts struct {
	Accounts []Account `json:"items"`
	errorInformation
}

// Account information
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

type errorInformation struct {
	IsError      bool   `json:"isError"`
	ErrorType    string `json:"errorType"`
	ErrorMessage string `json:"errorMessage"`
	TraceID      string `json:"traceId"`
}

// Credentials holds login information
type Credentials struct {
	Apikey string `json:"apikey"`
	Secret string `json:"secret"`
	UserID string `json:"userid"`
}

type tokenResponse struct {
	Token string `json:"access_token"`
}

// Transaction information
type Transaction struct {
	TransactionID      string  `json:"transactionId"`
	AccountingDate     string  `json:"accountingDate"`
	InterestDate       string  `json:"interestDate"`
	OtherAccountNumber string  `json:"otherAccountNumber"`
	Amount             float64 `json:"amount"`
	Text               string  `json:"text"`
	Source             string  `json:"source"`
}

type transactions struct {
	Transactions []Transaction `json:"items"`
	errorInformation
}

// APIConnection is the Api client
type APIConnection struct {
	cred  Credentials
	token string
}

func (conn *APIConnection) getToken() string {
	if conn.token == "" {
		postdata := url.Values{}
		postdata.Add("grant_type", "client_credentials")
		req, _ := http.NewRequest("POST", identityserver, strings.NewReader(postdata.Encode()))
		req.Header.Add("Content-type", "application/x-www-form-urlencoded; charset=utf-8")
		req.SetBasicAuth(conn.cred.Apikey, conn.cred.Secret)
		cli := &http.Client{}
		resp, _ := cli.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		var t tokenResponse
		json.Unmarshal(body, &t)
		conn.token = t.Token
	}
	return conn.token
}

type apirequest struct {
	target  string
	params  map[string]string
	headers map[string]string
}

func newAPIRequest() apirequest {
	var r apirequest
	r.params = map[string]string{}
	r.headers = map[string]string{}
	return r
}

func (conn *APIConnection) makeAPIRequest(r apirequest) []byte {
	req, _ := http.NewRequest("GET", r.target, nil)
	req.Header.Add("Authorization", "Bearer "+conn.getToken())

	req.Header.Add("customerId", conn.cred.UserID)
	for key, value := range r.headers {
		req.Header.Add(key, value)
	}
	cli := &http.Client{}
	resp, _ := cli.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return body
}

func (conn *APIConnection) GetAccounts() []Account {
	r := newAPIRequest()
	r.target = apiAccounts
	var a accounts
	json.Unmarshal(conn.makeAPIRequest(r), &a)
	return a.Accounts
}

func (conn *APIConnection) GetTransactions(accountid string) []Transaction {
	r := newAPIRequest()
	r.target = fmt.Sprintf(apiTransactions, accountid)
	var t transactions
	json.Unmarshal(conn.makeAPIRequest(r), &t)
	return t.Transactions
}

func NewApiConnection(cred Credentials) APIConnection {
	var conn APIConnection
	conn.cred = cred
	return conn
}
