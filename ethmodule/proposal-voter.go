package ethmodule

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type HandlerFunc func(m *relayer.Message) (*Proposal, error)

func NewMessageHandler(evmCaller EVMClient, bridgeAddress common.Address) *MessageHandler {
	return &MessageHandler{
		bridgeAddress: bridgeAddress,
		evmCaller:     evmCaller,
	}
}

type MessageHandler struct {
	evmCaller     EVMClient
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
	prop, err := handleMessage(m)
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

type Sender interface {
}

type Proposal struct {
	Source         uint8  // Source where message was initiated
	Destination    uint8  // Destination chain of message
	DepositNonce   uint64 // Nonce for the deposit
	ResourceId     [32]byte
	Payload        []interface{} // data associated with event sequence
	Data           []byte
	DataHash       common.Hash
	HandlerAddress common.Address
}

func (p *Proposal) Status() (relayer.ProposalStatus, error) {
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

func (p *Proposal) Execute(sender Sender) error {
	//executeProposal(uint8 chainID, uint64 depositNonce, bytes data, bytes32 resourceID) returns()
	data, err := buildDataUnsafe(
		[]byte("executeProposal(uint8,uint64,bytes,bytes32)"),
		big.NewInt(int64(p.Source)).Bytes(),
		big.NewInt(int64(p.DepositNonce)).Bytes(),
		p.Data,
		p.ResourceId[:])
	if err != nil {
		return err
	}

	err = sender.SignAndSendTransaction(data)
	if err != nil {
		return err
	}
	return nil
}

func (p *Proposal) Vote() error {
	//voteProposal(uint8 chainID, uint64 depositNonce, bytes32 resourceID, bytes32 dataHash)

	dataHash := createProposalDataHash(data, handlerContract, m.MPParams, m.SVParams)

	data, err := buildDataUnsafe(
		[]byte("voteProposal(uint8,uint64,bytes,bytes32,bytes32)"),
		big.NewInt(int64(proposal.Source)).Bytes(),
		big.NewInt(int64(proposal.DepositNonce)).Bytes(),
		proposal.ResourceId[:],
		proposal.DataHash,
	)

	err = r.SignAndSendTransaction(data)
	if err != nil {
		return err
	}
}

func (p *Proposal) VotedBy() bool {

}

func ERC20ProposalHandler(m *relayer.Message, handlerAddr string) (*Proposal, error) {
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

	caddress := common.HexToAddress(handlerAddr)
	return &Proposal{
		Source:         m.Source,
		DepositNonce:   m.DepositNonce,
		ResourceId:     m.ResourceId,
		Data:           data,
		DataHash:       crypto.Keccak256Hash(append(caddress.Bytes(), data...)),
		HandlerAddress: common.HexToAddress(handlerAddr),
	}, nil
}

func ERC721ProposalHandler(msg *relayer.Message, handlerAddr string) (*Proposal, error) {
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
	caddress := common.HexToAddress(handlerAddr)
	return &Proposal{
		Source:         msg.Source,
		DepositNonce:   msg.DepositNonce,
		ResourceId:     msg.ResourceId,
		Data:           data.Bytes(),
		DataHash:       crypto.Keccak256Hash(append(caddress.Bytes(), data.Bytes()...)),
		HandlerAddress: common.HexToAddress(handlerAddr),
	}, nil
}

func GenericProposalHandler(msg *relayer.Message, handlerAddr string) (*Proposal, error) {
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
	caddress := common.HexToAddress(handlerAddr)
	return &Proposal{
		Source:         msg.Source,
		DepositNonce:   msg.DepositNonce,
		ResourceId:     msg.ResourceId,
		Data:           data.Bytes(),
		DataHash:       crypto.Keccak256Hash(append(caddress.Bytes(), data.Bytes()...)),
		HandlerAddress: common.HexToAddress(handlerAddr),
	}, nil
}
