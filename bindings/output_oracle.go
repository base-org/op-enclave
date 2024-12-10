// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// TypesOutputProposal is an auto generated low-level Go binding around an user-defined struct.
type TypesOutputProposal struct {
	OutputRoot    [32]byte
	Timestamp     *big.Int
	L2BlockNumber *big.Int
}

// OutputOracleMetaData contains all meta data concerning the OutputOracle contract.
var OutputOracleMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_systemConfigGlobal\",\"type\":\"address\",\"internalType\":\"contractSystemConfigGlobal\"},{\"name\":\"_maxOutputCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"configHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"enableProofs\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getL2Output\",\"inputs\":[{\"name\":\"_l2OutputIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTypes.OutputProposal\",\"components\":[{\"name\":\"outputRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"l2BlockNumber\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getL2OutputAfter\",\"inputs\":[{\"name\":\"_l2BlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTypes.OutputProposal\",\"components\":[{\"name\":\"outputRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"l2BlockNumber\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getL2OutputIndexAfter\",\"inputs\":[{\"name\":\"_l2BlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_systemConfig\",\"type\":\"address\",\"internalType\":\"contractSystemConfigOwnable\"},{\"name\":\"_configHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_genesisOutputRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_proofsEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"latestBlockNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestL2Output\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTypes.OutputProposal\",\"components\":[{\"name\":\"outputRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"l2BlockNumber\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestOutputIndex\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxOutputCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextOutputIndex\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proofsEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposeL2Output\",\"inputs\":[{\"name\":\"_outputRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_l2BlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_l1BlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proposer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"systemConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractSystemConfigOwnable\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"systemConfigGlobal\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractSystemConfigGlobal\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OutputProposed\",\"inputs\":[{\"name\":\"outputRoot\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"l2OutputIndex\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"l2BlockNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"l1Timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
}

// OutputOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use OutputOracleMetaData.ABI instead.
var OutputOracleABI = OutputOracleMetaData.ABI

// OutputOracle is an auto generated Go binding around an Ethereum contract.
type OutputOracle struct {
	OutputOracleCaller     // Read-only binding to the contract
	OutputOracleTransactor // Write-only binding to the contract
	OutputOracleFilterer   // Log filterer for contract events
}

// OutputOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type OutputOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OutputOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OutputOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OutputOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OutputOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OutputOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OutputOracleSession struct {
	Contract     *OutputOracle     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OutputOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OutputOracleCallerSession struct {
	Contract *OutputOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// OutputOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OutputOracleTransactorSession struct {
	Contract     *OutputOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// OutputOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type OutputOracleRaw struct {
	Contract *OutputOracle // Generic contract binding to access the raw methods on
}

// OutputOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OutputOracleCallerRaw struct {
	Contract *OutputOracleCaller // Generic read-only contract binding to access the raw methods on
}

// OutputOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OutputOracleTransactorRaw struct {
	Contract *OutputOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOutputOracle creates a new instance of OutputOracle, bound to a specific deployed contract.
func NewOutputOracle(address common.Address, backend bind.ContractBackend) (*OutputOracle, error) {
	contract, err := bindOutputOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OutputOracle{OutputOracleCaller: OutputOracleCaller{contract: contract}, OutputOracleTransactor: OutputOracleTransactor{contract: contract}, OutputOracleFilterer: OutputOracleFilterer{contract: contract}}, nil
}

// NewOutputOracleCaller creates a new read-only instance of OutputOracle, bound to a specific deployed contract.
func NewOutputOracleCaller(address common.Address, caller bind.ContractCaller) (*OutputOracleCaller, error) {
	contract, err := bindOutputOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OutputOracleCaller{contract: contract}, nil
}

// NewOutputOracleTransactor creates a new write-only instance of OutputOracle, bound to a specific deployed contract.
func NewOutputOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OutputOracleTransactor, error) {
	contract, err := bindOutputOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OutputOracleTransactor{contract: contract}, nil
}

// NewOutputOracleFilterer creates a new log filterer instance of OutputOracle, bound to a specific deployed contract.
func NewOutputOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OutputOracleFilterer, error) {
	contract, err := bindOutputOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OutputOracleFilterer{contract: contract}, nil
}

// bindOutputOracle binds a generic wrapper to an already deployed contract.
func bindOutputOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OutputOracleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OutputOracle *OutputOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OutputOracle.Contract.OutputOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OutputOracle *OutputOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OutputOracle.Contract.OutputOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OutputOracle *OutputOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OutputOracle.Contract.OutputOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OutputOracle *OutputOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OutputOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OutputOracle *OutputOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OutputOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OutputOracle *OutputOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OutputOracle.Contract.contract.Transact(opts, method, params...)
}

// ConfigHash is a free data retrieval call binding the contract method 0xe1f1176d.
//
// Solidity: function configHash() view returns(bytes32)
func (_OutputOracle *OutputOracleCaller) ConfigHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "configHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ConfigHash is a free data retrieval call binding the contract method 0xe1f1176d.
//
// Solidity: function configHash() view returns(bytes32)
func (_OutputOracle *OutputOracleSession) ConfigHash() ([32]byte, error) {
	return _OutputOracle.Contract.ConfigHash(&_OutputOracle.CallOpts)
}

// ConfigHash is a free data retrieval call binding the contract method 0xe1f1176d.
//
// Solidity: function configHash() view returns(bytes32)
func (_OutputOracle *OutputOracleCallerSession) ConfigHash() ([32]byte, error) {
	return _OutputOracle.Contract.ConfigHash(&_OutputOracle.CallOpts)
}

// GetL2Output is a free data retrieval call binding the contract method 0xa25ae557.
//
// Solidity: function getL2Output(uint256 _l2OutputIndex) view returns((bytes32,uint128,uint128))
func (_OutputOracle *OutputOracleCaller) GetL2Output(opts *bind.CallOpts, _l2OutputIndex *big.Int) (TypesOutputProposal, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "getL2Output", _l2OutputIndex)

	if err != nil {
		return *new(TypesOutputProposal), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesOutputProposal)).(*TypesOutputProposal)

	return out0, err

}

// GetL2Output is a free data retrieval call binding the contract method 0xa25ae557.
//
// Solidity: function getL2Output(uint256 _l2OutputIndex) view returns((bytes32,uint128,uint128))
func (_OutputOracle *OutputOracleSession) GetL2Output(_l2OutputIndex *big.Int) (TypesOutputProposal, error) {
	return _OutputOracle.Contract.GetL2Output(&_OutputOracle.CallOpts, _l2OutputIndex)
}

// GetL2Output is a free data retrieval call binding the contract method 0xa25ae557.
//
// Solidity: function getL2Output(uint256 _l2OutputIndex) view returns((bytes32,uint128,uint128))
func (_OutputOracle *OutputOracleCallerSession) GetL2Output(_l2OutputIndex *big.Int) (TypesOutputProposal, error) {
	return _OutputOracle.Contract.GetL2Output(&_OutputOracle.CallOpts, _l2OutputIndex)
}

// GetL2OutputAfter is a free data retrieval call binding the contract method 0xcf8e5cf0.
//
// Solidity: function getL2OutputAfter(uint256 _l2BlockNumber) view returns((bytes32,uint128,uint128))
func (_OutputOracle *OutputOracleCaller) GetL2OutputAfter(opts *bind.CallOpts, _l2BlockNumber *big.Int) (TypesOutputProposal, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "getL2OutputAfter", _l2BlockNumber)

	if err != nil {
		return *new(TypesOutputProposal), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesOutputProposal)).(*TypesOutputProposal)

	return out0, err

}

// GetL2OutputAfter is a free data retrieval call binding the contract method 0xcf8e5cf0.
//
// Solidity: function getL2OutputAfter(uint256 _l2BlockNumber) view returns((bytes32,uint128,uint128))
func (_OutputOracle *OutputOracleSession) GetL2OutputAfter(_l2BlockNumber *big.Int) (TypesOutputProposal, error) {
	return _OutputOracle.Contract.GetL2OutputAfter(&_OutputOracle.CallOpts, _l2BlockNumber)
}

// GetL2OutputAfter is a free data retrieval call binding the contract method 0xcf8e5cf0.
//
// Solidity: function getL2OutputAfter(uint256 _l2BlockNumber) view returns((bytes32,uint128,uint128))
func (_OutputOracle *OutputOracleCallerSession) GetL2OutputAfter(_l2BlockNumber *big.Int) (TypesOutputProposal, error) {
	return _OutputOracle.Contract.GetL2OutputAfter(&_OutputOracle.CallOpts, _l2BlockNumber)
}

// GetL2OutputIndexAfter is a free data retrieval call binding the contract method 0x7f006420.
//
// Solidity: function getL2OutputIndexAfter(uint256 _l2BlockNumber) view returns(uint256)
func (_OutputOracle *OutputOracleCaller) GetL2OutputIndexAfter(opts *bind.CallOpts, _l2BlockNumber *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "getL2OutputIndexAfter", _l2BlockNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetL2OutputIndexAfter is a free data retrieval call binding the contract method 0x7f006420.
//
// Solidity: function getL2OutputIndexAfter(uint256 _l2BlockNumber) view returns(uint256)
func (_OutputOracle *OutputOracleSession) GetL2OutputIndexAfter(_l2BlockNumber *big.Int) (*big.Int, error) {
	return _OutputOracle.Contract.GetL2OutputIndexAfter(&_OutputOracle.CallOpts, _l2BlockNumber)
}

// GetL2OutputIndexAfter is a free data retrieval call binding the contract method 0x7f006420.
//
// Solidity: function getL2OutputIndexAfter(uint256 _l2BlockNumber) view returns(uint256)
func (_OutputOracle *OutputOracleCallerSession) GetL2OutputIndexAfter(_l2BlockNumber *big.Int) (*big.Int, error) {
	return _OutputOracle.Contract.GetL2OutputIndexAfter(&_OutputOracle.CallOpts, _l2BlockNumber)
}

// LatestBlockNumber is a free data retrieval call binding the contract method 0x4599c788.
//
// Solidity: function latestBlockNumber() view returns(uint256)
func (_OutputOracle *OutputOracleCaller) LatestBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "latestBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestBlockNumber is a free data retrieval call binding the contract method 0x4599c788.
//
// Solidity: function latestBlockNumber() view returns(uint256)
func (_OutputOracle *OutputOracleSession) LatestBlockNumber() (*big.Int, error) {
	return _OutputOracle.Contract.LatestBlockNumber(&_OutputOracle.CallOpts)
}

// LatestBlockNumber is a free data retrieval call binding the contract method 0x4599c788.
//
// Solidity: function latestBlockNumber() view returns(uint256)
func (_OutputOracle *OutputOracleCallerSession) LatestBlockNumber() (*big.Int, error) {
	return _OutputOracle.Contract.LatestBlockNumber(&_OutputOracle.CallOpts)
}

// LatestL2Output is a free data retrieval call binding the contract method 0xc885bbb6.
//
// Solidity: function latestL2Output() view returns((bytes32,uint128,uint128))
func (_OutputOracle *OutputOracleCaller) LatestL2Output(opts *bind.CallOpts) (TypesOutputProposal, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "latestL2Output")

	if err != nil {
		return *new(TypesOutputProposal), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesOutputProposal)).(*TypesOutputProposal)

	return out0, err

}

// LatestL2Output is a free data retrieval call binding the contract method 0xc885bbb6.
//
// Solidity: function latestL2Output() view returns((bytes32,uint128,uint128))
func (_OutputOracle *OutputOracleSession) LatestL2Output() (TypesOutputProposal, error) {
	return _OutputOracle.Contract.LatestL2Output(&_OutputOracle.CallOpts)
}

// LatestL2Output is a free data retrieval call binding the contract method 0xc885bbb6.
//
// Solidity: function latestL2Output() view returns((bytes32,uint128,uint128))
func (_OutputOracle *OutputOracleCallerSession) LatestL2Output() (TypesOutputProposal, error) {
	return _OutputOracle.Contract.LatestL2Output(&_OutputOracle.CallOpts)
}

// LatestOutputIndex is a free data retrieval call binding the contract method 0x69f16eec.
//
// Solidity: function latestOutputIndex() view returns(uint256)
func (_OutputOracle *OutputOracleCaller) LatestOutputIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "latestOutputIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestOutputIndex is a free data retrieval call binding the contract method 0x69f16eec.
//
// Solidity: function latestOutputIndex() view returns(uint256)
func (_OutputOracle *OutputOracleSession) LatestOutputIndex() (*big.Int, error) {
	return _OutputOracle.Contract.LatestOutputIndex(&_OutputOracle.CallOpts)
}

// LatestOutputIndex is a free data retrieval call binding the contract method 0x69f16eec.
//
// Solidity: function latestOutputIndex() view returns(uint256)
func (_OutputOracle *OutputOracleCallerSession) LatestOutputIndex() (*big.Int, error) {
	return _OutputOracle.Contract.LatestOutputIndex(&_OutputOracle.CallOpts)
}

// MaxOutputCount is a free data retrieval call binding the contract method 0xcc23c381.
//
// Solidity: function maxOutputCount() view returns(uint256)
func (_OutputOracle *OutputOracleCaller) MaxOutputCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "maxOutputCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxOutputCount is a free data retrieval call binding the contract method 0xcc23c381.
//
// Solidity: function maxOutputCount() view returns(uint256)
func (_OutputOracle *OutputOracleSession) MaxOutputCount() (*big.Int, error) {
	return _OutputOracle.Contract.MaxOutputCount(&_OutputOracle.CallOpts)
}

// MaxOutputCount is a free data retrieval call binding the contract method 0xcc23c381.
//
// Solidity: function maxOutputCount() view returns(uint256)
func (_OutputOracle *OutputOracleCallerSession) MaxOutputCount() (*big.Int, error) {
	return _OutputOracle.Contract.MaxOutputCount(&_OutputOracle.CallOpts)
}

// NextOutputIndex is a free data retrieval call binding the contract method 0x6abcf563.
//
// Solidity: function nextOutputIndex() view returns(uint256)
func (_OutputOracle *OutputOracleCaller) NextOutputIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "nextOutputIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextOutputIndex is a free data retrieval call binding the contract method 0x6abcf563.
//
// Solidity: function nextOutputIndex() view returns(uint256)
func (_OutputOracle *OutputOracleSession) NextOutputIndex() (*big.Int, error) {
	return _OutputOracle.Contract.NextOutputIndex(&_OutputOracle.CallOpts)
}

// NextOutputIndex is a free data retrieval call binding the contract method 0x6abcf563.
//
// Solidity: function nextOutputIndex() view returns(uint256)
func (_OutputOracle *OutputOracleCallerSession) NextOutputIndex() (*big.Int, error) {
	return _OutputOracle.Contract.NextOutputIndex(&_OutputOracle.CallOpts)
}

// ProofsEnabled is a free data retrieval call binding the contract method 0x25f88172.
//
// Solidity: function proofsEnabled() view returns(bool)
func (_OutputOracle *OutputOracleCaller) ProofsEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "proofsEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ProofsEnabled is a free data retrieval call binding the contract method 0x25f88172.
//
// Solidity: function proofsEnabled() view returns(bool)
func (_OutputOracle *OutputOracleSession) ProofsEnabled() (bool, error) {
	return _OutputOracle.Contract.ProofsEnabled(&_OutputOracle.CallOpts)
}

// ProofsEnabled is a free data retrieval call binding the contract method 0x25f88172.
//
// Solidity: function proofsEnabled() view returns(bool)
func (_OutputOracle *OutputOracleCallerSession) ProofsEnabled() (bool, error) {
	return _OutputOracle.Contract.ProofsEnabled(&_OutputOracle.CallOpts)
}

// Proposer is a free data retrieval call binding the contract method 0xa8e4fb90.
//
// Solidity: function proposer() view returns(address)
func (_OutputOracle *OutputOracleCaller) Proposer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "proposer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Proposer is a free data retrieval call binding the contract method 0xa8e4fb90.
//
// Solidity: function proposer() view returns(address)
func (_OutputOracle *OutputOracleSession) Proposer() (common.Address, error) {
	return _OutputOracle.Contract.Proposer(&_OutputOracle.CallOpts)
}

// Proposer is a free data retrieval call binding the contract method 0xa8e4fb90.
//
// Solidity: function proposer() view returns(address)
func (_OutputOracle *OutputOracleCallerSession) Proposer() (common.Address, error) {
	return _OutputOracle.Contract.Proposer(&_OutputOracle.CallOpts)
}

// SystemConfig is a free data retrieval call binding the contract method 0x33d7e2bd.
//
// Solidity: function systemConfig() view returns(address)
func (_OutputOracle *OutputOracleCaller) SystemConfig(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "systemConfig")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SystemConfig is a free data retrieval call binding the contract method 0x33d7e2bd.
//
// Solidity: function systemConfig() view returns(address)
func (_OutputOracle *OutputOracleSession) SystemConfig() (common.Address, error) {
	return _OutputOracle.Contract.SystemConfig(&_OutputOracle.CallOpts)
}

// SystemConfig is a free data retrieval call binding the contract method 0x33d7e2bd.
//
// Solidity: function systemConfig() view returns(address)
func (_OutputOracle *OutputOracleCallerSession) SystemConfig() (common.Address, error) {
	return _OutputOracle.Contract.SystemConfig(&_OutputOracle.CallOpts)
}

// SystemConfigGlobal is a free data retrieval call binding the contract method 0xcd92b3fe.
//
// Solidity: function systemConfigGlobal() view returns(address)
func (_OutputOracle *OutputOracleCaller) SystemConfigGlobal(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "systemConfigGlobal")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SystemConfigGlobal is a free data retrieval call binding the contract method 0xcd92b3fe.
//
// Solidity: function systemConfigGlobal() view returns(address)
func (_OutputOracle *OutputOracleSession) SystemConfigGlobal() (common.Address, error) {
	return _OutputOracle.Contract.SystemConfigGlobal(&_OutputOracle.CallOpts)
}

// SystemConfigGlobal is a free data retrieval call binding the contract method 0xcd92b3fe.
//
// Solidity: function systemConfigGlobal() view returns(address)
func (_OutputOracle *OutputOracleCallerSession) SystemConfigGlobal() (common.Address, error) {
	return _OutputOracle.Contract.SystemConfigGlobal(&_OutputOracle.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_OutputOracle *OutputOracleCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OutputOracle.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_OutputOracle *OutputOracleSession) Version() (string, error) {
	return _OutputOracle.Contract.Version(&_OutputOracle.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_OutputOracle *OutputOracleCallerSession) Version() (string, error) {
	return _OutputOracle.Contract.Version(&_OutputOracle.CallOpts)
}

// EnableProofs is a paid mutator transaction binding the contract method 0x60df09b2.
//
// Solidity: function enableProofs() returns()
func (_OutputOracle *OutputOracleTransactor) EnableProofs(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OutputOracle.contract.Transact(opts, "enableProofs")
}

// EnableProofs is a paid mutator transaction binding the contract method 0x60df09b2.
//
// Solidity: function enableProofs() returns()
func (_OutputOracle *OutputOracleSession) EnableProofs() (*types.Transaction, error) {
	return _OutputOracle.Contract.EnableProofs(&_OutputOracle.TransactOpts)
}

// EnableProofs is a paid mutator transaction binding the contract method 0x60df09b2.
//
// Solidity: function enableProofs() returns()
func (_OutputOracle *OutputOracleTransactorSession) EnableProofs() (*types.Transaction, error) {
	return _OutputOracle.Contract.EnableProofs(&_OutputOracle.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xb8205148.
//
// Solidity: function initialize(address _systemConfig, bytes32 _configHash, bytes32 _genesisOutputRoot, bool _proofsEnabled) returns()
func (_OutputOracle *OutputOracleTransactor) Initialize(opts *bind.TransactOpts, _systemConfig common.Address, _configHash [32]byte, _genesisOutputRoot [32]byte, _proofsEnabled bool) (*types.Transaction, error) {
	return _OutputOracle.contract.Transact(opts, "initialize", _systemConfig, _configHash, _genesisOutputRoot, _proofsEnabled)
}

// Initialize is a paid mutator transaction binding the contract method 0xb8205148.
//
// Solidity: function initialize(address _systemConfig, bytes32 _configHash, bytes32 _genesisOutputRoot, bool _proofsEnabled) returns()
func (_OutputOracle *OutputOracleSession) Initialize(_systemConfig common.Address, _configHash [32]byte, _genesisOutputRoot [32]byte, _proofsEnabled bool) (*types.Transaction, error) {
	return _OutputOracle.Contract.Initialize(&_OutputOracle.TransactOpts, _systemConfig, _configHash, _genesisOutputRoot, _proofsEnabled)
}

// Initialize is a paid mutator transaction binding the contract method 0xb8205148.
//
// Solidity: function initialize(address _systemConfig, bytes32 _configHash, bytes32 _genesisOutputRoot, bool _proofsEnabled) returns()
func (_OutputOracle *OutputOracleTransactorSession) Initialize(_systemConfig common.Address, _configHash [32]byte, _genesisOutputRoot [32]byte, _proofsEnabled bool) (*types.Transaction, error) {
	return _OutputOracle.Contract.Initialize(&_OutputOracle.TransactOpts, _systemConfig, _configHash, _genesisOutputRoot, _proofsEnabled)
}

// ProposeL2Output is a paid mutator transaction binding the contract method 0x9ad84880.
//
// Solidity: function proposeL2Output(bytes32 _outputRoot, uint256 _l2BlockNumber, uint256 _l1BlockNumber, bytes _signature) returns()
func (_OutputOracle *OutputOracleTransactor) ProposeL2Output(opts *bind.TransactOpts, _outputRoot [32]byte, _l2BlockNumber *big.Int, _l1BlockNumber *big.Int, _signature []byte) (*types.Transaction, error) {
	return _OutputOracle.contract.Transact(opts, "proposeL2Output", _outputRoot, _l2BlockNumber, _l1BlockNumber, _signature)
}

// ProposeL2Output is a paid mutator transaction binding the contract method 0x9ad84880.
//
// Solidity: function proposeL2Output(bytes32 _outputRoot, uint256 _l2BlockNumber, uint256 _l1BlockNumber, bytes _signature) returns()
func (_OutputOracle *OutputOracleSession) ProposeL2Output(_outputRoot [32]byte, _l2BlockNumber *big.Int, _l1BlockNumber *big.Int, _signature []byte) (*types.Transaction, error) {
	return _OutputOracle.Contract.ProposeL2Output(&_OutputOracle.TransactOpts, _outputRoot, _l2BlockNumber, _l1BlockNumber, _signature)
}

// ProposeL2Output is a paid mutator transaction binding the contract method 0x9ad84880.
//
// Solidity: function proposeL2Output(bytes32 _outputRoot, uint256 _l2BlockNumber, uint256 _l1BlockNumber, bytes _signature) returns()
func (_OutputOracle *OutputOracleTransactorSession) ProposeL2Output(_outputRoot [32]byte, _l2BlockNumber *big.Int, _l1BlockNumber *big.Int, _signature []byte) (*types.Transaction, error) {
	return _OutputOracle.Contract.ProposeL2Output(&_OutputOracle.TransactOpts, _outputRoot, _l2BlockNumber, _l1BlockNumber, _signature)
}

// OutputOracleInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the OutputOracle contract.
type OutputOracleInitializedIterator struct {
	Event *OutputOracleInitialized // Event containing the contract specifics and raw log

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
func (it *OutputOracleInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OutputOracleInitialized)
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
		it.Event = new(OutputOracleInitialized)
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
func (it *OutputOracleInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OutputOracleInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OutputOracleInitialized represents a Initialized event raised by the OutputOracle contract.
type OutputOracleInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_OutputOracle *OutputOracleFilterer) FilterInitialized(opts *bind.FilterOpts) (*OutputOracleInitializedIterator, error) {

	logs, sub, err := _OutputOracle.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &OutputOracleInitializedIterator{contract: _OutputOracle.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_OutputOracle *OutputOracleFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *OutputOracleInitialized) (event.Subscription, error) {

	logs, sub, err := _OutputOracle.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OutputOracleInitialized)
				if err := _OutputOracle.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_OutputOracle *OutputOracleFilterer) ParseInitialized(log types.Log) (*OutputOracleInitialized, error) {
	event := new(OutputOracleInitialized)
	if err := _OutputOracle.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OutputOracleOutputProposedIterator is returned from FilterOutputProposed and is used to iterate over the raw logs and unpacked data for OutputProposed events raised by the OutputOracle contract.
type OutputOracleOutputProposedIterator struct {
	Event *OutputOracleOutputProposed // Event containing the contract specifics and raw log

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
func (it *OutputOracleOutputProposedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OutputOracleOutputProposed)
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
		it.Event = new(OutputOracleOutputProposed)
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
func (it *OutputOracleOutputProposedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OutputOracleOutputProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OutputOracleOutputProposed represents a OutputProposed event raised by the OutputOracle contract.
type OutputOracleOutputProposed struct {
	OutputRoot    [32]byte
	L2OutputIndex *big.Int
	L2BlockNumber *big.Int
	L1Timestamp   *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOutputProposed is a free log retrieval operation binding the contract event 0xa7aaf2512769da4e444e3de247be2564225c2e7a8f74cfe528e46e17d24868e2.
//
// Solidity: event OutputProposed(bytes32 indexed outputRoot, uint256 indexed l2OutputIndex, uint256 indexed l2BlockNumber, uint256 l1Timestamp)
func (_OutputOracle *OutputOracleFilterer) FilterOutputProposed(opts *bind.FilterOpts, outputRoot [][32]byte, l2OutputIndex []*big.Int, l2BlockNumber []*big.Int) (*OutputOracleOutputProposedIterator, error) {

	var outputRootRule []interface{}
	for _, outputRootItem := range outputRoot {
		outputRootRule = append(outputRootRule, outputRootItem)
	}
	var l2OutputIndexRule []interface{}
	for _, l2OutputIndexItem := range l2OutputIndex {
		l2OutputIndexRule = append(l2OutputIndexRule, l2OutputIndexItem)
	}
	var l2BlockNumberRule []interface{}
	for _, l2BlockNumberItem := range l2BlockNumber {
		l2BlockNumberRule = append(l2BlockNumberRule, l2BlockNumberItem)
	}

	logs, sub, err := _OutputOracle.contract.FilterLogs(opts, "OutputProposed", outputRootRule, l2OutputIndexRule, l2BlockNumberRule)
	if err != nil {
		return nil, err
	}
	return &OutputOracleOutputProposedIterator{contract: _OutputOracle.contract, event: "OutputProposed", logs: logs, sub: sub}, nil
}

// WatchOutputProposed is a free log subscription operation binding the contract event 0xa7aaf2512769da4e444e3de247be2564225c2e7a8f74cfe528e46e17d24868e2.
//
// Solidity: event OutputProposed(bytes32 indexed outputRoot, uint256 indexed l2OutputIndex, uint256 indexed l2BlockNumber, uint256 l1Timestamp)
func (_OutputOracle *OutputOracleFilterer) WatchOutputProposed(opts *bind.WatchOpts, sink chan<- *OutputOracleOutputProposed, outputRoot [][32]byte, l2OutputIndex []*big.Int, l2BlockNumber []*big.Int) (event.Subscription, error) {

	var outputRootRule []interface{}
	for _, outputRootItem := range outputRoot {
		outputRootRule = append(outputRootRule, outputRootItem)
	}
	var l2OutputIndexRule []interface{}
	for _, l2OutputIndexItem := range l2OutputIndex {
		l2OutputIndexRule = append(l2OutputIndexRule, l2OutputIndexItem)
	}
	var l2BlockNumberRule []interface{}
	for _, l2BlockNumberItem := range l2BlockNumber {
		l2BlockNumberRule = append(l2BlockNumberRule, l2BlockNumberItem)
	}

	logs, sub, err := _OutputOracle.contract.WatchLogs(opts, "OutputProposed", outputRootRule, l2OutputIndexRule, l2BlockNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OutputOracleOutputProposed)
				if err := _OutputOracle.contract.UnpackLog(event, "OutputProposed", log); err != nil {
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

// ParseOutputProposed is a log parse operation binding the contract event 0xa7aaf2512769da4e444e3de247be2564225c2e7a8f74cfe528e46e17d24868e2.
//
// Solidity: event OutputProposed(bytes32 indexed outputRoot, uint256 indexed l2OutputIndex, uint256 indexed l2BlockNumber, uint256 l1Timestamp)
func (_OutputOracle *OutputOracleFilterer) ParseOutputProposed(log types.Log) (*OutputOracleOutputProposed, error) {
	event := new(OutputOracleOutputProposed)
	if err := _OutputOracle.contract.UnpackLog(event, "OutputProposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
