package blockchain

import (
    "log"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "crypto/ecdsa"
    "context"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "math/big"
    "github.com/BloxBerg-UTFPR/API-Blockchain/cmd/api/contract"
)

func Blockchain() {
	client := ConnectToEthereum()
    accountAddress := "0xAf013c22ef0b034C996E762A770cc18206C9f235"
    privateKeyAddress := "51f75a1c57cb2d558d6c6e4b0e9af357f6683cfcbb90d26351150467bc8da96a"

    auth := getAccountAuth(client, privateKeyAddress)

    //deploying smart contract
    ownerAddress := common.HexToAddress(accountAddress)
	address, tx, instance, err := contract.DeployContract(auth, client, ownerAddress)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Println(address.Hex())

	_, _ = instance, tx
	fmt.Println("instance->", instance)
	fmt.Println("tx->", tx.Hash().Hex())
    
}

func getAccountAuth(client *ethclient.Client, privateKeyAddress string) *bind.TransactOpts {

	privateKey, err := crypto.HexToECDSA(privateKeyAddress)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("invalid key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}
	fmt.Println("nounce=", nonce)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		panic(err)
	}

    // Suggest a gas price
    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        log.Fatalf("Failed to get gas price: %v", err)
    }

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	return auth
}

func ConnectToEthereum() *ethclient.Client {
    client, err := ethclient.Dial("http://127.0.0.1:7545")
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }
    return client
}