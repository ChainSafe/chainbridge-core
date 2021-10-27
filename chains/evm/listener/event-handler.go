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
	if len(calldata) == 0 {
		err := errors.New("missing calldata")
		return nil, err
	}

	amount := calldata[:32]
	recipientAddress := calldata[65:]

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
