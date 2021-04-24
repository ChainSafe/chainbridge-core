package writer

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func ERC20ProposalHandler(m *relayer.Message, handlerAddr string) (*relayer.Proposal, error) {
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
	b := bytes.Buffer{}
	b.Write(common.LeftPadBytes(amount, 32)) // amount (uint256)
	recipientLen := big.NewInt(int64(len(recipient))).Bytes()
	b.Write(common.LeftPadBytes(recipientLen, 32))
	b.Write(recipient)
	caddress := common.HexToAddress(handlerAddr)
	return &relayer.Proposal{
		Data:           b.Bytes(),
		DataHash:       crypto.Keccak256Hash(append(caddress.Bytes(), b.Bytes()...)),
		HandlerAddress: common.HexToAddress(handlerAddr),
	}, nil
}

func ERC721ProposalHandler(msg *relayer.Message, handlerAddr string) (*relayer.Proposal, error) {
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
	return &relayer.Proposal{
		Data:           data.Bytes(),
		DataHash:       crypto.Keccak256Hash(append(caddress.Bytes(), data.Bytes()...)),
		HandlerAddress: common.HexToAddress(handlerAddr),
	}, nil
}

func GenericProposalHandler(msg *relayer.Message, handlerAddr string) (*relayer.Proposal, error) {
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
	return &relayer.Proposal{
		Data:           data.Bytes(),
		DataHash:       crypto.Keccak256Hash(append(caddress.Bytes(), data.Bytes()...)),
		HandlerAddress: common.HexToAddress(handlerAddr),
	}, nil
}
