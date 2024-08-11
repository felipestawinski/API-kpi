package handlers

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/BloxBerg-UTFPR/API-Blockchain/cmd/api/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// welcomeHandler handles welcome requests for logged-in users
func BlockchainInteraction(w http.ResponseWriter, r *http.Request) {
	// separar em métodos de leitura e escrita
	//Make sure to remove the prefix '0x' of private key 
	//Eg.: 0xb13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83 -> b13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83
    accountAddress := "0x42A2069C3F18DCd7e3a61276A8401bE431958239"
	contractAddress := "0x5BcC1E2133119cb819c97418B73C770727001FBd"
	teste := "a5058f2733cf8c615730962e75ed4af277df4983a9243298cd1c77afa621681f"

	UserAuthorized(w, r)

	action := r.PathValue("method")
	fmt.Println("method: ", action)

	ethclient := ConnectToEthereum()
	conn, err := contract.NewContract(common.HexToAddress(contractAddress), ethclient)
		if err != nil {
			panic(err)
		}
	
	if action == "" {
		http.Error(w, "Action should be provided", http.StatusUnauthorized)
		return
	} else if action == "name" {
		res, err := conn.Name(&bind.CallOpts{})
		if err != nil {
			panic (err)
		}
		fmt.Fprintf(w, res)
		w.WriteHeader(http.StatusOK)
		fmt.Println(res)
	} else if action == "nonces" {
		res, err := conn.Nonces(&bind.CallOpts{}, common.HexToAddress((accountAddress)))
		if err != nil {
			panic (err)
		}
		fmt.Println("Nonces", res)
		w.WriteHeader(http.StatusOK)
	} else if action == "owner" {
		res, err := conn.Owner(&bind.CallOpts{})
		if err != nil {
			panic (err)
		}
		fmt.Println("Owner: ", res)
		w.WriteHeader(http.StatusOK)
		fmt.Println(res)
	} else if action == "paused" {
		res, err := conn.Paused(&bind.CallOpts{})
		if err != nil {
			panic (err)
		}
		fmt.Println("Paused: ", res)
		w.WriteHeader(http.StatusOK)
		fmt.Println(res)
	} else if action == "pause" {
		auth := getAccountAuth(ethclient, teste)
		res, err := conn.Pause(auth)

		if err != nil {
			panic (err)
		}
		fmt.Println("Pause: ", res)
		w.WriteHeader(http.StatusOK)
		fmt.Println(res)
	} else if action == "unpause" {
		auth := getAccountAuth(ethclient, teste)
		res, err := conn.Unpause(auth)

		if err != nil {
			panic (err)
		}
		fmt.Println("Pause: ", res)
		w.WriteHeader(http.StatusOK)
		fmt.Println(res)
	}else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Unknown method:" + action)
		return
	}

}

func ConnectToEthereum() *ethclient.Client {
    client, err := ethclient.Dial("https://polygon-amoy.g.alchemy.com/v2/5N0az1E28WsPnYk7g37TnNSNs3QZwKMr")
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }
    return client
}

func getAccountAuth(client *ethclient.Client, accountAddress string) *bind.TransactOpts {
	//Remove the '0x' prefix if present
    if strings.HasPrefix(accountAddress, "0x") {
        accountAddress = accountAddress[2:]
    }
	privateKey, err := crypto.HexToECDSA(accountAddress)
	
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("invalid key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//fetch the last use nonce of account
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
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = big.NewInt(1000000)
	return auth
}