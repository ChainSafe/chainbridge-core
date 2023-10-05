package deposithandlers

import (
	"errors"

	"github.com/ChainSafe/sygma-core/chains/evm/eventhandlers"

	"github.com/ChainSafe/sygma-core/types"
	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/common"
)

type arbitraryFunction func(Config interface{}) error
type DepositHandlers map[common.Address]eventhandlers.DepositHandler
type DepositHandlerFunc func(sourceID, destID uint8, nonce uint64, resourceID types.ResourceID, calldata, handlerResponse []byte) (*types.Message, error)
type HandlerMatcher interface {
	GetHandlerAddressForResourceID(resourceID types.ResourceID) (common.Address, error)
}

type ETHDepositHandler struct {
	handlerMatcher  HandlerMatcher
	depositHandlers DepositHandlers
}

// NewETHDepositHandler creates an instance of ETHDepositHandler that contains
// handler functions for processing deposit events
func NewETHDepositHandler(handlerMatcher HandlerMatcher) *ETHDepositHandler {
	return &ETHDepositHandler{
		handlerMatcher:  handlerMatcher,
		depositHandlers: make(map[common.Address]eventhandlers.DepositHandler),
	}
}

func (e *ETHDepositHandler) HandleDeposit(sourceID, destID uint8, depositNonce uint64, resourceID types.ResourceID, calldata, handlerResponse []byte) (*types.Message, error) {
	handlerAddr, err := e.handlerMatcher.GetHandlerAddressForResourceID(resourceID)
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
	return hf.HandleDeposit, nil
}

// RegisterDepositHandler registers an event handler by associating a handler function to a specified address
func (e *ETHDepositHandler) RegisterDepositHandler(handlerAddress string, handler eventhandlers.DepositHandler) {
	if handlerAddress == "" {
		return
	}

	log.Debug().Msgf("Registered deposit handler for address %s", handlerAddress)
	e.depositHandlers[common.HexToAddress(handlerAddress)] = handler
}
