// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package RLPReaderHelper

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

// RLPReaderHelperABI is the input ABI used to generate the binding from.
const RLPReaderHelperABI = "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"isList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"itemLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"rlpLen\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"payloadLen\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"numItems\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"toRlpBytes\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"toBytes\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"toUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"toUintStrict\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"toAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"toBoolean\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"bytesToString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"toIterator\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"nestedIteration\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"rlpHeader\",\"type\":\"bytes\"}],\"name\":\"toBlockHeader\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"parentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sha3Uncles\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"transactionsRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"receiptsRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"difficulty\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"customDestructure\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"customNestedDestructure\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"item\",\"type\":\"bytes\"}],\"name\":\"customNestedToRlpBytes\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true}]"

// RLPReaderHelperBin is the compiled bytecode used for deploying new contracts.
var RLPReaderHelperBin = "0x608060405234801561001057600080fd5b50611773806100206000396000f3fe608060405234801561001057600080fd5b50600436106101165760003560e01c806374c4c494116100a2578063aac1954511610071578063aac1954514610a71578063d03ca5a314610b15578063dd208ea114610bbb578063e7617bf814610c5f578063e79fbcec1461069957610116565b806374c4c494146107e157806376193f8c146108855780638109672f146109295780638529d576146109cd57610116565b806320f91849116100e957806320f918491461040a5780632d888869146104c057806333c82d181461058057806339614e4f146106995780636393b0fb1461073d57610116565b806304ba34271461011b578063059a97fb146101d357806315943cb61461029f57806319921bd014610366575b600080fd5b6101bf6004803603602081101561013157600080fd5b810190602081018135600160201b81111561014b57600080fd5b82018360208201111561015d57600080fd5b803590602001918460018302840111600160201b8311171561017e57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610d59945050505050565b604080519115158252519081900360200190f35b610277600480360360208110156101e957600080fd5b810190602081018135600160201b81111561020357600080fd5b82018360208201111561021557600080fd5b803590602001918460018302840111600160201b8311171561023657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610d80945050505050565b604080516001600160a01b039094168452911515602084015282820152519081900360600190f35b610343600480360360208110156102b557600080fd5b810190602081018135600160201b8111156102cf57600080fd5b8201836020820111156102e157600080fd5b803590602001918460018302840111600160201b8311171561030257600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610dff945050505050565b604080516001600160a01b03909316835260208301919091528051918290030190f35b6101bf6004803603602081101561037c57600080fd5b810190602081018135600160201b81111561039657600080fd5b8201836020820111156103a857600080fd5b803590602001918460018302840111600160201b831117156103c957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610e5b945050505050565b6104ae6004803603602081101561042057600080fd5b810190602081018135600160201b81111561043a57600080fd5b82018360208201111561044c57600080fd5b803590602001918460018302840111600160201b8311171561046d57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610e79945050505050565b60408051918252519081900360200190f35b610564600480360360208110156104d657600080fd5b810190602081018135600160201b8111156104f057600080fd5b82018360208201111561050257600080fd5b803590602001918460018302840111600160201b8311171561052357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610e97945050505050565b604080516001600160a01b039092168252519081900360200190f35b6106246004803603602081101561059657600080fd5b810190602081018135600160201b8111156105b057600080fd5b8201836020820111156105c257600080fd5b803590602001918460018302840111600160201b831117156105e357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610eb5945050505050565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561065e578181015183820152602001610646565b50505050905090810190601f16801561068b5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610624600480360360208110156106af57600080fd5b810190602081018135600160201b8111156106c957600080fd5b8201836020820111156106db57600080fd5b803590602001918460018302840111600160201b831117156106fc57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610ed3945050505050565b6104ae6004803603602081101561075357600080fd5b810190602081018135600160201b81111561076d57600080fd5b82018360208201111561077f57600080fd5b803590602001918460018302840111600160201b831117156107a057600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610ef1945050505050565b610624600480360360208110156107f757600080fd5b810190602081018135600160201b81111561081157600080fd5b82018360208201111561082357600080fd5b803590602001918460018302840111600160201b8311171561084457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610f0f945050505050565b6106246004803603602081101561089b57600080fd5b810190602081018135600160201b8111156108b557600080fd5b8201836020820111156108c757600080fd5b803590602001918460018302840111600160201b831117156108e857600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610f3d945050505050565b6104ae6004803603602081101561093f57600080fd5b810190602081018135600160201b81111561095957600080fd5b82018360208201111561096b57600080fd5b803590602001918460018302840111600160201b8311171561098c57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610fb3945050505050565b6104ae600480360360208110156109e357600080fd5b810190602081018135600160201b8111156109fd57600080fd5b820183602082011115610a0f57600080fd5b803590602001918460018302840111600160201b83111715610a3057600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610fc2945050505050565b6104ae60048036036020811015610a8757600080fd5b810190602081018135600160201b811115610aa157600080fd5b820183602082011115610ab357600080fd5b803590602001918460018302840111600160201b83111715610ad457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610fe0945050505050565b610bb960048036036020811015610b2b57600080fd5b810190602081018135600160201b811115610b4557600080fd5b820183602082011115610b5757600080fd5b803590602001918460018302840111600160201b83111715610b7857600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610ffe945050505050565b005b6104ae60048036036020811015610bd157600080fd5b810190602081018135600160201b811115610beb57600080fd5b820183602082011115610bfd57600080fd5b803590602001918460018302840111600160201b83111715610c1e57600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955061100e945050505050565b610d0360048036036020811015610c7557600080fd5b810190602081018135600160201b811115610c8f57600080fd5b820183602082011115610ca157600080fd5b803590602001918460018302840111600160201b83111715610cc257600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611026945050505050565b604080519b8c5260208c019a909a528a8a019890985260608a0196909652608089019490945260a088019290925260c087015260e086015261010085015261012084015261014083015251908190036101600190f35b6000610d63611703565b610d6c836111b7565b9050610d77816111dc565b9150505b919050565b60008060006060610d98610d93866111b7565b61120f565b9050610db781600081518110610daa57fe5b60200260200101516112e0565b610dd482600181518110610dc757fe5b60200260200101516111dc565b610df183600281518110610de457fe5b6020026020010151611300565b935093509350509193909250565b6000806060610e10610d93856111b7565b9050610e2f81600081518110610e2257fe5b602002602001015161120f565b9050610e4181600081518110610daa57fe5b610e5182600181518110610de457fe5b9250925050915091565b6000610e65611703565b610e6e836111b7565b9050610d7781611360565b6000610e83611703565b610e8c836111b7565b9050610d778161139a565b6000610ea1611703565b610eaa836111b7565b9050610d77816112e0565b6060610ebf611703565b610ec8836111b7565b9050610d77816113b7565b6060610edd611703565b610ee6836111b7565b9050610d778161141f565b6000610efb611703565b610f04836111b7565b9050610d778161148f565b606080610f1e610d93846111b7565b9050610d7781600081518110610f3057fe5b60200260200101516113b7565b6060610f4761171d565b610f58610f53846111b7565b6114a8565b9050610f6261171d565b610f6e610f53836114f2565b90506060610f83610f7e836114f2565b61141f565b9050610f8e83611540565b15610f9857600080fd5b610fa182611540565b15610fab57600080fd5b949350505050565b600060208201610d7781611561565b6000610fcc611703565b610fd5836111b7565b9050610d7781611300565b6000610fea611703565b610ff3836111b7565b9050610d77816115f5565b61100a610f53826111b7565b5050565b6000606061101e610d93846111b7565b519392505050565b600080600080600080600080600080600061103f61171d565b61104b610f538e6111b7565b905060005b61105982611540565b156111a6578061107b5761107461106f836114f2565b611300565b9c5061119e565b80600114156110975761109061106f836114f2565b9b5061119e565b80600314156110b3576110ac61106f836114f2565b9a5061119e565b80600414156110cf576110c861106f836114f2565b995061119e565b80600514156110eb576110e461106f836114f2565b985061119e565b80600714156111075761110061106f836114f2565b975061119e565b80600814156111235761111c61106f836114f2565b965061119e565b806009141561113f5761113861106f836114f2565b955061119e565b80600a141561115b5761115461106f836114f2565b945061119e565b80600b14156111775761117061106f836114f2565b935061119e565b80600e14156111935761118c61106f836114f2565b925061119e565b61119c826114f2565b505b600101611050565b505091939597999b90929496989a50565b6111bf611703565b506040805180820190915281518152602082810190820152919050565b80516000906001146111ed57600080fd5b6020820151805160001a908115611205576001610fab565b6000949350505050565b606061121a82611360565b61122357600080fd5b600061122e836115f9565b905060608160405190808252806020026020018201604052801561126c57816020015b611259611703565b8152602001906001900390816112515790505b509050600061127e8560200151611655565b60208601510190506000805b848110156112d55761129b83611561565b91506040518060400160405280838152602001848152508482815181106112be57fe5b60209081029190910101529181019160010161128a565b509195945050505050565b80516000906015146112f157600080fd5b6112fa82611300565b92915050565b80516000901580159061131557508151602110155b61131e57600080fd5b600061132d8360200151611655565b8351602080860151830180519394509184900392919083101561135757826020036101000a820491505b50949350505050565b805160009061137157506000610d7b565b6020820151805160001a9060c082101561139057600092505050610d7b565b5060019392505050565b80516000906021146113ab57600080fd5b50602001516001015190565b60608082600001516040519080825280601f01601f1916602001820160405280156113e9576020820181803683370190505b5090508051600014156113fd579050610d7b565b600081602001905061141884602001518286600001516116b8565b5092915050565b805160609061142d57600080fd5b600061143c8360200151611655565b83516040805191839003808352601f19601f8201168301602001909152919250606090828015611473576020820181803683370190505b50905060008160200190506113578487602001510182856116b8565b600061149e8260200151611655565b8251039050919050565b6114b061171d565b6114b982611360565b6114c257600080fd5b60006114d18360200151611655565b60208085015160408051808201909152868152920190820152915050919050565b6114fa611703565b61150382611540565b61150c57600080fd5b6020820151600061151c82611561565b80830160209586015260408051808201909152908152938401919091525090919050565b600061154a611703565b505080518051602091820151919092015191011190565b80516000908190811a608081101561157c5760019150611418565b60b881101561159157607e1981019150611418565b60c08110156115be5760b78103600185019450806020036101000a85510460018201810193505050611418565b60f88110156115d35760be1981019150611418565b60019390930151602084900360f7016101000a900490920160f5190192915050565b5190565b805160009061160a57506000610d7b565b6000809050600061161e8460200151611655565b602085015185519181019250015b8082101561164c5761163d82611561565b6001909301929091019061162c565b50909392505050565b8051600090811a608081101561166f576000915050610d7b565b60b881108061168a575060c0811080159061168a575060f881105b15611699576001915050610d7b565b60c08110156116ad5760b519019050610d7b565b60f519019050610d7b565b806116c2576116fe565b5b602081106116e2578251825260209283019290910190601f19016116c3565b8251825160208390036101000a60001901801990921691161782525b505050565b604051806040016040528060008152602001600081525090565b6040518060400160405280611730611703565b815260200160008152509056fea26469706673582212203c9ca2c4d246586948efe90c8e898253fada78d4fdea0da01e52864a62512f6664736f6c63430006040033"

// DeployRLPReaderHelper deploys a new Ethereum contract, binding an instance of RLPReaderHelper to it.
func DeployRLPReaderHelper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RLPReaderHelper, error) {
	parsed, err := abi.JSON(strings.NewReader(RLPReaderHelperABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RLPReaderHelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RLPReaderHelper{RLPReaderHelperCaller: RLPReaderHelperCaller{contract: contract}, RLPReaderHelperTransactor: RLPReaderHelperTransactor{contract: contract}, RLPReaderHelperFilterer: RLPReaderHelperFilterer{contract: contract}}, nil
}

// RLPReaderHelper is an auto generated Go binding around an Ethereum contract.
type RLPReaderHelper struct {
	RLPReaderHelperCaller     // Read-only binding to the contract
	RLPReaderHelperTransactor // Write-only binding to the contract
	RLPReaderHelperFilterer   // Log filterer for contract events
}

// RLPReaderHelperCaller is an auto generated read-only Go binding around an Ethereum contract.
type RLPReaderHelperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RLPReaderHelperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RLPReaderHelperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RLPReaderHelperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RLPReaderHelperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RLPReaderHelperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RLPReaderHelperSession struct {
	Contract     *RLPReaderHelper  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RLPReaderHelperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RLPReaderHelperCallerSession struct {
	Contract *RLPReaderHelperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// RLPReaderHelperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RLPReaderHelperTransactorSession struct {
	Contract     *RLPReaderHelperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// RLPReaderHelperRaw is an auto generated low-level Go binding around an Ethereum contract.
type RLPReaderHelperRaw struct {
	Contract *RLPReaderHelper // Generic contract binding to access the raw methods on
}

// RLPReaderHelperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RLPReaderHelperCallerRaw struct {
	Contract *RLPReaderHelperCaller // Generic read-only contract binding to access the raw methods on
}

// RLPReaderHelperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RLPReaderHelperTransactorRaw struct {
	Contract *RLPReaderHelperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRLPReaderHelper creates a new instance of RLPReaderHelper, bound to a specific deployed contract.
func NewRLPReaderHelper(address common.Address, backend bind.ContractBackend) (*RLPReaderHelper, error) {
	contract, err := bindRLPReaderHelper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RLPReaderHelper{RLPReaderHelperCaller: RLPReaderHelperCaller{contract: contract}, RLPReaderHelperTransactor: RLPReaderHelperTransactor{contract: contract}, RLPReaderHelperFilterer: RLPReaderHelperFilterer{contract: contract}}, nil
}

// NewRLPReaderHelperCaller creates a new read-only instance of RLPReaderHelper, bound to a specific deployed contract.
func NewRLPReaderHelperCaller(address common.Address, caller bind.ContractCaller) (*RLPReaderHelperCaller, error) {
	contract, err := bindRLPReaderHelper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RLPReaderHelperCaller{contract: contract}, nil
}

// NewRLPReaderHelperTransactor creates a new write-only instance of RLPReaderHelper, bound to a specific deployed contract.
func NewRLPReaderHelperTransactor(address common.Address, transactor bind.ContractTransactor) (*RLPReaderHelperTransactor, error) {
	contract, err := bindRLPReaderHelper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RLPReaderHelperTransactor{contract: contract}, nil
}

// NewRLPReaderHelperFilterer creates a new log filterer instance of RLPReaderHelper, bound to a specific deployed contract.
func NewRLPReaderHelperFilterer(address common.Address, filterer bind.ContractFilterer) (*RLPReaderHelperFilterer, error) {
	contract, err := bindRLPReaderHelper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RLPReaderHelperFilterer{contract: contract}, nil
}

// bindRLPReaderHelper binds a generic wrapper to an already deployed contract.
func bindRLPReaderHelper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RLPReaderHelperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// ParseRLPReaderHelperABI parses the ABI
func ParseRLPReaderHelperABI() (*abi.ABI, error) {
	parsed, err := abi.JSON(strings.NewReader(RLPReaderHelperABI))
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RLPReaderHelper *RLPReaderHelperRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RLPReaderHelper.Contract.RLPReaderHelperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RLPReaderHelper *RLPReaderHelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RLPReaderHelper.Contract.RLPReaderHelperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RLPReaderHelper *RLPReaderHelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RLPReaderHelper.Contract.RLPReaderHelperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RLPReaderHelper *RLPReaderHelperCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RLPReaderHelper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RLPReaderHelper *RLPReaderHelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RLPReaderHelper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RLPReaderHelper *RLPReaderHelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RLPReaderHelper.Contract.contract.Transact(opts, method, params...)
}

// BytesToString is a free data retrieval call binding the contract method 0x39614e4f.
//
// Solidity: function bytesToString(bytes item) constant returns(string)
func (_RLPReaderHelper *RLPReaderHelperCaller) BytesToString(opts *bind.CallOpts, item []byte) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "bytesToString", item)
	return *ret0, err
}

// BytesToString is a free data retrieval call binding the contract method 0x39614e4f.
//
// Solidity: function bytesToString(bytes item) constant returns(string)
func (_RLPReaderHelper *RLPReaderHelperSession) BytesToString(item []byte) (string, error) {
	return _RLPReaderHelper.Contract.BytesToString(&_RLPReaderHelper.CallOpts, item)
}

// BytesToString is a free data retrieval call binding the contract method 0x39614e4f.
//
// Solidity: function bytesToString(bytes item) constant returns(string)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) BytesToString(item []byte) (string, error) {
	return _RLPReaderHelper.Contract.BytesToString(&_RLPReaderHelper.CallOpts, item)
}

// CustomDestructure is a free data retrieval call binding the contract method 0x059a97fb.
//
// Solidity: function customDestructure(bytes item) constant returns(address, bool, uint256)
func (_RLPReaderHelper *RLPReaderHelperCaller) CustomDestructure(opts *bind.CallOpts, item []byte) (common.Address, bool, *big.Int, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new(bool)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _RLPReaderHelper.contract.Call(opts, out, "customDestructure", item)
	return *ret0, *ret1, *ret2, err
}

// CustomDestructure is a free data retrieval call binding the contract method 0x059a97fb.
//
// Solidity: function customDestructure(bytes item) constant returns(address, bool, uint256)
func (_RLPReaderHelper *RLPReaderHelperSession) CustomDestructure(item []byte) (common.Address, bool, *big.Int, error) {
	return _RLPReaderHelper.Contract.CustomDestructure(&_RLPReaderHelper.CallOpts, item)
}

// CustomDestructure is a free data retrieval call binding the contract method 0x059a97fb.
//
// Solidity: function customDestructure(bytes item) constant returns(address, bool, uint256)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) CustomDestructure(item []byte) (common.Address, bool, *big.Int, error) {
	return _RLPReaderHelper.Contract.CustomDestructure(&_RLPReaderHelper.CallOpts, item)
}

// CustomNestedDestructure is a free data retrieval call binding the contract method 0x15943cb6.
//
// Solidity: function customNestedDestructure(bytes item) constant returns(address, uint256)
func (_RLPReaderHelper *RLPReaderHelperCaller) CustomNestedDestructure(opts *bind.CallOpts, item []byte) (common.Address, *big.Int, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _RLPReaderHelper.contract.Call(opts, out, "customNestedDestructure", item)
	return *ret0, *ret1, err
}

// CustomNestedDestructure is a free data retrieval call binding the contract method 0x15943cb6.
//
// Solidity: function customNestedDestructure(bytes item) constant returns(address, uint256)
func (_RLPReaderHelper *RLPReaderHelperSession) CustomNestedDestructure(item []byte) (common.Address, *big.Int, error) {
	return _RLPReaderHelper.Contract.CustomNestedDestructure(&_RLPReaderHelper.CallOpts, item)
}

// CustomNestedDestructure is a free data retrieval call binding the contract method 0x15943cb6.
//
// Solidity: function customNestedDestructure(bytes item) constant returns(address, uint256)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) CustomNestedDestructure(item []byte) (common.Address, *big.Int, error) {
	return _RLPReaderHelper.Contract.CustomNestedDestructure(&_RLPReaderHelper.CallOpts, item)
}

// CustomNestedToRlpBytes is a free data retrieval call binding the contract method 0x74c4c494.
//
// Solidity: function customNestedToRlpBytes(bytes item) constant returns(bytes)
func (_RLPReaderHelper *RLPReaderHelperCaller) CustomNestedToRlpBytes(opts *bind.CallOpts, item []byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "customNestedToRlpBytes", item)
	return *ret0, err
}

// CustomNestedToRlpBytes is a free data retrieval call binding the contract method 0x74c4c494.
//
// Solidity: function customNestedToRlpBytes(bytes item) constant returns(bytes)
func (_RLPReaderHelper *RLPReaderHelperSession) CustomNestedToRlpBytes(item []byte) ([]byte, error) {
	return _RLPReaderHelper.Contract.CustomNestedToRlpBytes(&_RLPReaderHelper.CallOpts, item)
}

// CustomNestedToRlpBytes is a free data retrieval call binding the contract method 0x74c4c494.
//
// Solidity: function customNestedToRlpBytes(bytes item) constant returns(bytes)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) CustomNestedToRlpBytes(item []byte) ([]byte, error) {
	return _RLPReaderHelper.Contract.CustomNestedToRlpBytes(&_RLPReaderHelper.CallOpts, item)
}

// IsList is a free data retrieval call binding the contract method 0x19921bd0.
//
// Solidity: function isList(bytes item) constant returns(bool)
func (_RLPReaderHelper *RLPReaderHelperCaller) IsList(opts *bind.CallOpts, item []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "isList", item)
	return *ret0, err
}

// IsList is a free data retrieval call binding the contract method 0x19921bd0.
//
// Solidity: function isList(bytes item) constant returns(bool)
func (_RLPReaderHelper *RLPReaderHelperSession) IsList(item []byte) (bool, error) {
	return _RLPReaderHelper.Contract.IsList(&_RLPReaderHelper.CallOpts, item)
}

// IsList is a free data retrieval call binding the contract method 0x19921bd0.
//
// Solidity: function isList(bytes item) constant returns(bool)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) IsList(item []byte) (bool, error) {
	return _RLPReaderHelper.Contract.IsList(&_RLPReaderHelper.CallOpts, item)
}

// ItemLength is a free data retrieval call binding the contract method 0x8109672f.
//
// Solidity: function itemLength(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperCaller) ItemLength(opts *bind.CallOpts, item []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "itemLength", item)
	return *ret0, err
}

// ItemLength is a free data retrieval call binding the contract method 0x8109672f.
//
// Solidity: function itemLength(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperSession) ItemLength(item []byte) (*big.Int, error) {
	return _RLPReaderHelper.Contract.ItemLength(&_RLPReaderHelper.CallOpts, item)
}

// ItemLength is a free data retrieval call binding the contract method 0x8109672f.
//
// Solidity: function itemLength(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) ItemLength(item []byte) (*big.Int, error) {
	return _RLPReaderHelper.Contract.ItemLength(&_RLPReaderHelper.CallOpts, item)
}

// NestedIteration is a free data retrieval call binding the contract method 0x76193f8c.
//
// Solidity: function nestedIteration(bytes item) constant returns(string)
func (_RLPReaderHelper *RLPReaderHelperCaller) NestedIteration(opts *bind.CallOpts, item []byte) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "nestedIteration", item)
	return *ret0, err
}

// NestedIteration is a free data retrieval call binding the contract method 0x76193f8c.
//
// Solidity: function nestedIteration(bytes item) constant returns(string)
func (_RLPReaderHelper *RLPReaderHelperSession) NestedIteration(item []byte) (string, error) {
	return _RLPReaderHelper.Contract.NestedIteration(&_RLPReaderHelper.CallOpts, item)
}

// NestedIteration is a free data retrieval call binding the contract method 0x76193f8c.
//
// Solidity: function nestedIteration(bytes item) constant returns(string)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) NestedIteration(item []byte) (string, error) {
	return _RLPReaderHelper.Contract.NestedIteration(&_RLPReaderHelper.CallOpts, item)
}

// NumItems is a free data retrieval call binding the contract method 0xdd208ea1.
//
// Solidity: function numItems(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperCaller) NumItems(opts *bind.CallOpts, item []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "numItems", item)
	return *ret0, err
}

// NumItems is a free data retrieval call binding the contract method 0xdd208ea1.
//
// Solidity: function numItems(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperSession) NumItems(item []byte) (*big.Int, error) {
	return _RLPReaderHelper.Contract.NumItems(&_RLPReaderHelper.CallOpts, item)
}

// NumItems is a free data retrieval call binding the contract method 0xdd208ea1.
//
// Solidity: function numItems(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) NumItems(item []byte) (*big.Int, error) {
	return _RLPReaderHelper.Contract.NumItems(&_RLPReaderHelper.CallOpts, item)
}

// PayloadLen is a free data retrieval call binding the contract method 0x6393b0fb.
//
// Solidity: function payloadLen(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperCaller) PayloadLen(opts *bind.CallOpts, item []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "payloadLen", item)
	return *ret0, err
}

// PayloadLen is a free data retrieval call binding the contract method 0x6393b0fb.
//
// Solidity: function payloadLen(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperSession) PayloadLen(item []byte) (*big.Int, error) {
	return _RLPReaderHelper.Contract.PayloadLen(&_RLPReaderHelper.CallOpts, item)
}

// PayloadLen is a free data retrieval call binding the contract method 0x6393b0fb.
//
// Solidity: function payloadLen(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) PayloadLen(item []byte) (*big.Int, error) {
	return _RLPReaderHelper.Contract.PayloadLen(&_RLPReaderHelper.CallOpts, item)
}

// RlpLen is a free data retrieval call binding the contract method 0xaac19545.
//
// Solidity: function rlpLen(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperCaller) RlpLen(opts *bind.CallOpts, item []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "rlpLen", item)
	return *ret0, err
}

// RlpLen is a free data retrieval call binding the contract method 0xaac19545.
//
// Solidity: function rlpLen(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperSession) RlpLen(item []byte) (*big.Int, error) {
	return _RLPReaderHelper.Contract.RlpLen(&_RLPReaderHelper.CallOpts, item)
}

// RlpLen is a free data retrieval call binding the contract method 0xaac19545.
//
// Solidity: function rlpLen(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) RlpLen(item []byte) (*big.Int, error) {
	return _RLPReaderHelper.Contract.RlpLen(&_RLPReaderHelper.CallOpts, item)
}

// ToAddress is a free data retrieval call binding the contract method 0x2d888869.
//
// Solidity: function toAddress(bytes item) constant returns(address)
func (_RLPReaderHelper *RLPReaderHelperCaller) ToAddress(opts *bind.CallOpts, item []byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "toAddress", item)
	return *ret0, err
}

// ToAddress is a free data retrieval call binding the contract method 0x2d888869.
//
// Solidity: function toAddress(bytes item) constant returns(address)
func (_RLPReaderHelper *RLPReaderHelperSession) ToAddress(item []byte) (common.Address, error) {
	return _RLPReaderHelper.Contract.ToAddress(&_RLPReaderHelper.CallOpts, item)
}

// ToAddress is a free data retrieval call binding the contract method 0x2d888869.
//
// Solidity: function toAddress(bytes item) constant returns(address)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) ToAddress(item []byte) (common.Address, error) {
	return _RLPReaderHelper.Contract.ToAddress(&_RLPReaderHelper.CallOpts, item)
}

// ToBlockHeader is a free data retrieval call binding the contract method 0xe7617bf8.
//
// Solidity: function toBlockHeader(bytes rlpHeader) constant returns(bytes32 parentHash, bytes32 sha3Uncles, bytes32 stateRoot, bytes32 transactionsRoot, bytes32 receiptsRoot, uint256 difficulty, uint256 number, uint256 gasLimit, uint256 gasUsed, uint256 timestamp, uint256 nonce)
func (_RLPReaderHelper *RLPReaderHelperCaller) ToBlockHeader(opts *bind.CallOpts, rlpHeader []byte) (struct {
	ParentHash       [32]byte
	Sha3Uncles       [32]byte
	StateRoot        [32]byte
	TransactionsRoot [32]byte
	ReceiptsRoot     [32]byte
	Difficulty       *big.Int
	Number           *big.Int
	GasLimit         *big.Int
	GasUsed          *big.Int
	Timestamp        *big.Int
	Nonce            *big.Int
}, error) {
	ret := new(struct {
		ParentHash       [32]byte
		Sha3Uncles       [32]byte
		StateRoot        [32]byte
		TransactionsRoot [32]byte
		ReceiptsRoot     [32]byte
		Difficulty       *big.Int
		Number           *big.Int
		GasLimit         *big.Int
		GasUsed          *big.Int
		Timestamp        *big.Int
		Nonce            *big.Int
	})
	out := ret
	err := _RLPReaderHelper.contract.Call(opts, out, "toBlockHeader", rlpHeader)
	return *ret, err
}

// ToBlockHeader is a free data retrieval call binding the contract method 0xe7617bf8.
//
// Solidity: function toBlockHeader(bytes rlpHeader) constant returns(bytes32 parentHash, bytes32 sha3Uncles, bytes32 stateRoot, bytes32 transactionsRoot, bytes32 receiptsRoot, uint256 difficulty, uint256 number, uint256 gasLimit, uint256 gasUsed, uint256 timestamp, uint256 nonce)
func (_RLPReaderHelper *RLPReaderHelperSession) ToBlockHeader(rlpHeader []byte) (struct {
	ParentHash       [32]byte
	Sha3Uncles       [32]byte
	StateRoot        [32]byte
	TransactionsRoot [32]byte
	ReceiptsRoot     [32]byte
	Difficulty       *big.Int
	Number           *big.Int
	GasLimit         *big.Int
	GasUsed          *big.Int
	Timestamp        *big.Int
	Nonce            *big.Int
}, error) {
	return _RLPReaderHelper.Contract.ToBlockHeader(&_RLPReaderHelper.CallOpts, rlpHeader)
}

// ToBlockHeader is a free data retrieval call binding the contract method 0xe7617bf8.
//
// Solidity: function toBlockHeader(bytes rlpHeader) constant returns(bytes32 parentHash, bytes32 sha3Uncles, bytes32 stateRoot, bytes32 transactionsRoot, bytes32 receiptsRoot, uint256 difficulty, uint256 number, uint256 gasLimit, uint256 gasUsed, uint256 timestamp, uint256 nonce)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) ToBlockHeader(rlpHeader []byte) (struct {
	ParentHash       [32]byte
	Sha3Uncles       [32]byte
	StateRoot        [32]byte
	TransactionsRoot [32]byte
	ReceiptsRoot     [32]byte
	Difficulty       *big.Int
	Number           *big.Int
	GasLimit         *big.Int
	GasUsed          *big.Int
	Timestamp        *big.Int
	Nonce            *big.Int
}, error) {
	return _RLPReaderHelper.Contract.ToBlockHeader(&_RLPReaderHelper.CallOpts, rlpHeader)
}

// ToBoolean is a free data retrieval call binding the contract method 0x04ba3427.
//
// Solidity: function toBoolean(bytes item) constant returns(bool)
func (_RLPReaderHelper *RLPReaderHelperCaller) ToBoolean(opts *bind.CallOpts, item []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "toBoolean", item)
	return *ret0, err
}

// ToBoolean is a free data retrieval call binding the contract method 0x04ba3427.
//
// Solidity: function toBoolean(bytes item) constant returns(bool)
func (_RLPReaderHelper *RLPReaderHelperSession) ToBoolean(item []byte) (bool, error) {
	return _RLPReaderHelper.Contract.ToBoolean(&_RLPReaderHelper.CallOpts, item)
}

// ToBoolean is a free data retrieval call binding the contract method 0x04ba3427.
//
// Solidity: function toBoolean(bytes item) constant returns(bool)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) ToBoolean(item []byte) (bool, error) {
	return _RLPReaderHelper.Contract.ToBoolean(&_RLPReaderHelper.CallOpts, item)
}

// ToBytes is a free data retrieval call binding the contract method 0xe79fbcec.
//
// Solidity: function toBytes(bytes item) constant returns(bytes)
func (_RLPReaderHelper *RLPReaderHelperCaller) ToBytes(opts *bind.CallOpts, item []byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "toBytes", item)
	return *ret0, err
}

// ToBytes is a free data retrieval call binding the contract method 0xe79fbcec.
//
// Solidity: function toBytes(bytes item) constant returns(bytes)
func (_RLPReaderHelper *RLPReaderHelperSession) ToBytes(item []byte) ([]byte, error) {
	return _RLPReaderHelper.Contract.ToBytes(&_RLPReaderHelper.CallOpts, item)
}

// ToBytes is a free data retrieval call binding the contract method 0xe79fbcec.
//
// Solidity: function toBytes(bytes item) constant returns(bytes)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) ToBytes(item []byte) ([]byte, error) {
	return _RLPReaderHelper.Contract.ToBytes(&_RLPReaderHelper.CallOpts, item)
}

// ToIterator is a free data retrieval call binding the contract method 0xd03ca5a3.
//
// Solidity: function toIterator(bytes item) constant returns()
func (_RLPReaderHelper *RLPReaderHelperCaller) ToIterator(opts *bind.CallOpts, item []byte) error {
	var ()
	out := &[]interface{}{}
	err := _RLPReaderHelper.contract.Call(opts, out, "toIterator", item)
	return err
}

// ToIterator is a free data retrieval call binding the contract method 0xd03ca5a3.
//
// Solidity: function toIterator(bytes item) constant returns()
func (_RLPReaderHelper *RLPReaderHelperSession) ToIterator(item []byte) error {
	return _RLPReaderHelper.Contract.ToIterator(&_RLPReaderHelper.CallOpts, item)
}

// ToIterator is a free data retrieval call binding the contract method 0xd03ca5a3.
//
// Solidity: function toIterator(bytes item) constant returns()
func (_RLPReaderHelper *RLPReaderHelperCallerSession) ToIterator(item []byte) error {
	return _RLPReaderHelper.Contract.ToIterator(&_RLPReaderHelper.CallOpts, item)
}

// ToRlpBytes is a free data retrieval call binding the contract method 0x33c82d18.
//
// Solidity: function toRlpBytes(bytes item) constant returns(bytes)
func (_RLPReaderHelper *RLPReaderHelperCaller) ToRlpBytes(opts *bind.CallOpts, item []byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "toRlpBytes", item)
	return *ret0, err
}

// ToRlpBytes is a free data retrieval call binding the contract method 0x33c82d18.
//
// Solidity: function toRlpBytes(bytes item) constant returns(bytes)
func (_RLPReaderHelper *RLPReaderHelperSession) ToRlpBytes(item []byte) ([]byte, error) {
	return _RLPReaderHelper.Contract.ToRlpBytes(&_RLPReaderHelper.CallOpts, item)
}

// ToRlpBytes is a free data retrieval call binding the contract method 0x33c82d18.
//
// Solidity: function toRlpBytes(bytes item) constant returns(bytes)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) ToRlpBytes(item []byte) ([]byte, error) {
	return _RLPReaderHelper.Contract.ToRlpBytes(&_RLPReaderHelper.CallOpts, item)
}

// ToUint is a free data retrieval call binding the contract method 0x8529d576.
//
// Solidity: function toUint(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperCaller) ToUint(opts *bind.CallOpts, item []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "toUint", item)
	return *ret0, err
}

// ToUint is a free data retrieval call binding the contract method 0x8529d576.
//
// Solidity: function toUint(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperSession) ToUint(item []byte) (*big.Int, error) {
	return _RLPReaderHelper.Contract.ToUint(&_RLPReaderHelper.CallOpts, item)
}

// ToUint is a free data retrieval call binding the contract method 0x8529d576.
//
// Solidity: function toUint(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) ToUint(item []byte) (*big.Int, error) {
	return _RLPReaderHelper.Contract.ToUint(&_RLPReaderHelper.CallOpts, item)
}

// ToUintStrict is a free data retrieval call binding the contract method 0x20f91849.
//
// Solidity: function toUintStrict(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperCaller) ToUintStrict(opts *bind.CallOpts, item []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RLPReaderHelper.contract.Call(opts, out, "toUintStrict", item)
	return *ret0, err
}

// ToUintStrict is a free data retrieval call binding the contract method 0x20f91849.
//
// Solidity: function toUintStrict(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperSession) ToUintStrict(item []byte) (*big.Int, error) {
	return _RLPReaderHelper.Contract.ToUintStrict(&_RLPReaderHelper.CallOpts, item)
}

// ToUintStrict is a free data retrieval call binding the contract method 0x20f91849.
//
// Solidity: function toUintStrict(bytes item) constant returns(uint256)
func (_RLPReaderHelper *RLPReaderHelperCallerSession) ToUintStrict(item []byte) (*big.Int, error) {
	return _RLPReaderHelper.Contract.ToUintStrict(&_RLPReaderHelper.CallOpts, item)
}

// TryParseLog attempts to parse a log. Returns the parsed log, evenName and whether it was succesfull
func (_RLPReaderHelper *RLPReaderHelperFilterer) TryParseLog(log types.Log) (eventName string, event interface{}, ok bool, err error) {
	eventName, ok, err = _RLPReaderHelper.contract.LogEventName(log)
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
