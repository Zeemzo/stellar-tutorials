package main

import (
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
)

func main() {
    // Generate a new randomly generated address
    pair, err := keypair.Random()
    if err != nil {
        log.Fatal(err)
	}

	log.Println("Public Key:", pair.Address())
    log.Println("Secret Key:", pair.Seed())

    // Create and fund the address on TestNet, using friendbot
    client := horizonclient.DefaultTestNetClient
    client.Fund(pair.Address())

    // Get information about the account we just created
    accountRequest := horizonclient.AccountRequest{AccountID: pair.Address()}
    hAccount0, err := client.AccountDetail(accountRequest)
    if err != nil {
        log.Fatal(err)
    }

    
	log.Println("Account ID:", hAccount0.AccountID)
    log.Println("Account Balance:", hAccount0.Balances[0].Balance)
}
