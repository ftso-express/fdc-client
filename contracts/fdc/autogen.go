// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package fdc

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

// FdcMetaData contains all meta data concerning the Fdc contract.
var FdcMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIGovernanceSettings\",\"name\":\"_governanceSettings\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_initialGovernance\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_addressUpdater\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"AttestationRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"authorizedAmountWei\",\"type\":\"uint256\"}],\"name\":\"DailyAuthorizedInflationSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"allowedAfterTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encodedCall\",\"type\":\"bytes\"}],\"name\":\"GovernanceCallTimelocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"initialGovernance\",\"type\":\"address\"}],\"name\":\"GovernanceInitialised\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"governanceSettings\",\"type\":\"address\"}],\"name\":\"GovernedProductionModeEntered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountReceivedWei\",\"type\":\"uint256\"}],\"name\":\"InflationReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint24\",\"name\":\"rewardEpochId\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"InflationRewardsOffered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"TimelockedGovernanceCallCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"TimelockedGovernanceCallExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_type\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"source\",\"type\":\"bytes32\"}],\"name\":\"TypeAndSourceFeeRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_type\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"source\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"TypeAndSourceFeeSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_selector\",\"type\":\"bytes4\"}],\"name\":\"cancelGovernanceCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dailyAuthorizedInflation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_selector\",\"type\":\"bytes4\"}],\"name\":\"executeGovernanceCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"flareSystemsManager\",\"outputs\":[{\"internalType\":\"contractIIFlareSystemsManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAddressUpdater\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_addressUpdater\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContractName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExpectedBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getInflationAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"getRequestFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenPoolSupplyData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_lockedFundsWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalInflationAuthorizedWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalClaimedWei\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"governance\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"governanceSettings\",\"outputs\":[{\"internalType\":\"contractIGovernanceSettings\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIGovernanceSettings\",\"name\":\"_governanceSettings\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_initialGovernance\",\"type\":\"address\"}],\"name\":\"initialise\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"isExecutor\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastInflationAuthorizationReceivedTs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastInflationReceivedTs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"productionMode\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"receiveInflation\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_type\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_source\",\"type\":\"bytes32\"}],\"name\":\"removeTypeAndSourceFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_types\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_sources\",\"type\":\"bytes32[]\"}],\"name\":\"removeTypeAndSourceFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"requestAttestation\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardManager\",\"outputs\":[{\"internalType\":\"contractIIRewardManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_toAuthorizeWei\",\"type\":\"uint256\"}],\"name\":\"setDailyAuthorizedInflation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_type\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_source\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"setTypeAndSourceFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_types\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_sources\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_fees\",\"type\":\"uint256[]\"}],\"name\":\"setTypeAndSourceFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"switchToProductionMode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"}],\"name\":\"timelockedCalls\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"allowedAfterTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"encodedCall\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalInflationAuthorizedWei\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalInflationReceivedWei\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalInflationRewardsOfferedWei\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint24\",\"name\":\"_currentRewardEpochId\",\"type\":\"uint24\"},{\"internalType\":\"uint64\",\"name\":\"_currentRewardEpochExpectedEndTs\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_rewardEpochDurationSeconds\",\"type\":\"uint64\"}],\"name\":\"triggerRewardEpochSwitchover\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"typeAndSource\",\"type\":\"bytes32\"}],\"name\":\"typeAndSourceFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_contractNameHashes\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"_contractAddresses\",\"type\":\"address[]\"}],\"name\":\"updateContractAddresses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// FdcABI is the input ABI used to generate the binding from.
// Deprecated: Use FdcMetaData.ABI instead.
var FdcABI = FdcMetaData.ABI

// Fdc is an auto generated Go binding around an Ethereum contract.
type Fdc struct {
	FdcCaller     // Read-only binding to the contract
	FdcTransactor // Write-only binding to the contract
	FdcFilterer   // Log filterer for contract events
}

// FdcCaller is an auto generated read-only Go binding around an Ethereum contract.
type FdcCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FdcTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FdcTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FdcFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FdcFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FdcSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FdcSession struct {
	Contract     *Fdc              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FdcCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FdcCallerSession struct {
	Contract *FdcCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FdcTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FdcTransactorSession struct {
	Contract     *FdcTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FdcRaw is an auto generated low-level Go binding around an Ethereum contract.
type FdcRaw struct {
	Contract *Fdc // Generic contract binding to access the raw methods on
}

// FdcCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FdcCallerRaw struct {
	Contract *FdcCaller // Generic read-only contract binding to access the raw methods on
}

// FdcTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FdcTransactorRaw struct {
	Contract *FdcTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFdc creates a new instance of Fdc, bound to a specific deployed contract.
func NewFdc(address common.Address, backend bind.ContractBackend) (*Fdc, error) {
	contract, err := bindFdc(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Fdc{FdcCaller: FdcCaller{contract: contract}, FdcTransactor: FdcTransactor{contract: contract}, FdcFilterer: FdcFilterer{contract: contract}}, nil
}

// NewFdcCaller creates a new read-only instance of Fdc, bound to a specific deployed contract.
func NewFdcCaller(address common.Address, caller bind.ContractCaller) (*FdcCaller, error) {
	contract, err := bindFdc(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FdcCaller{contract: contract}, nil
}

// NewFdcTransactor creates a new write-only instance of Fdc, bound to a specific deployed contract.
func NewFdcTransactor(address common.Address, transactor bind.ContractTransactor) (*FdcTransactor, error) {
	contract, err := bindFdc(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FdcTransactor{contract: contract}, nil
}

// NewFdcFilterer creates a new log filterer instance of Fdc, bound to a specific deployed contract.
func NewFdcFilterer(address common.Address, filterer bind.ContractFilterer) (*FdcFilterer, error) {
	contract, err := bindFdc(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FdcFilterer{contract: contract}, nil
}

// bindFdc binds a generic wrapper to an already deployed contract.
func bindFdc(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FdcMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Fdc *FdcRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Fdc.Contract.FdcCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Fdc *FdcRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fdc.Contract.FdcTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Fdc *FdcRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Fdc.Contract.FdcTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Fdc *FdcCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Fdc.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Fdc *FdcTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fdc.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Fdc *FdcTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Fdc.Contract.contract.Transact(opts, method, params...)
}

// DailyAuthorizedInflation is a free data retrieval call binding the contract method 0x708e34ce.
//
// Solidity: function dailyAuthorizedInflation() view returns(uint256)
func (_Fdc *FdcCaller) DailyAuthorizedInflation(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "dailyAuthorizedInflation")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DailyAuthorizedInflation is a free data retrieval call binding the contract method 0x708e34ce.
//
// Solidity: function dailyAuthorizedInflation() view returns(uint256)
func (_Fdc *FdcSession) DailyAuthorizedInflation() (*big.Int, error) {
	return _Fdc.Contract.DailyAuthorizedInflation(&_Fdc.CallOpts)
}

// DailyAuthorizedInflation is a free data retrieval call binding the contract method 0x708e34ce.
//
// Solidity: function dailyAuthorizedInflation() view returns(uint256)
func (_Fdc *FdcCallerSession) DailyAuthorizedInflation() (*big.Int, error) {
	return _Fdc.Contract.DailyAuthorizedInflation(&_Fdc.CallOpts)
}

// FlareSystemsManager is a free data retrieval call binding the contract method 0xfaae7fc9.
//
// Solidity: function flareSystemsManager() view returns(address)
func (_Fdc *FdcCaller) FlareSystemsManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "flareSystemsManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FlareSystemsManager is a free data retrieval call binding the contract method 0xfaae7fc9.
//
// Solidity: function flareSystemsManager() view returns(address)
func (_Fdc *FdcSession) FlareSystemsManager() (common.Address, error) {
	return _Fdc.Contract.FlareSystemsManager(&_Fdc.CallOpts)
}

// FlareSystemsManager is a free data retrieval call binding the contract method 0xfaae7fc9.
//
// Solidity: function flareSystemsManager() view returns(address)
func (_Fdc *FdcCallerSession) FlareSystemsManager() (common.Address, error) {
	return _Fdc.Contract.FlareSystemsManager(&_Fdc.CallOpts)
}

// GetAddressUpdater is a free data retrieval call binding the contract method 0x5267a15d.
//
// Solidity: function getAddressUpdater() view returns(address _addressUpdater)
func (_Fdc *FdcCaller) GetAddressUpdater(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "getAddressUpdater")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAddressUpdater is a free data retrieval call binding the contract method 0x5267a15d.
//
// Solidity: function getAddressUpdater() view returns(address _addressUpdater)
func (_Fdc *FdcSession) GetAddressUpdater() (common.Address, error) {
	return _Fdc.Contract.GetAddressUpdater(&_Fdc.CallOpts)
}

// GetAddressUpdater is a free data retrieval call binding the contract method 0x5267a15d.
//
// Solidity: function getAddressUpdater() view returns(address _addressUpdater)
func (_Fdc *FdcCallerSession) GetAddressUpdater() (common.Address, error) {
	return _Fdc.Contract.GetAddressUpdater(&_Fdc.CallOpts)
}

// GetContractName is a free data retrieval call binding the contract method 0xf5f5ba72.
//
// Solidity: function getContractName() pure returns(string)
func (_Fdc *FdcCaller) GetContractName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "getContractName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetContractName is a free data retrieval call binding the contract method 0xf5f5ba72.
//
// Solidity: function getContractName() pure returns(string)
func (_Fdc *FdcSession) GetContractName() (string, error) {
	return _Fdc.Contract.GetContractName(&_Fdc.CallOpts)
}

// GetContractName is a free data retrieval call binding the contract method 0xf5f5ba72.
//
// Solidity: function getContractName() pure returns(string)
func (_Fdc *FdcCallerSession) GetContractName() (string, error) {
	return _Fdc.Contract.GetContractName(&_Fdc.CallOpts)
}

// GetExpectedBalance is a free data retrieval call binding the contract method 0xaf04cd3b.
//
// Solidity: function getExpectedBalance() view returns(uint256)
func (_Fdc *FdcCaller) GetExpectedBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "getExpectedBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetExpectedBalance is a free data retrieval call binding the contract method 0xaf04cd3b.
//
// Solidity: function getExpectedBalance() view returns(uint256)
func (_Fdc *FdcSession) GetExpectedBalance() (*big.Int, error) {
	return _Fdc.Contract.GetExpectedBalance(&_Fdc.CallOpts)
}

// GetExpectedBalance is a free data retrieval call binding the contract method 0xaf04cd3b.
//
// Solidity: function getExpectedBalance() view returns(uint256)
func (_Fdc *FdcCallerSession) GetExpectedBalance() (*big.Int, error) {
	return _Fdc.Contract.GetExpectedBalance(&_Fdc.CallOpts)
}

// GetInflationAddress is a free data retrieval call binding the contract method 0xed39d3f8.
//
// Solidity: function getInflationAddress() view returns(address)
func (_Fdc *FdcCaller) GetInflationAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "getInflationAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetInflationAddress is a free data retrieval call binding the contract method 0xed39d3f8.
//
// Solidity: function getInflationAddress() view returns(address)
func (_Fdc *FdcSession) GetInflationAddress() (common.Address, error) {
	return _Fdc.Contract.GetInflationAddress(&_Fdc.CallOpts)
}

// GetInflationAddress is a free data retrieval call binding the contract method 0xed39d3f8.
//
// Solidity: function getInflationAddress() view returns(address)
func (_Fdc *FdcCallerSession) GetInflationAddress() (common.Address, error) {
	return _Fdc.Contract.GetInflationAddress(&_Fdc.CallOpts)
}

// GetRequestFee is a free data retrieval call binding the contract method 0x0a0f2476.
//
// Solidity: function getRequestFee(bytes _data) view returns(uint256)
func (_Fdc *FdcCaller) GetRequestFee(opts *bind.CallOpts, _data []byte) (*big.Int, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "getRequestFee", _data)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRequestFee is a free data retrieval call binding the contract method 0x0a0f2476.
//
// Solidity: function getRequestFee(bytes _data) view returns(uint256)
func (_Fdc *FdcSession) GetRequestFee(_data []byte) (*big.Int, error) {
	return _Fdc.Contract.GetRequestFee(&_Fdc.CallOpts, _data)
}

// GetRequestFee is a free data retrieval call binding the contract method 0x0a0f2476.
//
// Solidity: function getRequestFee(bytes _data) view returns(uint256)
func (_Fdc *FdcCallerSession) GetRequestFee(_data []byte) (*big.Int, error) {
	return _Fdc.Contract.GetRequestFee(&_Fdc.CallOpts, _data)
}

// GetTokenPoolSupplyData is a free data retrieval call binding the contract method 0x2dafdbbf.
//
// Solidity: function getTokenPoolSupplyData() view returns(uint256 _lockedFundsWei, uint256 _totalInflationAuthorizedWei, uint256 _totalClaimedWei)
func (_Fdc *FdcCaller) GetTokenPoolSupplyData(opts *bind.CallOpts) (struct {
	LockedFundsWei              *big.Int
	TotalInflationAuthorizedWei *big.Int
	TotalClaimedWei             *big.Int
}, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "getTokenPoolSupplyData")

	outstruct := new(struct {
		LockedFundsWei              *big.Int
		TotalInflationAuthorizedWei *big.Int
		TotalClaimedWei             *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.LockedFundsWei = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalInflationAuthorizedWei = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TotalClaimedWei = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetTokenPoolSupplyData is a free data retrieval call binding the contract method 0x2dafdbbf.
//
// Solidity: function getTokenPoolSupplyData() view returns(uint256 _lockedFundsWei, uint256 _totalInflationAuthorizedWei, uint256 _totalClaimedWei)
func (_Fdc *FdcSession) GetTokenPoolSupplyData() (struct {
	LockedFundsWei              *big.Int
	TotalInflationAuthorizedWei *big.Int
	TotalClaimedWei             *big.Int
}, error) {
	return _Fdc.Contract.GetTokenPoolSupplyData(&_Fdc.CallOpts)
}

// GetTokenPoolSupplyData is a free data retrieval call binding the contract method 0x2dafdbbf.
//
// Solidity: function getTokenPoolSupplyData() view returns(uint256 _lockedFundsWei, uint256 _totalInflationAuthorizedWei, uint256 _totalClaimedWei)
func (_Fdc *FdcCallerSession) GetTokenPoolSupplyData() (struct {
	LockedFundsWei              *big.Int
	TotalInflationAuthorizedWei *big.Int
	TotalClaimedWei             *big.Int
}, error) {
	return _Fdc.Contract.GetTokenPoolSupplyData(&_Fdc.CallOpts)
}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() view returns(address)
func (_Fdc *FdcCaller) Governance(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "governance")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() view returns(address)
func (_Fdc *FdcSession) Governance() (common.Address, error) {
	return _Fdc.Contract.Governance(&_Fdc.CallOpts)
}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() view returns(address)
func (_Fdc *FdcCallerSession) Governance() (common.Address, error) {
	return _Fdc.Contract.Governance(&_Fdc.CallOpts)
}

// GovernanceSettings is a free data retrieval call binding the contract method 0x62354e03.
//
// Solidity: function governanceSettings() view returns(address)
func (_Fdc *FdcCaller) GovernanceSettings(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "governanceSettings")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GovernanceSettings is a free data retrieval call binding the contract method 0x62354e03.
//
// Solidity: function governanceSettings() view returns(address)
func (_Fdc *FdcSession) GovernanceSettings() (common.Address, error) {
	return _Fdc.Contract.GovernanceSettings(&_Fdc.CallOpts)
}

// GovernanceSettings is a free data retrieval call binding the contract method 0x62354e03.
//
// Solidity: function governanceSettings() view returns(address)
func (_Fdc *FdcCallerSession) GovernanceSettings() (common.Address, error) {
	return _Fdc.Contract.GovernanceSettings(&_Fdc.CallOpts)
}

// IsExecutor is a free data retrieval call binding the contract method 0xdebfda30.
//
// Solidity: function isExecutor(address _address) view returns(bool)
func (_Fdc *FdcCaller) IsExecutor(opts *bind.CallOpts, _address common.Address) (bool, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "isExecutor", _address)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsExecutor is a free data retrieval call binding the contract method 0xdebfda30.
//
// Solidity: function isExecutor(address _address) view returns(bool)
func (_Fdc *FdcSession) IsExecutor(_address common.Address) (bool, error) {
	return _Fdc.Contract.IsExecutor(&_Fdc.CallOpts, _address)
}

// IsExecutor is a free data retrieval call binding the contract method 0xdebfda30.
//
// Solidity: function isExecutor(address _address) view returns(bool)
func (_Fdc *FdcCallerSession) IsExecutor(_address common.Address) (bool, error) {
	return _Fdc.Contract.IsExecutor(&_Fdc.CallOpts, _address)
}

// LastInflationAuthorizationReceivedTs is a free data retrieval call binding the contract method 0x473252c4.
//
// Solidity: function lastInflationAuthorizationReceivedTs() view returns(uint256)
func (_Fdc *FdcCaller) LastInflationAuthorizationReceivedTs(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "lastInflationAuthorizationReceivedTs")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastInflationAuthorizationReceivedTs is a free data retrieval call binding the contract method 0x473252c4.
//
// Solidity: function lastInflationAuthorizationReceivedTs() view returns(uint256)
func (_Fdc *FdcSession) LastInflationAuthorizationReceivedTs() (*big.Int, error) {
	return _Fdc.Contract.LastInflationAuthorizationReceivedTs(&_Fdc.CallOpts)
}

// LastInflationAuthorizationReceivedTs is a free data retrieval call binding the contract method 0x473252c4.
//
// Solidity: function lastInflationAuthorizationReceivedTs() view returns(uint256)
func (_Fdc *FdcCallerSession) LastInflationAuthorizationReceivedTs() (*big.Int, error) {
	return _Fdc.Contract.LastInflationAuthorizationReceivedTs(&_Fdc.CallOpts)
}

// LastInflationReceivedTs is a free data retrieval call binding the contract method 0x12afcf0b.
//
// Solidity: function lastInflationReceivedTs() view returns(uint256)
func (_Fdc *FdcCaller) LastInflationReceivedTs(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "lastInflationReceivedTs")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastInflationReceivedTs is a free data retrieval call binding the contract method 0x12afcf0b.
//
// Solidity: function lastInflationReceivedTs() view returns(uint256)
func (_Fdc *FdcSession) LastInflationReceivedTs() (*big.Int, error) {
	return _Fdc.Contract.LastInflationReceivedTs(&_Fdc.CallOpts)
}

// LastInflationReceivedTs is a free data retrieval call binding the contract method 0x12afcf0b.
//
// Solidity: function lastInflationReceivedTs() view returns(uint256)
func (_Fdc *FdcCallerSession) LastInflationReceivedTs() (*big.Int, error) {
	return _Fdc.Contract.LastInflationReceivedTs(&_Fdc.CallOpts)
}

// ProductionMode is a free data retrieval call binding the contract method 0xe17f212e.
//
// Solidity: function productionMode() view returns(bool)
func (_Fdc *FdcCaller) ProductionMode(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "productionMode")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ProductionMode is a free data retrieval call binding the contract method 0xe17f212e.
//
// Solidity: function productionMode() view returns(bool)
func (_Fdc *FdcSession) ProductionMode() (bool, error) {
	return _Fdc.Contract.ProductionMode(&_Fdc.CallOpts)
}

// ProductionMode is a free data retrieval call binding the contract method 0xe17f212e.
//
// Solidity: function productionMode() view returns(bool)
func (_Fdc *FdcCallerSession) ProductionMode() (bool, error) {
	return _Fdc.Contract.ProductionMode(&_Fdc.CallOpts)
}

// RewardManager is a free data retrieval call binding the contract method 0x0f4ef8a6.
//
// Solidity: function rewardManager() view returns(address)
func (_Fdc *FdcCaller) RewardManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "rewardManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardManager is a free data retrieval call binding the contract method 0x0f4ef8a6.
//
// Solidity: function rewardManager() view returns(address)
func (_Fdc *FdcSession) RewardManager() (common.Address, error) {
	return _Fdc.Contract.RewardManager(&_Fdc.CallOpts)
}

// RewardManager is a free data retrieval call binding the contract method 0x0f4ef8a6.
//
// Solidity: function rewardManager() view returns(address)
func (_Fdc *FdcCallerSession) RewardManager() (common.Address, error) {
	return _Fdc.Contract.RewardManager(&_Fdc.CallOpts)
}

// TimelockedCalls is a free data retrieval call binding the contract method 0x74e6310e.
//
// Solidity: function timelockedCalls(bytes4 selector) view returns(uint256 allowedAfterTimestamp, bytes encodedCall)
func (_Fdc *FdcCaller) TimelockedCalls(opts *bind.CallOpts, selector [4]byte) (struct {
	AllowedAfterTimestamp *big.Int
	EncodedCall           []byte
}, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "timelockedCalls", selector)

	outstruct := new(struct {
		AllowedAfterTimestamp *big.Int
		EncodedCall           []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AllowedAfterTimestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.EncodedCall = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

// TimelockedCalls is a free data retrieval call binding the contract method 0x74e6310e.
//
// Solidity: function timelockedCalls(bytes4 selector) view returns(uint256 allowedAfterTimestamp, bytes encodedCall)
func (_Fdc *FdcSession) TimelockedCalls(selector [4]byte) (struct {
	AllowedAfterTimestamp *big.Int
	EncodedCall           []byte
}, error) {
	return _Fdc.Contract.TimelockedCalls(&_Fdc.CallOpts, selector)
}

// TimelockedCalls is a free data retrieval call binding the contract method 0x74e6310e.
//
// Solidity: function timelockedCalls(bytes4 selector) view returns(uint256 allowedAfterTimestamp, bytes encodedCall)
func (_Fdc *FdcCallerSession) TimelockedCalls(selector [4]byte) (struct {
	AllowedAfterTimestamp *big.Int
	EncodedCall           []byte
}, error) {
	return _Fdc.Contract.TimelockedCalls(&_Fdc.CallOpts, selector)
}

// TotalInflationAuthorizedWei is a free data retrieval call binding the contract method 0xd0c1c393.
//
// Solidity: function totalInflationAuthorizedWei() view returns(uint256)
func (_Fdc *FdcCaller) TotalInflationAuthorizedWei(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "totalInflationAuthorizedWei")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalInflationAuthorizedWei is a free data retrieval call binding the contract method 0xd0c1c393.
//
// Solidity: function totalInflationAuthorizedWei() view returns(uint256)
func (_Fdc *FdcSession) TotalInflationAuthorizedWei() (*big.Int, error) {
	return _Fdc.Contract.TotalInflationAuthorizedWei(&_Fdc.CallOpts)
}

// TotalInflationAuthorizedWei is a free data retrieval call binding the contract method 0xd0c1c393.
//
// Solidity: function totalInflationAuthorizedWei() view returns(uint256)
func (_Fdc *FdcCallerSession) TotalInflationAuthorizedWei() (*big.Int, error) {
	return _Fdc.Contract.TotalInflationAuthorizedWei(&_Fdc.CallOpts)
}

// TotalInflationReceivedWei is a free data retrieval call binding the contract method 0xa5555aea.
//
// Solidity: function totalInflationReceivedWei() view returns(uint256)
func (_Fdc *FdcCaller) TotalInflationReceivedWei(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "totalInflationReceivedWei")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalInflationReceivedWei is a free data retrieval call binding the contract method 0xa5555aea.
//
// Solidity: function totalInflationReceivedWei() view returns(uint256)
func (_Fdc *FdcSession) TotalInflationReceivedWei() (*big.Int, error) {
	return _Fdc.Contract.TotalInflationReceivedWei(&_Fdc.CallOpts)
}

// TotalInflationReceivedWei is a free data retrieval call binding the contract method 0xa5555aea.
//
// Solidity: function totalInflationReceivedWei() view returns(uint256)
func (_Fdc *FdcCallerSession) TotalInflationReceivedWei() (*big.Int, error) {
	return _Fdc.Contract.TotalInflationReceivedWei(&_Fdc.CallOpts)
}

// TotalInflationRewardsOfferedWei is a free data retrieval call binding the contract method 0xbd76b69c.
//
// Solidity: function totalInflationRewardsOfferedWei() view returns(uint256)
func (_Fdc *FdcCaller) TotalInflationRewardsOfferedWei(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "totalInflationRewardsOfferedWei")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalInflationRewardsOfferedWei is a free data retrieval call binding the contract method 0xbd76b69c.
//
// Solidity: function totalInflationRewardsOfferedWei() view returns(uint256)
func (_Fdc *FdcSession) TotalInflationRewardsOfferedWei() (*big.Int, error) {
	return _Fdc.Contract.TotalInflationRewardsOfferedWei(&_Fdc.CallOpts)
}

// TotalInflationRewardsOfferedWei is a free data retrieval call binding the contract method 0xbd76b69c.
//
// Solidity: function totalInflationRewardsOfferedWei() view returns(uint256)
func (_Fdc *FdcCallerSession) TotalInflationRewardsOfferedWei() (*big.Int, error) {
	return _Fdc.Contract.TotalInflationRewardsOfferedWei(&_Fdc.CallOpts)
}

// TypeAndSourceFees is a free data retrieval call binding the contract method 0xb9e70b39.
//
// Solidity: function typeAndSourceFees(bytes32 typeAndSource) view returns(uint256 fee)
func (_Fdc *FdcCaller) TypeAndSourceFees(opts *bind.CallOpts, typeAndSource [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Fdc.contract.Call(opts, &out, "typeAndSourceFees", typeAndSource)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TypeAndSourceFees is a free data retrieval call binding the contract method 0xb9e70b39.
//
// Solidity: function typeAndSourceFees(bytes32 typeAndSource) view returns(uint256 fee)
func (_Fdc *FdcSession) TypeAndSourceFees(typeAndSource [32]byte) (*big.Int, error) {
	return _Fdc.Contract.TypeAndSourceFees(&_Fdc.CallOpts, typeAndSource)
}

// TypeAndSourceFees is a free data retrieval call binding the contract method 0xb9e70b39.
//
// Solidity: function typeAndSourceFees(bytes32 typeAndSource) view returns(uint256 fee)
func (_Fdc *FdcCallerSession) TypeAndSourceFees(typeAndSource [32]byte) (*big.Int, error) {
	return _Fdc.Contract.TypeAndSourceFees(&_Fdc.CallOpts, typeAndSource)
}

// CancelGovernanceCall is a paid mutator transaction binding the contract method 0x67fc4029.
//
// Solidity: function cancelGovernanceCall(bytes4 _selector) returns()
func (_Fdc *FdcTransactor) CancelGovernanceCall(opts *bind.TransactOpts, _selector [4]byte) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "cancelGovernanceCall", _selector)
}

// CancelGovernanceCall is a paid mutator transaction binding the contract method 0x67fc4029.
//
// Solidity: function cancelGovernanceCall(bytes4 _selector) returns()
func (_Fdc *FdcSession) CancelGovernanceCall(_selector [4]byte) (*types.Transaction, error) {
	return _Fdc.Contract.CancelGovernanceCall(&_Fdc.TransactOpts, _selector)
}

// CancelGovernanceCall is a paid mutator transaction binding the contract method 0x67fc4029.
//
// Solidity: function cancelGovernanceCall(bytes4 _selector) returns()
func (_Fdc *FdcTransactorSession) CancelGovernanceCall(_selector [4]byte) (*types.Transaction, error) {
	return _Fdc.Contract.CancelGovernanceCall(&_Fdc.TransactOpts, _selector)
}

// ExecuteGovernanceCall is a paid mutator transaction binding the contract method 0x5ff27079.
//
// Solidity: function executeGovernanceCall(bytes4 _selector) returns()
func (_Fdc *FdcTransactor) ExecuteGovernanceCall(opts *bind.TransactOpts, _selector [4]byte) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "executeGovernanceCall", _selector)
}

// ExecuteGovernanceCall is a paid mutator transaction binding the contract method 0x5ff27079.
//
// Solidity: function executeGovernanceCall(bytes4 _selector) returns()
func (_Fdc *FdcSession) ExecuteGovernanceCall(_selector [4]byte) (*types.Transaction, error) {
	return _Fdc.Contract.ExecuteGovernanceCall(&_Fdc.TransactOpts, _selector)
}

// ExecuteGovernanceCall is a paid mutator transaction binding the contract method 0x5ff27079.
//
// Solidity: function executeGovernanceCall(bytes4 _selector) returns()
func (_Fdc *FdcTransactorSession) ExecuteGovernanceCall(_selector [4]byte) (*types.Transaction, error) {
	return _Fdc.Contract.ExecuteGovernanceCall(&_Fdc.TransactOpts, _selector)
}

// Initialise is a paid mutator transaction binding the contract method 0xef88bf13.
//
// Solidity: function initialise(address _governanceSettings, address _initialGovernance) returns()
func (_Fdc *FdcTransactor) Initialise(opts *bind.TransactOpts, _governanceSettings common.Address, _initialGovernance common.Address) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "initialise", _governanceSettings, _initialGovernance)
}

// Initialise is a paid mutator transaction binding the contract method 0xef88bf13.
//
// Solidity: function initialise(address _governanceSettings, address _initialGovernance) returns()
func (_Fdc *FdcSession) Initialise(_governanceSettings common.Address, _initialGovernance common.Address) (*types.Transaction, error) {
	return _Fdc.Contract.Initialise(&_Fdc.TransactOpts, _governanceSettings, _initialGovernance)
}

// Initialise is a paid mutator transaction binding the contract method 0xef88bf13.
//
// Solidity: function initialise(address _governanceSettings, address _initialGovernance) returns()
func (_Fdc *FdcTransactorSession) Initialise(_governanceSettings common.Address, _initialGovernance common.Address) (*types.Transaction, error) {
	return _Fdc.Contract.Initialise(&_Fdc.TransactOpts, _governanceSettings, _initialGovernance)
}

// ReceiveInflation is a paid mutator transaction binding the contract method 0x06201f1d.
//
// Solidity: function receiveInflation() payable returns()
func (_Fdc *FdcTransactor) ReceiveInflation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "receiveInflation")
}

// ReceiveInflation is a paid mutator transaction binding the contract method 0x06201f1d.
//
// Solidity: function receiveInflation() payable returns()
func (_Fdc *FdcSession) ReceiveInflation() (*types.Transaction, error) {
	return _Fdc.Contract.ReceiveInflation(&_Fdc.TransactOpts)
}

// ReceiveInflation is a paid mutator transaction binding the contract method 0x06201f1d.
//
// Solidity: function receiveInflation() payable returns()
func (_Fdc *FdcTransactorSession) ReceiveInflation() (*types.Transaction, error) {
	return _Fdc.Contract.ReceiveInflation(&_Fdc.TransactOpts)
}

// RemoveTypeAndSourceFee is a paid mutator transaction binding the contract method 0xda42a778.
//
// Solidity: function removeTypeAndSourceFee(bytes32 _type, bytes32 _source) returns()
func (_Fdc *FdcTransactor) RemoveTypeAndSourceFee(opts *bind.TransactOpts, _type [32]byte, _source [32]byte) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "removeTypeAndSourceFee", _type, _source)
}

// RemoveTypeAndSourceFee is a paid mutator transaction binding the contract method 0xda42a778.
//
// Solidity: function removeTypeAndSourceFee(bytes32 _type, bytes32 _source) returns()
func (_Fdc *FdcSession) RemoveTypeAndSourceFee(_type [32]byte, _source [32]byte) (*types.Transaction, error) {
	return _Fdc.Contract.RemoveTypeAndSourceFee(&_Fdc.TransactOpts, _type, _source)
}

// RemoveTypeAndSourceFee is a paid mutator transaction binding the contract method 0xda42a778.
//
// Solidity: function removeTypeAndSourceFee(bytes32 _type, bytes32 _source) returns()
func (_Fdc *FdcTransactorSession) RemoveTypeAndSourceFee(_type [32]byte, _source [32]byte) (*types.Transaction, error) {
	return _Fdc.Contract.RemoveTypeAndSourceFee(&_Fdc.TransactOpts, _type, _source)
}

// RemoveTypeAndSourceFees is a paid mutator transaction binding the contract method 0xe2f6bca8.
//
// Solidity: function removeTypeAndSourceFees(bytes32[] _types, bytes32[] _sources) returns()
func (_Fdc *FdcTransactor) RemoveTypeAndSourceFees(opts *bind.TransactOpts, _types [][32]byte, _sources [][32]byte) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "removeTypeAndSourceFees", _types, _sources)
}

// RemoveTypeAndSourceFees is a paid mutator transaction binding the contract method 0xe2f6bca8.
//
// Solidity: function removeTypeAndSourceFees(bytes32[] _types, bytes32[] _sources) returns()
func (_Fdc *FdcSession) RemoveTypeAndSourceFees(_types [][32]byte, _sources [][32]byte) (*types.Transaction, error) {
	return _Fdc.Contract.RemoveTypeAndSourceFees(&_Fdc.TransactOpts, _types, _sources)
}

// RemoveTypeAndSourceFees is a paid mutator transaction binding the contract method 0xe2f6bca8.
//
// Solidity: function removeTypeAndSourceFees(bytes32[] _types, bytes32[] _sources) returns()
func (_Fdc *FdcTransactorSession) RemoveTypeAndSourceFees(_types [][32]byte, _sources [][32]byte) (*types.Transaction, error) {
	return _Fdc.Contract.RemoveTypeAndSourceFees(&_Fdc.TransactOpts, _types, _sources)
}

// RequestAttestation is a paid mutator transaction binding the contract method 0x6238f354.
//
// Solidity: function requestAttestation(bytes _data) payable returns()
func (_Fdc *FdcTransactor) RequestAttestation(opts *bind.TransactOpts, _data []byte) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "requestAttestation", _data)
}

// RequestAttestation is a paid mutator transaction binding the contract method 0x6238f354.
//
// Solidity: function requestAttestation(bytes _data) payable returns()
func (_Fdc *FdcSession) RequestAttestation(_data []byte) (*types.Transaction, error) {
	return _Fdc.Contract.RequestAttestation(&_Fdc.TransactOpts, _data)
}

// RequestAttestation is a paid mutator transaction binding the contract method 0x6238f354.
//
// Solidity: function requestAttestation(bytes _data) payable returns()
func (_Fdc *FdcTransactorSession) RequestAttestation(_data []byte) (*types.Transaction, error) {
	return _Fdc.Contract.RequestAttestation(&_Fdc.TransactOpts, _data)
}

// SetDailyAuthorizedInflation is a paid mutator transaction binding the contract method 0xe2739563.
//
// Solidity: function setDailyAuthorizedInflation(uint256 _toAuthorizeWei) returns()
func (_Fdc *FdcTransactor) SetDailyAuthorizedInflation(opts *bind.TransactOpts, _toAuthorizeWei *big.Int) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "setDailyAuthorizedInflation", _toAuthorizeWei)
}

// SetDailyAuthorizedInflation is a paid mutator transaction binding the contract method 0xe2739563.
//
// Solidity: function setDailyAuthorizedInflation(uint256 _toAuthorizeWei) returns()
func (_Fdc *FdcSession) SetDailyAuthorizedInflation(_toAuthorizeWei *big.Int) (*types.Transaction, error) {
	return _Fdc.Contract.SetDailyAuthorizedInflation(&_Fdc.TransactOpts, _toAuthorizeWei)
}

// SetDailyAuthorizedInflation is a paid mutator transaction binding the contract method 0xe2739563.
//
// Solidity: function setDailyAuthorizedInflation(uint256 _toAuthorizeWei) returns()
func (_Fdc *FdcTransactorSession) SetDailyAuthorizedInflation(_toAuthorizeWei *big.Int) (*types.Transaction, error) {
	return _Fdc.Contract.SetDailyAuthorizedInflation(&_Fdc.TransactOpts, _toAuthorizeWei)
}

// SetTypeAndSourceFee is a paid mutator transaction binding the contract method 0x5b1bacb7.
//
// Solidity: function setTypeAndSourceFee(bytes32 _type, bytes32 _source, uint256 _fee) returns()
func (_Fdc *FdcTransactor) SetTypeAndSourceFee(opts *bind.TransactOpts, _type [32]byte, _source [32]byte, _fee *big.Int) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "setTypeAndSourceFee", _type, _source, _fee)
}

// SetTypeAndSourceFee is a paid mutator transaction binding the contract method 0x5b1bacb7.
//
// Solidity: function setTypeAndSourceFee(bytes32 _type, bytes32 _source, uint256 _fee) returns()
func (_Fdc *FdcSession) SetTypeAndSourceFee(_type [32]byte, _source [32]byte, _fee *big.Int) (*types.Transaction, error) {
	return _Fdc.Contract.SetTypeAndSourceFee(&_Fdc.TransactOpts, _type, _source, _fee)
}

// SetTypeAndSourceFee is a paid mutator transaction binding the contract method 0x5b1bacb7.
//
// Solidity: function setTypeAndSourceFee(bytes32 _type, bytes32 _source, uint256 _fee) returns()
func (_Fdc *FdcTransactorSession) SetTypeAndSourceFee(_type [32]byte, _source [32]byte, _fee *big.Int) (*types.Transaction, error) {
	return _Fdc.Contract.SetTypeAndSourceFee(&_Fdc.TransactOpts, _type, _source, _fee)
}

// SetTypeAndSourceFees is a paid mutator transaction binding the contract method 0x3da84157.
//
// Solidity: function setTypeAndSourceFees(bytes32[] _types, bytes32[] _sources, uint256[] _fees) returns()
func (_Fdc *FdcTransactor) SetTypeAndSourceFees(opts *bind.TransactOpts, _types [][32]byte, _sources [][32]byte, _fees []*big.Int) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "setTypeAndSourceFees", _types, _sources, _fees)
}

// SetTypeAndSourceFees is a paid mutator transaction binding the contract method 0x3da84157.
//
// Solidity: function setTypeAndSourceFees(bytes32[] _types, bytes32[] _sources, uint256[] _fees) returns()
func (_Fdc *FdcSession) SetTypeAndSourceFees(_types [][32]byte, _sources [][32]byte, _fees []*big.Int) (*types.Transaction, error) {
	return _Fdc.Contract.SetTypeAndSourceFees(&_Fdc.TransactOpts, _types, _sources, _fees)
}

// SetTypeAndSourceFees is a paid mutator transaction binding the contract method 0x3da84157.
//
// Solidity: function setTypeAndSourceFees(bytes32[] _types, bytes32[] _sources, uint256[] _fees) returns()
func (_Fdc *FdcTransactorSession) SetTypeAndSourceFees(_types [][32]byte, _sources [][32]byte, _fees []*big.Int) (*types.Transaction, error) {
	return _Fdc.Contract.SetTypeAndSourceFees(&_Fdc.TransactOpts, _types, _sources, _fees)
}

// SwitchToProductionMode is a paid mutator transaction binding the contract method 0xf5a98383.
//
// Solidity: function switchToProductionMode() returns()
func (_Fdc *FdcTransactor) SwitchToProductionMode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "switchToProductionMode")
}

// SwitchToProductionMode is a paid mutator transaction binding the contract method 0xf5a98383.
//
// Solidity: function switchToProductionMode() returns()
func (_Fdc *FdcSession) SwitchToProductionMode() (*types.Transaction, error) {
	return _Fdc.Contract.SwitchToProductionMode(&_Fdc.TransactOpts)
}

// SwitchToProductionMode is a paid mutator transaction binding the contract method 0xf5a98383.
//
// Solidity: function switchToProductionMode() returns()
func (_Fdc *FdcTransactorSession) SwitchToProductionMode() (*types.Transaction, error) {
	return _Fdc.Contract.SwitchToProductionMode(&_Fdc.TransactOpts)
}

// TriggerRewardEpochSwitchover is a paid mutator transaction binding the contract method 0x91f25679.
//
// Solidity: function triggerRewardEpochSwitchover(uint24 _currentRewardEpochId, uint64 _currentRewardEpochExpectedEndTs, uint64 _rewardEpochDurationSeconds) returns()
func (_Fdc *FdcTransactor) TriggerRewardEpochSwitchover(opts *bind.TransactOpts, _currentRewardEpochId *big.Int, _currentRewardEpochExpectedEndTs uint64, _rewardEpochDurationSeconds uint64) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "triggerRewardEpochSwitchover", _currentRewardEpochId, _currentRewardEpochExpectedEndTs, _rewardEpochDurationSeconds)
}

// TriggerRewardEpochSwitchover is a paid mutator transaction binding the contract method 0x91f25679.
//
// Solidity: function triggerRewardEpochSwitchover(uint24 _currentRewardEpochId, uint64 _currentRewardEpochExpectedEndTs, uint64 _rewardEpochDurationSeconds) returns()
func (_Fdc *FdcSession) TriggerRewardEpochSwitchover(_currentRewardEpochId *big.Int, _currentRewardEpochExpectedEndTs uint64, _rewardEpochDurationSeconds uint64) (*types.Transaction, error) {
	return _Fdc.Contract.TriggerRewardEpochSwitchover(&_Fdc.TransactOpts, _currentRewardEpochId, _currentRewardEpochExpectedEndTs, _rewardEpochDurationSeconds)
}

// TriggerRewardEpochSwitchover is a paid mutator transaction binding the contract method 0x91f25679.
//
// Solidity: function triggerRewardEpochSwitchover(uint24 _currentRewardEpochId, uint64 _currentRewardEpochExpectedEndTs, uint64 _rewardEpochDurationSeconds) returns()
func (_Fdc *FdcTransactorSession) TriggerRewardEpochSwitchover(_currentRewardEpochId *big.Int, _currentRewardEpochExpectedEndTs uint64, _rewardEpochDurationSeconds uint64) (*types.Transaction, error) {
	return _Fdc.Contract.TriggerRewardEpochSwitchover(&_Fdc.TransactOpts, _currentRewardEpochId, _currentRewardEpochExpectedEndTs, _rewardEpochDurationSeconds)
}

// UpdateContractAddresses is a paid mutator transaction binding the contract method 0xb00c0b76.
//
// Solidity: function updateContractAddresses(bytes32[] _contractNameHashes, address[] _contractAddresses) returns()
func (_Fdc *FdcTransactor) UpdateContractAddresses(opts *bind.TransactOpts, _contractNameHashes [][32]byte, _contractAddresses []common.Address) (*types.Transaction, error) {
	return _Fdc.contract.Transact(opts, "updateContractAddresses", _contractNameHashes, _contractAddresses)
}

// UpdateContractAddresses is a paid mutator transaction binding the contract method 0xb00c0b76.
//
// Solidity: function updateContractAddresses(bytes32[] _contractNameHashes, address[] _contractAddresses) returns()
func (_Fdc *FdcSession) UpdateContractAddresses(_contractNameHashes [][32]byte, _contractAddresses []common.Address) (*types.Transaction, error) {
	return _Fdc.Contract.UpdateContractAddresses(&_Fdc.TransactOpts, _contractNameHashes, _contractAddresses)
}

// UpdateContractAddresses is a paid mutator transaction binding the contract method 0xb00c0b76.
//
// Solidity: function updateContractAddresses(bytes32[] _contractNameHashes, address[] _contractAddresses) returns()
func (_Fdc *FdcTransactorSession) UpdateContractAddresses(_contractNameHashes [][32]byte, _contractAddresses []common.Address) (*types.Transaction, error) {
	return _Fdc.Contract.UpdateContractAddresses(&_Fdc.TransactOpts, _contractNameHashes, _contractAddresses)
}

// FdcAttestationRequestIterator is returned from FilterAttestationRequest and is used to iterate over the raw logs and unpacked data for AttestationRequest events raised by the Fdc contract.
type FdcAttestationRequestIterator struct {
	Event *FdcAttestationRequest // Event containing the contract specifics and raw log

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
func (it *FdcAttestationRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcAttestationRequest)
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
		it.Event = new(FdcAttestationRequest)
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
func (it *FdcAttestationRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcAttestationRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcAttestationRequest represents a AttestationRequest event raised by the Fdc contract.
type FdcAttestationRequest struct {
	Data []byte
	Fee  *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAttestationRequest is a free log retrieval operation binding the contract event 0x251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9.
//
// Solidity: event AttestationRequest(bytes data, uint256 fee)
func (_Fdc *FdcFilterer) FilterAttestationRequest(opts *bind.FilterOpts) (*FdcAttestationRequestIterator, error) {

	logs, sub, err := _Fdc.contract.FilterLogs(opts, "AttestationRequest")
	if err != nil {
		return nil, err
	}
	return &FdcAttestationRequestIterator{contract: _Fdc.contract, event: "AttestationRequest", logs: logs, sub: sub}, nil
}

// WatchAttestationRequest is a free log subscription operation binding the contract event 0x251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9.
//
// Solidity: event AttestationRequest(bytes data, uint256 fee)
func (_Fdc *FdcFilterer) WatchAttestationRequest(opts *bind.WatchOpts, sink chan<- *FdcAttestationRequest) (event.Subscription, error) {

	logs, sub, err := _Fdc.contract.WatchLogs(opts, "AttestationRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcAttestationRequest)
				if err := _Fdc.contract.UnpackLog(event, "AttestationRequest", log); err != nil {
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

// ParseAttestationRequest is a log parse operation binding the contract event 0x251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9.
//
// Solidity: event AttestationRequest(bytes data, uint256 fee)
func (_Fdc *FdcFilterer) ParseAttestationRequest(log types.Log) (*FdcAttestationRequest, error) {
	event := new(FdcAttestationRequest)
	if err := _Fdc.contract.UnpackLog(event, "AttestationRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcDailyAuthorizedInflationSetIterator is returned from FilterDailyAuthorizedInflationSet and is used to iterate over the raw logs and unpacked data for DailyAuthorizedInflationSet events raised by the Fdc contract.
type FdcDailyAuthorizedInflationSetIterator struct {
	Event *FdcDailyAuthorizedInflationSet // Event containing the contract specifics and raw log

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
func (it *FdcDailyAuthorizedInflationSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcDailyAuthorizedInflationSet)
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
		it.Event = new(FdcDailyAuthorizedInflationSet)
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
func (it *FdcDailyAuthorizedInflationSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcDailyAuthorizedInflationSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcDailyAuthorizedInflationSet represents a DailyAuthorizedInflationSet event raised by the Fdc contract.
type FdcDailyAuthorizedInflationSet struct {
	AuthorizedAmountWei *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterDailyAuthorizedInflationSet is a free log retrieval operation binding the contract event 0x187f32a0f765499f15b3bb52ed0aebf6015059f230f2ace7e701e60a47669595.
//
// Solidity: event DailyAuthorizedInflationSet(uint256 authorizedAmountWei)
func (_Fdc *FdcFilterer) FilterDailyAuthorizedInflationSet(opts *bind.FilterOpts) (*FdcDailyAuthorizedInflationSetIterator, error) {

	logs, sub, err := _Fdc.contract.FilterLogs(opts, "DailyAuthorizedInflationSet")
	if err != nil {
		return nil, err
	}
	return &FdcDailyAuthorizedInflationSetIterator{contract: _Fdc.contract, event: "DailyAuthorizedInflationSet", logs: logs, sub: sub}, nil
}

// WatchDailyAuthorizedInflationSet is a free log subscription operation binding the contract event 0x187f32a0f765499f15b3bb52ed0aebf6015059f230f2ace7e701e60a47669595.
//
// Solidity: event DailyAuthorizedInflationSet(uint256 authorizedAmountWei)
func (_Fdc *FdcFilterer) WatchDailyAuthorizedInflationSet(opts *bind.WatchOpts, sink chan<- *FdcDailyAuthorizedInflationSet) (event.Subscription, error) {

	logs, sub, err := _Fdc.contract.WatchLogs(opts, "DailyAuthorizedInflationSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcDailyAuthorizedInflationSet)
				if err := _Fdc.contract.UnpackLog(event, "DailyAuthorizedInflationSet", log); err != nil {
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

// ParseDailyAuthorizedInflationSet is a log parse operation binding the contract event 0x187f32a0f765499f15b3bb52ed0aebf6015059f230f2ace7e701e60a47669595.
//
// Solidity: event DailyAuthorizedInflationSet(uint256 authorizedAmountWei)
func (_Fdc *FdcFilterer) ParseDailyAuthorizedInflationSet(log types.Log) (*FdcDailyAuthorizedInflationSet, error) {
	event := new(FdcDailyAuthorizedInflationSet)
	if err := _Fdc.contract.UnpackLog(event, "DailyAuthorizedInflationSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcGovernanceCallTimelockedIterator is returned from FilterGovernanceCallTimelocked and is used to iterate over the raw logs and unpacked data for GovernanceCallTimelocked events raised by the Fdc contract.
type FdcGovernanceCallTimelockedIterator struct {
	Event *FdcGovernanceCallTimelocked // Event containing the contract specifics and raw log

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
func (it *FdcGovernanceCallTimelockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcGovernanceCallTimelocked)
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
		it.Event = new(FdcGovernanceCallTimelocked)
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
func (it *FdcGovernanceCallTimelockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcGovernanceCallTimelockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcGovernanceCallTimelocked represents a GovernanceCallTimelocked event raised by the Fdc contract.
type FdcGovernanceCallTimelocked struct {
	Selector              [4]byte
	AllowedAfterTimestamp *big.Int
	EncodedCall           []byte
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterGovernanceCallTimelocked is a free log retrieval operation binding the contract event 0xed948300a3694aa01d4a6b258bfd664350193d770c0b51f8387277f6d83ea3b6.
//
// Solidity: event GovernanceCallTimelocked(bytes4 selector, uint256 allowedAfterTimestamp, bytes encodedCall)
func (_Fdc *FdcFilterer) FilterGovernanceCallTimelocked(opts *bind.FilterOpts) (*FdcGovernanceCallTimelockedIterator, error) {

	logs, sub, err := _Fdc.contract.FilterLogs(opts, "GovernanceCallTimelocked")
	if err != nil {
		return nil, err
	}
	return &FdcGovernanceCallTimelockedIterator{contract: _Fdc.contract, event: "GovernanceCallTimelocked", logs: logs, sub: sub}, nil
}

// WatchGovernanceCallTimelocked is a free log subscription operation binding the contract event 0xed948300a3694aa01d4a6b258bfd664350193d770c0b51f8387277f6d83ea3b6.
//
// Solidity: event GovernanceCallTimelocked(bytes4 selector, uint256 allowedAfterTimestamp, bytes encodedCall)
func (_Fdc *FdcFilterer) WatchGovernanceCallTimelocked(opts *bind.WatchOpts, sink chan<- *FdcGovernanceCallTimelocked) (event.Subscription, error) {

	logs, sub, err := _Fdc.contract.WatchLogs(opts, "GovernanceCallTimelocked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcGovernanceCallTimelocked)
				if err := _Fdc.contract.UnpackLog(event, "GovernanceCallTimelocked", log); err != nil {
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

// ParseGovernanceCallTimelocked is a log parse operation binding the contract event 0xed948300a3694aa01d4a6b258bfd664350193d770c0b51f8387277f6d83ea3b6.
//
// Solidity: event GovernanceCallTimelocked(bytes4 selector, uint256 allowedAfterTimestamp, bytes encodedCall)
func (_Fdc *FdcFilterer) ParseGovernanceCallTimelocked(log types.Log) (*FdcGovernanceCallTimelocked, error) {
	event := new(FdcGovernanceCallTimelocked)
	if err := _Fdc.contract.UnpackLog(event, "GovernanceCallTimelocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcGovernanceInitialisedIterator is returned from FilterGovernanceInitialised and is used to iterate over the raw logs and unpacked data for GovernanceInitialised events raised by the Fdc contract.
type FdcGovernanceInitialisedIterator struct {
	Event *FdcGovernanceInitialised // Event containing the contract specifics and raw log

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
func (it *FdcGovernanceInitialisedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcGovernanceInitialised)
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
		it.Event = new(FdcGovernanceInitialised)
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
func (it *FdcGovernanceInitialisedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcGovernanceInitialisedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcGovernanceInitialised represents a GovernanceInitialised event raised by the Fdc contract.
type FdcGovernanceInitialised struct {
	InitialGovernance common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterGovernanceInitialised is a free log retrieval operation binding the contract event 0x9789733827840833afc031fb2ef9ab6894271f77bad2085687cf4ae5c7bee4db.
//
// Solidity: event GovernanceInitialised(address initialGovernance)
func (_Fdc *FdcFilterer) FilterGovernanceInitialised(opts *bind.FilterOpts) (*FdcGovernanceInitialisedIterator, error) {

	logs, sub, err := _Fdc.contract.FilterLogs(opts, "GovernanceInitialised")
	if err != nil {
		return nil, err
	}
	return &FdcGovernanceInitialisedIterator{contract: _Fdc.contract, event: "GovernanceInitialised", logs: logs, sub: sub}, nil
}

// WatchGovernanceInitialised is a free log subscription operation binding the contract event 0x9789733827840833afc031fb2ef9ab6894271f77bad2085687cf4ae5c7bee4db.
//
// Solidity: event GovernanceInitialised(address initialGovernance)
func (_Fdc *FdcFilterer) WatchGovernanceInitialised(opts *bind.WatchOpts, sink chan<- *FdcGovernanceInitialised) (event.Subscription, error) {

	logs, sub, err := _Fdc.contract.WatchLogs(opts, "GovernanceInitialised")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcGovernanceInitialised)
				if err := _Fdc.contract.UnpackLog(event, "GovernanceInitialised", log); err != nil {
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

// ParseGovernanceInitialised is a log parse operation binding the contract event 0x9789733827840833afc031fb2ef9ab6894271f77bad2085687cf4ae5c7bee4db.
//
// Solidity: event GovernanceInitialised(address initialGovernance)
func (_Fdc *FdcFilterer) ParseGovernanceInitialised(log types.Log) (*FdcGovernanceInitialised, error) {
	event := new(FdcGovernanceInitialised)
	if err := _Fdc.contract.UnpackLog(event, "GovernanceInitialised", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcGovernedProductionModeEnteredIterator is returned from FilterGovernedProductionModeEntered and is used to iterate over the raw logs and unpacked data for GovernedProductionModeEntered events raised by the Fdc contract.
type FdcGovernedProductionModeEnteredIterator struct {
	Event *FdcGovernedProductionModeEntered // Event containing the contract specifics and raw log

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
func (it *FdcGovernedProductionModeEnteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcGovernedProductionModeEntered)
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
		it.Event = new(FdcGovernedProductionModeEntered)
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
func (it *FdcGovernedProductionModeEnteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcGovernedProductionModeEnteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcGovernedProductionModeEntered represents a GovernedProductionModeEntered event raised by the Fdc contract.
type FdcGovernedProductionModeEntered struct {
	GovernanceSettings common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterGovernedProductionModeEntered is a free log retrieval operation binding the contract event 0x83af113638b5422f9e977cebc0aaf0eaf2188eb9a8baae7f9d46c42b33a1560c.
//
// Solidity: event GovernedProductionModeEntered(address governanceSettings)
func (_Fdc *FdcFilterer) FilterGovernedProductionModeEntered(opts *bind.FilterOpts) (*FdcGovernedProductionModeEnteredIterator, error) {

	logs, sub, err := _Fdc.contract.FilterLogs(opts, "GovernedProductionModeEntered")
	if err != nil {
		return nil, err
	}
	return &FdcGovernedProductionModeEnteredIterator{contract: _Fdc.contract, event: "GovernedProductionModeEntered", logs: logs, sub: sub}, nil
}

// WatchGovernedProductionModeEntered is a free log subscription operation binding the contract event 0x83af113638b5422f9e977cebc0aaf0eaf2188eb9a8baae7f9d46c42b33a1560c.
//
// Solidity: event GovernedProductionModeEntered(address governanceSettings)
func (_Fdc *FdcFilterer) WatchGovernedProductionModeEntered(opts *bind.WatchOpts, sink chan<- *FdcGovernedProductionModeEntered) (event.Subscription, error) {

	logs, sub, err := _Fdc.contract.WatchLogs(opts, "GovernedProductionModeEntered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcGovernedProductionModeEntered)
				if err := _Fdc.contract.UnpackLog(event, "GovernedProductionModeEntered", log); err != nil {
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

// ParseGovernedProductionModeEntered is a log parse operation binding the contract event 0x83af113638b5422f9e977cebc0aaf0eaf2188eb9a8baae7f9d46c42b33a1560c.
//
// Solidity: event GovernedProductionModeEntered(address governanceSettings)
func (_Fdc *FdcFilterer) ParseGovernedProductionModeEntered(log types.Log) (*FdcGovernedProductionModeEntered, error) {
	event := new(FdcGovernedProductionModeEntered)
	if err := _Fdc.contract.UnpackLog(event, "GovernedProductionModeEntered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcInflationReceivedIterator is returned from FilterInflationReceived and is used to iterate over the raw logs and unpacked data for InflationReceived events raised by the Fdc contract.
type FdcInflationReceivedIterator struct {
	Event *FdcInflationReceived // Event containing the contract specifics and raw log

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
func (it *FdcInflationReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcInflationReceived)
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
		it.Event = new(FdcInflationReceived)
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
func (it *FdcInflationReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcInflationReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcInflationReceived represents a InflationReceived event raised by the Fdc contract.
type FdcInflationReceived struct {
	AmountReceivedWei *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterInflationReceived is a free log retrieval operation binding the contract event 0x95c4e29cc99bc027cfc3cd719d6fd973d5f0317061885fbb322b9b17d8d35d37.
//
// Solidity: event InflationReceived(uint256 amountReceivedWei)
func (_Fdc *FdcFilterer) FilterInflationReceived(opts *bind.FilterOpts) (*FdcInflationReceivedIterator, error) {

	logs, sub, err := _Fdc.contract.FilterLogs(opts, "InflationReceived")
	if err != nil {
		return nil, err
	}
	return &FdcInflationReceivedIterator{contract: _Fdc.contract, event: "InflationReceived", logs: logs, sub: sub}, nil
}

// WatchInflationReceived is a free log subscription operation binding the contract event 0x95c4e29cc99bc027cfc3cd719d6fd973d5f0317061885fbb322b9b17d8d35d37.
//
// Solidity: event InflationReceived(uint256 amountReceivedWei)
func (_Fdc *FdcFilterer) WatchInflationReceived(opts *bind.WatchOpts, sink chan<- *FdcInflationReceived) (event.Subscription, error) {

	logs, sub, err := _Fdc.contract.WatchLogs(opts, "InflationReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcInflationReceived)
				if err := _Fdc.contract.UnpackLog(event, "InflationReceived", log); err != nil {
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

// ParseInflationReceived is a log parse operation binding the contract event 0x95c4e29cc99bc027cfc3cd719d6fd973d5f0317061885fbb322b9b17d8d35d37.
//
// Solidity: event InflationReceived(uint256 amountReceivedWei)
func (_Fdc *FdcFilterer) ParseInflationReceived(log types.Log) (*FdcInflationReceived, error) {
	event := new(FdcInflationReceived)
	if err := _Fdc.contract.UnpackLog(event, "InflationReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcInflationRewardsOfferedIterator is returned from FilterInflationRewardsOffered and is used to iterate over the raw logs and unpacked data for InflationRewardsOffered events raised by the Fdc contract.
type FdcInflationRewardsOfferedIterator struct {
	Event *FdcInflationRewardsOffered // Event containing the contract specifics and raw log

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
func (it *FdcInflationRewardsOfferedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcInflationRewardsOffered)
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
		it.Event = new(FdcInflationRewardsOffered)
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
func (it *FdcInflationRewardsOfferedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcInflationRewardsOfferedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcInflationRewardsOffered represents a InflationRewardsOffered event raised by the Fdc contract.
type FdcInflationRewardsOffered struct {
	RewardEpochId *big.Int
	Amount        *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterInflationRewardsOffered is a free log retrieval operation binding the contract event 0x27b4386cf0cd6c0d9f8d1368d66972efd9dee44649ec21a62c46cef37646c7e6.
//
// Solidity: event InflationRewardsOffered(uint24 indexed rewardEpochId, uint256 amount)
func (_Fdc *FdcFilterer) FilterInflationRewardsOffered(opts *bind.FilterOpts, rewardEpochId []*big.Int) (*FdcInflationRewardsOfferedIterator, error) {

	var rewardEpochIdRule []interface{}
	for _, rewardEpochIdItem := range rewardEpochId {
		rewardEpochIdRule = append(rewardEpochIdRule, rewardEpochIdItem)
	}

	logs, sub, err := _Fdc.contract.FilterLogs(opts, "InflationRewardsOffered", rewardEpochIdRule)
	if err != nil {
		return nil, err
	}
	return &FdcInflationRewardsOfferedIterator{contract: _Fdc.contract, event: "InflationRewardsOffered", logs: logs, sub: sub}, nil
}

// WatchInflationRewardsOffered is a free log subscription operation binding the contract event 0x27b4386cf0cd6c0d9f8d1368d66972efd9dee44649ec21a62c46cef37646c7e6.
//
// Solidity: event InflationRewardsOffered(uint24 indexed rewardEpochId, uint256 amount)
func (_Fdc *FdcFilterer) WatchInflationRewardsOffered(opts *bind.WatchOpts, sink chan<- *FdcInflationRewardsOffered, rewardEpochId []*big.Int) (event.Subscription, error) {

	var rewardEpochIdRule []interface{}
	for _, rewardEpochIdItem := range rewardEpochId {
		rewardEpochIdRule = append(rewardEpochIdRule, rewardEpochIdItem)
	}

	logs, sub, err := _Fdc.contract.WatchLogs(opts, "InflationRewardsOffered", rewardEpochIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcInflationRewardsOffered)
				if err := _Fdc.contract.UnpackLog(event, "InflationRewardsOffered", log); err != nil {
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

// ParseInflationRewardsOffered is a log parse operation binding the contract event 0x27b4386cf0cd6c0d9f8d1368d66972efd9dee44649ec21a62c46cef37646c7e6.
//
// Solidity: event InflationRewardsOffered(uint24 indexed rewardEpochId, uint256 amount)
func (_Fdc *FdcFilterer) ParseInflationRewardsOffered(log types.Log) (*FdcInflationRewardsOffered, error) {
	event := new(FdcInflationRewardsOffered)
	if err := _Fdc.contract.UnpackLog(event, "InflationRewardsOffered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcTimelockedGovernanceCallCanceledIterator is returned from FilterTimelockedGovernanceCallCanceled and is used to iterate over the raw logs and unpacked data for TimelockedGovernanceCallCanceled events raised by the Fdc contract.
type FdcTimelockedGovernanceCallCanceledIterator struct {
	Event *FdcTimelockedGovernanceCallCanceled // Event containing the contract specifics and raw log

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
func (it *FdcTimelockedGovernanceCallCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcTimelockedGovernanceCallCanceled)
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
		it.Event = new(FdcTimelockedGovernanceCallCanceled)
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
func (it *FdcTimelockedGovernanceCallCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcTimelockedGovernanceCallCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcTimelockedGovernanceCallCanceled represents a TimelockedGovernanceCallCanceled event raised by the Fdc contract.
type FdcTimelockedGovernanceCallCanceled struct {
	Selector  [4]byte
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTimelockedGovernanceCallCanceled is a free log retrieval operation binding the contract event 0x7735b2391c38a81419c513e30ca578db7158eadd7101511b23e221c654d19cf8.
//
// Solidity: event TimelockedGovernanceCallCanceled(bytes4 selector, uint256 timestamp)
func (_Fdc *FdcFilterer) FilterTimelockedGovernanceCallCanceled(opts *bind.FilterOpts) (*FdcTimelockedGovernanceCallCanceledIterator, error) {

	logs, sub, err := _Fdc.contract.FilterLogs(opts, "TimelockedGovernanceCallCanceled")
	if err != nil {
		return nil, err
	}
	return &FdcTimelockedGovernanceCallCanceledIterator{contract: _Fdc.contract, event: "TimelockedGovernanceCallCanceled", logs: logs, sub: sub}, nil
}

// WatchTimelockedGovernanceCallCanceled is a free log subscription operation binding the contract event 0x7735b2391c38a81419c513e30ca578db7158eadd7101511b23e221c654d19cf8.
//
// Solidity: event TimelockedGovernanceCallCanceled(bytes4 selector, uint256 timestamp)
func (_Fdc *FdcFilterer) WatchTimelockedGovernanceCallCanceled(opts *bind.WatchOpts, sink chan<- *FdcTimelockedGovernanceCallCanceled) (event.Subscription, error) {

	logs, sub, err := _Fdc.contract.WatchLogs(opts, "TimelockedGovernanceCallCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcTimelockedGovernanceCallCanceled)
				if err := _Fdc.contract.UnpackLog(event, "TimelockedGovernanceCallCanceled", log); err != nil {
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

// ParseTimelockedGovernanceCallCanceled is a log parse operation binding the contract event 0x7735b2391c38a81419c513e30ca578db7158eadd7101511b23e221c654d19cf8.
//
// Solidity: event TimelockedGovernanceCallCanceled(bytes4 selector, uint256 timestamp)
func (_Fdc *FdcFilterer) ParseTimelockedGovernanceCallCanceled(log types.Log) (*FdcTimelockedGovernanceCallCanceled, error) {
	event := new(FdcTimelockedGovernanceCallCanceled)
	if err := _Fdc.contract.UnpackLog(event, "TimelockedGovernanceCallCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcTimelockedGovernanceCallExecutedIterator is returned from FilterTimelockedGovernanceCallExecuted and is used to iterate over the raw logs and unpacked data for TimelockedGovernanceCallExecuted events raised by the Fdc contract.
type FdcTimelockedGovernanceCallExecutedIterator struct {
	Event *FdcTimelockedGovernanceCallExecuted // Event containing the contract specifics and raw log

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
func (it *FdcTimelockedGovernanceCallExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcTimelockedGovernanceCallExecuted)
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
		it.Event = new(FdcTimelockedGovernanceCallExecuted)
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
func (it *FdcTimelockedGovernanceCallExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcTimelockedGovernanceCallExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcTimelockedGovernanceCallExecuted represents a TimelockedGovernanceCallExecuted event raised by the Fdc contract.
type FdcTimelockedGovernanceCallExecuted struct {
	Selector  [4]byte
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTimelockedGovernanceCallExecuted is a free log retrieval operation binding the contract event 0xa7326b57fc9cfe267aaea5e7f0b01757154d265620a0585819416ee9ddd2c438.
//
// Solidity: event TimelockedGovernanceCallExecuted(bytes4 selector, uint256 timestamp)
func (_Fdc *FdcFilterer) FilterTimelockedGovernanceCallExecuted(opts *bind.FilterOpts) (*FdcTimelockedGovernanceCallExecutedIterator, error) {

	logs, sub, err := _Fdc.contract.FilterLogs(opts, "TimelockedGovernanceCallExecuted")
	if err != nil {
		return nil, err
	}
	return &FdcTimelockedGovernanceCallExecutedIterator{contract: _Fdc.contract, event: "TimelockedGovernanceCallExecuted", logs: logs, sub: sub}, nil
}

// WatchTimelockedGovernanceCallExecuted is a free log subscription operation binding the contract event 0xa7326b57fc9cfe267aaea5e7f0b01757154d265620a0585819416ee9ddd2c438.
//
// Solidity: event TimelockedGovernanceCallExecuted(bytes4 selector, uint256 timestamp)
func (_Fdc *FdcFilterer) WatchTimelockedGovernanceCallExecuted(opts *bind.WatchOpts, sink chan<- *FdcTimelockedGovernanceCallExecuted) (event.Subscription, error) {

	logs, sub, err := _Fdc.contract.WatchLogs(opts, "TimelockedGovernanceCallExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcTimelockedGovernanceCallExecuted)
				if err := _Fdc.contract.UnpackLog(event, "TimelockedGovernanceCallExecuted", log); err != nil {
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

// ParseTimelockedGovernanceCallExecuted is a log parse operation binding the contract event 0xa7326b57fc9cfe267aaea5e7f0b01757154d265620a0585819416ee9ddd2c438.
//
// Solidity: event TimelockedGovernanceCallExecuted(bytes4 selector, uint256 timestamp)
func (_Fdc *FdcFilterer) ParseTimelockedGovernanceCallExecuted(log types.Log) (*FdcTimelockedGovernanceCallExecuted, error) {
	event := new(FdcTimelockedGovernanceCallExecuted)
	if err := _Fdc.contract.UnpackLog(event, "TimelockedGovernanceCallExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcTypeAndSourceFeeRemovedIterator is returned from FilterTypeAndSourceFeeRemoved and is used to iterate over the raw logs and unpacked data for TypeAndSourceFeeRemoved events raised by the Fdc contract.
type FdcTypeAndSourceFeeRemovedIterator struct {
	Event *FdcTypeAndSourceFeeRemoved // Event containing the contract specifics and raw log

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
func (it *FdcTypeAndSourceFeeRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcTypeAndSourceFeeRemoved)
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
		it.Event = new(FdcTypeAndSourceFeeRemoved)
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
func (it *FdcTypeAndSourceFeeRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcTypeAndSourceFeeRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcTypeAndSourceFeeRemoved represents a TypeAndSourceFeeRemoved event raised by the Fdc contract.
type FdcTypeAndSourceFeeRemoved struct {
	Type   [32]byte
	Source [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTypeAndSourceFeeRemoved is a free log retrieval operation binding the contract event 0xa186aa06a988f2cd41162e05c078a9f69e58eb1a7233361c2604a1a49721d521.
//
// Solidity: event TypeAndSourceFeeRemoved(bytes32 indexed _type, bytes32 indexed source)
func (_Fdc *FdcFilterer) FilterTypeAndSourceFeeRemoved(opts *bind.FilterOpts, _type [][32]byte, source [][32]byte) (*FdcTypeAndSourceFeeRemovedIterator, error) {

	var _typeRule []interface{}
	for _, _typeItem := range _type {
		_typeRule = append(_typeRule, _typeItem)
	}
	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _Fdc.contract.FilterLogs(opts, "TypeAndSourceFeeRemoved", _typeRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return &FdcTypeAndSourceFeeRemovedIterator{contract: _Fdc.contract, event: "TypeAndSourceFeeRemoved", logs: logs, sub: sub}, nil
}

// WatchTypeAndSourceFeeRemoved is a free log subscription operation binding the contract event 0xa186aa06a988f2cd41162e05c078a9f69e58eb1a7233361c2604a1a49721d521.
//
// Solidity: event TypeAndSourceFeeRemoved(bytes32 indexed _type, bytes32 indexed source)
func (_Fdc *FdcFilterer) WatchTypeAndSourceFeeRemoved(opts *bind.WatchOpts, sink chan<- *FdcTypeAndSourceFeeRemoved, _type [][32]byte, source [][32]byte) (event.Subscription, error) {

	var _typeRule []interface{}
	for _, _typeItem := range _type {
		_typeRule = append(_typeRule, _typeItem)
	}
	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _Fdc.contract.WatchLogs(opts, "TypeAndSourceFeeRemoved", _typeRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcTypeAndSourceFeeRemoved)
				if err := _Fdc.contract.UnpackLog(event, "TypeAndSourceFeeRemoved", log); err != nil {
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

// ParseTypeAndSourceFeeRemoved is a log parse operation binding the contract event 0xa186aa06a988f2cd41162e05c078a9f69e58eb1a7233361c2604a1a49721d521.
//
// Solidity: event TypeAndSourceFeeRemoved(bytes32 indexed _type, bytes32 indexed source)
func (_Fdc *FdcFilterer) ParseTypeAndSourceFeeRemoved(log types.Log) (*FdcTypeAndSourceFeeRemoved, error) {
	event := new(FdcTypeAndSourceFeeRemoved)
	if err := _Fdc.contract.UnpackLog(event, "TypeAndSourceFeeRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcTypeAndSourceFeeSetIterator is returned from FilterTypeAndSourceFeeSet and is used to iterate over the raw logs and unpacked data for TypeAndSourceFeeSet events raised by the Fdc contract.
type FdcTypeAndSourceFeeSetIterator struct {
	Event *FdcTypeAndSourceFeeSet // Event containing the contract specifics and raw log

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
func (it *FdcTypeAndSourceFeeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcTypeAndSourceFeeSet)
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
		it.Event = new(FdcTypeAndSourceFeeSet)
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
func (it *FdcTypeAndSourceFeeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcTypeAndSourceFeeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcTypeAndSourceFeeSet represents a TypeAndSourceFeeSet event raised by the Fdc contract.
type FdcTypeAndSourceFeeSet struct {
	Type   [32]byte
	Source [32]byte
	Fee    *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTypeAndSourceFeeSet is a free log retrieval operation binding the contract event 0x53fbcb0688fd3904cd33601e03a5eba0c6bba217d8c9312333ea0c116ea91382.
//
// Solidity: event TypeAndSourceFeeSet(bytes32 indexed _type, bytes32 indexed source, uint256 fee)
func (_Fdc *FdcFilterer) FilterTypeAndSourceFeeSet(opts *bind.FilterOpts, _type [][32]byte, source [][32]byte) (*FdcTypeAndSourceFeeSetIterator, error) {

	var _typeRule []interface{}
	for _, _typeItem := range _type {
		_typeRule = append(_typeRule, _typeItem)
	}
	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _Fdc.contract.FilterLogs(opts, "TypeAndSourceFeeSet", _typeRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return &FdcTypeAndSourceFeeSetIterator{contract: _Fdc.contract, event: "TypeAndSourceFeeSet", logs: logs, sub: sub}, nil
}

// WatchTypeAndSourceFeeSet is a free log subscription operation binding the contract event 0x53fbcb0688fd3904cd33601e03a5eba0c6bba217d8c9312333ea0c116ea91382.
//
// Solidity: event TypeAndSourceFeeSet(bytes32 indexed _type, bytes32 indexed source, uint256 fee)
func (_Fdc *FdcFilterer) WatchTypeAndSourceFeeSet(opts *bind.WatchOpts, sink chan<- *FdcTypeAndSourceFeeSet, _type [][32]byte, source [][32]byte) (event.Subscription, error) {

	var _typeRule []interface{}
	for _, _typeItem := range _type {
		_typeRule = append(_typeRule, _typeItem)
	}
	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _Fdc.contract.WatchLogs(opts, "TypeAndSourceFeeSet", _typeRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcTypeAndSourceFeeSet)
				if err := _Fdc.contract.UnpackLog(event, "TypeAndSourceFeeSet", log); err != nil {
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

// ParseTypeAndSourceFeeSet is a log parse operation binding the contract event 0x53fbcb0688fd3904cd33601e03a5eba0c6bba217d8c9312333ea0c116ea91382.
//
// Solidity: event TypeAndSourceFeeSet(bytes32 indexed _type, bytes32 indexed source, uint256 fee)
func (_Fdc *FdcFilterer) ParseTypeAndSourceFeeSet(log types.Log) (*FdcTypeAndSourceFeeSet, error) {
	event := new(FdcTypeAndSourceFeeSet)
	if err := _Fdc.contract.UnpackLog(event, "TypeAndSourceFeeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
