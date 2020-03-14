package main

import (
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Issuer Credentials - replace with your respective keys
	IssuerPublicKey := "GC2HSDISNE676F23TSAXE3WJWGVS3PUG7WHH3KJDS2BYOIQWP4XJYVVA"
	IssuerSecretKey := "SAWYJEMWLENAGPLTEYAH55EHHYYZ2QU6Q76ITUSS3ROHYITH2J5UJ7SC"
	IssuerKeypair, _ := keypair.ParseFull(IssuerSecretKey)

	// Distributor Credentials - replace with your respective keys
	DistributorPublicKey := "GCGZDP5FLFJMVPY4WRPPNA5I6B7XV2M3JCZAWTRZHIJ7XXNKUH26NYJX"
	DistributorSecretKey := "SBUHPUUDUG4HGLTJBBIPNZ7UY65V4SVWYLIUUH3MXYOJHLSKDXLSP3NM"
	DistributorKeypair, _ := keypair.ParseFull(DistributorSecretKey)

	// User Credentials - replace with your respective keys
	// UserPublicKey := "GC3J24BHVNCV6OQIFOEDC6IPYQQ6MJZFP5XPFAGZN6L3T3C4COOFBCMX"
	UserSecretKey := "SCNUELCHDINKKRW2NKOUYJLK4OFD7NE26MWQ3HYGPIXDAJ5JFOULRGEK"
	UserKeypair, _ := keypair.ParseFull(UserSecretKey)

	// Set Horizon Client to Testnet
	client := horizonclient.DefaultTestNetClient

	// Distributor Creates Trustline to Asset Issuer
	log.Println("STEP 1:Distributor Creates Trustline to Asset Issuer")
	ChangeTrustResponse := <-ChangeTrust(DistributorKeypair, IssuerPublicKey, client)
	log.Println("Transaction Hash: ", ChangeTrustResponse)

	// Issuer Funds Distributor Account with total coin allocation
	log.Println("STEP 2:Issuer Funds Distributor Account with total coin allocation")
	FundDistributorResponse := <-FundDistributor(IssuerKeypair, DistributorPublicKey, client)
	log.Println("Transaction Hash: ", FundDistributorResponse)

	// Issuer locks itself by nullifying Master Threshold
	log.Println("STEP 3:Issuer locks itself by nullifying Master Threshold")
	LockIsuerResponse := <-LockIsuer(IssuerKeypair, client)
	log.Println("Transaction Hash: ", LockIsuerResponse)

	// User Creates Trustline to Asset Issuer
	log.Println("STEP 4:User Creates Trustline to Asset Issuer")
	ChangeTrustResponse2 := <-ChangeTrust(UserKeypair, IssuerPublicKey, client)
	log.Println("Transaction Hash: ", ChangeTrustResponse2)

	// User Buys the Coin in exchange for XLMS
	log.Println("STEP 5:User Buys the Coin in exchange for XLMS")
	OfferCoinExchangeResponse := <-OfferCoinExchange(UserKeypair, DistributorKeypair, IssuerPublicKey, client)
	log.Println("Transaction Hash: ", OfferCoinExchangeResponse)

}

//ChangeTrust creates and submits the create trust operation
func ChangeTrust(DistributorKeypair *keypair.Full, IssuerPublicKey string,
	client *horizonclient.Client) <-chan string {

	res := make(chan string)

	go func() {
		defer close(res)

		// Get information about the Distributor account
		accountRequest := horizonclient.AccountRequest{AccountID: DistributorKeypair.Address()}
		Account, err := client.AccountDetail(accountRequest)
		if err != nil {
			log.Fatal(err)
		}

		// Construct the operation
		changeTrustOp := txnbuild.ChangeTrust{
			Line: txnbuild.CreditAsset{
				Code:   "DM",
				Issuer: IssuerPublicKey,
			},
			Limit:         "1000000",
			SourceAccount: &Account,
		}

		// Construct the transaction that will carry the operation
		tx := txnbuild.Transaction{
			SourceAccount: &Account,
			Operations:    []txnbuild.Operation{&changeTrustOp},
			Network:       network.TestNetworkPassphrase,
			Timebounds:    txnbuild.NewInfiniteTimeout(),
		}

		// Sign the transaction, serialise it to XDR, and base 64 encode it
		txeBase64, err := tx.BuildSignEncode(DistributorKeypair)
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

//FundDistributor creates and submits the payment operation to fund the distributor
func FundDistributor(IssuerKeypair *keypair.Full, DistributorPublicKey string,
	client *horizonclient.Client) <-chan string {

	res := make(chan string)

	go func() {
		defer close(res)
		// Get information about the Distributor account
		accountRequest := horizonclient.AccountRequest{AccountID: IssuerKeypair.Address()}
		Account, err := client.AccountDetail(accountRequest)
		if err != nil {
			log.Fatal(err)
		}

		// Construct the operation
		paymentOp := txnbuild.Payment{
			Destination: DistributorPublicKey,
			Amount:      "1000000",
			Asset: txnbuild.CreditAsset{
				Code:   "DM",
				Issuer: IssuerKeypair.Address(),
			},
			SourceAccount: &Account,
		}

		// Construct the transaction that will carry the operation
		tx := txnbuild.Transaction{
			SourceAccount: &Account,
			Operations:    []txnbuild.Operation{&paymentOp},
			Network:       network.TestNetworkPassphrase,
			Timebounds:    txnbuild.NewInfiniteTimeout(),
		}

		// Sign the transaction, serialise it to XDR, and base 64 encode it
		txeBase64, err := tx.BuildSignEncode(IssuerKeypair)
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

//LockIsuer creates and submits the set option operation to lock account
func LockIsuer(IssuerKeypair *keypair.Full, client *horizonclient.Client) <-chan string {
	res := make(chan string)

	go func() {
		// Get information about the Issuer account
		accountRequest := horizonclient.AccountRequest{AccountID: IssuerKeypair.Address()}
		Account, err := client.AccountDetail(accountRequest)
		if err != nil {
			log.Fatal(err)
		}

		// Construct the operation
		setOptionsOpp := txnbuild.SetOptions{MasterWeight: txnbuild.NewThreshold(0)}

		// Construct the transaction that will carry the operation
		tx := txnbuild.Transaction{
			SourceAccount: &Account,
			Operations:    []txnbuild.Operation{&setOptionsOpp},
			Network:       network.TestNetworkPassphrase,
			Timebounds:    txnbuild.NewInfiniteTimeout(),
		}

		// Sign the transaction, serialise it to XDR, and base 64 encode it
		txeBase64, err := tx.BuildSignEncode(IssuerKeypair)
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

//OfferCoinExchange creates and submits the coin exchange payment multi-signature transaction
func OfferCoinExchange(UserKeypair *keypair.Full, DistributorKeypair *keypair.Full,
	IssuerPublicKey string, client *horizonclient.Client) <-chan string {
	res := make(chan string)

	go func() {
		defer close(res)

		// Get information about the User's account
		UserAccountRequest := horizonclient.AccountRequest{AccountID: UserKeypair.Address()}
		UserAccount, err := client.AccountDetail(UserAccountRequest)
		if err != nil {
			log.Fatal(err)
		}

		// Get information about the DistributorRequest's account
		DistributorRequest := horizonclient.AccountRequest{AccountID: DistributorKeypair.Address()}
		DistributorAccount, err := client.AccountDetail(DistributorRequest)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("User's Wallet Before")
		//View account balance
		for _, Bal := range UserAccount.Balances {
			if Bal.Asset.Code == "" {
				log.Println("Asset Type:"+Bal.Asset.Type, "Asset Code:XLM", "Asset Balance:"+Bal.Balance)
			} else {
				log.Println("Asset Type:"+Bal.Asset.Type, "Asset Code:"+Bal.Asset.Code, "Asset Balance:"+Bal.Balance)
			}
		}

		//Assign Native Asset Interface
		nativeAsset := txnbuild.NativeAsset{}

		// Construct the operation
		paymentOp1 := txnbuild.Payment{
			SourceAccount: &UserAccount,
			Destination:   DistributorKeypair.Address(),
			Amount:        "5",
			Asset:         nativeAsset,
		}

		paymentOp2 := txnbuild.Payment{
			SourceAccount: &DistributorAccount,
			Destination:   UserKeypair.Address(),
			Amount:        "1",
			Asset: txnbuild.CreditAsset{
				Code:   "DM",
				Issuer: IssuerPublicKey,
			},
		}

		// Construct the transaction that will carry the operation
		tx := txnbuild.Transaction{
			SourceAccount: &UserAccount,
			Operations:    []txnbuild.Operation{&paymentOp1, &paymentOp2},
			Network:       network.TestNetworkPassphrase,
			Timebounds:    txnbuild.NewInfiniteTimeout(),
		}

		// User signs the transaction, serialise it to XDR, and base 64 encode it
		txeBase64A, err := tx.BuildSignEncode(UserKeypair)
		if err != nil {
			log.Fatal("Error Building transaction:", err)
		}

		//Distributor receives xdr and opens as xdr
		txn, errXDR := txnbuild.TransactionFromXDR(txeBase64A)
		if errXDR != nil {
			log.Fatal("Error decoding xdr:", errXDR)
		}

		txn.Network = network.TestNetworkPassphrase
		// Distributor signs the transaction if right, serialises it to XDR, and base 64 encode it
		errS := txn.Sign(DistributorKeypair)
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

		// log.Println("Transaction Hash: ", resp.Hash)

		UserAccount, err = client.AccountDetail(UserAccountRequest)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("User's Wallet After")
		//View User account balance
		for _, Bal := range UserAccount.Balances {
			if Bal.Asset.Code == "" {
				log.Println("Asset Type:"+Bal.Asset.Type, "Asset Code:XLM", "Asset Balance:"+Bal.Balance)
			} else {
				log.Println("Asset Type:"+Bal.Asset.Type, "Asset Code:"+Bal.Asset.Code, "Asset Balance:"+Bal.Balance)
			}
		}

		res <- resp.Hash

	}()

	return res

}
