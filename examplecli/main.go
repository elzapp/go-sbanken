package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/elzapp/go-sbanken"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(log.WarnLevel)
}

func main() {
	var creds sbanken.Credentials
	configfile, _ := os.Open(os.Args[1])
	bconfig, _ := ioutil.ReadAll(configfile)
	json.Unmarshal(bconfig, &creds)
	connection := sbanken.NewAPIConnection(creds)
	accounts := connection.GetAccounts()
	fmt.Println("╔══ Account overview ══════════════════════════════════════════════╗")
	fmt.Printf("║ %-25s%11s    % 10s    % 10s ║\n", "Name", "Number", "Balance", "Available")
	for _, acc := range accounts {
		fmt.Printf("║ %-25s%11s kr % 10.2f kr % 10.2f ║\n", acc.Name, acc.AccountNumber, acc.Balance, acc.Available)
	}

	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	efakturas, err := connection.GetAllEFakturas()
	if err != nil {
		log.Errorf("%s", err.Error())
	} else {
		for _, ef := range efakturas {
			fmt.Printf("%s", ef.IssuerName)
		}
	}
}
