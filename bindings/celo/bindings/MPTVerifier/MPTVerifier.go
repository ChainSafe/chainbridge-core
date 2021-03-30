// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package MPTVerifier

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MPTVerifierABI is the input ABI used to generate the binding from.
const MPTVerifierABI = "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rootHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"mptPath\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rlpStack\",\"type\":\"bytes\"}],\"name\":\"validateMPTProof\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rootHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"mptPath\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rlpStack\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"valueHash\",\"type\":\"bytes32\"}],\"name\":\"validateMPTValueInclusion\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isIncluded\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true}]"

// MPTVerifierBin is the compiled bytecode used for deploying new contracts.
var MPTVerifierBin = "0x608060405234801561001057600080fd5b506110d9806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806393d70bf01461003b578063b941472d14610175575b600080fd5b6101006004803603606081101561005157600080fd5b81359190810190604081016020820135600160201b81111561007257600080fd5b82018360208201111561008457600080fd5b803590602001918460018302840111600160201b831117156100a557600080fd5b919390929091602081019035600160201b8111156100c257600080fd5b8201836020820111156100d457600080fd5b803590602001918460018302840111600160201b831117156100f557600080fd5b50909250905061024e565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561013a578181015183820152602001610122565b50505050905090810190601f1680156101675780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61023a6004803603608081101561018b57600080fd5b81359190810190604081016020820135600160201b8111156101ac57600080fd5b8201836020820111156101be57600080fd5b803590602001918460018302840111600160201b831117156101df57600080fd5b919390929091602081019035600160201b8111156101fc57600080fd5b82018360208201111561020e57600080fd5b803590602001918460018302840111600160201b8311171561022f57600080fd5b9193509150356102ce565b604080519115158252519081900360200190f35b60606102c48686868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050604080516020601f8a01819004810282018101909252888152925088915087908190840183828082843760009201919091525061035d92505050565b9695505050505050565b600060606103468888888080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050604080516020601f8c018190048102820181019092528a815292508a915089908190840183828082843760009201919091525061035d92505050565b805160209091012083149150509695505050505050565b606061037a848461037561037086610384565b6103aa565b61047b565b90505b9392505050565b61038c611089565b5060408051808201909152815181526020828101908201525b919050565b60606103b582610a86565b6103be57600080fd5b60006103c983610ac0565b905060608160405190808252806020026020018201604052801561040757816020015b6103f4611089565b8152602001906001900390816103ec5790505b50905060006104198560200151610b1c565b60208601510190506000805b848110156104705761043683610b7f565b915060405180604001604052808381526020018481525084828151811061045957fe5b602090810291909101015291810191600101610425565b509195945050505050565b60606000601081838061048c611089565b87516104db577f56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b4218a146104be57600080fd5b5050604080516000815260208101909152945061037d9350505050565b60005b8851811015610a7857856010146105125761050b8387815181106104fe57fe5b60200260200101516103aa565b92506105fb565b61052e89828151811061052157fe5b6020026020010151610c18565b9350801580156105445750835160208501208b14155b1561058b576040805162461bcd60e51b81526020600482015260126024820152710a4dedee840d0c2e6d040dad2e6dac2e8c6d60731b604482015290519081900360640190fd5b80158015906105a2575061059e84610c79565b8514155b156105e9576040805162461bcd60e51b815260206004820152601260248201527109cdec8ca40d0c2e6d040dad2e6dac2e8c6d60731b604482015290519081900360640190fd5b6105f88982815181106104fe57fe5b92505b82516002141561083a576000606061062e6106298660008151811061061c57fe5b6020026020010151610d2d565b610da6565b909250905060006106408a8e84610e84565b9050808a01995081518110156106cd5760018c51038410156106a9576040805162461bcd60e51b815260206004820152601860248201527f446976657267656e636520696e20746865206d6964646c650000000000000000604482015290519081900360640190fd5b60408051600080825260208201909252905b509a505050505050505050505061037d565b82156107715760018c5103841015610721576040805162461bcd60e51b81526020600482015260126024820152714c65616620696e20746865206d6964646c6560701b604482015290519081900360640190fd5b8c518a10156107405760408051600080825260208201909252906106bb565b8560018151811061074d57fe5b6020026020010151945061076085610d2d565b9a505050505050505050505061037d565b886010148015610784575060018c510384145b156107cd576040805162461bcd60e51b8152602060048201526014602482015273115e1d195b9cda5bdb881a5b881d1a1948195b9960621b604482015290519081900360640190fd5b6107ea866001815181106107dd57fe5b6020026020010151610a86565b61080c576107fe8660018151811061061c57fe5b805190602001209750610832565b61081c8660018151811061052157fe5b8051602090910120600199509750600019909301925b505050610a70565b825160111415610a2d57895187146109b25760008a888151811061085a57fe5b01602001516001989098019760f81c9050601081106108b1576040805162461bcd60e51b815260206004820152600e60248201526d496e76616c6964206e6962626c6560901b604482015290519081900360640190fd5b6108d0848260ff16815181106108c357fe5b6020026020010151610ef1565b1561094c5760018a5103821461092d576040805162461bcd60e51b815260206004820152601860248201527f456d707479206c65616620696e20746865206d6964646c650000000000000000604482015290519081900360640190fd5b5050604080516000815260208101909152965061037d95505050505050565b61095e848260ff16815181106107dd57fe5b61098257610974848260ff168151811061061c57fe5b8051906020012095506109ac565b610994848260ff168151811061052157fe5b805160209091012060ff821697509550600019909101905b50610a28565b60018951038114610a0a576040805162461bcd60e51b815260206004820152601a60248201527f4b657920636f6e73756d656420696e20746865206d6964646c65000000000000604482015290519081900360640190fd5b610a1a8360108151811061061c57fe5b97505050505050505061037d565b610a70565b6040805162461bcd60e51b8152602060048201526013602482015272092dcecc2d8d2c840dcdec8ca40d8cadccee8d606b1b604482015290519081900360640190fd5b6001016104de565b505050505050509392505050565b8051600090610a97575060006103a5565b6020820151805160001a9060c0821015610ab6576000925050506103a5565b5060019392505050565b8051600090610ad1575060006103a5565b60008090506000610ae58460200151610b1c565b602085015185519181019250015b80821015610b1357610b0482610b7f565b60019093019290910190610af3565b50909392505050565b8051600090811a6080811015610b365760009150506103a5565b60b8811080610b51575060c08110801590610b51575060f881105b15610b605760019150506103a5565b60c0811015610b745760b5190190506103a5565b60f5190190506103a5565b80516000908190811a6080811015610b9a5760019150610c11565b60b8811015610baf57607e1981019150610c11565b60c0811015610bdc5760b78103600185019450806020036101000a85510460018201810193505050610c11565b60f8811015610bf15760be1981019150610c11565b60f78103600185019450806020036101000a855104600182018101935050505b5092915050565b60608082600001516040519080825280601f01601f191660200182016040528015610c4a576020820181803683370190505b509050805160001415610c5e5790506103a5565b6000816020019050610c118460200151828660000151610f14565b6000602082511015610c925750805160208201206103a5565b816040516020018082805190602001908083835b60208310610cc55780518252601f199092019160209182019101610ca6565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405160208183030381529060405280519060200120604051602001808281526020019150506040516020818303038152906040528051906020012090506103a5565b8051606090610d3b57600080fd5b6000610d4a8360200151610b1c565b83516040805191839003808352601f19601f8201168301602001909152919250606090828015610d81576020820181803683370190505b5090506000816020019050610d9d848760200151018285610f14565b50949350505050565b600060606000835111610db857600080fd5b6000600484600081518110610dc957fe5b60209101015160f81c901c600f169050600081610dec5750600092506002610e6e565b8160011415610e015750600092506001610e6e565b8160021415610e165750600192506002610e6e565b8160031415610e2a57506001925082610e6e565b6040805162461bcd60e51b8152602060048201526014602482015273496e76616c6964206669727374206e6962626c6560601b604482015290519081900360640190fd5b83610e798683610f5f565b935093505050915091565b6000805b8351858201108015610e9a5750825181105b1561037a57828181518110610eab57fe5b602001015160f81c60f81b6001600160f81b0319168486830181518110610ece57fe5b01602001516001600160f81b03191614610ee957905061037d565b600101610e88565b8051600090600114610f05575060006103a5565b50602001515160001a60801490565b80610f1e57610f5a565b5b60208110610f3e578251825260209283019290910190601f1901610f1f565b8251825160208390036101000a60001901801990921691161782525b505050565b60606000835111610f6f57600080fd5b825160020280831115610f8157600080fd5b8281039050806040519080825280601f01601f191660200182016040528015610fb1576020820181803683370190505b5091506000835b828501811015611076576002810661101c576004866002830481518110610fdb57fe5b602001015160f81c60f81b60f81c60ff16901c600f1660f81b84838151811061100057fe5b60200101906001600160f81b031916908160001a90535061106a565b600086600283048151811061102d57fe5b602001015160f81c60f81b60f81c60ff16901c600f1660f81b84838151811061105257fe5b60200101906001600160f81b031916908160001a9053505b60019182019101610fb8565b508251811461108157fe5b505092915050565b60405180604001604052806000815260200160008152509056fea26469706673582212203fab88038a1211aaffd7c7ddcb12cb077e39b2932ac7d8920f9c56fb8024baac64736f6c63430006040033"

// DeployMPTVerifier deploys a new Ethereum contract, binding an instance of MPTVerifier to it.
func DeployMPTVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MPTVerifier, error) {
	parsed, err := abi.JSON(strings.NewReader(MPTVerifierABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MPTVerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MPTVerifier{MPTVerifierCaller: MPTVerifierCaller{contract: contract}, MPTVerifierTransactor: MPTVerifierTransactor{contract: contract}, MPTVerifierFilterer: MPTVerifierFilterer{contract: contract}}, nil
}

// MPTVerifier is an auto generated Go binding around an Ethereum contract.
type MPTVerifier struct {
	MPTVerifierCaller     // Read-only binding to the contract
	MPTVerifierTransactor // Write-only binding to the contract
	MPTVerifierFilterer   // Log filterer for contract events
}

// MPTVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type MPTVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MPTVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MPTVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MPTVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MPTVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MPTVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MPTVerifierSession struct {
	Contract     *MPTVerifier      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MPTVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MPTVerifierCallerSession struct {
	Contract *MPTVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// MPTVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MPTVerifierTransactorSession struct {
	Contract     *MPTVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// MPTVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type MPTVerifierRaw struct {
	Contract *MPTVerifier // Generic contract binding to access the raw methods on
}

// MPTVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MPTVerifierCallerRaw struct {
	Contract *MPTVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// MPTVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MPTVerifierTransactorRaw struct {
	Contract *MPTVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMPTVerifier creates a new instance of MPTVerifier, bound to a specific deployed contract.
func NewMPTVerifier(address common.Address, backend bind.ContractBackend) (*MPTVerifier, error) {
	contract, err := bindMPTVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MPTVerifier{MPTVerifierCaller: MPTVerifierCaller{contract: contract}, MPTVerifierTransactor: MPTVerifierTransactor{contract: contract}, MPTVerifierFilterer: MPTVerifierFilterer{contract: contract}}, nil
}

// NewMPTVerifierCaller creates a new read-only instance of MPTVerifier, bound to a specific deployed contract.
func NewMPTVerifierCaller(address common.Address, caller bind.ContractCaller) (*MPTVerifierCaller, error) {
	contract, err := bindMPTVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MPTVerifierCaller{contract: contract}, nil
}

// NewMPTVerifierTransactor creates a new write-only instance of MPTVerifier, bound to a specific deployed contract.
func NewMPTVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*MPTVerifierTransactor, error) {
	contract, err := bindMPTVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MPTVerifierTransactor{contract: contract}, nil
}

// NewMPTVerifierFilterer creates a new log filterer instance of MPTVerifier, bound to a specific deployed contract.
func NewMPTVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*MPTVerifierFilterer, error) {
	contract, err := bindMPTVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MPTVerifierFilterer{contract: contract}, nil
}

// bindMPTVerifier binds a generic wrapper to an already deployed contract.
func bindMPTVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MPTVerifierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// ParseMPTVerifierABI parses the ABI
func ParseMPTVerifierABI() (*abi.ABI, error) {
	parsed, err := abi.JSON(strings.NewReader(MPTVerifierABI))
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MPTVerifier *MPTVerifierRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MPTVerifier.Contract.MPTVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MPTVerifier *MPTVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MPTVerifier.Contract.MPTVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MPTVerifier *MPTVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MPTVerifier.Contract.MPTVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MPTVerifier *MPTVerifierCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MPTVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MPTVerifier *MPTVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MPTVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MPTVerifier *MPTVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MPTVerifier.Contract.contract.Transact(opts, method, params...)
}

// ValidateMPTProof is a free data retrieval call binding the contract method 0x93d70bf0.
//
// Solidity: function validateMPTProof(bytes32 rootHash, bytes mptPath, bytes rlpStack) constant returns(bytes value)
func (_MPTVerifier *MPTVerifierCaller) ValidateMPTProof(opts *bind.CallOpts, rootHash [32]byte, mptPath []byte, rlpStack []byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _MPTVerifier.contract.Call(opts, out, "validateMPTProof", rootHash, mptPath, rlpStack)
	return *ret0, err
}

// ValidateMPTProof is a free data retrieval call binding the contract method 0x93d70bf0.
//
// Solidity: function validateMPTProof(bytes32 rootHash, bytes mptPath, bytes rlpStack) constant returns(bytes value)
func (_MPTVerifier *MPTVerifierSession) ValidateMPTProof(rootHash [32]byte, mptPath []byte, rlpStack []byte) ([]byte, error) {
	return _MPTVerifier.Contract.ValidateMPTProof(&_MPTVerifier.CallOpts, rootHash, mptPath, rlpStack)
}

// ValidateMPTProof is a free data retrieval call binding the contract method 0x93d70bf0.
//
// Solidity: function validateMPTProof(bytes32 rootHash, bytes mptPath, bytes rlpStack) constant returns(bytes value)
func (_MPTVerifier *MPTVerifierCallerSession) ValidateMPTProof(rootHash [32]byte, mptPath []byte, rlpStack []byte) ([]byte, error) {
	return _MPTVerifier.Contract.ValidateMPTProof(&_MPTVerifier.CallOpts, rootHash, mptPath, rlpStack)
}

// ValidateMPTValueInclusion is a free data retrieval call binding the contract method 0xb941472d.
//
// Solidity: function validateMPTValueInclusion(bytes32 rootHash, bytes mptPath, bytes rlpStack, bytes32 valueHash) constant returns(bool isIncluded)
func (_MPTVerifier *MPTVerifierCaller) ValidateMPTValueInclusion(opts *bind.CallOpts, rootHash [32]byte, mptPath []byte, rlpStack []byte, valueHash [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _MPTVerifier.contract.Call(opts, out, "validateMPTValueInclusion", rootHash, mptPath, rlpStack, valueHash)
	return *ret0, err
}

// ValidateMPTValueInclusion is a free data retrieval call binding the contract method 0xb941472d.
//
// Solidity: function validateMPTValueInclusion(bytes32 rootHash, bytes mptPath, bytes rlpStack, bytes32 valueHash) constant returns(bool isIncluded)
func (_MPTVerifier *MPTVerifierSession) ValidateMPTValueInclusion(rootHash [32]byte, mptPath []byte, rlpStack []byte, valueHash [32]byte) (bool, error) {
	return _MPTVerifier.Contract.ValidateMPTValueInclusion(&_MPTVerifier.CallOpts, rootHash, mptPath, rlpStack, valueHash)
}

// ValidateMPTValueInclusion is a free data retrieval call binding the contract method 0xb941472d.
//
// Solidity: function validateMPTValueInclusion(bytes32 rootHash, bytes mptPath, bytes rlpStack, bytes32 valueHash) constant returns(bool isIncluded)
func (_MPTVerifier *MPTVerifierCallerSession) ValidateMPTValueInclusion(rootHash [32]byte, mptPath []byte, rlpStack []byte, valueHash [32]byte) (bool, error) {
	return _MPTVerifier.Contract.ValidateMPTValueInclusion(&_MPTVerifier.CallOpts, rootHash, mptPath, rlpStack, valueHash)
}

// TryParseLog attempts to parse a log. Returns the parsed log, evenName and whether it was succesfull
func (_MPTVerifier *MPTVerifierFilterer) TryParseLog(log types.Log) (eventName string, event interface{}, ok bool, err error) {
	eventName, ok, err = _MPTVerifier.contract.LogEventName(log)
	if err != nil || !ok {
		return "", nil, false, err
	}

	switch eventName {
	}
	if err != nil {
		return "", nil, false, err
	}

	return eventName, event, ok, nil
}
