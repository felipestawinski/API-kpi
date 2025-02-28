package handlers

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"encoding/json"
	"math/big"
	"net/http"
	"reflect"
	"strings"

	"github.com/BloxBerg-UTFPR/API-Blockchain/cmd/api/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
)

// var readMethods = map[string]bool{
// 	"DOMAIN_SEPARATOR": true,
// 	"Allowance":        true,
// 	"BalanceOf":        true,
// 	"Decimals":         true,
// 	"Eip712Domain":     true,
// 	"Name":             true,
// 	"Nonces":           true,
// 	"Owner":            true,
// 	"Paused":           true,
// 	"Symbol":           true,
// 	"TotalSupply":      true,
// }

var readMethods = map[string]bool{
    "GetAccessLevel": true,
    "GetMetaData": true,
}
//Create a function to interact with the blockchain

func BlockchainInteraction(w http.ResponseWriter, r *http.Request) {
	//Make sure to remove the prefix '0x' of private key 
	//Eg.: 0xb13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83 -> b13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83

	//contrato vitao
	contractAddress := "0x473f8eA5Ce1F35acf7Eb61A6D4b74C8f5cf2f362"
	
	//UserAuthorized(w, r)

	action := r.PathValue("method")
	if action == "" {
		http.Error(w, "Action should be provided", http.StatusBadRequest)
		return
	}
	
	ethclient := ConnectToEthereum()
	conn, err := contract.NewContract(common.HexToAddress(contractAddress), ethclient)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to instantiate contract: %v", err), http.StatusInternalServerError)
	}

	var result interface{}
	var txn *types.Transaction

	// Parse query parameters
	params := make([]interface{}, 0)
	for _, values := range r.URL.Query() {
		if len(values) > 0 {
			params = append(params, values[0])
		}
	}

	if readMethods[action] {
		result, err = callReadMethod(conn, action, params...)
	} else {
		//accountAddress := "0x42A2069C3F18DCd7e3a61276A8401bE431958239"
		//privateKey := "a5058f2733cf8c615730962e75ed4af277df4983a9243298cd1c77afa621681f"

		//privateKey vitao
		privateKey := "65ec240a4866e5f3aa86aef6da44daea4eed19172b9991c6244ba87865de955f"
		auth := getAccountAuth(ethclient, privateKey)
		txn, err = callWriteMethod(conn, action, auth, params...)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling method %s: %v", action, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if txn != nil {
		json.NewEncoder(w).Encode(map[string]string{"txHash": txn.Hash().Hex()})
	} else {
		json.NewEncoder(w).Encode(result)
	}
	
}

func callReadMethod(conn *contract.Contract, methodName string, params ...interface{}) (interface{}, error) {
	method := reflect.ValueOf(conn).MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("method %s not found", methodName)
	}

	inputs := make([]reflect.Value, len(params)+1)
	inputs[0] = reflect.ValueOf(&bind.CallOpts{})
	methodType := method.Type()
	for i := 0; i < len(params); i++ {
		paramType := methodType.In(i + 1)
		convertedParam, err := convertParam(params[i], paramType)
		if err != nil {
			return nil, fmt.Errorf("error converting parameter %d: %v", i, err)
		}
		inputs[i+1] = reflect.ValueOf(convertedParam)
	}


	results := method.Call(inputs)
	if len(results) == 0 {
		return nil, fmt.Errorf("method %s returned no results", methodName)
	}

	if err, ok := results[len(results)-1].Interface().(error); ok && err != nil {
		return nil, err
	}

	return results[0].Interface(), nil
}

func callWriteMethod(conn *contract.Contract, methodName string, auth *bind.TransactOpts, params ...interface{}) (*types.Transaction, error) {
	method := reflect.ValueOf(conn).MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("method %s not found", methodName)
	}

	inputs := make([]reflect.Value, len(params)+1)
	inputs[0] = reflect.ValueOf(auth)

	methodType := method.Type()
	for i := 0; i < len(params); i++ {
		paramType := methodType.In(i + 1)
		convertedParam, err := convertParam(params[i], paramType)
		if err != nil {
			return nil, fmt.Errorf("error converting parameter %d: %v", i, err)
		}
		inputs[i+1] = reflect.ValueOf(convertedParam)
	}

	results := method.Call(inputs)
	if len(results) == 0 {
		return nil, fmt.Errorf("method %s returned no results", methodName)
	}

	if err, ok := results[len(results)-1].Interface().(error); ok && err != nil {
		return nil, err
	}

	return results[0].Interface().(*types.Transaction), nil
}

func convertParam(param interface{}, targetType reflect.Type) (interface{}, error) {
	switch targetType {
	case reflect.TypeOf(common.Address{}):
		if str, ok := param.(string); ok {
			return common.HexToAddress(str), nil
		}
	case reflect.TypeOf(&big.Int{}):
		if str, ok := param.(string); ok {
			n, ok := new(big.Int).SetString(str, 10)
			if !ok {
				return nil, fmt.Errorf("invalid big.Int string: %s", str)
			}
			return n, nil
		}
	}
	return param, nil
}


func ConnectToEthereum() *ethclient.Client {
    //client, err := ethclient.Dial("https://polygon-amoy.g.alchemy.com/v2/5N0az1E28WsPnYk7g37TnNSNs3QZwKMr")

	client, err := ethclient.Dial("https://polygon-amoy.g.alchemy.com/v2/jsocrRMAC5988zy4JCKk1BY8_hWHBMHB")
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
	auth.GasPrice = big.NewInt(25000000000)
	return auth
}