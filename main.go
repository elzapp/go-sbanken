package sbanken

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const dateFormat = "2006-01-02T15:04:05-07:00" //2019-03-06T00:00:00+01:00
const identityserver = "https://auth.sbanken.no/identityserver/connect/token"
const apiAccounts = "https://api.sbanken.no/exec.bank/api/v1/Accounts"
const apiTransactions = "https://api.sbanken.no/exec.bank/api/v1/Transactions/%s"

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
	TransactionID        string      `json:"transactionId"`
	AccountingDate       string      `json:"accountingDate"`
	InterestDate         string      `json:"interestDate"`
	OtherAccountNumber   string      `json:"otherAccountNumber"`
	TransactionType      string      `json:"transactionType"`
	TransactionTypeCode  int64       `json:"transactionTypeCode"`
	TransactionTypeText  string      `json:"transactionTypeText"`
	IsReservation        bool        `json:"isReservation"`
	CardDetailsSpecified bool        `json:"cardDetailsSpecified"`
	Amount               float64     `json:"amount"`
	Text                 string      `json:"text"`
	Source               string      `json:"source"`
	CardDetails          cardDetails `json:"cardDetails"`
}

type cardDetails struct {
	CardNumber                  string  `json:"cardNumber"`
	CurrencyAmount              float64 `json:"currencyAmount"`
	CurrencyRate                float64 `json:"currencyRate"`
	MerchantCategoryCode        string  `json:"merchantCategoryCode"`
	MerchantCategoryDescription string  `json:"merchantCategoryDescription"`
	MerchantCity                string  `json:"merchantCity"`
	MerchantName                string  `json:"merchantName"`
	OriginalCurrencyCode        string  `json:"originalCurrencyCode"`
	PurchaseDate                string  `json:"purchaseDate"`
	TransactionID               string  `json:"transactionId"`
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
	makeAPIRequest func(r apirequest) ([]byte, error)
}

// HasToken returns true if this session has been authenticated
func (conn *APIConnection) HasToken() bool {
	if conn.token == "" {
		return false
	}
	return true
}

func (conn *APIConnection) getToken() (string, error) {
	if conn.token == "" {
		log.Debug("Getting token")
		postdata := url.Values{}
		postdata.Add("grant_type", "client_credentials")
		req, err := http.NewRequest("POST", identityserver, strings.NewReader(postdata.Encode()))
		if err != nil {
			return "", fmt.Errorf("Failed to create request %w", err)
		}
		req.Header.Add("Content-type", "application/x-www-form-urlencoded; charset=utf-8")
		req.Header.Add("User-Agent", "github.com/elzapp/go-sbanken")
		req.SetBasicAuth(conn.cred.Apikey, url.QueryEscape(conn.cred.Secret))
		cli := &http.Client{}
		resp, err := cli.Do(req)
		if err != nil {
			return "", fmt.Errorf("Failed to get token: %w", err)
		}
		if resp.StatusCode == 400 {
			return "", fmt.Errorf("Got \"%s\" while requesting token, check that your secret is valid", resp.Status)
		} else if resp.StatusCode > 399 {
			return "", fmt.Errorf("Got \"%s\" while requesting token (%+v)", resp.Status, resp)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var t tokenResponse
		json.Unmarshal(body, &t)
		if t.Token == "" {
			log.Errorf("Received empty token from identityserver (%s)", body)
		}
		conn.token = t.Token
	}
	return conn.token, nil
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
	resp, err := conn.makeAPIRequest(r)
	if err != nil {
		log.Error(err)
	}
	json.Unmarshal(resp, &a)
	return a.Accounts
}

// GetTransactions returns the latest transactions on a given account
// using the default limits set by Sbanken
func (conn *APIConnection) GetTransactions(accountid string) []Transaction {
	r := newAPIRequest()
	r.target = fmt.Sprintf(apiTransactions, accountid)
	var t transactions
	resp, _ := conn.makeAPIRequest(r)
	json.Unmarshal(resp, &t)
	return t.Transactions
}

// GetTransactionsSince returns the latest transactions on a given account
// for a given period. The period must be less than, or equal to 366 days.
// At this point this will only return the last 1000 transactions in the
// period
func (conn *APIConnection) GetTransactionsSince(accountid string, startDate string) []Transaction {
	r := newAPIRequest()
	r.target = fmt.Sprintf(apiTransactions, accountid)
	sd := time.Now()
	sd = sd.AddDate(-1, 0, 0)
	r.params["startDate"] = sd.Format("2006-01-02")
	r.params["length"] = "1000"

	var t transactions
	resp, _ := conn.makeAPIRequest(r)
	json.Unmarshal(resp, &t)
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
	conn.makeAPIRequest = func(r apirequest) ([]byte, error) {
		token, err := conn.getToken()
		if err != nil {
			return []byte{}, err
		}
		req, err := http.NewRequest("GET", r.target, nil)
		if err != nil {
			return []byte{}, fmt.Errorf("Failed to create request towards %s (%w)", r.target, err)
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("User-Agent", "github.com/elzapp/go-sbanken")

		req.Header.Add("customerId", conn.cred.UserID)
		for key, value := range r.headers {
			req.Header.Add(key, value)
		}
		log.Debugf("Requesting %+v using these headers: %+v", r, req.Header)
		if len(r.params) > 0 {
			q := req.URL.Query()
			for key, value := range r.params {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
		}
		cli := &http.Client{Timeout: time.Second * 10}
		resp, err := cli.Do(req)
		if err == nil {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return []byte{}, err
			}
			return body, nil
		}
		return []byte{}, fmt.Errorf("Unhandled error while requesting {%+v} %w", r, err)

	}
	return conn
}
