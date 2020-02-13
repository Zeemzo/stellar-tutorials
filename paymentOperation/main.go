package main

import (
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Sender Credentials - replace with your respective keys
	
	senderPublicKey := "GBLPKM32JA4KTYAEWO55XSVLJNGYIABKHP6KYROKA4ZEUCGAARVLCRYU"
	senderSecretKey := "SA2ORJEER4H2QECTRCYM67CM6IFC2INOWE2OVCNWQWX3IGULTKJC3KQL"
	pair, _ := keypair.ParseFull(senderSecretKey)

	// Receiver Address - replace with your respective keys
	receiverPublicKey := "GBJZEJXDHK56LHHGRNS2BNHOPFUIHSQSJPI2ACT7HF3JA5XAOKG5D4Y5"

	// Set Horizon Client to Testnet
	client := horizonclient.DefaultTestNetClient

	// Get information about the sender account
	accountRequest := horizonclient.AccountRequest{AccountID: senderPublicKey}
	hAccount0, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
	}

	//View account balance
	log.Println("Balance Before: ", hAccount0.Balances[0].Asset, hAccount0.Balances[0].Balance)

	//Assign Native Asset Interface
	asset := txnbuild.NativeAsset{}

	// Construct the operation
	paymentOp := txnbuild.Payment{
		Destination: receiverPublicKey,
		Amount:      "100",
		Asset:       asset,
	}

	// Construct the transaction that will carry the operation
	tx := txnbuild.Transaction{
		SourceAccount: &hAccount0,
		Operations:    []txnbuild.Operation{&paymentOp},
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
	log.Println("Balance After: ", hAccount0.Balances[0].Asset, hAccount0.Balances[0].Balance)

}
