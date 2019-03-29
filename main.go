package sbanken

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const dateFormat = "2006-01-02T15:04:05-07:00" //2019-03-06T00:00:00+01:00
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

// GetInterestDate returns the interest date as a Time struct
func (t *Transaction) GetInterestDate() time.Time {
	r, _ := time.Parse(dateFormat, t.InterestDate)
	return r
}

// GetAccountingDate returns the interest date as a Time struct
func (t *Transaction) GetAccountingDate() time.Time {
	r, _ := time.Parse(dateFormat, t.AccountingDate)
	return r
}

type transactions struct {
	Transactions []Transaction `json:"items"`
	errorInformation
}

// APIConnection is the Api client
type APIConnection struct {
	cred           Credentials
	token          string
	makeAPIRequest func(r apirequest) []byte
}

func (conn *APIConnection) HasToken() bool {
	if conn.token == "" {
		return false
	}
	return true
}

func (conn *APIConnection) getToken() string {
	if conn.token == "" {
		postdata := url.Values{}
		postdata.Add("grant_type", "client_credentials")
		req, _ := http.NewRequest("POST", identityserver, strings.NewReader(postdata.Encode()))
		req.Header.Add("Content-type", "application/x-www-form-urlencoded; charset=utf-8")
		req.Header.Add("User-Agent", "github.com/elzapp/go-sbanken")
		req.SetBasicAuth(conn.cred.Apikey, conn.cred.Secret)
		cli := &http.Client{}
		resp, err := cli.Do(req)
		if err != nil {
			fmt.Printf("%+v", err)
			return ""
		}
		fmt.Printf("%+v", resp)
		defer resp.Body.Close()
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

// GetAccounts returns a list of all your bank accounts
// See the Account struct for details
func (conn *APIConnection) GetAccounts() []Account {
	r := newAPIRequest()
	r.target = apiAccounts
	var a accounts
	json.Unmarshal(conn.makeAPIRequest(r), &a)
	return a.Accounts
}

// GetTransactions returns the latest transactions on a given account
// using the default limits set by Sbanken
func (conn *APIConnection) GetTransactions(accountid string) []Transaction {
	r := newAPIRequest()
	r.target = fmt.Sprintf(apiTransactions, accountid)
	var t transactions
	json.Unmarshal(conn.makeAPIRequest(r), &t)
	return t.Transactions
}

// NewAPIConnection creates an API connection for you
// This is your starting point, supply it with a
// Credentials struct, which you easily can read from a
// JSON file.
//
// The returned APIConnection struct contains all the
// methods to communicate with the public Sbanken API
func NewAPIConnection(cred Credentials) APIConnection {
	var conn APIConnection
	conn.cred = cred
	conn.makeAPIRequest = func(r apirequest) []byte {
		token := conn.getToken()
		req, _ := http.NewRequest("GET", r.target, nil)
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("User-Agent", "github.com/elzapp/go-sbanken")

		req.Header.Add("customerId", conn.cred.UserID)
		for key, value := range r.headers {
			req.Header.Add(key, value)
		}
		cli := &http.Client{Timeout: time.Second * 10}
		resp, err := cli.Do(req)
		if err == nil {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
			return body
		} else {

			fmt.Printf("%+v %d\n", r, resp.StatusCode)
		}
		return []byte{}
	}
	return conn
}
