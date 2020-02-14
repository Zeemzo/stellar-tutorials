package main

import (
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"

	"strconv"
)

func main() {
	// Alice's Credentials - replace with your respective keys
	senderPublicKey := "GBLPKM32JA4KTYAEWO55XSVLJNGYIABKHP6KYROKA4ZEUCGAARVLCRYU"
	senderSecretKey := "SA2ORJEER4H2QECTRCYM67CM6IFC2INOWE2OVCNWQWX3IGULTKJC3KQL"
	aPair, _ := keypair.ParseFull(senderSecretKey)

	// Bob's Credentials - replace with your respective keys
	receiverPublicKey := "GBJZEJXDHK56LHHGRNS2BNHOPFUIHSQSJPI2ACT7HF3JA5XAOKG5D4Y5"
	receiverSecretKey := "SA5TXFONM6G6BKXW7JK66KFEPBEZ3LHAG5ZKZDXGNAKRMP3IETRGGVW6"
	bPair, _ := keypair.ParseFull(receiverSecretKey)

	// Set Horizon Client to Testnet
	client := horizonclient.DefaultTestNetClient

	// Get information about the Alice's account
	aliceAccountRequest := horizonclient.AccountRequest{AccountID: senderPublicKey}
	aliceAccount, err := client.AccountDetail(aliceAccountRequest)
	if err != nil {
		log.Fatal(err)
	}

	// Get information about the Bob's account
	bobAccountRequest := horizonclient.AccountRequest{AccountID: receiverPublicKey}
	bobAccount, err := client.AccountDetail(bobAccountRequest)
	if err != nil {
		log.Fatal(err)
	}

	//View account balance
	log.Println("Alice's Balance Before:", aliceAccount.Balances[0].Asset.Type, aliceAccount.Balances[0].Balance)
	log.Println("Bob's Balance Before:", bobAccount.Balances[0].Asset.Type, bobAccount.Balances[0].Balance)

	//Assign Native Asset Interface
	asset := txnbuild.NativeAsset{}

	//Alice's sequence
	aliceSequence, err := strconv.ParseInt(aliceAccount.Sequence, 10, 64)
	if err != nil {
		log.Println(err)
	}

	//Bob's sequence
	bobSequence, err := strconv.ParseInt(bobAccount.Sequence, 10, 64)
	if err != nil {
		log.Println(err)
	}

	// Construct the operation
	paymentOp1 := txnbuild.Payment{
		SourceAccount: &txnbuild.SimpleAccount{AccountID: senderPublicKey, Sequence: aliceSequence},
		Destination:   receiverPublicKey,
		Amount:        "5",
		Asset:         asset,
	}

	paymentOp2 := txnbuild.Payment{
		SourceAccount: &txnbuild.SimpleAccount{AccountID: receiverPublicKey, Sequence: bobSequence},
		Destination:   senderPublicKey,
		Amount:        "10",
		Asset:         asset,
	}

	// Construct the transaction that will carry the operation
	tx := txnbuild.Transaction{
		SourceAccount: &txnbuild.SimpleAccount{AccountID: senderPublicKey, Sequence: aliceSequence},
		Operations:    []txnbuild.Operation{&paymentOp1, &paymentOp2},
		Network:       network.TestNetworkPassphrase,
		Timebounds:    txnbuild.NewInfiniteTimeout(),
	}

	// Alice signs the transaction, serialise it to XDR, and base 64 encode it
	txeBase64A, err := tx.BuildSignEncode(aPair)
	if err != nil {
		log.Fatal("Error Building transaction:", err)
	}

	//Bob receives xdr and opens as xdr
	txn, errXDR := txnbuild.TransactionFromXDR(txeBase64A)
	if errXDR != nil {
		log.Fatal("Error decoding xdr:", errXDR)
	}

	txn.Network=network.TestNetworkPassphrase
	// Bob signs the transaction, serialise it to XDR, and base 64 encode it
	errS:= txn.Sign(bPair)
	if errS != nil {
		log.Fatal("Error submitting transaction:", errS)
	}

	txeBase64B,err := txn.Base64()
	if err != nil {
		log.Fatal("Error encoding transaction:", err)
	}

	// Submit the transaction
	resp, err := client.SubmitTransactionXDR(txeBase64B)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Fatal("Error submitting transaction:", hError)
	}

	log.Println("Transaction Hash: ", resp.Hash)

	// Get information about the Alice's account
	aliceAccount, err = client.AccountDetail(aliceAccountRequest)
	if err != nil {
		log.Fatal(err)
	}

	// Get information about the Bob's account
	bobAccount, err = client.AccountDetail(bobAccountRequest)
	if err != nil {
		log.Fatal(err)
	}

	//View account balance
	log.Println("Alice's Balance After:", aliceAccount.Balances[0].Asset.Type, aliceAccount.Balances[0].Balance)
	log.Println("Bob's Balance After:", bobAccount.Balances[0].Asset.Type, bobAccount.Balances[0].Balance)
}
