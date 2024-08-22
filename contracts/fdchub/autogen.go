// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package fdchub

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

// IFdcInflationConfigurationsFdcConfiguration is an auto generated low-level Go binding around an user-defined struct.
type IFdcInflationConfigurationsFdcConfiguration struct {
	AttestationType      [32]byte
	Source               [32]byte
	InflationShare       *big.Int
	MinRequestsThreshold uint8
	Mode                 *big.Int
}

// FdcHubMetaData contains all meta data concerning the FdcHub contract.
var FdcHubMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIGovernanceSettings\",\"name\":\"_governanceSettings\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_initialGovernance\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_addressUpdater\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_requestsOffsetSeconds\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"AttestationRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"authorizedAmountWei\",\"type\":\"uint256\"}],\"name\":\"DailyAuthorizedInflationSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"allowedAfterTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encodedCall\",\"type\":\"bytes\"}],\"name\":\"GovernanceCallTimelocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"initialGovernance\",\"type\":\"address\"}],\"name\":\"GovernanceInitialised\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"governanceSettings\",\"type\":\"address\"}],\"name\":\"GovernedProductionModeEntered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountReceivedWei\",\"type\":\"uint256\"}],\"name\":\"InflationReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint24\",\"name\":\"rewardEpochId\",\"type\":\"uint24\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"attestationType\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"source\",\"type\":\"bytes32\"},{\"internalType\":\"uint24\",\"name\":\"inflationShare\",\"type\":\"uint24\"},{\"internalType\":\"uint8\",\"name\":\"minRequestsThreshold\",\"type\":\"uint8\"},{\"internalType\":\"uint224\",\"name\":\"mode\",\"type\":\"uint224\"}],\"indexed\":false,\"internalType\":\"structIFdcInflationConfigurations.FdcConfiguration[]\",\"name\":\"fdcConfigurations\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"InflationRewardsOffered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"requestsOffsetSeconds\",\"type\":\"uint8\"}],\"name\":\"RequestsOffsetSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"TimelockedGovernanceCallCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"TimelockedGovernanceCallExecuted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_selector\",\"type\":\"bytes4\"}],\"name\":\"cancelGovernanceCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dailyAuthorizedInflation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_selector\",\"type\":\"bytes4\"}],\"name\":\"executeGovernanceCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fdcInflationConfigurations\",\"outputs\":[{\"internalType\":\"contractIFdcInflationConfigurations\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fdcRequestFeeConfigurations\",\"outputs\":[{\"internalType\":\"contractIFdcRequestFeeConfigurations\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"flareSystemsManager\",\"outputs\":[{\"internalType\":\"contractIIFlareSystemsManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAddressUpdater\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_addressUpdater\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContractName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExpectedBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getInflationAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenPoolSupplyData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_lockedFundsWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalInflationAuthorizedWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalClaimedWei\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"governance\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"governanceSettings\",\"outputs\":[{\"internalType\":\"contractIGovernanceSettings\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIGovernanceSettings\",\"name\":\"_governanceSettings\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_initialGovernance\",\"type\":\"address\"}],\"name\":\"initialise\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"isExecutor\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastInflationAuthorizationReceivedTs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastInflationReceivedTs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"productionMode\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"receiveInflation\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"requestAttestation\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestsOffsetSeconds\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardManager\",\"outputs\":[{\"internalType\":\"contractIIRewardManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_toAuthorizeWei\",\"type\":\"uint256\"}],\"name\":\"setDailyAuthorizedInflation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_requestsOffsetSeconds\",\"type\":\"uint8\"}],\"name\":\"setRequestsOffset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"switchToProductionMode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"}],\"name\":\"timelockedCalls\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"allowedAfterTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"encodedCall\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalInflationAuthorizedWei\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalInflationReceivedWei\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalInflationRewardsOfferedWei\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint24\",\"name\":\"_currentRewardEpochId\",\"type\":\"uint24\"},{\"internalType\":\"uint64\",\"name\":\"_currentRewardEpochExpectedEndTs\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_rewardEpochDurationSeconds\",\"type\":\"uint64\"}],\"name\":\"triggerRewardEpochSwitchover\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_contractNameHashes\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"_contractAddresses\",\"type\":\"address[]\"}],\"name\":\"updateContractAddresses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// FdcHubABI is the input ABI used to generate the binding from.
// Deprecated: Use FdcHubMetaData.ABI instead.
var FdcHubABI = FdcHubMetaData.ABI

// FdcHub is an auto generated Go binding around an Ethereum contract.
type FdcHub struct {
	FdcHubCaller     // Read-only binding to the contract
	FdcHubTransactor // Write-only binding to the contract
	FdcHubFilterer   // Log filterer for contract events
}

// FdcHubCaller is an auto generated read-only Go binding around an Ethereum contract.
type FdcHubCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FdcHubTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FdcHubTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FdcHubFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FdcHubFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FdcHubSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FdcHubSession struct {
	Contract     *FdcHub           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FdcHubCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FdcHubCallerSession struct {
	Contract *FdcHubCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FdcHubTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FdcHubTransactorSession struct {
	Contract     *FdcHubTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FdcHubRaw is an auto generated low-level Go binding around an Ethereum contract.
type FdcHubRaw struct {
	Contract *FdcHub // Generic contract binding to access the raw methods on
}

// FdcHubCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FdcHubCallerRaw struct {
	Contract *FdcHubCaller // Generic read-only contract binding to access the raw methods on
}

// FdcHubTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FdcHubTransactorRaw struct {
	Contract *FdcHubTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFdcHub creates a new instance of FdcHub, bound to a specific deployed contract.
func NewFdcHub(address common.Address, backend bind.ContractBackend) (*FdcHub, error) {
	contract, err := bindFdcHub(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FdcHub{FdcHubCaller: FdcHubCaller{contract: contract}, FdcHubTransactor: FdcHubTransactor{contract: contract}, FdcHubFilterer: FdcHubFilterer{contract: contract}}, nil
}

// NewFdcHubCaller creates a new read-only instance of FdcHub, bound to a specific deployed contract.
func NewFdcHubCaller(address common.Address, caller bind.ContractCaller) (*FdcHubCaller, error) {
	contract, err := bindFdcHub(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FdcHubCaller{contract: contract}, nil
}

// NewFdcHubTransactor creates a new write-only instance of FdcHub, bound to a specific deployed contract.
func NewFdcHubTransactor(address common.Address, transactor bind.ContractTransactor) (*FdcHubTransactor, error) {
	contract, err := bindFdcHub(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FdcHubTransactor{contract: contract}, nil
}

// NewFdcHubFilterer creates a new log filterer instance of FdcHub, bound to a specific deployed contract.
func NewFdcHubFilterer(address common.Address, filterer bind.ContractFilterer) (*FdcHubFilterer, error) {
	contract, err := bindFdcHub(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FdcHubFilterer{contract: contract}, nil
}

// bindFdcHub binds a generic wrapper to an already deployed contract.
func bindFdcHub(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FdcHubMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FdcHub *FdcHubRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FdcHub.Contract.FdcHubCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FdcHub *FdcHubRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FdcHub.Contract.FdcHubTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FdcHub *FdcHubRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FdcHub.Contract.FdcHubTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FdcHub *FdcHubCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FdcHub.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FdcHub *FdcHubTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FdcHub.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FdcHub *FdcHubTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FdcHub.Contract.contract.Transact(opts, method, params...)
}

// DailyAuthorizedInflation is a free data retrieval call binding the contract method 0x708e34ce.
//
// Solidity: function dailyAuthorizedInflation() view returns(uint256)
func (_FdcHub *FdcHubCaller) DailyAuthorizedInflation(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "dailyAuthorizedInflation")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DailyAuthorizedInflation is a free data retrieval call binding the contract method 0x708e34ce.
//
// Solidity: function dailyAuthorizedInflation() view returns(uint256)
func (_FdcHub *FdcHubSession) DailyAuthorizedInflation() (*big.Int, error) {
	return _FdcHub.Contract.DailyAuthorizedInflation(&_FdcHub.CallOpts)
}

// DailyAuthorizedInflation is a free data retrieval call binding the contract method 0x708e34ce.
//
// Solidity: function dailyAuthorizedInflation() view returns(uint256)
func (_FdcHub *FdcHubCallerSession) DailyAuthorizedInflation() (*big.Int, error) {
	return _FdcHub.Contract.DailyAuthorizedInflation(&_FdcHub.CallOpts)
}

// FdcInflationConfigurations is a free data retrieval call binding the contract method 0x4c5a1d28.
//
// Solidity: function fdcInflationConfigurations() view returns(address)
func (_FdcHub *FdcHubCaller) FdcInflationConfigurations(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "fdcInflationConfigurations")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FdcInflationConfigurations is a free data retrieval call binding the contract method 0x4c5a1d28.
//
// Solidity: function fdcInflationConfigurations() view returns(address)
func (_FdcHub *FdcHubSession) FdcInflationConfigurations() (common.Address, error) {
	return _FdcHub.Contract.FdcInflationConfigurations(&_FdcHub.CallOpts)
}

// FdcInflationConfigurations is a free data retrieval call binding the contract method 0x4c5a1d28.
//
// Solidity: function fdcInflationConfigurations() view returns(address)
func (_FdcHub *FdcHubCallerSession) FdcInflationConfigurations() (common.Address, error) {
	return _FdcHub.Contract.FdcInflationConfigurations(&_FdcHub.CallOpts)
}

// FdcRequestFeeConfigurations is a free data retrieval call binding the contract method 0x116ea702.
//
// Solidity: function fdcRequestFeeConfigurations() view returns(address)
func (_FdcHub *FdcHubCaller) FdcRequestFeeConfigurations(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "fdcRequestFeeConfigurations")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FdcRequestFeeConfigurations is a free data retrieval call binding the contract method 0x116ea702.
//
// Solidity: function fdcRequestFeeConfigurations() view returns(address)
func (_FdcHub *FdcHubSession) FdcRequestFeeConfigurations() (common.Address, error) {
	return _FdcHub.Contract.FdcRequestFeeConfigurations(&_FdcHub.CallOpts)
}

// FdcRequestFeeConfigurations is a free data retrieval call binding the contract method 0x116ea702.
//
// Solidity: function fdcRequestFeeConfigurations() view returns(address)
func (_FdcHub *FdcHubCallerSession) FdcRequestFeeConfigurations() (common.Address, error) {
	return _FdcHub.Contract.FdcRequestFeeConfigurations(&_FdcHub.CallOpts)
}

// FlareSystemsManager is a free data retrieval call binding the contract method 0xfaae7fc9.
//
// Solidity: function flareSystemsManager() view returns(address)
func (_FdcHub *FdcHubCaller) FlareSystemsManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "flareSystemsManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FlareSystemsManager is a free data retrieval call binding the contract method 0xfaae7fc9.
//
// Solidity: function flareSystemsManager() view returns(address)
func (_FdcHub *FdcHubSession) FlareSystemsManager() (common.Address, error) {
	return _FdcHub.Contract.FlareSystemsManager(&_FdcHub.CallOpts)
}

// FlareSystemsManager is a free data retrieval call binding the contract method 0xfaae7fc9.
//
// Solidity: function flareSystemsManager() view returns(address)
func (_FdcHub *FdcHubCallerSession) FlareSystemsManager() (common.Address, error) {
	return _FdcHub.Contract.FlareSystemsManager(&_FdcHub.CallOpts)
}

// GetAddressUpdater is a free data retrieval call binding the contract method 0x5267a15d.
//
// Solidity: function getAddressUpdater() view returns(address _addressUpdater)
func (_FdcHub *FdcHubCaller) GetAddressUpdater(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "getAddressUpdater")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAddressUpdater is a free data retrieval call binding the contract method 0x5267a15d.
//
// Solidity: function getAddressUpdater() view returns(address _addressUpdater)
func (_FdcHub *FdcHubSession) GetAddressUpdater() (common.Address, error) {
	return _FdcHub.Contract.GetAddressUpdater(&_FdcHub.CallOpts)
}

// GetAddressUpdater is a free data retrieval call binding the contract method 0x5267a15d.
//
// Solidity: function getAddressUpdater() view returns(address _addressUpdater)
func (_FdcHub *FdcHubCallerSession) GetAddressUpdater() (common.Address, error) {
	return _FdcHub.Contract.GetAddressUpdater(&_FdcHub.CallOpts)
}

// GetContractName is a free data retrieval call binding the contract method 0xf5f5ba72.
//
// Solidity: function getContractName() pure returns(string)
func (_FdcHub *FdcHubCaller) GetContractName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "getContractName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetContractName is a free data retrieval call binding the contract method 0xf5f5ba72.
//
// Solidity: function getContractName() pure returns(string)
func (_FdcHub *FdcHubSession) GetContractName() (string, error) {
	return _FdcHub.Contract.GetContractName(&_FdcHub.CallOpts)
}

// GetContractName is a free data retrieval call binding the contract method 0xf5f5ba72.
//
// Solidity: function getContractName() pure returns(string)
func (_FdcHub *FdcHubCallerSession) GetContractName() (string, error) {
	return _FdcHub.Contract.GetContractName(&_FdcHub.CallOpts)
}

// GetExpectedBalance is a free data retrieval call binding the contract method 0xaf04cd3b.
//
// Solidity: function getExpectedBalance() view returns(uint256)
func (_FdcHub *FdcHubCaller) GetExpectedBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "getExpectedBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetExpectedBalance is a free data retrieval call binding the contract method 0xaf04cd3b.
//
// Solidity: function getExpectedBalance() view returns(uint256)
func (_FdcHub *FdcHubSession) GetExpectedBalance() (*big.Int, error) {
	return _FdcHub.Contract.GetExpectedBalance(&_FdcHub.CallOpts)
}

// GetExpectedBalance is a free data retrieval call binding the contract method 0xaf04cd3b.
//
// Solidity: function getExpectedBalance() view returns(uint256)
func (_FdcHub *FdcHubCallerSession) GetExpectedBalance() (*big.Int, error) {
	return _FdcHub.Contract.GetExpectedBalance(&_FdcHub.CallOpts)
}

// GetInflationAddress is a free data retrieval call binding the contract method 0xed39d3f8.
//
// Solidity: function getInflationAddress() view returns(address)
func (_FdcHub *FdcHubCaller) GetInflationAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "getInflationAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetInflationAddress is a free data retrieval call binding the contract method 0xed39d3f8.
//
// Solidity: function getInflationAddress() view returns(address)
func (_FdcHub *FdcHubSession) GetInflationAddress() (common.Address, error) {
	return _FdcHub.Contract.GetInflationAddress(&_FdcHub.CallOpts)
}

// GetInflationAddress is a free data retrieval call binding the contract method 0xed39d3f8.
//
// Solidity: function getInflationAddress() view returns(address)
func (_FdcHub *FdcHubCallerSession) GetInflationAddress() (common.Address, error) {
	return _FdcHub.Contract.GetInflationAddress(&_FdcHub.CallOpts)
}

// GetTokenPoolSupplyData is a free data retrieval call binding the contract method 0x2dafdbbf.
//
// Solidity: function getTokenPoolSupplyData() view returns(uint256 _lockedFundsWei, uint256 _totalInflationAuthorizedWei, uint256 _totalClaimedWei)
func (_FdcHub *FdcHubCaller) GetTokenPoolSupplyData(opts *bind.CallOpts) (struct {
	LockedFundsWei              *big.Int
	TotalInflationAuthorizedWei *big.Int
	TotalClaimedWei             *big.Int
}, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "getTokenPoolSupplyData")

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
func (_FdcHub *FdcHubSession) GetTokenPoolSupplyData() (struct {
	LockedFundsWei              *big.Int
	TotalInflationAuthorizedWei *big.Int
	TotalClaimedWei             *big.Int
}, error) {
	return _FdcHub.Contract.GetTokenPoolSupplyData(&_FdcHub.CallOpts)
}

// GetTokenPoolSupplyData is a free data retrieval call binding the contract method 0x2dafdbbf.
//
// Solidity: function getTokenPoolSupplyData() view returns(uint256 _lockedFundsWei, uint256 _totalInflationAuthorizedWei, uint256 _totalClaimedWei)
func (_FdcHub *FdcHubCallerSession) GetTokenPoolSupplyData() (struct {
	LockedFundsWei              *big.Int
	TotalInflationAuthorizedWei *big.Int
	TotalClaimedWei             *big.Int
}, error) {
	return _FdcHub.Contract.GetTokenPoolSupplyData(&_FdcHub.CallOpts)
}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() view returns(address)
func (_FdcHub *FdcHubCaller) Governance(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "governance")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() view returns(address)
func (_FdcHub *FdcHubSession) Governance() (common.Address, error) {
	return _FdcHub.Contract.Governance(&_FdcHub.CallOpts)
}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() view returns(address)
func (_FdcHub *FdcHubCallerSession) Governance() (common.Address, error) {
	return _FdcHub.Contract.Governance(&_FdcHub.CallOpts)
}

// GovernanceSettings is a free data retrieval call binding the contract method 0x62354e03.
//
// Solidity: function governanceSettings() view returns(address)
func (_FdcHub *FdcHubCaller) GovernanceSettings(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "governanceSettings")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GovernanceSettings is a free data retrieval call binding the contract method 0x62354e03.
//
// Solidity: function governanceSettings() view returns(address)
func (_FdcHub *FdcHubSession) GovernanceSettings() (common.Address, error) {
	return _FdcHub.Contract.GovernanceSettings(&_FdcHub.CallOpts)
}

// GovernanceSettings is a free data retrieval call binding the contract method 0x62354e03.
//
// Solidity: function governanceSettings() view returns(address)
func (_FdcHub *FdcHubCallerSession) GovernanceSettings() (common.Address, error) {
	return _FdcHub.Contract.GovernanceSettings(&_FdcHub.CallOpts)
}

// IsExecutor is a free data retrieval call binding the contract method 0xdebfda30.
//
// Solidity: function isExecutor(address _address) view returns(bool)
func (_FdcHub *FdcHubCaller) IsExecutor(opts *bind.CallOpts, _address common.Address) (bool, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "isExecutor", _address)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsExecutor is a free data retrieval call binding the contract method 0xdebfda30.
//
// Solidity: function isExecutor(address _address) view returns(bool)
func (_FdcHub *FdcHubSession) IsExecutor(_address common.Address) (bool, error) {
	return _FdcHub.Contract.IsExecutor(&_FdcHub.CallOpts, _address)
}

// IsExecutor is a free data retrieval call binding the contract method 0xdebfda30.
//
// Solidity: function isExecutor(address _address) view returns(bool)
func (_FdcHub *FdcHubCallerSession) IsExecutor(_address common.Address) (bool, error) {
	return _FdcHub.Contract.IsExecutor(&_FdcHub.CallOpts, _address)
}

// LastInflationAuthorizationReceivedTs is a free data retrieval call binding the contract method 0x473252c4.
//
// Solidity: function lastInflationAuthorizationReceivedTs() view returns(uint256)
func (_FdcHub *FdcHubCaller) LastInflationAuthorizationReceivedTs(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "lastInflationAuthorizationReceivedTs")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastInflationAuthorizationReceivedTs is a free data retrieval call binding the contract method 0x473252c4.
//
// Solidity: function lastInflationAuthorizationReceivedTs() view returns(uint256)
func (_FdcHub *FdcHubSession) LastInflationAuthorizationReceivedTs() (*big.Int, error) {
	return _FdcHub.Contract.LastInflationAuthorizationReceivedTs(&_FdcHub.CallOpts)
}

// LastInflationAuthorizationReceivedTs is a free data retrieval call binding the contract method 0x473252c4.
//
// Solidity: function lastInflationAuthorizationReceivedTs() view returns(uint256)
func (_FdcHub *FdcHubCallerSession) LastInflationAuthorizationReceivedTs() (*big.Int, error) {
	return _FdcHub.Contract.LastInflationAuthorizationReceivedTs(&_FdcHub.CallOpts)
}

// LastInflationReceivedTs is a free data retrieval call binding the contract method 0x12afcf0b.
//
// Solidity: function lastInflationReceivedTs() view returns(uint256)
func (_FdcHub *FdcHubCaller) LastInflationReceivedTs(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "lastInflationReceivedTs")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastInflationReceivedTs is a free data retrieval call binding the contract method 0x12afcf0b.
//
// Solidity: function lastInflationReceivedTs() view returns(uint256)
func (_FdcHub *FdcHubSession) LastInflationReceivedTs() (*big.Int, error) {
	return _FdcHub.Contract.LastInflationReceivedTs(&_FdcHub.CallOpts)
}

// LastInflationReceivedTs is a free data retrieval call binding the contract method 0x12afcf0b.
//
// Solidity: function lastInflationReceivedTs() view returns(uint256)
func (_FdcHub *FdcHubCallerSession) LastInflationReceivedTs() (*big.Int, error) {
	return _FdcHub.Contract.LastInflationReceivedTs(&_FdcHub.CallOpts)
}

// ProductionMode is a free data retrieval call binding the contract method 0xe17f212e.
//
// Solidity: function productionMode() view returns(bool)
func (_FdcHub *FdcHubCaller) ProductionMode(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "productionMode")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ProductionMode is a free data retrieval call binding the contract method 0xe17f212e.
//
// Solidity: function productionMode() view returns(bool)
func (_FdcHub *FdcHubSession) ProductionMode() (bool, error) {
	return _FdcHub.Contract.ProductionMode(&_FdcHub.CallOpts)
}

// ProductionMode is a free data retrieval call binding the contract method 0xe17f212e.
//
// Solidity: function productionMode() view returns(bool)
func (_FdcHub *FdcHubCallerSession) ProductionMode() (bool, error) {
	return _FdcHub.Contract.ProductionMode(&_FdcHub.CallOpts)
}

// RequestsOffsetSeconds is a free data retrieval call binding the contract method 0x94d019f1.
//
// Solidity: function requestsOffsetSeconds() view returns(uint8)
func (_FdcHub *FdcHubCaller) RequestsOffsetSeconds(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "requestsOffsetSeconds")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// RequestsOffsetSeconds is a free data retrieval call binding the contract method 0x94d019f1.
//
// Solidity: function requestsOffsetSeconds() view returns(uint8)
func (_FdcHub *FdcHubSession) RequestsOffsetSeconds() (uint8, error) {
	return _FdcHub.Contract.RequestsOffsetSeconds(&_FdcHub.CallOpts)
}

// RequestsOffsetSeconds is a free data retrieval call binding the contract method 0x94d019f1.
//
// Solidity: function requestsOffsetSeconds() view returns(uint8)
func (_FdcHub *FdcHubCallerSession) RequestsOffsetSeconds() (uint8, error) {
	return _FdcHub.Contract.RequestsOffsetSeconds(&_FdcHub.CallOpts)
}

// RewardManager is a free data retrieval call binding the contract method 0x0f4ef8a6.
//
// Solidity: function rewardManager() view returns(address)
func (_FdcHub *FdcHubCaller) RewardManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "rewardManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardManager is a free data retrieval call binding the contract method 0x0f4ef8a6.
//
// Solidity: function rewardManager() view returns(address)
func (_FdcHub *FdcHubSession) RewardManager() (common.Address, error) {
	return _FdcHub.Contract.RewardManager(&_FdcHub.CallOpts)
}

// RewardManager is a free data retrieval call binding the contract method 0x0f4ef8a6.
//
// Solidity: function rewardManager() view returns(address)
func (_FdcHub *FdcHubCallerSession) RewardManager() (common.Address, error) {
	return _FdcHub.Contract.RewardManager(&_FdcHub.CallOpts)
}

// TimelockedCalls is a free data retrieval call binding the contract method 0x74e6310e.
//
// Solidity: function timelockedCalls(bytes4 selector) view returns(uint256 allowedAfterTimestamp, bytes encodedCall)
func (_FdcHub *FdcHubCaller) TimelockedCalls(opts *bind.CallOpts, selector [4]byte) (struct {
	AllowedAfterTimestamp *big.Int
	EncodedCall           []byte
}, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "timelockedCalls", selector)

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
func (_FdcHub *FdcHubSession) TimelockedCalls(selector [4]byte) (struct {
	AllowedAfterTimestamp *big.Int
	EncodedCall           []byte
}, error) {
	return _FdcHub.Contract.TimelockedCalls(&_FdcHub.CallOpts, selector)
}

// TimelockedCalls is a free data retrieval call binding the contract method 0x74e6310e.
//
// Solidity: function timelockedCalls(bytes4 selector) view returns(uint256 allowedAfterTimestamp, bytes encodedCall)
func (_FdcHub *FdcHubCallerSession) TimelockedCalls(selector [4]byte) (struct {
	AllowedAfterTimestamp *big.Int
	EncodedCall           []byte
}, error) {
	return _FdcHub.Contract.TimelockedCalls(&_FdcHub.CallOpts, selector)
}

// TotalInflationAuthorizedWei is a free data retrieval call binding the contract method 0xd0c1c393.
//
// Solidity: function totalInflationAuthorizedWei() view returns(uint256)
func (_FdcHub *FdcHubCaller) TotalInflationAuthorizedWei(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "totalInflationAuthorizedWei")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalInflationAuthorizedWei is a free data retrieval call binding the contract method 0xd0c1c393.
//
// Solidity: function totalInflationAuthorizedWei() view returns(uint256)
func (_FdcHub *FdcHubSession) TotalInflationAuthorizedWei() (*big.Int, error) {
	return _FdcHub.Contract.TotalInflationAuthorizedWei(&_FdcHub.CallOpts)
}

// TotalInflationAuthorizedWei is a free data retrieval call binding the contract method 0xd0c1c393.
//
// Solidity: function totalInflationAuthorizedWei() view returns(uint256)
func (_FdcHub *FdcHubCallerSession) TotalInflationAuthorizedWei() (*big.Int, error) {
	return _FdcHub.Contract.TotalInflationAuthorizedWei(&_FdcHub.CallOpts)
}

// TotalInflationReceivedWei is a free data retrieval call binding the contract method 0xa5555aea.
//
// Solidity: function totalInflationReceivedWei() view returns(uint256)
func (_FdcHub *FdcHubCaller) TotalInflationReceivedWei(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "totalInflationReceivedWei")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalInflationReceivedWei is a free data retrieval call binding the contract method 0xa5555aea.
//
// Solidity: function totalInflationReceivedWei() view returns(uint256)
func (_FdcHub *FdcHubSession) TotalInflationReceivedWei() (*big.Int, error) {
	return _FdcHub.Contract.TotalInflationReceivedWei(&_FdcHub.CallOpts)
}

// TotalInflationReceivedWei is a free data retrieval call binding the contract method 0xa5555aea.
//
// Solidity: function totalInflationReceivedWei() view returns(uint256)
func (_FdcHub *FdcHubCallerSession) TotalInflationReceivedWei() (*big.Int, error) {
	return _FdcHub.Contract.TotalInflationReceivedWei(&_FdcHub.CallOpts)
}

// TotalInflationRewardsOfferedWei is a free data retrieval call binding the contract method 0xbd76b69c.
//
// Solidity: function totalInflationRewardsOfferedWei() view returns(uint256)
func (_FdcHub *FdcHubCaller) TotalInflationRewardsOfferedWei(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FdcHub.contract.Call(opts, &out, "totalInflationRewardsOfferedWei")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalInflationRewardsOfferedWei is a free data retrieval call binding the contract method 0xbd76b69c.
//
// Solidity: function totalInflationRewardsOfferedWei() view returns(uint256)
func (_FdcHub *FdcHubSession) TotalInflationRewardsOfferedWei() (*big.Int, error) {
	return _FdcHub.Contract.TotalInflationRewardsOfferedWei(&_FdcHub.CallOpts)
}

// TotalInflationRewardsOfferedWei is a free data retrieval call binding the contract method 0xbd76b69c.
//
// Solidity: function totalInflationRewardsOfferedWei() view returns(uint256)
func (_FdcHub *FdcHubCallerSession) TotalInflationRewardsOfferedWei() (*big.Int, error) {
	return _FdcHub.Contract.TotalInflationRewardsOfferedWei(&_FdcHub.CallOpts)
}

// CancelGovernanceCall is a paid mutator transaction binding the contract method 0x67fc4029.
//
// Solidity: function cancelGovernanceCall(bytes4 _selector) returns()
func (_FdcHub *FdcHubTransactor) CancelGovernanceCall(opts *bind.TransactOpts, _selector [4]byte) (*types.Transaction, error) {
	return _FdcHub.contract.Transact(opts, "cancelGovernanceCall", _selector)
}

// CancelGovernanceCall is a paid mutator transaction binding the contract method 0x67fc4029.
//
// Solidity: function cancelGovernanceCall(bytes4 _selector) returns()
func (_FdcHub *FdcHubSession) CancelGovernanceCall(_selector [4]byte) (*types.Transaction, error) {
	return _FdcHub.Contract.CancelGovernanceCall(&_FdcHub.TransactOpts, _selector)
}

// CancelGovernanceCall is a paid mutator transaction binding the contract method 0x67fc4029.
//
// Solidity: function cancelGovernanceCall(bytes4 _selector) returns()
func (_FdcHub *FdcHubTransactorSession) CancelGovernanceCall(_selector [4]byte) (*types.Transaction, error) {
	return _FdcHub.Contract.CancelGovernanceCall(&_FdcHub.TransactOpts, _selector)
}

// ExecuteGovernanceCall is a paid mutator transaction binding the contract method 0x5ff27079.
//
// Solidity: function executeGovernanceCall(bytes4 _selector) returns()
func (_FdcHub *FdcHubTransactor) ExecuteGovernanceCall(opts *bind.TransactOpts, _selector [4]byte) (*types.Transaction, error) {
	return _FdcHub.contract.Transact(opts, "executeGovernanceCall", _selector)
}

// ExecuteGovernanceCall is a paid mutator transaction binding the contract method 0x5ff27079.
//
// Solidity: function executeGovernanceCall(bytes4 _selector) returns()
func (_FdcHub *FdcHubSession) ExecuteGovernanceCall(_selector [4]byte) (*types.Transaction, error) {
	return _FdcHub.Contract.ExecuteGovernanceCall(&_FdcHub.TransactOpts, _selector)
}

// ExecuteGovernanceCall is a paid mutator transaction binding the contract method 0x5ff27079.
//
// Solidity: function executeGovernanceCall(bytes4 _selector) returns()
func (_FdcHub *FdcHubTransactorSession) ExecuteGovernanceCall(_selector [4]byte) (*types.Transaction, error) {
	return _FdcHub.Contract.ExecuteGovernanceCall(&_FdcHub.TransactOpts, _selector)
}

// Initialise is a paid mutator transaction binding the contract method 0xef88bf13.
//
// Solidity: function initialise(address _governanceSettings, address _initialGovernance) returns()
func (_FdcHub *FdcHubTransactor) Initialise(opts *bind.TransactOpts, _governanceSettings common.Address, _initialGovernance common.Address) (*types.Transaction, error) {
	return _FdcHub.contract.Transact(opts, "initialise", _governanceSettings, _initialGovernance)
}

// Initialise is a paid mutator transaction binding the contract method 0xef88bf13.
//
// Solidity: function initialise(address _governanceSettings, address _initialGovernance) returns()
func (_FdcHub *FdcHubSession) Initialise(_governanceSettings common.Address, _initialGovernance common.Address) (*types.Transaction, error) {
	return _FdcHub.Contract.Initialise(&_FdcHub.TransactOpts, _governanceSettings, _initialGovernance)
}

// Initialise is a paid mutator transaction binding the contract method 0xef88bf13.
//
// Solidity: function initialise(address _governanceSettings, address _initialGovernance) returns()
func (_FdcHub *FdcHubTransactorSession) Initialise(_governanceSettings common.Address, _initialGovernance common.Address) (*types.Transaction, error) {
	return _FdcHub.Contract.Initialise(&_FdcHub.TransactOpts, _governanceSettings, _initialGovernance)
}

// ReceiveInflation is a paid mutator transaction binding the contract method 0x06201f1d.
//
// Solidity: function receiveInflation() payable returns()
func (_FdcHub *FdcHubTransactor) ReceiveInflation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FdcHub.contract.Transact(opts, "receiveInflation")
}

// ReceiveInflation is a paid mutator transaction binding the contract method 0x06201f1d.
//
// Solidity: function receiveInflation() payable returns()
func (_FdcHub *FdcHubSession) ReceiveInflation() (*types.Transaction, error) {
	return _FdcHub.Contract.ReceiveInflation(&_FdcHub.TransactOpts)
}

// ReceiveInflation is a paid mutator transaction binding the contract method 0x06201f1d.
//
// Solidity: function receiveInflation() payable returns()
func (_FdcHub *FdcHubTransactorSession) ReceiveInflation() (*types.Transaction, error) {
	return _FdcHub.Contract.ReceiveInflation(&_FdcHub.TransactOpts)
}

// RequestAttestation is a paid mutator transaction binding the contract method 0x6238f354.
//
// Solidity: function requestAttestation(bytes _data) payable returns()
func (_FdcHub *FdcHubTransactor) RequestAttestation(opts *bind.TransactOpts, _data []byte) (*types.Transaction, error) {
	return _FdcHub.contract.Transact(opts, "requestAttestation", _data)
}

// RequestAttestation is a paid mutator transaction binding the contract method 0x6238f354.
//
// Solidity: function requestAttestation(bytes _data) payable returns()
func (_FdcHub *FdcHubSession) RequestAttestation(_data []byte) (*types.Transaction, error) {
	return _FdcHub.Contract.RequestAttestation(&_FdcHub.TransactOpts, _data)
}

// RequestAttestation is a paid mutator transaction binding the contract method 0x6238f354.
//
// Solidity: function requestAttestation(bytes _data) payable returns()
func (_FdcHub *FdcHubTransactorSession) RequestAttestation(_data []byte) (*types.Transaction, error) {
	return _FdcHub.Contract.RequestAttestation(&_FdcHub.TransactOpts, _data)
}

// SetDailyAuthorizedInflation is a paid mutator transaction binding the contract method 0xe2739563.
//
// Solidity: function setDailyAuthorizedInflation(uint256 _toAuthorizeWei) returns()
func (_FdcHub *FdcHubTransactor) SetDailyAuthorizedInflation(opts *bind.TransactOpts, _toAuthorizeWei *big.Int) (*types.Transaction, error) {
	return _FdcHub.contract.Transact(opts, "setDailyAuthorizedInflation", _toAuthorizeWei)
}

// SetDailyAuthorizedInflation is a paid mutator transaction binding the contract method 0xe2739563.
//
// Solidity: function setDailyAuthorizedInflation(uint256 _toAuthorizeWei) returns()
func (_FdcHub *FdcHubSession) SetDailyAuthorizedInflation(_toAuthorizeWei *big.Int) (*types.Transaction, error) {
	return _FdcHub.Contract.SetDailyAuthorizedInflation(&_FdcHub.TransactOpts, _toAuthorizeWei)
}

// SetDailyAuthorizedInflation is a paid mutator transaction binding the contract method 0xe2739563.
//
// Solidity: function setDailyAuthorizedInflation(uint256 _toAuthorizeWei) returns()
func (_FdcHub *FdcHubTransactorSession) SetDailyAuthorizedInflation(_toAuthorizeWei *big.Int) (*types.Transaction, error) {
	return _FdcHub.Contract.SetDailyAuthorizedInflation(&_FdcHub.TransactOpts, _toAuthorizeWei)
}

// SetRequestsOffset is a paid mutator transaction binding the contract method 0xbfda6086.
//
// Solidity: function setRequestsOffset(uint8 _requestsOffsetSeconds) returns()
func (_FdcHub *FdcHubTransactor) SetRequestsOffset(opts *bind.TransactOpts, _requestsOffsetSeconds uint8) (*types.Transaction, error) {
	return _FdcHub.contract.Transact(opts, "setRequestsOffset", _requestsOffsetSeconds)
}

// SetRequestsOffset is a paid mutator transaction binding the contract method 0xbfda6086.
//
// Solidity: function setRequestsOffset(uint8 _requestsOffsetSeconds) returns()
func (_FdcHub *FdcHubSession) SetRequestsOffset(_requestsOffsetSeconds uint8) (*types.Transaction, error) {
	return _FdcHub.Contract.SetRequestsOffset(&_FdcHub.TransactOpts, _requestsOffsetSeconds)
}

// SetRequestsOffset is a paid mutator transaction binding the contract method 0xbfda6086.
//
// Solidity: function setRequestsOffset(uint8 _requestsOffsetSeconds) returns()
func (_FdcHub *FdcHubTransactorSession) SetRequestsOffset(_requestsOffsetSeconds uint8) (*types.Transaction, error) {
	return _FdcHub.Contract.SetRequestsOffset(&_FdcHub.TransactOpts, _requestsOffsetSeconds)
}

// SwitchToProductionMode is a paid mutator transaction binding the contract method 0xf5a98383.
//
// Solidity: function switchToProductionMode() returns()
func (_FdcHub *FdcHubTransactor) SwitchToProductionMode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FdcHub.contract.Transact(opts, "switchToProductionMode")
}

// SwitchToProductionMode is a paid mutator transaction binding the contract method 0xf5a98383.
//
// Solidity: function switchToProductionMode() returns()
func (_FdcHub *FdcHubSession) SwitchToProductionMode() (*types.Transaction, error) {
	return _FdcHub.Contract.SwitchToProductionMode(&_FdcHub.TransactOpts)
}

// SwitchToProductionMode is a paid mutator transaction binding the contract method 0xf5a98383.
//
// Solidity: function switchToProductionMode() returns()
func (_FdcHub *FdcHubTransactorSession) SwitchToProductionMode() (*types.Transaction, error) {
	return _FdcHub.Contract.SwitchToProductionMode(&_FdcHub.TransactOpts)
}

// TriggerRewardEpochSwitchover is a paid mutator transaction binding the contract method 0x91f25679.
//
// Solidity: function triggerRewardEpochSwitchover(uint24 _currentRewardEpochId, uint64 _currentRewardEpochExpectedEndTs, uint64 _rewardEpochDurationSeconds) returns()
func (_FdcHub *FdcHubTransactor) TriggerRewardEpochSwitchover(opts *bind.TransactOpts, _currentRewardEpochId *big.Int, _currentRewardEpochExpectedEndTs uint64, _rewardEpochDurationSeconds uint64) (*types.Transaction, error) {
	return _FdcHub.contract.Transact(opts, "triggerRewardEpochSwitchover", _currentRewardEpochId, _currentRewardEpochExpectedEndTs, _rewardEpochDurationSeconds)
}

// TriggerRewardEpochSwitchover is a paid mutator transaction binding the contract method 0x91f25679.
//
// Solidity: function triggerRewardEpochSwitchover(uint24 _currentRewardEpochId, uint64 _currentRewardEpochExpectedEndTs, uint64 _rewardEpochDurationSeconds) returns()
func (_FdcHub *FdcHubSession) TriggerRewardEpochSwitchover(_currentRewardEpochId *big.Int, _currentRewardEpochExpectedEndTs uint64, _rewardEpochDurationSeconds uint64) (*types.Transaction, error) {
	return _FdcHub.Contract.TriggerRewardEpochSwitchover(&_FdcHub.TransactOpts, _currentRewardEpochId, _currentRewardEpochExpectedEndTs, _rewardEpochDurationSeconds)
}

// TriggerRewardEpochSwitchover is a paid mutator transaction binding the contract method 0x91f25679.
//
// Solidity: function triggerRewardEpochSwitchover(uint24 _currentRewardEpochId, uint64 _currentRewardEpochExpectedEndTs, uint64 _rewardEpochDurationSeconds) returns()
func (_FdcHub *FdcHubTransactorSession) TriggerRewardEpochSwitchover(_currentRewardEpochId *big.Int, _currentRewardEpochExpectedEndTs uint64, _rewardEpochDurationSeconds uint64) (*types.Transaction, error) {
	return _FdcHub.Contract.TriggerRewardEpochSwitchover(&_FdcHub.TransactOpts, _currentRewardEpochId, _currentRewardEpochExpectedEndTs, _rewardEpochDurationSeconds)
}

// UpdateContractAddresses is a paid mutator transaction binding the contract method 0xb00c0b76.
//
// Solidity: function updateContractAddresses(bytes32[] _contractNameHashes, address[] _contractAddresses) returns()
func (_FdcHub *FdcHubTransactor) UpdateContractAddresses(opts *bind.TransactOpts, _contractNameHashes [][32]byte, _contractAddresses []common.Address) (*types.Transaction, error) {
	return _FdcHub.contract.Transact(opts, "updateContractAddresses", _contractNameHashes, _contractAddresses)
}

// UpdateContractAddresses is a paid mutator transaction binding the contract method 0xb00c0b76.
//
// Solidity: function updateContractAddresses(bytes32[] _contractNameHashes, address[] _contractAddresses) returns()
func (_FdcHub *FdcHubSession) UpdateContractAddresses(_contractNameHashes [][32]byte, _contractAddresses []common.Address) (*types.Transaction, error) {
	return _FdcHub.Contract.UpdateContractAddresses(&_FdcHub.TransactOpts, _contractNameHashes, _contractAddresses)
}

// UpdateContractAddresses is a paid mutator transaction binding the contract method 0xb00c0b76.
//
// Solidity: function updateContractAddresses(bytes32[] _contractNameHashes, address[] _contractAddresses) returns()
func (_FdcHub *FdcHubTransactorSession) UpdateContractAddresses(_contractNameHashes [][32]byte, _contractAddresses []common.Address) (*types.Transaction, error) {
	return _FdcHub.Contract.UpdateContractAddresses(&_FdcHub.TransactOpts, _contractNameHashes, _contractAddresses)
}

// FdcHubAttestationRequestIterator is returned from FilterAttestationRequest and is used to iterate over the raw logs and unpacked data for AttestationRequest events raised by the FdcHub contract.
type FdcHubAttestationRequestIterator struct {
	Event *FdcHubAttestationRequest // Event containing the contract specifics and raw log

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
func (it *FdcHubAttestationRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcHubAttestationRequest)
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
		it.Event = new(FdcHubAttestationRequest)
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
func (it *FdcHubAttestationRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcHubAttestationRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcHubAttestationRequest represents a AttestationRequest event raised by the FdcHub contract.
type FdcHubAttestationRequest struct {
	Data []byte
	Fee  *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAttestationRequest is a free log retrieval operation binding the contract event 0x251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9.
//
// Solidity: event AttestationRequest(bytes data, uint256 fee)
func (_FdcHub *FdcHubFilterer) FilterAttestationRequest(opts *bind.FilterOpts) (*FdcHubAttestationRequestIterator, error) {

	logs, sub, err := _FdcHub.contract.FilterLogs(opts, "AttestationRequest")
	if err != nil {
		return nil, err
	}
	return &FdcHubAttestationRequestIterator{contract: _FdcHub.contract, event: "AttestationRequest", logs: logs, sub: sub}, nil
}

// WatchAttestationRequest is a free log subscription operation binding the contract event 0x251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9.
//
// Solidity: event AttestationRequest(bytes data, uint256 fee)
func (_FdcHub *FdcHubFilterer) WatchAttestationRequest(opts *bind.WatchOpts, sink chan<- *FdcHubAttestationRequest) (event.Subscription, error) {

	logs, sub, err := _FdcHub.contract.WatchLogs(opts, "AttestationRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcHubAttestationRequest)
				if err := _FdcHub.contract.UnpackLog(event, "AttestationRequest", log); err != nil {
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
func (_FdcHub *FdcHubFilterer) ParseAttestationRequest(log types.Log) (*FdcHubAttestationRequest, error) {
	event := new(FdcHubAttestationRequest)
	if err := _FdcHub.contract.UnpackLog(event, "AttestationRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcHubDailyAuthorizedInflationSetIterator is returned from FilterDailyAuthorizedInflationSet and is used to iterate over the raw logs and unpacked data for DailyAuthorizedInflationSet events raised by the FdcHub contract.
type FdcHubDailyAuthorizedInflationSetIterator struct {
	Event *FdcHubDailyAuthorizedInflationSet // Event containing the contract specifics and raw log

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
func (it *FdcHubDailyAuthorizedInflationSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcHubDailyAuthorizedInflationSet)
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
		it.Event = new(FdcHubDailyAuthorizedInflationSet)
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
func (it *FdcHubDailyAuthorizedInflationSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcHubDailyAuthorizedInflationSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcHubDailyAuthorizedInflationSet represents a DailyAuthorizedInflationSet event raised by the FdcHub contract.
type FdcHubDailyAuthorizedInflationSet struct {
	AuthorizedAmountWei *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterDailyAuthorizedInflationSet is a free log retrieval operation binding the contract event 0x187f32a0f765499f15b3bb52ed0aebf6015059f230f2ace7e701e60a47669595.
//
// Solidity: event DailyAuthorizedInflationSet(uint256 authorizedAmountWei)
func (_FdcHub *FdcHubFilterer) FilterDailyAuthorizedInflationSet(opts *bind.FilterOpts) (*FdcHubDailyAuthorizedInflationSetIterator, error) {

	logs, sub, err := _FdcHub.contract.FilterLogs(opts, "DailyAuthorizedInflationSet")
	if err != nil {
		return nil, err
	}
	return &FdcHubDailyAuthorizedInflationSetIterator{contract: _FdcHub.contract, event: "DailyAuthorizedInflationSet", logs: logs, sub: sub}, nil
}

// WatchDailyAuthorizedInflationSet is a free log subscription operation binding the contract event 0x187f32a0f765499f15b3bb52ed0aebf6015059f230f2ace7e701e60a47669595.
//
// Solidity: event DailyAuthorizedInflationSet(uint256 authorizedAmountWei)
func (_FdcHub *FdcHubFilterer) WatchDailyAuthorizedInflationSet(opts *bind.WatchOpts, sink chan<- *FdcHubDailyAuthorizedInflationSet) (event.Subscription, error) {

	logs, sub, err := _FdcHub.contract.WatchLogs(opts, "DailyAuthorizedInflationSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcHubDailyAuthorizedInflationSet)
				if err := _FdcHub.contract.UnpackLog(event, "DailyAuthorizedInflationSet", log); err != nil {
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
func (_FdcHub *FdcHubFilterer) ParseDailyAuthorizedInflationSet(log types.Log) (*FdcHubDailyAuthorizedInflationSet, error) {
	event := new(FdcHubDailyAuthorizedInflationSet)
	if err := _FdcHub.contract.UnpackLog(event, "DailyAuthorizedInflationSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcHubGovernanceCallTimelockedIterator is returned from FilterGovernanceCallTimelocked and is used to iterate over the raw logs and unpacked data for GovernanceCallTimelocked events raised by the FdcHub contract.
type FdcHubGovernanceCallTimelockedIterator struct {
	Event *FdcHubGovernanceCallTimelocked // Event containing the contract specifics and raw log

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
func (it *FdcHubGovernanceCallTimelockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcHubGovernanceCallTimelocked)
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
		it.Event = new(FdcHubGovernanceCallTimelocked)
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
func (it *FdcHubGovernanceCallTimelockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcHubGovernanceCallTimelockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcHubGovernanceCallTimelocked represents a GovernanceCallTimelocked event raised by the FdcHub contract.
type FdcHubGovernanceCallTimelocked struct {
	Selector              [4]byte
	AllowedAfterTimestamp *big.Int
	EncodedCall           []byte
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterGovernanceCallTimelocked is a free log retrieval operation binding the contract event 0xed948300a3694aa01d4a6b258bfd664350193d770c0b51f8387277f6d83ea3b6.
//
// Solidity: event GovernanceCallTimelocked(bytes4 selector, uint256 allowedAfterTimestamp, bytes encodedCall)
func (_FdcHub *FdcHubFilterer) FilterGovernanceCallTimelocked(opts *bind.FilterOpts) (*FdcHubGovernanceCallTimelockedIterator, error) {

	logs, sub, err := _FdcHub.contract.FilterLogs(opts, "GovernanceCallTimelocked")
	if err != nil {
		return nil, err
	}
	return &FdcHubGovernanceCallTimelockedIterator{contract: _FdcHub.contract, event: "GovernanceCallTimelocked", logs: logs, sub: sub}, nil
}

// WatchGovernanceCallTimelocked is a free log subscription operation binding the contract event 0xed948300a3694aa01d4a6b258bfd664350193d770c0b51f8387277f6d83ea3b6.
//
// Solidity: event GovernanceCallTimelocked(bytes4 selector, uint256 allowedAfterTimestamp, bytes encodedCall)
func (_FdcHub *FdcHubFilterer) WatchGovernanceCallTimelocked(opts *bind.WatchOpts, sink chan<- *FdcHubGovernanceCallTimelocked) (event.Subscription, error) {

	logs, sub, err := _FdcHub.contract.WatchLogs(opts, "GovernanceCallTimelocked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcHubGovernanceCallTimelocked)
				if err := _FdcHub.contract.UnpackLog(event, "GovernanceCallTimelocked", log); err != nil {
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
func (_FdcHub *FdcHubFilterer) ParseGovernanceCallTimelocked(log types.Log) (*FdcHubGovernanceCallTimelocked, error) {
	event := new(FdcHubGovernanceCallTimelocked)
	if err := _FdcHub.contract.UnpackLog(event, "GovernanceCallTimelocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcHubGovernanceInitialisedIterator is returned from FilterGovernanceInitialised and is used to iterate over the raw logs and unpacked data for GovernanceInitialised events raised by the FdcHub contract.
type FdcHubGovernanceInitialisedIterator struct {
	Event *FdcHubGovernanceInitialised // Event containing the contract specifics and raw log

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
func (it *FdcHubGovernanceInitialisedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcHubGovernanceInitialised)
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
		it.Event = new(FdcHubGovernanceInitialised)
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
func (it *FdcHubGovernanceInitialisedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcHubGovernanceInitialisedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcHubGovernanceInitialised represents a GovernanceInitialised event raised by the FdcHub contract.
type FdcHubGovernanceInitialised struct {
	InitialGovernance common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterGovernanceInitialised is a free log retrieval operation binding the contract event 0x9789733827840833afc031fb2ef9ab6894271f77bad2085687cf4ae5c7bee4db.
//
// Solidity: event GovernanceInitialised(address initialGovernance)
func (_FdcHub *FdcHubFilterer) FilterGovernanceInitialised(opts *bind.FilterOpts) (*FdcHubGovernanceInitialisedIterator, error) {

	logs, sub, err := _FdcHub.contract.FilterLogs(opts, "GovernanceInitialised")
	if err != nil {
		return nil, err
	}
	return &FdcHubGovernanceInitialisedIterator{contract: _FdcHub.contract, event: "GovernanceInitialised", logs: logs, sub: sub}, nil
}

// WatchGovernanceInitialised is a free log subscription operation binding the contract event 0x9789733827840833afc031fb2ef9ab6894271f77bad2085687cf4ae5c7bee4db.
//
// Solidity: event GovernanceInitialised(address initialGovernance)
func (_FdcHub *FdcHubFilterer) WatchGovernanceInitialised(opts *bind.WatchOpts, sink chan<- *FdcHubGovernanceInitialised) (event.Subscription, error) {

	logs, sub, err := _FdcHub.contract.WatchLogs(opts, "GovernanceInitialised")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcHubGovernanceInitialised)
				if err := _FdcHub.contract.UnpackLog(event, "GovernanceInitialised", log); err != nil {
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
func (_FdcHub *FdcHubFilterer) ParseGovernanceInitialised(log types.Log) (*FdcHubGovernanceInitialised, error) {
	event := new(FdcHubGovernanceInitialised)
	if err := _FdcHub.contract.UnpackLog(event, "GovernanceInitialised", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcHubGovernedProductionModeEnteredIterator is returned from FilterGovernedProductionModeEntered and is used to iterate over the raw logs and unpacked data for GovernedProductionModeEntered events raised by the FdcHub contract.
type FdcHubGovernedProductionModeEnteredIterator struct {
	Event *FdcHubGovernedProductionModeEntered // Event containing the contract specifics and raw log

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
func (it *FdcHubGovernedProductionModeEnteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcHubGovernedProductionModeEntered)
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
		it.Event = new(FdcHubGovernedProductionModeEntered)
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
func (it *FdcHubGovernedProductionModeEnteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcHubGovernedProductionModeEnteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcHubGovernedProductionModeEntered represents a GovernedProductionModeEntered event raised by the FdcHub contract.
type FdcHubGovernedProductionModeEntered struct {
	GovernanceSettings common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterGovernedProductionModeEntered is a free log retrieval operation binding the contract event 0x83af113638b5422f9e977cebc0aaf0eaf2188eb9a8baae7f9d46c42b33a1560c.
//
// Solidity: event GovernedProductionModeEntered(address governanceSettings)
func (_FdcHub *FdcHubFilterer) FilterGovernedProductionModeEntered(opts *bind.FilterOpts) (*FdcHubGovernedProductionModeEnteredIterator, error) {

	logs, sub, err := _FdcHub.contract.FilterLogs(opts, "GovernedProductionModeEntered")
	if err != nil {
		return nil, err
	}
	return &FdcHubGovernedProductionModeEnteredIterator{contract: _FdcHub.contract, event: "GovernedProductionModeEntered", logs: logs, sub: sub}, nil
}

// WatchGovernedProductionModeEntered is a free log subscription operation binding the contract event 0x83af113638b5422f9e977cebc0aaf0eaf2188eb9a8baae7f9d46c42b33a1560c.
//
// Solidity: event GovernedProductionModeEntered(address governanceSettings)
func (_FdcHub *FdcHubFilterer) WatchGovernedProductionModeEntered(opts *bind.WatchOpts, sink chan<- *FdcHubGovernedProductionModeEntered) (event.Subscription, error) {

	logs, sub, err := _FdcHub.contract.WatchLogs(opts, "GovernedProductionModeEntered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcHubGovernedProductionModeEntered)
				if err := _FdcHub.contract.UnpackLog(event, "GovernedProductionModeEntered", log); err != nil {
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
func (_FdcHub *FdcHubFilterer) ParseGovernedProductionModeEntered(log types.Log) (*FdcHubGovernedProductionModeEntered, error) {
	event := new(FdcHubGovernedProductionModeEntered)
	if err := _FdcHub.contract.UnpackLog(event, "GovernedProductionModeEntered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcHubInflationReceivedIterator is returned from FilterInflationReceived and is used to iterate over the raw logs and unpacked data for InflationReceived events raised by the FdcHub contract.
type FdcHubInflationReceivedIterator struct {
	Event *FdcHubInflationReceived // Event containing the contract specifics and raw log

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
func (it *FdcHubInflationReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcHubInflationReceived)
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
		it.Event = new(FdcHubInflationReceived)
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
func (it *FdcHubInflationReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcHubInflationReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcHubInflationReceived represents a InflationReceived event raised by the FdcHub contract.
type FdcHubInflationReceived struct {
	AmountReceivedWei *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterInflationReceived is a free log retrieval operation binding the contract event 0x95c4e29cc99bc027cfc3cd719d6fd973d5f0317061885fbb322b9b17d8d35d37.
//
// Solidity: event InflationReceived(uint256 amountReceivedWei)
func (_FdcHub *FdcHubFilterer) FilterInflationReceived(opts *bind.FilterOpts) (*FdcHubInflationReceivedIterator, error) {

	logs, sub, err := _FdcHub.contract.FilterLogs(opts, "InflationReceived")
	if err != nil {
		return nil, err
	}
	return &FdcHubInflationReceivedIterator{contract: _FdcHub.contract, event: "InflationReceived", logs: logs, sub: sub}, nil
}

// WatchInflationReceived is a free log subscription operation binding the contract event 0x95c4e29cc99bc027cfc3cd719d6fd973d5f0317061885fbb322b9b17d8d35d37.
//
// Solidity: event InflationReceived(uint256 amountReceivedWei)
func (_FdcHub *FdcHubFilterer) WatchInflationReceived(opts *bind.WatchOpts, sink chan<- *FdcHubInflationReceived) (event.Subscription, error) {

	logs, sub, err := _FdcHub.contract.WatchLogs(opts, "InflationReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcHubInflationReceived)
				if err := _FdcHub.contract.UnpackLog(event, "InflationReceived", log); err != nil {
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
func (_FdcHub *FdcHubFilterer) ParseInflationReceived(log types.Log) (*FdcHubInflationReceived, error) {
	event := new(FdcHubInflationReceived)
	if err := _FdcHub.contract.UnpackLog(event, "InflationReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcHubInflationRewardsOfferedIterator is returned from FilterInflationRewardsOffered and is used to iterate over the raw logs and unpacked data for InflationRewardsOffered events raised by the FdcHub contract.
type FdcHubInflationRewardsOfferedIterator struct {
	Event *FdcHubInflationRewardsOffered // Event containing the contract specifics and raw log

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
func (it *FdcHubInflationRewardsOfferedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcHubInflationRewardsOffered)
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
		it.Event = new(FdcHubInflationRewardsOffered)
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
func (it *FdcHubInflationRewardsOfferedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcHubInflationRewardsOfferedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcHubInflationRewardsOffered represents a InflationRewardsOffered event raised by the FdcHub contract.
type FdcHubInflationRewardsOffered struct {
	RewardEpochId     *big.Int
	FdcConfigurations []IFdcInflationConfigurationsFdcConfiguration
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterInflationRewardsOffered is a free log retrieval operation binding the contract event 0xedcf03eed469135e307ec8dc425dc2c49560d3014b724a532f6f468fcc975df8.
//
// Solidity: event InflationRewardsOffered(uint24 indexed rewardEpochId, (bytes32,bytes32,uint24,uint8,uint224)[] fdcConfigurations, uint256 amount)
func (_FdcHub *FdcHubFilterer) FilterInflationRewardsOffered(opts *bind.FilterOpts, rewardEpochId []*big.Int) (*FdcHubInflationRewardsOfferedIterator, error) {

	var rewardEpochIdRule []interface{}
	for _, rewardEpochIdItem := range rewardEpochId {
		rewardEpochIdRule = append(rewardEpochIdRule, rewardEpochIdItem)
	}

	logs, sub, err := _FdcHub.contract.FilterLogs(opts, "InflationRewardsOffered", rewardEpochIdRule)
	if err != nil {
		return nil, err
	}
	return &FdcHubInflationRewardsOfferedIterator{contract: _FdcHub.contract, event: "InflationRewardsOffered", logs: logs, sub: sub}, nil
}

// WatchInflationRewardsOffered is a free log subscription operation binding the contract event 0xedcf03eed469135e307ec8dc425dc2c49560d3014b724a532f6f468fcc975df8.
//
// Solidity: event InflationRewardsOffered(uint24 indexed rewardEpochId, (bytes32,bytes32,uint24,uint8,uint224)[] fdcConfigurations, uint256 amount)
func (_FdcHub *FdcHubFilterer) WatchInflationRewardsOffered(opts *bind.WatchOpts, sink chan<- *FdcHubInflationRewardsOffered, rewardEpochId []*big.Int) (event.Subscription, error) {

	var rewardEpochIdRule []interface{}
	for _, rewardEpochIdItem := range rewardEpochId {
		rewardEpochIdRule = append(rewardEpochIdRule, rewardEpochIdItem)
	}

	logs, sub, err := _FdcHub.contract.WatchLogs(opts, "InflationRewardsOffered", rewardEpochIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcHubInflationRewardsOffered)
				if err := _FdcHub.contract.UnpackLog(event, "InflationRewardsOffered", log); err != nil {
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

// ParseInflationRewardsOffered is a log parse operation binding the contract event 0xedcf03eed469135e307ec8dc425dc2c49560d3014b724a532f6f468fcc975df8.
//
// Solidity: event InflationRewardsOffered(uint24 indexed rewardEpochId, (bytes32,bytes32,uint24,uint8,uint224)[] fdcConfigurations, uint256 amount)
func (_FdcHub *FdcHubFilterer) ParseInflationRewardsOffered(log types.Log) (*FdcHubInflationRewardsOffered, error) {
	event := new(FdcHubInflationRewardsOffered)
	if err := _FdcHub.contract.UnpackLog(event, "InflationRewardsOffered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcHubRequestsOffsetSetIterator is returned from FilterRequestsOffsetSet and is used to iterate over the raw logs and unpacked data for RequestsOffsetSet events raised by the FdcHub contract.
type FdcHubRequestsOffsetSetIterator struct {
	Event *FdcHubRequestsOffsetSet // Event containing the contract specifics and raw log

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
func (it *FdcHubRequestsOffsetSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcHubRequestsOffsetSet)
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
		it.Event = new(FdcHubRequestsOffsetSet)
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
func (it *FdcHubRequestsOffsetSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcHubRequestsOffsetSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcHubRequestsOffsetSet represents a RequestsOffsetSet event raised by the FdcHub contract.
type FdcHubRequestsOffsetSet struct {
	RequestsOffsetSeconds uint8
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterRequestsOffsetSet is a free log retrieval operation binding the contract event 0x5d5d031078427b5a6de8e1f2973a0edacfca00cc17bb10fc99b4c726df1f1f4c.
//
// Solidity: event RequestsOffsetSet(uint8 requestsOffsetSeconds)
func (_FdcHub *FdcHubFilterer) FilterRequestsOffsetSet(opts *bind.FilterOpts) (*FdcHubRequestsOffsetSetIterator, error) {

	logs, sub, err := _FdcHub.contract.FilterLogs(opts, "RequestsOffsetSet")
	if err != nil {
		return nil, err
	}
	return &FdcHubRequestsOffsetSetIterator{contract: _FdcHub.contract, event: "RequestsOffsetSet", logs: logs, sub: sub}, nil
}

// WatchRequestsOffsetSet is a free log subscription operation binding the contract event 0x5d5d031078427b5a6de8e1f2973a0edacfca00cc17bb10fc99b4c726df1f1f4c.
//
// Solidity: event RequestsOffsetSet(uint8 requestsOffsetSeconds)
func (_FdcHub *FdcHubFilterer) WatchRequestsOffsetSet(opts *bind.WatchOpts, sink chan<- *FdcHubRequestsOffsetSet) (event.Subscription, error) {

	logs, sub, err := _FdcHub.contract.WatchLogs(opts, "RequestsOffsetSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcHubRequestsOffsetSet)
				if err := _FdcHub.contract.UnpackLog(event, "RequestsOffsetSet", log); err != nil {
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

// ParseRequestsOffsetSet is a log parse operation binding the contract event 0x5d5d031078427b5a6de8e1f2973a0edacfca00cc17bb10fc99b4c726df1f1f4c.
//
// Solidity: event RequestsOffsetSet(uint8 requestsOffsetSeconds)
func (_FdcHub *FdcHubFilterer) ParseRequestsOffsetSet(log types.Log) (*FdcHubRequestsOffsetSet, error) {
	event := new(FdcHubRequestsOffsetSet)
	if err := _FdcHub.contract.UnpackLog(event, "RequestsOffsetSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcHubTimelockedGovernanceCallCanceledIterator is returned from FilterTimelockedGovernanceCallCanceled and is used to iterate over the raw logs and unpacked data for TimelockedGovernanceCallCanceled events raised by the FdcHub contract.
type FdcHubTimelockedGovernanceCallCanceledIterator struct {
	Event *FdcHubTimelockedGovernanceCallCanceled // Event containing the contract specifics and raw log

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
func (it *FdcHubTimelockedGovernanceCallCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcHubTimelockedGovernanceCallCanceled)
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
		it.Event = new(FdcHubTimelockedGovernanceCallCanceled)
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
func (it *FdcHubTimelockedGovernanceCallCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcHubTimelockedGovernanceCallCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcHubTimelockedGovernanceCallCanceled represents a TimelockedGovernanceCallCanceled event raised by the FdcHub contract.
type FdcHubTimelockedGovernanceCallCanceled struct {
	Selector  [4]byte
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTimelockedGovernanceCallCanceled is a free log retrieval operation binding the contract event 0x7735b2391c38a81419c513e30ca578db7158eadd7101511b23e221c654d19cf8.
//
// Solidity: event TimelockedGovernanceCallCanceled(bytes4 selector, uint256 timestamp)
func (_FdcHub *FdcHubFilterer) FilterTimelockedGovernanceCallCanceled(opts *bind.FilterOpts) (*FdcHubTimelockedGovernanceCallCanceledIterator, error) {

	logs, sub, err := _FdcHub.contract.FilterLogs(opts, "TimelockedGovernanceCallCanceled")
	if err != nil {
		return nil, err
	}
	return &FdcHubTimelockedGovernanceCallCanceledIterator{contract: _FdcHub.contract, event: "TimelockedGovernanceCallCanceled", logs: logs, sub: sub}, nil
}

// WatchTimelockedGovernanceCallCanceled is a free log subscription operation binding the contract event 0x7735b2391c38a81419c513e30ca578db7158eadd7101511b23e221c654d19cf8.
//
// Solidity: event TimelockedGovernanceCallCanceled(bytes4 selector, uint256 timestamp)
func (_FdcHub *FdcHubFilterer) WatchTimelockedGovernanceCallCanceled(opts *bind.WatchOpts, sink chan<- *FdcHubTimelockedGovernanceCallCanceled) (event.Subscription, error) {

	logs, sub, err := _FdcHub.contract.WatchLogs(opts, "TimelockedGovernanceCallCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcHubTimelockedGovernanceCallCanceled)
				if err := _FdcHub.contract.UnpackLog(event, "TimelockedGovernanceCallCanceled", log); err != nil {
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
func (_FdcHub *FdcHubFilterer) ParseTimelockedGovernanceCallCanceled(log types.Log) (*FdcHubTimelockedGovernanceCallCanceled, error) {
	event := new(FdcHubTimelockedGovernanceCallCanceled)
	if err := _FdcHub.contract.UnpackLog(event, "TimelockedGovernanceCallCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FdcHubTimelockedGovernanceCallExecutedIterator is returned from FilterTimelockedGovernanceCallExecuted and is used to iterate over the raw logs and unpacked data for TimelockedGovernanceCallExecuted events raised by the FdcHub contract.
type FdcHubTimelockedGovernanceCallExecutedIterator struct {
	Event *FdcHubTimelockedGovernanceCallExecuted // Event containing the contract specifics and raw log

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
func (it *FdcHubTimelockedGovernanceCallExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FdcHubTimelockedGovernanceCallExecuted)
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
		it.Event = new(FdcHubTimelockedGovernanceCallExecuted)
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
func (it *FdcHubTimelockedGovernanceCallExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FdcHubTimelockedGovernanceCallExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FdcHubTimelockedGovernanceCallExecuted represents a TimelockedGovernanceCallExecuted event raised by the FdcHub contract.
type FdcHubTimelockedGovernanceCallExecuted struct {
	Selector  [4]byte
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTimelockedGovernanceCallExecuted is a free log retrieval operation binding the contract event 0xa7326b57fc9cfe267aaea5e7f0b01757154d265620a0585819416ee9ddd2c438.
//
// Solidity: event TimelockedGovernanceCallExecuted(bytes4 selector, uint256 timestamp)
func (_FdcHub *FdcHubFilterer) FilterTimelockedGovernanceCallExecuted(opts *bind.FilterOpts) (*FdcHubTimelockedGovernanceCallExecutedIterator, error) {

	logs, sub, err := _FdcHub.contract.FilterLogs(opts, "TimelockedGovernanceCallExecuted")
	if err != nil {
		return nil, err
	}
	return &FdcHubTimelockedGovernanceCallExecutedIterator{contract: _FdcHub.contract, event: "TimelockedGovernanceCallExecuted", logs: logs, sub: sub}, nil
}

// WatchTimelockedGovernanceCallExecuted is a free log subscription operation binding the contract event 0xa7326b57fc9cfe267aaea5e7f0b01757154d265620a0585819416ee9ddd2c438.
//
// Solidity: event TimelockedGovernanceCallExecuted(bytes4 selector, uint256 timestamp)
func (_FdcHub *FdcHubFilterer) WatchTimelockedGovernanceCallExecuted(opts *bind.WatchOpts, sink chan<- *FdcHubTimelockedGovernanceCallExecuted) (event.Subscription, error) {

	logs, sub, err := _FdcHub.contract.WatchLogs(opts, "TimelockedGovernanceCallExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FdcHubTimelockedGovernanceCallExecuted)
				if err := _FdcHub.contract.UnpackLog(event, "TimelockedGovernanceCallExecuted", log); err != nil {
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
func (_FdcHub *FdcHubFilterer) ParseTimelockedGovernanceCallExecuted(log types.Log) (*FdcHubTimelockedGovernanceCallExecuted, error) {
	event := new(FdcHubTimelockedGovernanceCallExecuted)
	if err := _FdcHub.contract.UnpackLog(event, "TimelockedGovernanceCallExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
