package main

import (
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// First User's Credentials - replace with your respective keys
	// FirstPublicKey := "GBV5VZTKXO72JAB46PHAKTZCVBCETVUKOGXOI3XY27EVFXRJW52TLTG3"
	FirstSecretKey := "SCR2COELWOG6YCY4WILZCJ7AMID5272PUXBA2OBC7B6F7F5E2LKXJIR3"
	FirstKeypair, _ := keypair.ParseFull(FirstSecretKey)

	// Second User's Credentials - replace with your respective keys
	SecondPublicKey := "GD46OP3JSAHY2NX5EVFG3TLBAHUYUJDUASSA5ZMQYRETLDUF5LCJQA2J"
	SecondSecretKey := "SDP75EO6MCYBVFAD7TAUWBLGQ3NEFUBGOIQSQ3AUK3UQHMDYLNK2SLVP"
	SecondKeypair, _ := keypair.ParseFull(SecondSecretKey)

	// Set Horizon Client to Testnet
	client := horizonclient.DefaultTestNetClient

	// Set Options Signer type to add another signer
	log.Println("STEP 1:Set Options Signer type to add another signer")
	AddSignerResponse := <-AddSigner(FirstKeypair, SecondPublicKey, client)
	log.Println("Transaction Hash: ", AddSignerResponse)

	// Test a transaction using second account
	log.Println("STEP 2:Test a transaction using second account")
	TestTransactionResponse := <-TestTransaction(SecondKeypair, client)
	log.Println("Transaction Hash: ", TestTransactionResponse)

	// Set Options Low Threshold to 2
	log.Println("STEP 3:Set Options Low Threshold to 2")
	SetThresholdResponse := <-SetThreshold(FirstKeypair, client)
	log.Println("Transaction Hash: ", SetThresholdResponse)

	// Test a transaction using Multisignature
	log.Println("STEP 4: Test a transaction using Multisignature")
	MultisignatureTransactionResponse := <-MultisignatureTransaction(FirstKeypair, SecondKeypair,
		client)
	log.Println("Transaction Hash: ", MultisignatureTransactionResponse)

}

//AddSigner Set Options Signer type to add an another signer
func AddSigner(FirstKeypair *keypair.Full, SecondPublicKey string,
	client *horizonclient.Client) <-chan string {

	res := make(chan string)

	go func() {
		defer close(res)

		// Get information about the Distributor account
		accountRequest := horizonclient.AccountRequest{AccountID: FirstKeypair.Address()}
		Account, err := client.AccountDetail(accountRequest)
		if err != nil {
			log.Fatal(err)
		}

		// Construct the operation
		setOptionsOp := txnbuild.SetOptions{
			Signer: &txnbuild.Signer{
				Address: SecondPublicKey,
				Weight:  1,
			},
		}

		// Construct the transaction that will carry the operation
		tx := txnbuild.Transaction{
			SourceAccount: &Account,
			Operations:    []txnbuild.Operation{&setOptionsOp},
			Network:       network.TestNetworkPassphrase,
			Timebounds:    txnbuild.NewInfiniteTimeout(),
		}

		// Sign the transaction, serialise it to XDR, and base 64 encode it
		txeBase64, err := tx.BuildSignEncode(FirstKeypair)
		if err != nil {
			hError := err.(*horizonclient.Error)
			log.Fatal("Error submitting transaction:", hError)
		}
		log.Println("txeBase64: ", txeBase64)

		// Submit the transaction
		resp, err := client.SubmitTransactionXDR(txeBase64)
		if err != nil {
			hError := err.(*horizonclient.Error)
			log.Fatal("Error submitting transaction:", hError)
		}

		res <- resp.Hash
	}()

	return res

}

//SetThreshold Set Options Signer type to add an another signer
func SetThreshold(FirstKeypair *keypair.Full,
	client *horizonclient.Client) <-chan string {

	res := make(chan string)

	go func() {
		defer close(res)

		// Get information about the First account
		accountRequest := horizonclient.AccountRequest{AccountID: FirstKeypair.Address()}
		Account, err := client.AccountDetail(accountRequest)
		if err != nil {
			log.Fatal(err)
		}

		threshold := txnbuild.NewThreshold(2)
		// Construct the operation
		setOptionsOp := txnbuild.SetOptions{
			LowThreshold: threshold,
		}

		// Construct the transaction that will carry the operation
		tx := txnbuild.Transaction{
			SourceAccount: &Account,
			Operations:    []txnbuild.Operation{&setOptionsOp},
			Network:       network.TestNetworkPassphrase,
			Timebounds:    txnbuild.NewInfiniteTimeout(),
		}

		// Sign the transaction, serialise it to XDR, and base 64 encode it
		txeBase64, err := tx.BuildSignEncode(FirstKeypair)
		if err != nil {
			hError := err.(*horizonclient.Error)
			log.Fatal("Error submitting transaction:", hError)
		}
		log.Println("txeBase64: ", txeBase64)

		// Submit the transaction
		resp, err := client.SubmitTransactionXDR(txeBase64)
		if err != nil {
			hError := err.(*horizonclient.Error)
			log.Fatal("Error submitting transaction:", hError)
		}

		res <- resp.Hash
	}()

	return res

}

//TestTransaction Tests the account ability to sign a tranasaction depending on it's threshold
func TestTransaction(Keypair *keypair.Full, client *horizonclient.Client) <-chan string {

	res := make(chan string)

	go func() {
		defer close(res)
		// Get information about the account
		accountRequest := horizonclient.AccountRequest{AccountID: Keypair.Address()}
		Account, err := client.AccountDetail(accountRequest)
		if err != nil {
			log.Fatal(err)
		}

		// Construct the operation
		manageDataOp := txnbuild.ManageData{
			Name:  "name",
			Value: []byte("Joint Account"),
		}
		// Construct the transaction that will carry the operation
		tx := txnbuild.Transaction{
			SourceAccount: &Account,
			Operations:    []txnbuild.Operation{&manageDataOp},
			Network:       network.TestNetworkPassphrase,
			Timebounds:    txnbuild.NewInfiniteTimeout(),
		}

		// Sign the transaction, serialise it to XDR, and base 64 encode it
		txeBase64, err := tx.BuildSignEncode(Keypair)
		if err != nil {
			hError := err.(*horizonclient.Error)
			log.Fatal("Error submitting transaction:", hError)
		}
		// log.Println("txeBase64: ", txeBase64)

		// Submit the transaction
		resp, err := client.SubmitTransactionXDR(txeBase64)
		if err != nil {
			hError := err.(*horizonclient.Error)
			log.Fatal("Error submitting transaction:", hError)
		}

		res <- resp.Hash
	}()
	return res

}

//MultisignatureTransaction Tests multisign tranasaction to have minimum threshold
func MultisignatureTransaction(FirstKeypair *keypair.Full, SecondKeypair *keypair.Full, client *horizonclient.Client) <-chan string {
	res := make(chan string)

	go func() {
		defer close(res)

		// Get information about the First account
		FirstAccountRequest := horizonclient.AccountRequest{AccountID: FirstKeypair.Address()}
		FirstAccount, err := client.AccountDetail(FirstAccountRequest)
		if err != nil {
			log.Fatal(err)
		}

		// Construct the operation
		manageDataOp := txnbuild.ManageData{
			Name:  "name",
			Value: []byte("Joint Account"),
		}

		// Construct the transaction that will carry the operation
		tx := txnbuild.Transaction{
			SourceAccount: &FirstAccount,
			Operations:    []txnbuild.Operation{&manageDataOp},
			Network:       network.TestNetworkPassphrase,
			Timebounds:    txnbuild.NewInfiniteTimeout(),
		}

		// First Signer signs the transaction, serialise it to XDR, and base 64 encode it
		txeBase64A, err := tx.BuildSignEncode(FirstKeypair)
		if err != nil {
			log.Fatal("Error Building transaction:", err)
		}

		//Second Signer receives xdr and opens as xdr
		txn, errXDR := txnbuild.TransactionFromXDR(txeBase64A)
		if errXDR != nil {
			log.Fatal("Error decoding xdr:", errXDR)
		}

		txn.Network = network.TestNetworkPassphrase
		// Second Signer signs the transaction if right, serialises it to XDR, and base 64 encode it
		errS := txn.Sign(SecondKeypair)
		if errS != nil {
			log.Fatal("Error submitting transaction:", errS)
		}

		txeBase64B, err := txn.Base64()
		if err != nil {
			log.Fatal("Error encoding transaction:", err)
		}

		// Submit the transaction
		resp, err := client.SubmitTransactionXDR(txeBase64B)
		if err != nil {
			hError := err.(*horizonclient.Error)
			log.Fatal("Error submitting transaction:", hError)
		}

		res <- resp.Hash

	}()

	return res

}
