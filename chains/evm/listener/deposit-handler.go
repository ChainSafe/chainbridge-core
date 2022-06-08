package listener

import (
	"errors"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/chainbridge-core/types"
	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/common"
)

type DepositHandlers map[common.Address]DepositHandlerFunc
type DepositHandlerFunc func(sourceID, destId uint8, nonce uint64, resourceID types.ResourceID, calldata, handlerResponse []byte) (*message.Message, error)

type ETHDepositHandler struct {
	bridgeContract  bridge.BridgeContract
	depositHandlers DepositHandlers
}

// NewETHDepositHandler creates an instance of ETHDepositHandler that contains
// handler functions for processing deposit events
func NewETHDepositHandler(bridgeContract bridge.BridgeContract) *ETHDepositHandler {
	return &ETHDepositHandler{
		bridgeContract:  bridgeContract,
		depositHandlers: make(map[common.Address]DepositHandlerFunc),
	}
}

func (e *ETHDepositHandler) HandleDeposit(sourceID, destID uint8, depositNonce uint64, resourceID types.ResourceID, calldata, handlerResponse []byte) (*message.Message, error) {
	handlerAddr, err := e.bridgeContract.GetHandlerAddressForResourceID(resourceID)
	if err != nil {
		return nil, err
	}

	depositHandler, err := e.matchAddressWithHandlerFunc(handlerAddr)
	if err != nil {
		return nil, err
	}

	return depositHandler(sourceID, destID, depositNonce, resourceID, calldata, handlerResponse)
}

// matchAddressWithHandlerFunc matches a handler address with an associated handler function
func (e *ETHDepositHandler) matchAddressWithHandlerFunc(handlerAddress common.Address) (DepositHandlerFunc, error) {
	hf, ok := e.depositHandlers[handlerAddress]
	if !ok {
		return nil, errors.New("no corresponding deposit handler for this address exists")
	}
	return hf, nil
}

// RegisterDepositHandler registers an event handler by associating a handler function to a specified address
func (e *ETHDepositHandler) RegisterDepositHandler(handlerAddress string, handler DepositHandlerFunc) {
	if handlerAddress == "" {
		return
	}

	log.Info().Msgf("Registered deposit handler for address %s", handlerAddress)
	e.depositHandlers[common.HexToAddress(handlerAddress)] = handler
}

// Erc20DepositHandler converts data pulled from event logs into message
// handlerResponse can be an empty slice
func Erc20DepositHandler(sourceID, destId uint8, nonce uint64, resourceID types.ResourceID, calldata, handlerResponse []byte) (*message.Message, error) {
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

	// 32-64 is recipient address length
	recipientAddressLength := big.NewInt(0).SetBytes(calldata[32:64])

	// 64 - (64 + recipient address length) is recipient address
	recipientAddress := calldata[64:(64 + recipientAddressLength.Int64())]

	// if there is priority data, parse it and use it
	payload := []interface{}{
		amount,
		recipientAddress,
	}

	// arbitrary metadata that will be most likely be used by the relayer
	var metadata message.Metadata
	if 64+recipientAddressLength.Int64() < int64(len(calldata)) {
		priorityLength := big.NewInt(0).SetBytes(calldata[(64 + recipientAddressLength.Int64()):((64 + recipientAddressLength.Int64()) + 1)])

		// (64 + recipient address length + 1) - ((64 + recipient address length + 1) + priority length) is priority data
		priority := calldata[(64 + recipientAddressLength.Int64() + 1):((64 + recipientAddressLength.Int64()) + 1 + priorityLength.Int64())]
		metadata.Priority = priority[0]
	}
	return message.NewMessage(sourceID, destId, nonce, resourceID, message.FungibleTransfer, payload, metadata), nil
}

// GenericDepositHandler converts data pulled from generic deposit event logs into message
func GenericDepositHandler(sourceID, destId uint8, nonce uint64, resourceID types.ResourceID, calldata, handlerResponse []byte) (*message.Message, error) {
	if len(calldata) < 32 {
		err := errors.New("invalid calldata length: less than 32 bytes")
		return nil, err
	}

	// first 32 bytes are metadata length
	metadataLen := big.NewInt(0).SetBytes(calldata[:32])
	metadata := calldata[32 : 32+metadataLen.Int64()]
	payload := []interface{}{
		metadata,
	}

	// generic handler has specific payload length and doesn't support arbitrary metadata
	meta := message.Metadata{}
	return message.NewMessage(sourceID, destId, nonce, resourceID, message.GenericTransfer, payload, meta), nil
}

// Erc721DepositHandler converts data pulled from ERC721 deposit event logs into message
func Erc721DepositHandler(sourceID, destId uint8, nonce uint64, resourceID types.ResourceID, calldata, handlerResponse []byte) (*message.Message, error) {
	if len(calldata) < 64 {
		err := errors.New("invalid calldata length: less than 84 bytes")
		return nil, err
	}

	// first 32 bytes are tokenId
	tokenId := calldata[:32]

	// 32 - 64 is recipient address length
	recipientAddressLength := big.NewInt(0).SetBytes(calldata[32:64])

	// 64 - (64 + recipient address length) is recipient address
	recipientAddress := calldata[64:(64 + recipientAddressLength.Int64())]

	// (64 + recipient address length) - ((64 + recipient address length) + 32) is metadata length
	metadataLength := big.NewInt(0).SetBytes(
		calldata[(64 + recipientAddressLength.Int64()):((64 + recipientAddressLength.Int64()) + 32)],
	)
	// ((64 + recipient address length) + 32) - ((64 + recipient address length) + 32 + metadata length) is metadata
	var metadata []byte
	var metadataStart int64
	if metadataLength.Cmp(big.NewInt(0)) == 1 {
		metadataStart = (64 + recipientAddressLength.Int64()) + 32
		metadata = calldata[metadataStart : metadataStart+metadataLength.Int64()]
	}
	// arbitrary metadata that will be most likely be used by the relayer
	var meta message.Metadata

	payload := []interface{}{
		tokenId,
		recipientAddress,
		metadata,
	}

	if 64+recipientAddressLength.Int64()+32+metadataLength.Int64() < int64(len(calldata)) {
		// (metadataStart + metadataLength) - (metadataStart + metadataLength + 1) is priority length
		priorityLength := big.NewInt(0).SetBytes(calldata[(64 + recipientAddressLength.Int64() + 32 + metadataLength.Int64()):(64 + recipientAddressLength.Int64() + 32 + metadataLength.Int64() + 1)])
		// (metadataStart + metadataLength + 1) - (metadataStart + metadataLength + 1) + priority length) is priority data
		priority := calldata[(64 + recipientAddressLength.Int64() + 32 + metadataLength.Int64() + 1):(64 + recipientAddressLength.Int64() + 32 + metadataLength.Int64() + 1 + priorityLength.Int64())]
		meta.Priority = priority[0]
	}
	return message.NewMessage(sourceID, destId, nonce, resourceID, message.NonFungibleTransfer, payload, meta), nil
}
