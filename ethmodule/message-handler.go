package ethmodule

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"

	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type HandlerFunc func(m *relayer.Message, handlerAddr, bridgeAddress common.Address) (*Proposal, error)

func NewMessageHandler(evmCaller voter.EVMClient, bridgeAddress common.Address) *MessageHandler {
	return &MessageHandler{
		bridgeAddress: bridgeAddress,
		evmCaller:     evmCaller,
	}
}

type MessageHandler struct {
	evmCaller     voter.EVMClient
	handlers      map[common.Address]HandlerFunc
	bridgeAddress common.Address
}

func (mh *MessageHandler) HandleMessage(m *relayer.Message) (*Proposal, error) {
	// Matching resource ID with handler.
	addr, err := mh.matchResourceIDToHandlerAddress(m.ResourceId)
	// Based on handler that registered on BridgeContract
	handleMessage, err := mh.MatchAddressWithHandlerFunc(addr)
	if err != nil {
		return nil, err
	}
	log.Info().Str("type", string(m.Type)).Uint8("src", m.Source).Uint8("dst", m.Destination).Uint64("nonce", m.DepositNonce).Str("rId", fmt.Sprintf("%x", m.ResourceId)).Msg("Handling new message")
	prop, err := handleMessage(m, addr, mh.bridgeAddress)
	if err != nil {
		return nil, err
	}
	return prop, nil
}

func (mh *MessageHandler) matchResourceIDToHandlerAddress(rID [32]byte) (common.Address, error) {
	//_resourceIDToHandlerAddress(bytes32) view returns(address)
	input, err := buildDataUnsafe([]byte("_resourceIDToHandlerAddress(bytes32"), rID[:])
	if err != nil {
		return common.Address{}, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &mh.bridgeAddress, Data: input}
	out, err := mh.evmCaller.CallContract(context.TODO(), toCallArg(msg), nil)
	if err != nil {
		return common.Address{}, err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

func (mh *MessageHandler) MatchAddressWithHandlerFunc(addr common.Address) (HandlerFunc, error) {
	h, ok := mh.handlers[addr]
	if !ok {
		return nil, errors.New("no corresponding handler for this address exists")
	}
	return h, nil
}

func (mh *MessageHandler) RegisterProposalHandler(address common.Address, handler HandlerFunc) {
	mh.handlers[address] = handler
}

func ERC20MessageHandler(m *relayer.Message, handlerAddr, bridgeAddress common.Address) (*Proposal, error) {
	if len(m.Payload) != 2 {
		return nil, errors.New("malformed payload. Len  of payload should be 2")
	}
	amount, ok := m.Payload[0].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads amount format")
	}

	recipient, ok := m.Payload[1].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads recipient format")

	}
	var data []byte
	data = append(data, common.LeftPadBytes(amount, 32)...) // amount (uint256)

	recipientLen := big.NewInt(int64(len(recipient))).Bytes()
	data = append(data, common.LeftPadBytes(recipientLen, 32)...) // length of recipient (uint256)
	data = append(data, recipient...)                             // recipient ([]byte)

	return &Proposal{
		Source:         m.Source,
		DepositNonce:   m.DepositNonce,
		ResourceId:     m.ResourceId,
		Data:           data,
		HandlerAddress: handlerAddr,
		BridgeAddress:  bridgeAddress,
	}, nil
}

func ERC721MessageHandler(msg *relayer.Message, handlerAddr, bridgeAddress common.Address) (*Proposal, error) {
	if len(msg.Payload) != 3 {
		return nil, errors.New("malformed payload. Len  of payload should be 3")
	}
	tokenID, ok := msg.Payload[0].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads tokenID format")
	}
	recipient, ok := msg.Payload[1].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads recipient format")
	}
	metadata, ok := msg.Payload[2].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads metadata format")
	}

	data := bytes.Buffer{}
	data.Write(common.LeftPadBytes(tokenID, 32))

	recipientLen := big.NewInt(int64(len(recipient))).Bytes()
	data.Write(common.LeftPadBytes(recipientLen, 32))
	data.Write(recipient)

	metadataLen := big.NewInt(int64(len(metadata))).Bytes()
	data.Write(common.LeftPadBytes(metadataLen, 32))
	data.Write(metadata)
	return &Proposal{
		Source:         msg.Source,
		DepositNonce:   msg.DepositNonce,
		ResourceId:     msg.ResourceId,
		Data:           data.Bytes(),
		HandlerAddress: handlerAddr,
		BridgeAddress:  bridgeAddress,
	}, nil
}

func GenericMessageHandler(msg *relayer.Message, handlerAddr, bridgeAddress common.Address) (*Proposal, error) {
	if len(msg.Payload) != 1 {
		return nil, errors.New("malformed payload. Len  of payload should be 1")
	}
	metadata, ok := msg.Payload[0].([]byte)
	if !ok {
		return nil, errors.New("unable to convert metadata to []byte")
	}
	data := bytes.Buffer{}
	metadataLen := big.NewInt(int64(len(metadata))).Bytes()
	data.Write(common.LeftPadBytes(metadataLen, 32)) // length of metadata (uint256)
	data.Write(metadata)
	return &Proposal{
		Source:         msg.Source,
		DepositNonce:   msg.DepositNonce,
		ResourceId:     msg.ResourceId,
		Data:           data.Bytes(),
		HandlerAddress: handlerAddr,
		BridgeAddress:  bridgeAddress,
	}, nil
}
