package main

import (
	"encoding/base64"
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Credentials - replace with your respective keys

	publicKey := "GBLPKM32JA4KTYAEWO55XSVLJNGYIABKHP6KYROKA4ZEUCGAARVLCRYU"
	secretKey := "SA2ORJEER4H2QECTRCYM67CM6IFC2INOWE2OVCNWQWX3IGULTKJC3KQL"
	pair, _ := keypair.ParseFull(secretKey)

	// Set Horizon Client to Testnet
	client := horizonclient.DefaultTestNetClient

	// Get information about the sender account
	accountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	hAccount0, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
	}

	// Construct the operation
	manageDataOp := txnbuild.ManageData{
		Name:  "name",
		Value: []byte("Alice"),
	}

	// Construct the transaction that will carry the operation
	tx := txnbuild.Transaction{
		SourceAccount: &hAccount0,
		Operations:    []txnbuild.Operation{&manageDataOp},
		Network:       network.TestNetworkPassphrase,
		Timebounds:    txnbuild.NewInfiniteTimeout(),
	}

	// Sign the transaction, serialise it to XDR, and base 64 encode it
	txeBase64, err := tx.BuildSignEncode(pair)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Fatal("Error submitting transaction:", hError)
	}

	// Submit the transaction
	resp, err := client.SubmitTransactionXDR(txeBase64)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Fatal("Error submitting transaction:", hError)
	}

	log.Println("Transaction Hash: ", resp.Hash)

	hAccount0, err = client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
	}
	//View account value:name
	//decode from base64
	decoded, err := base64.StdEncoding.DecodeString(hAccount0.Data["name"])
	if err != nil {
		log.Println("decode error:", err)
	}

	log.Println("Added Data: ", string(decoded))
}
