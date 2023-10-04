package executor

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type HandlerMatcher interface {
	GetHandlerAddressForResourceID(resourceID types.ResourceID) (common.Address, error)
	ContractAddress() *common.Address
}

type MessageHandlerFunc func(m *types.Message, handlerAddr, bridgeAddress common.Address) (*types.Proposal, error)

// NewEVMMessageHandler creates an instance of EVMMessageHandler that contains
// message handler functions for converting deposit message into a chain specific
// proposal
func NewEVMMessageHandler(handlerMatcher HandlerMatcher) *EVMMessageHandler {
	return &EVMMessageHandler{
		handlerMatcher: handlerMatcher,
	}
}

type EVMMessageHandler struct {
	handlerMatcher HandlerMatcher
	handlers       map[common.Address]MessageHandlerFunc
}

func (mh *EVMMessageHandler) HandleMessage(m *types.Message) (*types.Proposal, error) {
	// Matching resource ID with handler.
	addr, err := mh.handlerMatcher.GetHandlerAddressForResourceID(m.ResourceId)
	if err != nil {
		return nil, err
	}
	// Based on handler that registered on BridgeContract
	handleMessage, err := mh.MatchAddressWithHandlerFunc(addr)
	if err != nil {
		return nil, err
	}
	log.Info().Str("type", string(m.Type)).Uint8("src", m.Source).Uint8("dst", m.Destination).Uint64("nonce", m.DepositNonce).Str("resourceID", fmt.Sprintf("%x", m.ResourceId)).Msg("Handling new message")
	prop, err := handleMessage(m, addr, *mh.handlerMatcher.ContractAddress())
	if err != nil {
		return nil, err
	}
	return prop, nil
}

func (mh *EVMMessageHandler) MatchAddressWithHandlerFunc(addr common.Address) (MessageHandlerFunc, error) {
	h, ok := mh.handlers[addr]
	if !ok {
		return nil, fmt.Errorf("no corresponding message handler for this address %s exists", addr.Hex())
	}
	return h, nil
}

// RegisterEventHandler registers an message handler by associating a handler function to a specified address
func (mh *EVMMessageHandler) RegisterMessageHandler(address string, handler MessageHandlerFunc) {
	if address == "" {
		return
	}
	if mh.handlers == nil {
		mh.handlers = make(map[common.Address]MessageHandlerFunc)
	}

	log.Debug().Msgf("Registered message handler for address %s", address)

	mh.handlers[common.HexToAddress(address)] = handler
}
