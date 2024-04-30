// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package hub

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

// HubMetaData contains all meta data concerning the Hub contract.
var HubMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"AttestationRequest\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"minimalFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"requestAttestation\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// HubABI is the input ABI used to generate the binding from.
// Deprecated: Use HubMetaData.ABI instead.
var HubABI = HubMetaData.ABI

// Hub is an auto generated Go binding around an Ethereum contract.
type Hub struct {
	HubCaller     // Read-only binding to the contract
	HubTransactor // Write-only binding to the contract
	HubFilterer   // Log filterer for contract events
}

// HubCaller is an auto generated read-only Go binding around an Ethereum contract.
type HubCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HubTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HubTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HubFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HubFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HubSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HubSession struct {
	Contract     *Hub              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HubCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HubCallerSession struct {
	Contract *HubCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// HubTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HubTransactorSession struct {
	Contract     *HubTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HubRaw is an auto generated low-level Go binding around an Ethereum contract.
type HubRaw struct {
	Contract *Hub // Generic contract binding to access the raw methods on
}

// HubCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HubCallerRaw struct {
	Contract *HubCaller // Generic read-only contract binding to access the raw methods on
}

// HubTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HubTransactorRaw struct {
	Contract *HubTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHub creates a new instance of Hub, bound to a specific deployed contract.
func NewHub(address common.Address, backend bind.ContractBackend) (*Hub, error) {
	contract, err := bindHub(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Hub{HubCaller: HubCaller{contract: contract}, HubTransactor: HubTransactor{contract: contract}, HubFilterer: HubFilterer{contract: contract}}, nil
}

// NewHubCaller creates a new read-only instance of Hub, bound to a specific deployed contract.
func NewHubCaller(address common.Address, caller bind.ContractCaller) (*HubCaller, error) {
	contract, err := bindHub(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HubCaller{contract: contract}, nil
}

// NewHubTransactor creates a new write-only instance of Hub, bound to a specific deployed contract.
func NewHubTransactor(address common.Address, transactor bind.ContractTransactor) (*HubTransactor, error) {
	contract, err := bindHub(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HubTransactor{contract: contract}, nil
}

// NewHubFilterer creates a new log filterer instance of Hub, bound to a specific deployed contract.
func NewHubFilterer(address common.Address, filterer bind.ContractFilterer) (*HubFilterer, error) {
	contract, err := bindHub(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HubFilterer{contract: contract}, nil
}

// bindHub binds a generic wrapper to an already deployed contract.
func bindHub(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HubMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Hub *HubRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Hub.Contract.HubCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Hub *HubRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Hub.Contract.HubTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Hub *HubRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Hub.Contract.HubTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Hub *HubCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Hub.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Hub *HubTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Hub.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Hub *HubTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Hub.Contract.contract.Transact(opts, method, params...)
}

// MinimalFee is a free data retrieval call binding the contract method 0xf61549da.
//
// Solidity: function minimalFee() view returns(uint256)
func (_Hub *HubCaller) MinimalFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Hub.contract.Call(opts, &out, "minimalFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinimalFee is a free data retrieval call binding the contract method 0xf61549da.
//
// Solidity: function minimalFee() view returns(uint256)
func (_Hub *HubSession) MinimalFee() (*big.Int, error) {
	return _Hub.Contract.MinimalFee(&_Hub.CallOpts)
}

// MinimalFee is a free data retrieval call binding the contract method 0xf61549da.
//
// Solidity: function minimalFee() view returns(uint256)
func (_Hub *HubCallerSession) MinimalFee() (*big.Int, error) {
	return _Hub.Contract.MinimalFee(&_Hub.CallOpts)
}

// RequestAttestation is a paid mutator transaction binding the contract method 0x6238f354.
//
// Solidity: function requestAttestation(bytes _data) payable returns()
func (_Hub *HubTransactor) RequestAttestation(opts *bind.TransactOpts, _data []byte) (*types.Transaction, error) {
	return _Hub.contract.Transact(opts, "requestAttestation", _data)
}

// RequestAttestation is a paid mutator transaction binding the contract method 0x6238f354.
//
// Solidity: function requestAttestation(bytes _data) payable returns()
func (_Hub *HubSession) RequestAttestation(_data []byte) (*types.Transaction, error) {
	return _Hub.Contract.RequestAttestation(&_Hub.TransactOpts, _data)
}

// RequestAttestation is a paid mutator transaction binding the contract method 0x6238f354.
//
// Solidity: function requestAttestation(bytes _data) payable returns()
func (_Hub *HubTransactorSession) RequestAttestation(_data []byte) (*types.Transaction, error) {
	return _Hub.Contract.RequestAttestation(&_Hub.TransactOpts, _data)
}

// HubAttestationRequestIterator is returned from FilterAttestationRequest and is used to iterate over the raw logs and unpacked data for AttestationRequest events raised by the Hub contract.
type HubAttestationRequestIterator struct {
	Event *HubAttestationRequest // Event containing the contract specifics and raw log

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
func (it *HubAttestationRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HubAttestationRequest)
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
		it.Event = new(HubAttestationRequest)
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
func (it *HubAttestationRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HubAttestationRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HubAttestationRequest represents a AttestationRequest event raised by the Hub contract.
type HubAttestationRequest struct {
	Data []byte
	Fee  *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAttestationRequest is a free log retrieval operation binding the contract event 0x251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9.
//
// Solidity: event AttestationRequest(bytes data, uint256 fee)
func (_Hub *HubFilterer) FilterAttestationRequest(opts *bind.FilterOpts) (*HubAttestationRequestIterator, error) {

	logs, sub, err := _Hub.contract.FilterLogs(opts, "AttestationRequest")
	if err != nil {
		return nil, err
	}
	return &HubAttestationRequestIterator{contract: _Hub.contract, event: "AttestationRequest", logs: logs, sub: sub}, nil
}

// WatchAttestationRequest is a free log subscription operation binding the contract event 0x251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9.
//
// Solidity: event AttestationRequest(bytes data, uint256 fee)
func (_Hub *HubFilterer) WatchAttestationRequest(opts *bind.WatchOpts, sink chan<- *HubAttestationRequest) (event.Subscription, error) {

	logs, sub, err := _Hub.contract.WatchLogs(opts, "AttestationRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HubAttestationRequest)
				if err := _Hub.contract.UnpackLog(event, "AttestationRequest", log); err != nil {
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
func (_Hub *HubFilterer) ParseAttestationRequest(log types.Log) (*HubAttestationRequest, error) {
	event := new(HubAttestationRequest)
	if err := _Hub.contract.UnpackLog(event, "AttestationRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
