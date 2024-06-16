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

func Blockchain( deploy string ) {
	client := ConnectToEthereum()
	
	//Make sure to remove the prefix '0x' of private key 
	//Eg.: 0xb13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83 -> b13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83
    accountAddress := "0x42A2069C3F18DCd7e3a61276A8401bE431958239"
    privateKeyAddress := "b13b0008c5cb0379f1d3a427bbfc75838a50eea795cf549c17b25e2c350c2e83"
	contractAddress := "0x5BcC1E2133119cb819c97418B73C770727001FBd"
	

	if (deploy == "deploy") {
		//deploying smart contract
		auth := getAccountAuth(client, privateKeyAddress)
		ownerAddress := common.HexToAddress(accountAddress)
		address, _, _, err := contract.DeployContract(auth, client, ownerAddress)
		if err != nil {
			log.Fatalf("Failed: %v", err)
		}
	
		fmt.Println(address.Hex())

		conn, err := contract.NewContract(address, client)
		if err != nil {
			panic(err)
		}

		//Balance of account
		reply, err := conn.BalanceOf(&bind.CallOpts{}, common.HexToAddress((accountAddress)))
			if err != nil {
				panic(err) 
			}
		fmt.Println("reply->", reply)
		
		//Pause TO DO: Fix 
		reply2, err := conn.Unpause(&bind.TransactOpts{})
			if err != nil {
				panic(err) 
			}
		fmt.Println("reply2->", reply2)

	} else {
		conn, err := contract.NewContract(common.HexToAddress(contractAddress), client)
		if err != nil {
			fmt.Println("aqui")
			panic(err)
		}

		reply, err := conn.BalanceOf(&bind.CallOpts{}, common.HexToAddress((accountAddress)))
			if err != nil {
				fmt.Println("aqui2")
				panic(err) 
			}
		fmt.Println("reply->", reply)


		// Transfer to account
		// num := big.NewInt(3) 
		// reply3, err := conn.Mint(&bind.TransactOpts{}, common.HexToAddress((accountAddress)), num)
		// if err != nil {
		// 	panic(err) 
		// }
		// fmt.Println("reply3->", reply3)

		reply4, err := conn.Name(&bind.CallOpts{})
		if err != nil {
			panic (err)
		}
		fmt.Println("reply4->", reply4)

		reply5, err := conn.Nonces(&bind.CallOpts{}, common.HexToAddress((accountAddress)))
		if err != nil {
			panic (err)
		}
		fmt.Println("reply5->", reply5)

		reply6, err := conn.Owner(&bind.CallOpts{})
		if err != nil {
			panic (err)
		}
		fmt.Println("reply6->", reply6)

		reply7, err := conn.Paused(&bind.CallOpts{})
		if err != nil {
			panic (err)
		}
		fmt.Println("reply7->", reply7)

	}
	
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
    client, err := ethclient.Dial("https://polygon-amoy.g.alchemy.com/v2/5N0az1E28WsPnYk7g37TnNSNs3QZwKMr")
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }
    return client
}