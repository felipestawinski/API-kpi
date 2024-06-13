package blockchain

import (
    "log"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/common"
    //"github.com/ethereum/go-ethereum/crypto"
    //"crypto/ecdsa"
    //"context"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    //"math/big"
    "github.com/BloxBerg-UTFPR/API-Blockchain/cmd/api/contract"
)

func Blockchain() {
	client := ConnectToEthereum()
    accountAddress := "0xD929775Aa54f87E641F176620332C9DF8f11146D"

	//Make sure to remove the prefix '0x' of private key 
	//Eg.: 0xb13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83 -> b13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83
    //privateKeyAddress := "b13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83"

    //auth := getAccountAuth(client, privateKeyAddress)

    //deploying smart contract
    // ownerAddress := common.HexToAddress(accountAddress)
	// address, _, _, err := contract.DeployContract(auth, client, ownerAddress)
	// if err != nil {
	// 	log.Fatalf("Failed: %v", err)
	// }

	//fmt.Println(address.Hex())

	contractAddress := "0x682862aA6d611D0D492Bdc1E0f1db170337cbD62"

	conn, err := contract.NewContract(common.HexToAddress(contractAddress), client)
	if err != nil {
		panic(err)
	}

	reply, err := conn.BalanceOf(&bind.CallOpts{}, common.HexToAddress((accountAddress)))
		if err != nil {
			panic(err) 
		}
	fmt.Println("reply->", reply)

	reply2, err := conn.Pause(&bind.TransactOpts{})
		if err != nil {
			panic(err) 
		}
	fmt.Println("reply2->", reply2)
}

// func getAccountAuth(client *ethclient.Client, privateKeyAddress string) *bind.TransactOpts {

// 	privateKey, err := crypto.HexToECDSA(privateKeyAddress)
// 	if err != nil {
// 		panic(err)
// 	}

// 	publicKey := privateKey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		panic("invalid key")
// 	}

// 	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
// 	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("nounce=", nonce)
// 	chainID, err := client.ChainID(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}

// 	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
// 	if err != nil {
// 		panic(err)
// 	}

//     // Suggest a gas price
//     gasPrice, err := client.SuggestGasPrice(context.Background())
//     if err != nil {
//         log.Fatalf("Failed to get gas price: %v", err)
//     }

// 	auth.Nonce = big.NewInt(int64(nonce))
// 	auth.Value = big.NewInt(0)      // in wei
// 	auth.GasLimit = uint64(3000000) // in units
// 	auth.GasPrice = gasPrice

// 	return auth
// }

func ConnectToEthereum() *ethclient.Client {
    client, err := ethclient.Dial("http://127.0.0.1:7545")
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }
    return client
}