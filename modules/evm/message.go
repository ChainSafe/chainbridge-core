package evm

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TransferType string

var FungibleTransfer TransferType = "FungibleTransfer"
var NonFungibleTransfer TransferType = "NonFungibleTransfer"
var GenericTransfer TransferType = "GenericTransfer"

type DefaultEVMMessage struct {
	Source       uint8        // Source where message was initiated
	Destination  uint8        // Destination chain of message
	Type         TransferType // type of bridge transfer
	DepositNonce uint64       // Nonce for the deposit
	ResourceId   [32]byte
	Payload      []interface{} // data associated with event sequence
}

func (m *DefaultEVMMessage) CreateProposalData() ([]byte, error) {
	var data []byte
	var err error
	switch m.Type {
	case FungibleTransfer:
		data, err = m.createERC20ProposalData()
	case NonFungibleTransfer:
		data, err = m.createErc721ProposalData()
	case GenericTransfer:
		data, err = m.createGenericDepositProposalData()
	default:
		return nil, errors.New(fmt.Sprintf("unknown message type received %s", m.Type))
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *DefaultEVMMessage) createERC20ProposalData() ([]byte, error) {
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
	return b.Bytes(), nil
}

func (m *DefaultEVMMessage) createErc721ProposalData() ([]byte, error) {
	if len(m.Payload) != 3 {
		return nil, errors.New("malformed payload. Len  of payload should be 3")
	}
	tokenID, ok := m.Payload[0].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads tokenID format")
	}
	recipient, ok := m.Payload[1].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads recipient format")
	}
	metadata, ok := m.Payload[2].([]byte)
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
	return data.Bytes(), nil
}

func (m *DefaultEVMMessage) createGenericDepositProposalData() ([]byte, error) {
	if len(m.Payload) != 1 {
		return nil, errors.New("malformed payload. Len  of payload should be 1")
	}
	metadata, ok := m.Payload[0].([]byte)
	if !ok {
		return nil, errors.New("unable to convert metadata to []byte")
	}
	data := bytes.Buffer{}
	metadataLen := big.NewInt(int64(len(metadata))).Bytes()
	data.Write(common.LeftPadBytes(metadataLen, 32)) // length of metadata (uint256)
	data.Write(metadata)
	return data.Bytes(), nil
}
