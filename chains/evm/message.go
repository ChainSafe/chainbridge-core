package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/status-im/keycard-go/hexutils"

	"github.com/ethereum/go-ethereum/common"
)

type EVMMessage struct {
	Source       uint8  // Source where message was initiated
	Destination  uint8  // Destination chain of message
	DepositNonce uint64 // Nonce for the deposit
	ResourceId   [32]byte
	Payload      []interface{} // data associated with event sequence
}

func (p EVMMessage) String() string {
	return fmt.Sprintf("evmMesssage from %v to %v nonce %v rID %v", p.Source, p.Destination, p.DepositNonce, hexutils.BytesToHex(p.ResourceId[:]))
}

func (m *EVMMessage) GetSource() uint8 {
	return m.GetSource()
}
func (m *EVMMessage) GetDestination() uint8 {
	return m.GetDestination()
}

func (m *EVMMessage) GetDepositNonce() uint64 {
	return m.GetDepositNonce()
}
func (m *EVMMessage) GetResourceID() [32]byte {
	return m.GetResourceID()
}
func (m *EVMMessage) GetPayload() []interface{} {
	return m.GetPayload()
}
func (m *EVMMessage) CreateProposalDataHash(data []byte) common.Hash {
	return crypto.Keccak256Hash(data)
}
