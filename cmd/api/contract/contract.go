// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumAccessControl.AccessLevel\",\"name\":\"accessLevel\",\"type\":\"uint8\"}],\"name\":\"AccessLevelUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"entity\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"poster\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"}],\"name\":\"DataPosted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_entity\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getMetaData\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"poster\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_entity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_uri\",\"type\":\"string\"}],\"name\":\"postData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"enumAccessControl.AccessLevel\",\"name\":\"_level\",\"type\":\"uint8\"}],\"name\":\"setPermanentAccessLevel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"enumAccessControl.AccessLevel\",\"name\":\"_level\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_durationBlocks\",\"type\":\"uint256\"}],\"name\":\"setTemporaryAccessLevel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"getAccessLevel\",\"outputs\":[{\"internalType\":\"enumAccessControl.AccessLevel\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b5060055f803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690836005811115606b57606a6074565b5b021790555060a1565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b611b71806100ae5f395ff3fe608060405234801561000f575f80fd5b5060043610610055575f3560e01c80631b9aa7e6146100595780635540dc2f1461008a5780635bc008a0146100ba5780636d54d4c7146100eb578063f46b02d114610107575b5f80fd5b610073600480360381019061006e91906111f8565b610123565b6040516100819291906112f1565b60405180910390f35b6100a4600480360381019061009f919061131f565b6104c3565b6040516100b191906113a4565b60405180910390f35b6100d460048036038101906100cf91906113e7565b610816565b6040516100e2929190611485565b60405180910390f35b610105600480360381019061010091906114cf565b610936565b005b610121600480360381019061011c919061151f565b610cc6565b005b5f606060015f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1690506001600581111561018857610187611412565b5b81600581111561019b5761019a611412565b5b14806101cb5750600360058111156101b6576101b5611412565b5b8160058111156101c9576101c8611412565b5b145b156102b45760015f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20544311156102b3575f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083600581111561027357610272611412565b5b02179055506040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102aa906115cd565b60405180910390fd5b5b8160058111156102c7576102c6611412565b5b8160058111156102da576102d9611412565b5b101561031b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103129061165b565b60405180910390fd5b60028660405161032b91906116b3565b908152602001604051809103902080549050851061037e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161037590611713565b60405180910390fd5b5f60028760405161038f91906116b3565b908152602001604051809103902086815481106103af576103ae611731565b5b905f5260205f2090600202016040518060400160405290815f82015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200160018201805461042a9061178b565b80601f01602080910402602001604051908101604052809291908181526020018280546104569061178b565b80156104a15780601f10610478576101008083540402835291602001916104a1565b820191905f5260205f20905b81548152906001019060200180831161048457829003601f168201915b5050505050815250509050805f01518160200151945094505050509250929050565b5f60035f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1690506001600581111561052657610525611412565b5b81600581111561053957610538611412565b5b148061056957506003600581111561055457610553611412565b5b81600581111561056757610566611412565b5b145b156106525760015f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054431115610651575f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083600581111561061157610610611412565b5b02179055506040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610648906115cd565b60405180910390fd5b5b81600581111561066557610664611412565b5b81600581111561067857610677611412565b5b10156106b9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106b09061165b565b60405180910390fd5b5f60405180604001604052803373ffffffffffffffffffffffffffffffffffffffff1681526020018681525090506002866040516106f791906116b3565b908152602001604051809103902081908060018154018082558091505060019003905f5260205f2090600202015f909190919091505f820151815f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160010190816107859190611958565b5050505f600160028860405161079b91906116b3565b9081526020016040518091039020805490506107b79190611a54565b905080876040516107c891906116b3565b60405180910390207f76141fea0627408a61a3a63bd9859d827b4b50e089ca7cb9d34fdb2a8777f20533896040516108019291906112f1565b60405180910390a38094505050505092915050565b5f805f805f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1690505f6001600581111561087957610878611412565b5b82600581111561088c5761088b611412565b5b14806108bc5750600360058111156108a7576108a6611412565b5b8260058111156108ba576108b9611412565b5b145b15610928575f60015f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905080431115610918575f8094509450505050610931565b43816109249190611a54565b9150505b81819350935050505b915091565b60055f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1690506001600581111561099857610997611412565b5b8160058111156109ab576109aa611412565b5b14806109db5750600360058111156109c6576109c5611412565b5b8160058111156109d9576109d8611412565b5b145b15610ac45760015f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054431115610ac3575f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690836005811115610a8357610a82611412565b5b02179055506040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610aba906115cd565b60405180910390fd5b5b816005811115610ad757610ad6611412565b5b816005811115610aea57610ae9611412565b5b1015610b2b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b229061165b565b60405180910390fd5b60016005811115610b3f57610b3e611412565b5b846005811115610b5257610b51611412565b5b1480610b82575060036005811115610b6d57610b6c611412565b5b846005811115610b8057610b7f611412565b5b145b610bc1576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bb890611ad1565b60405180910390fd5b835f808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690836005811115610c1f57610c1e611412565b5b02179055508243610c309190611aef565b60015f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508473ffffffffffffffffffffffffffffffffffffffff167fe7d80362d709ac8a3ab14b6ce9659e33aba3d2e123e4782476af81de61d5093e85604051610cb79190611b22565b60405180910390a25050505050565b60055f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905060016005811115610d2857610d27611412565b5b816005811115610d3b57610d3a611412565b5b1480610d6b575060036005811115610d5657610d55611412565b5b816005811115610d6957610d68611412565b5b145b15610e545760015f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054431115610e53575f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690836005811115610e1357610e12611412565b5b02179055506040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e4a906115cd565b60405180910390fd5b5b816005811115610e6757610e66611412565b5b816005811115610e7a57610e79611412565b5b1015610ebb576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610eb29061165b565b60405180910390fd5b60026005811115610ecf57610ece611412565b5b836005811115610ee257610ee1611412565b5b1480610f12575060046005811115610efd57610efc611412565b5b836005811115610f1057610f0f611412565b5b145b80610f405750600580811115610f2b57610f2a611412565b5b836005811115610f3e57610f3d611412565b5b145b610f7f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f7690611ad1565b60405180910390fd5b825f808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690836005811115610fdd57610fdc611412565b5b02179055505f60015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508373ffffffffffffffffffffffffffffffffffffffff167fe7d80362d709ac8a3ab14b6ce9659e33aba3d2e123e4782476af81de61d5093e8460405161106a9190611b22565b60405180910390a250505050565b5f604051905090565b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6110d782611091565b810181811067ffffffffffffffff821117156110f6576110f56110a1565b5b80604052505050565b5f611108611078565b905061111482826110ce565b919050565b5f67ffffffffffffffff821115611133576111326110a1565b5b61113c82611091565b9050602081019050919050565b828183375f83830152505050565b5f61116961116484611119565b6110ff565b9050828152602081018484840111156111855761118461108d565b5b611190848285611149565b509392505050565b5f82601f8301126111ac576111ab611089565b5b81356111bc848260208601611157565b91505092915050565b5f819050919050565b6111d7816111c5565b81146111e1575f80fd5b50565b5f813590506111f2816111ce565b92915050565b5f806040838503121561120e5761120d611081565b5b5f83013567ffffffffffffffff81111561122b5761122a611085565b5b61123785828601611198565b9250506020611248858286016111e4565b9150509250929050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61127b82611252565b9050919050565b61128b81611271565b82525050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f6112c382611291565b6112cd818561129b565b93506112dd8185602086016112ab565b6112e681611091565b840191505092915050565b5f6040820190506113045f830185611282565b818103602083015261131681846112b9565b90509392505050565b5f806040838503121561133557611334611081565b5b5f83013567ffffffffffffffff81111561135257611351611085565b5b61135e85828601611198565b925050602083013567ffffffffffffffff81111561137f5761137e611085565b5b61138b85828601611198565b9150509250929050565b61139e816111c5565b82525050565b5f6020820190506113b75f830184611395565b92915050565b6113c681611271565b81146113d0575f80fd5b50565b5f813590506113e1816113bd565b92915050565b5f602082840312156113fc576113fb611081565b5b5f611409848285016113d3565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b600681106114505761144f611412565b5b50565b5f8190506114608261143f565b919050565b5f61146f82611453565b9050919050565b61147f81611465565b82525050565b5f6040820190506114985f830185611476565b6114a56020830184611395565b9392505050565b600681106114b8575f80fd5b50565b5f813590506114c9816114ac565b92915050565b5f805f606084860312156114e6576114e5611081565b5b5f6114f3868287016113d3565b9350506020611504868287016114bb565b9250506040611515868287016111e4565b9150509250925092565b5f806040838503121561153557611534611081565b5b5f611542858286016113d3565b9250506020611553858286016114bb565b9150509250929050565b7f4163636573732044656e6965643a2054656d706f7261727920616363657373205f8201527f6578706972656400000000000000000000000000000000000000000000000000602082015250565b5f6115b760278361129b565b91506115c28261155d565b604082019050919050565b5f6020820190508181035f8301526115e4816115ab565b9050919050565b7f4163636573732044656e6965643a20496e73756666696369656e74207065726d5f8201527f697373696f6e7300000000000000000000000000000000000000000000000000602082015250565b5f61164560278361129b565b9150611650826115eb565b604082019050919050565b5f6020820190508181035f83015261167281611639565b9050919050565b5f81905092915050565b5f61168d82611291565b6116978185611679565b93506116a78185602086016112ab565b80840191505092915050565b5f6116be8284611683565b915081905092915050565b7f496e76616c6964204944000000000000000000000000000000000000000000005f82015250565b5f6116fd600a8361129b565b9150611708826116c9565b602082019050919050565b5f6020820190508181035f83015261172a816116f1565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806117a257607f821691505b6020821081036117b5576117b461175e565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026118177fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826117dc565b61182186836117dc565b95508019841693508086168417925050509392505050565b5f819050919050565b5f61185c611857611852846111c5565b611839565b6111c5565b9050919050565b5f819050919050565b61187583611842565b61188961188182611863565b8484546117e8565b825550505050565b5f90565b61189d611891565b6118a881848461186c565b505050565b5b818110156118cb576118c05f82611895565b6001810190506118ae565b5050565b601f821115611910576118e1816117bb565b6118ea846117cd565b810160208510156118f9578190505b61190d611905856117cd565b8301826118ad565b50505b505050565b5f82821c905092915050565b5f6119305f1984600802611915565b1980831691505092915050565b5f6119488383611921565b9150826002028217905092915050565b61196182611291565b67ffffffffffffffff81111561197a576119796110a1565b5b611984825461178b565b61198f8282856118cf565b5f60209050601f8311600181146119c0575f84156119ae578287015190505b6119b8858261193d565b865550611a1f565b601f1984166119ce866117bb565b5f5b828110156119f5578489015182556001820191506020850194506020810190506119d0565b86831015611a125784890151611a0e601f891682611921565b8355505b6001600288020188555050505b505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f611a5e826111c5565b9150611a69836111c5565b9250828203905081811115611a8157611a80611a27565b5b92915050565b7f496e76616c696420616363657373206c6576656c0000000000000000000000005f82015250565b5f611abb60148361129b565b9150611ac682611a87565b602082019050919050565b5f6020820190508181035f830152611ae881611aaf565b9050919050565b5f611af9826111c5565b9150611b04836111c5565b9250828201905080821115611b1c57611b1b611a27565b5b92915050565b5f602082019050611b355f830184611476565b9291505056fea264697066735822122066b33f850db37441b899984bf4ecb94eb3ac2b1476f992ab9fb226df5eb64f5564736f6c634300081a0033",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// GetAccessLevel is a free data retrieval call binding the contract method 0x5bc008a0.
//
// Solidity: function getAccessLevel(address _user) view returns(uint8, uint256)
func (_Contract *ContractCaller) GetAccessLevel(opts *bind.CallOpts, _user common.Address) (uint8, *big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getAccessLevel", _user)

	if err != nil {
		return *new(uint8), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetAccessLevel is a free data retrieval call binding the contract method 0x5bc008a0.
//
// Solidity: function getAccessLevel(address _user) view returns(uint8, uint256)
func (_Contract *ContractSession) GetAccessLevel(_user common.Address) (uint8, *big.Int, error) {
	return _Contract.Contract.GetAccessLevel(&_Contract.CallOpts, _user)
}

// GetAccessLevel is a free data retrieval call binding the contract method 0x5bc008a0.
//
// Solidity: function getAccessLevel(address _user) view returns(uint8, uint256)
func (_Contract *ContractCallerSession) GetAccessLevel(_user common.Address) (uint8, *big.Int, error) {
	return _Contract.Contract.GetAccessLevel(&_Contract.CallOpts, _user)
}

// GetMetaData is a paid mutator transaction binding the contract method 0x1b9aa7e6.
//
// Solidity: function getMetaData(string _entity, uint256 _id) returns(address poster, string uri)
func (_Contract *ContractTransactor) GetMetaData(opts *bind.TransactOpts, _entity string, _id *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "getMetaData", _entity, _id)
}

// GetMetaData is a paid mutator transaction binding the contract method 0x1b9aa7e6.
//
// Solidity: function getMetaData(string _entity, uint256 _id) returns(address poster, string uri)
func (_Contract *ContractSession) GetMetaData(_entity string, _id *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.GetMetaData(&_Contract.TransactOpts, _entity, _id)
}

// GetMetaData is a paid mutator transaction binding the contract method 0x1b9aa7e6.
//
// Solidity: function getMetaData(string _entity, uint256 _id) returns(address poster, string uri)
func (_Contract *ContractTransactorSession) GetMetaData(_entity string, _id *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.GetMetaData(&_Contract.TransactOpts, _entity, _id)
}

// PostData is a paid mutator transaction binding the contract method 0x5540dc2f.
//
// Solidity: function postData(string _entity, string _uri) returns(uint256)
func (_Contract *ContractTransactor) PostData(opts *bind.TransactOpts, _entity string, _uri string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "postData", _entity, _uri)
}

// PostData is a paid mutator transaction binding the contract method 0x5540dc2f.
//
// Solidity: function postData(string _entity, string _uri) returns(uint256)
func (_Contract *ContractSession) PostData(_entity string, _uri string) (*types.Transaction, error) {
	return _Contract.Contract.PostData(&_Contract.TransactOpts, _entity, _uri)
}

// PostData is a paid mutator transaction binding the contract method 0x5540dc2f.
//
// Solidity: function postData(string _entity, string _uri) returns(uint256)
func (_Contract *ContractTransactorSession) PostData(_entity string, _uri string) (*types.Transaction, error) {
	return _Contract.Contract.PostData(&_Contract.TransactOpts, _entity, _uri)
}

// SetPermanentAccessLevel is a paid mutator transaction binding the contract method 0xf46b02d1.
//
// Solidity: function setPermanentAccessLevel(address _user, uint8 _level) returns()
func (_Contract *ContractTransactor) SetPermanentAccessLevel(opts *bind.TransactOpts, _user common.Address, _level uint8) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setPermanentAccessLevel", _user, _level)
}

// SetPermanentAccessLevel is a paid mutator transaction binding the contract method 0xf46b02d1.
//
// Solidity: function setPermanentAccessLevel(address _user, uint8 _level) returns()
func (_Contract *ContractSession) SetPermanentAccessLevel(_user common.Address, _level uint8) (*types.Transaction, error) {
	return _Contract.Contract.SetPermanentAccessLevel(&_Contract.TransactOpts, _user, _level)
}

// SetPermanentAccessLevel is a paid mutator transaction binding the contract method 0xf46b02d1.
//
// Solidity: function setPermanentAccessLevel(address _user, uint8 _level) returns()
func (_Contract *ContractTransactorSession) SetPermanentAccessLevel(_user common.Address, _level uint8) (*types.Transaction, error) {
	return _Contract.Contract.SetPermanentAccessLevel(&_Contract.TransactOpts, _user, _level)
}

// SetTemporaryAccessLevel is a paid mutator transaction binding the contract method 0x6d54d4c7.
//
// Solidity: function setTemporaryAccessLevel(address _user, uint8 _level, uint256 _durationBlocks) returns()
func (_Contract *ContractTransactor) SetTemporaryAccessLevel(opts *bind.TransactOpts, _user common.Address, _level uint8, _durationBlocks *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setTemporaryAccessLevel", _user, _level, _durationBlocks)
}

// SetTemporaryAccessLevel is a paid mutator transaction binding the contract method 0x6d54d4c7.
//
// Solidity: function setTemporaryAccessLevel(address _user, uint8 _level, uint256 _durationBlocks) returns()
func (_Contract *ContractSession) SetTemporaryAccessLevel(_user common.Address, _level uint8, _durationBlocks *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetTemporaryAccessLevel(&_Contract.TransactOpts, _user, _level, _durationBlocks)
}

// SetTemporaryAccessLevel is a paid mutator transaction binding the contract method 0x6d54d4c7.
//
// Solidity: function setTemporaryAccessLevel(address _user, uint8 _level, uint256 _durationBlocks) returns()
func (_Contract *ContractTransactorSession) SetTemporaryAccessLevel(_user common.Address, _level uint8, _durationBlocks *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetTemporaryAccessLevel(&_Contract.TransactOpts, _user, _level, _durationBlocks)
}

// ContractAccessLevelUpdatedIterator is returned from FilterAccessLevelUpdated and is used to iterate over the raw logs and unpacked data for AccessLevelUpdated events raised by the Contract contract.
type ContractAccessLevelUpdatedIterator struct {
	Event *ContractAccessLevelUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAccessLevelUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAccessLevelUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAccessLevelUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAccessLevelUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAccessLevelUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAccessLevelUpdated represents a AccessLevelUpdated event raised by the Contract contract.
type ContractAccessLevelUpdated struct {
	User        common.Address
	AccessLevel uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAccessLevelUpdated is a free log retrieval operation binding the contract event 0xe7d80362d709ac8a3ab14b6ce9659e33aba3d2e123e4782476af81de61d5093e.
//
// Solidity: event AccessLevelUpdated(address indexed user, uint8 accessLevel)
func (_Contract *ContractFilterer) FilterAccessLevelUpdated(opts *bind.FilterOpts, user []common.Address) (*ContractAccessLevelUpdatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "AccessLevelUpdated", userRule)
	if err != nil {
		return nil, err
	}
	return &ContractAccessLevelUpdatedIterator{contract: _Contract.contract, event: "AccessLevelUpdated", logs: logs, sub: sub}, nil
}

// WatchAccessLevelUpdated is a free log subscription operation binding the contract event 0xe7d80362d709ac8a3ab14b6ce9659e33aba3d2e123e4782476af81de61d5093e.
//
// Solidity: event AccessLevelUpdated(address indexed user, uint8 accessLevel)
func (_Contract *ContractFilterer) WatchAccessLevelUpdated(opts *bind.WatchOpts, sink chan<- *ContractAccessLevelUpdated, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "AccessLevelUpdated", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAccessLevelUpdated)
				if err := _Contract.contract.UnpackLog(event, "AccessLevelUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAccessLevelUpdated is a log parse operation binding the contract event 0xe7d80362d709ac8a3ab14b6ce9659e33aba3d2e123e4782476af81de61d5093e.
//
// Solidity: event AccessLevelUpdated(address indexed user, uint8 accessLevel)
func (_Contract *ContractFilterer) ParseAccessLevelUpdated(log types.Log) (*ContractAccessLevelUpdated, error) {
	event := new(ContractAccessLevelUpdated)
	if err := _Contract.contract.UnpackLog(event, "AccessLevelUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDataPostedIterator is returned from FilterDataPosted and is used to iterate over the raw logs and unpacked data for DataPosted events raised by the Contract contract.
type ContractDataPostedIterator struct {
	Event *ContractDataPosted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractDataPostedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDataPosted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractDataPosted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractDataPostedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDataPostedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDataPosted represents a DataPosted event raised by the Contract contract.
type ContractDataPosted struct {
	Entity common.Hash
	Id     *big.Int
	Poster common.Address
	Uri    string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDataPosted is a free log retrieval operation binding the contract event 0x76141fea0627408a61a3a63bd9859d827b4b50e089ca7cb9d34fdb2a8777f205.
//
// Solidity: event DataPosted(string indexed entity, uint256 indexed id, address poster, string uri)
func (_Contract *ContractFilterer) FilterDataPosted(opts *bind.FilterOpts, entity []string, id []*big.Int) (*ContractDataPostedIterator, error) {

	var entityRule []interface{}
	for _, entityItem := range entity {
		entityRule = append(entityRule, entityItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "DataPosted", entityRule, idRule)
	if err != nil {
		return nil, err
	}
	return &ContractDataPostedIterator{contract: _Contract.contract, event: "DataPosted", logs: logs, sub: sub}, nil
}

// WatchDataPosted is a free log subscription operation binding the contract event 0x76141fea0627408a61a3a63bd9859d827b4b50e089ca7cb9d34fdb2a8777f205.
//
// Solidity: event DataPosted(string indexed entity, uint256 indexed id, address poster, string uri)
func (_Contract *ContractFilterer) WatchDataPosted(opts *bind.WatchOpts, sink chan<- *ContractDataPosted, entity []string, id []*big.Int) (event.Subscription, error) {

	var entityRule []interface{}
	for _, entityItem := range entity {
		entityRule = append(entityRule, entityItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "DataPosted", entityRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDataPosted)
				if err := _Contract.contract.UnpackLog(event, "DataPosted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDataPosted is a log parse operation binding the contract event 0x76141fea0627408a61a3a63bd9859d827b4b50e089ca7cb9d34fdb2a8777f205.
//
// Solidity: event DataPosted(string indexed entity, uint256 indexed id, address poster, string uri)
func (_Contract *ContractFilterer) ParseDataPosted(log types.Log) (*ContractDataPosted, error) {
	event := new(ContractDataPosted)
	if err := _Contract.contract.UnpackLog(event, "DataPosted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
