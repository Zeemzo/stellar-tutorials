package main

import (
	"log"

	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
)

func main() {

	//Credentials - replace with your respective keys
	publicKey := "GBLPKM32JA4KTYAEWO55XSVLJNGYIABKHP6KYROKA4ZEUCGAARVLCRYU"
	secretKey := "SA2ORJEER4H2QECTRCYM67CM6IFC2INOWE2OVCNWQWX3IGULTKJC3KQL"
	pair, _ := keypair.ParseFull(secretKey)
	// Get information about the account we just created
	client := horizonclient.DefaultTestNetClient

	accountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	hAccount0, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
	}

	// Generate a second randomly generated address
	kp1, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Public Key:", kp1.Address())
	log.Println("Secret Key:", kp1.Seed())

	// Construct the operation
	createAccountOp := txnbuild.CreateAccount{
		Destination: kp1.Address(),
		Amount:      "100",
	}

	// Construct the transaction that will carry the operation
	tx := txnbuild.Transaction{
		SourceAccount: &hAccount0,
		Operations:    []txnbuild.Operation{&createAccountOp},
		Timebounds:    txnbuild.NewTimeout(300),
		Network:       network.TestNetworkPassphrase,
	}

	
	// Sign the transaction, serialise it to XDR, and base 64 encode it
	txeBase64, err := tx.BuildSignEncode(pair)
	log.Println("Transaction base64: ", txeBase64)

	// Submit the transaction
	resp, err := client.SubmitTransactionXDR(txeBase64)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Fatal("Error submitting transaction:", hError)
	}

	log.Println("Transaction response: ", resp)
}
