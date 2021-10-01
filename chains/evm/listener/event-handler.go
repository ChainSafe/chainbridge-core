package listener

import (
	"context"
	"errors"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/relayer"
	internalTypes "github.com/ChainSafe/chainbridge-core/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type EventHandlers map[common.Address]EventHandlerFunc
type EventHandlerFunc func(sourceID, destId uint8, nonce uint64, resourceID internalTypes.ResourceID, calldata, handlerResponse []byte) (*relayer.Message, error)

type ETHEventHandler struct {
	bridgeAddress common.Address
	eventHandlers EventHandlers
	client        ChainClient
}

func NewETHEventHandler(address common.Address, client ChainClient) *ETHEventHandler {
	return &ETHEventHandler{
		bridgeAddress: address,
		client:        client,
	}
}

func (e *ETHEventHandler) HandleEvent(sourceID, destID uint8, depositNonce uint64, resourceID internalTypes.ResourceID, calldata, handlerResponse []byte) (*relayer.Message, error) {
	handlerAddr, err := e.matchResourceIDToHandlerAddress(resourceID)
	if err != nil {
		return nil, err
	}

	eventHandler, err := e.matchAddressWithHandlerFunc(handlerAddr)
	if err != nil {
		return nil, err
	}

	return eventHandler(sourceID, destID, depositNonce, resourceID, calldata, handlerResponse)
}

// matchResourceIDToHandlerAddress is a private method that matches a previously registered resource ID to its corresponding handler address
func (e *ETHEventHandler) matchResourceIDToHandlerAddress(resourceID internalTypes.ResourceID) (common.Address, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return common.Address{}, err
	}
	input, err := a.Pack("_resourceIDToHandlerAddress", resourceID)
	if err != nil {
		return common.Address{}, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &e.bridgeAddress, Data: input}
	out, err := e.client.CallContract(context.TODO(), toCallArg(msg), nil)
	if err != nil {
		return common.Address{}, err
	}
	res, err := a.Unpack("_resourceIDToHandlerAddress", out)
	if err != nil {
		return common.Address{}, err
	}
	out0 := *abi.ConvertType(res[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// matchAddressWithHandlerFunc is a private method that matches a handler address with an associated handler function
func (e *ETHEventHandler) matchAddressWithHandlerFunc(handlerAddress common.Address) (EventHandlerFunc, error) {
	hf, ok := e.eventHandlers[handlerAddress]
	if !ok {
		return nil, errors.New("no corresponding event handler for this address exists")
	}
	return hf, nil
}

// RegisterEventHandler is a public method that registers an event handler by associating a handler function to a specific address
func (e *ETHEventHandler) RegisterEventHandler(handlerAddress string, handler EventHandlerFunc) {
	if e.eventHandlers == nil {
		e.eventHandlers = make(map[common.Address]EventHandlerFunc)
	}
	e.eventHandlers[common.HexToAddress(handlerAddress)] = handler
}

func toCallArg(msg ethereum.CallMsg) map[string]interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}

// Erc20EventHandler converts data pulled from event logs into message
// handlerResponse can be an empty slice
func Erc20EventHandler(sourceID, destId uint8, nonce uint64, resourceID internalTypes.ResourceID, calldata, handlerResponse []byte) (*relayer.Message, error) {
	if len(calldata) < 84 {
		err := errors.New("invalid calldata length: less than 84 bytes")
		return nil, err
	}

	// @dev
	// amount: first 32 bytes of calldata
	amount := calldata[:32]

	// lenRecipientAddress: second 32 bytes of calldata [32:64]
	// does not need to be derived because it is being calculated
	// within ERC20MessageHandler
	// https://github.com/ChainSafe/chainbridge-core/blob/main/chains/evm/voter/message-handler.go#L108

	// recipientAddress: last 20 bytes of calldata
	recipientAddress := calldata[64:]

	return &relayer.Message{
		Source:       sourceID,
		Destination:  destId,
		DepositNonce: nonce,
		ResourceId:   resourceID,
		Type:         relayer.FungibleTransfer,
		Payload: []interface{}{
			amount,
			recipientAddress,
		},
	}, nil
}

func GenericEventHandler(sourceID, destId uint8, nonce uint64, handlerContractAddress common.Address, client ChainClient) (*relayer.Message, error) {
	type GenericHandlerDepositRecord struct {
		DestinationChainID uint8
		Depositer          common.Address
		ResourceID         [32]byte
		MetaData           []byte
	}

	rec, err := getDepositRecord(consts.GenericHandlerABI, nonce, destId, &handlerContractAddress, client)
	if err != nil {
		return nil, err
	}

	out0 := *abi.ConvertType(rec, new(GenericHandlerDepositRecord)).(*GenericHandlerDepositRecord)
	return &relayer.Message{
		Source:       sourceID,
		Destination:  destId,
		DepositNonce: nonce,
		ResourceId:   out0.ResourceID,
		Type:         relayer.GenericTransfer,
		Payload: []interface{}{
			out0.MetaData,
		},
	}, nil
}

func Erc721EventHandler(sourceID, destId uint8, nonce uint64, handlerContractAddress common.Address, client ChainClient) (*relayer.Message, error) {
	type Erc721HandlerDepositRecord struct {
		TokenAddress                   common.Address
		LenDestinationRecipientAddress uint8
		DestinationDomainID            uint8
		ResourceID                     [32]byte
		DestinationRecipientAddress    []byte
		Depositer                      common.Address
		TokenId                        [32]byte
		MetaData                       []byte
	}

	rec, err := getDepositRecord(consts.ERC721HandlerABI, nonce, destId, &handlerContractAddress, client)
	if err != nil {
		return nil, err
	}

	out0 := *abi.ConvertType(rec, new(Erc721HandlerDepositRecord)).(*Erc721HandlerDepositRecord)
	return &relayer.Message{
		Source:       sourceID,
		Destination:  destId,
		DepositNonce: nonce,
		ResourceId:   out0.ResourceID,
		Type:         relayer.NonFungibleTransfer,
		Payload: []interface{}{
			out0.TokenId,
			out0.DestinationRecipientAddress,
			out0.MetaData,
		},
	}, nil
}

func GenericEventHandler(sourceID, destId uint8, nonce uint64, handlerContractAddress common.Address, client ChainClient) (*relayer.Message, error) {
	definition := "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"initialResourceIDs\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"initialContractAddresses\",\"type\":\"address[]\"},{\"internalType\":\"bytes4[]\",\"name\":\"initialDepositFunctionSignatures\",\"type\":\"bytes4[]\"},{\"internalType\":\"bytes4[]\",\"name\":\"initialExecuteFunctionSignatures\",\"type\":\"bytes4[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"_bridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_contractAddressToDepositFunctionSignature\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_contractAddressToExecuteFunctionSignature\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_contractAddressToResourceID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_contractWhitelist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"_depositRecords\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"_destinationChainID\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"_depositer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_metaData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"_resourceIDToContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"destId\",\"type\":\"uint8\"}],\"name\":\"getDepositRecord\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"_destinationChainID\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"_depositer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_metaData\",\"type\":\"bytes\"}],\"internalType\":\"structGenericHandler.DepositRecord\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"depositFunctionSig\",\"type\":\"bytes4\"},{\"internalType\":\"bytes4\",\"name\":\"executeFunctionSig\",\"type\":\"bytes4\"}],\"name\":\"setResource\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"destinationChainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"depositer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"executeProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	type GenericHandlerDepositRecord struct {
		DestinationChainID uint8
		Depositer          common.Address
		ResourceID         [32]byte
		MetaData           []byte
	}

	rec, err := getDepositRecord(definition, nonce, destId, &handlerContractAddress, client)
	if err != nil {
		return nil, err
	}

	out0 := *abi.ConvertType(rec, new(GenericHandlerDepositRecord)).(*GenericHandlerDepositRecord)
	return &relayer.Message{
		Source:       sourceID,
		Destination:  destId,
		DepositNonce: nonce,
		ResourceId:   out0.ResourceID,
		Type:         relayer.GenericTransfer,
		Payload: []interface{}{
			out0.MetaData,
		},
	}, nil
}
