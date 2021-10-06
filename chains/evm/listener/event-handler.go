package listener

import (
	"context"
	"errors"
	"strings"

	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type EventHandlers map[common.Address]EventHandlerFunc
type EventHandlerFunc func(sourceID, destId uint8, nonce uint64, handlerContractAddress common.Address, caller ChainClient, resourceId [32]byte, data []byte) (*relayer.Message, error)

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

func (e *ETHEventHandler) HandleEvent(sourceID, destID uint8, depositNonce uint64, rID [32]byte, data []byte) (*relayer.Message, error) {
	addr, err := e.matchResourceIDToHandlerAddress(rID)
	if err != nil {
		return nil, err
	}

	eventHandler, err := e.matchAddressWithHandlerFunc(addr)
	if err != nil {
		return nil, err
	}

	return eventHandler(sourceID, destID, depositNonce, addr, e.client, rID, data)
}

func (e *ETHEventHandler) matchResourceIDToHandlerAddress(rID [32]byte) (common.Address, error) {
	definition := "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"_resourceIDToHandlerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return common.Address{}, err
	}
	input, err := a.Pack("_resourceIDToHandlerAddress", rID)
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

func (e *ETHEventHandler) matchAddressWithHandlerFunc(addr common.Address) (EventHandlerFunc, error) {
	hf, ok := e.eventHandlers[addr]
	if !ok {
		return nil, errors.New("no corresponding event handler for this address exists")
	}
	return hf, nil
}

func (e *ETHEventHandler) RegisterEventHandler(address string, handler EventHandlerFunc) {
	if e.eventHandlers == nil {
		e.eventHandlers = make(map[common.Address]EventHandlerFunc)
	}
	e.eventHandlers[common.HexToAddress(address)] = handler
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

// Erc20EventHandler converts data pulled from contract event logs into message
func Erc20EventHandler(sourceID, destId uint8, nonce uint64, handlerContractAddress common.Address, client ChainClient, resourceID [32]byte, calldata []byte) (*relayer.Message, error) {

	// TODO: parse calldata

	return &relayer.Message{
		Source:       sourceID,
		Destination:  destId,
		DepositNonce: nonce,
		ResourceId:   resourceID,
		Type:         relayer.FungibleTransfer,
		Payload: []interface{}{
			calldata, // amount?
			// add destination recipient address?
		},
	}, nil
}
