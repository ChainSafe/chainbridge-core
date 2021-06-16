package listener

import (
	"context"
	"errors"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"
)

type EventHandlers map[common.Address]EventHandlerFunc
type EventHandlerFunc func(sourceID, destId uint8, nonce uint64, handlerContractAddress common.Address, caller ChainClient) (*relayer.Message, error)

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

func (e *ETHEventHandler) HandleEvent(sourceID, destID uint8, depositNonce uint64, rID [32]byte) (*relayer.Message, error) {
	addr, err := e.matchResourceIDToHandlerAddress(rID)
	if err != nil {
		return nil, err
	}

	eventHandler, err := e.matchAddressWithHandlerFunc(addr)
	if err != nil {
		return nil, err
	}

	return eventHandler(sourceID, destID, depositNonce, addr, e.client)
}

func (e *ETHEventHandler) matchResourceIDToHandlerAddress(rID [32]byte) (common.Address, error) {
	//_resourceIDToHandlerAddress(bytes32) view returns(address)
	input, err := buildDataUnsafe([]byte("_resourceIDToHandlerAddress(bytes32"), rID[:])
	if err != nil {
		return common.Address{}, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &e.bridgeAddress, Data: input}
	out, err := e.client.CallContract(context.TODO(), toCallArg(msg), nil)
	if err != nil {
		return common.Address{}, err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

func (e *ETHEventHandler) matchAddressWithHandlerFunc(addr common.Address) (EventHandlerFunc, error) {
	hf, ok := e.eventHandlers[addr]
	if !ok {
		return nil, errors.New("no corresponding handler for this address exists")
	}
	return hf, nil
}

func (e *ETHEventHandler) RegisterHandlerFabric(address string, handler EventHandlerFunc) {
	if e.eventHandlers == nil {
		e.eventHandlers = make(map[common.Address]EventHandlerFunc)
	}
	e.eventHandlers[common.HexToAddress(address)] = handler
}

func buildDataUnsafe(method []byte, params ...[]byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(method)
	if err != nil {
		return nil, err
	}
	methodID := hash.Sum(nil)[:4]

	var data []byte
	data = append(data, methodID...)
	for _, v := range params {
		paddedParam := common.LeftPadBytes(v, 32)
		data = append(data, paddedParam...)
	}
	return data, nil
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

func Erc20EventHandler(sourceID, destId uint8, nonce uint64, handlerContractAddress common.Address, client ChainClient) (*relayer.Message, error) {
	type ERC20HandlerDepositRecord struct {
		TokenAddress                   common.Address
		LenDestinationRecipientAddress uint8
		DestinationChainID             uint8
		ResourceID                     [32]byte
		DestinationRecipientAddress    []byte
		Depositer                      common.Address
		Amount                         *big.Int
	}
	input, err := buildDataUnsafe([]byte("getDepositRecord(uint64,uint8"), big.NewInt(0).SetUint64(nonce).Bytes(), big.NewInt(0).SetUint64(uint64(destId)).Bytes())
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &handlerContractAddress, Data: input}
	res, err := client.CallContract(context.TODO(), toCallArg(msg), nil)
	if err != nil {
		return nil, err
	}
	out0 := *abi.ConvertType(res[0], new(ERC20HandlerDepositRecord)).(*ERC20HandlerDepositRecord)
	return &relayer.Message{
		Source:       sourceID,
		Destination:  destId,
		DepositNonce: nonce,
		ResourceId:   out0.ResourceID,
		Type:         relayer.FungibleTransfer,
		Payload: []interface{}{
			out0.Amount.Bytes(),
			out0.DestinationRecipientAddress,
		},
	}, nil
}
