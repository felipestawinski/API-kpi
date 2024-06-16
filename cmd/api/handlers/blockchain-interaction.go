package handlers

import (
	"fmt"
	"net/http"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/BloxBerg-UTFPR/API-Blockchain/cmd/api/contract"
)



// welcomeHandler handles welcome requests for logged-in users
func BlockchainInteraction(w http.ResponseWriter, r *http.Request) {

	//Make sure to remove the prefix '0x' of private key 
	//Eg.: 0xb13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83 -> b13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83
    accountAddress := "0x42A2069C3F18DCd7e3a61276A8401bE431958239"
    //privateKeyAddress := "b13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83"
	contractAddress := "0x5BcC1E2133119cb819c97418B73C770727001FBd"
	UserAuthorized(w, r)

	action := r.Header.Get("Action")

	ethclient := ConnectToEthereum()
	conn, err := contract.NewContract(common.HexToAddress(contractAddress), ethclient)
		if err != nil {
			panic(err)
		}
	
	res := ""

	fmt.Println("ACTION: ", action)

	if action == "" {
		http.Error(w, "Action should be provided", http.StatusUnauthorized)
		return
	}

	if action == "name" {
		res, err := conn.Name(&bind.CallOpts{})
		if err != nil {
			panic (err)
		}
		fmt.Println("Name", res)
	}

	if action == "nonces" {
		res, err := conn.Nonces(&bind.CallOpts{}, common.HexToAddress((accountAddress)))
		if err != nil {
			panic (err)
		}
		fmt.Println("Nonces", res)
	}

	if action == "owner" {
		res, err := conn.Owner(&bind.CallOpts{})
		if err != nil {
			panic (err)
		}
		fmt.Println("Owner: ", res)
	}

	if action == "paused" {
		res, err := conn.Paused(&bind.CallOpts{})
		if err != nil {
			panic (err)
		}
		fmt.Println("Paused: ", res)
	}

	
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Response: " + res)
}

func ConnectToEthereum() *ethclient.Client {
    client, err := ethclient.Dial("https://polygon-amoy.g.alchemy.com/v2/5N0az1E28WsPnYk7g37TnNSNs3QZwKMr")
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }
    return client
}