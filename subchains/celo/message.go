// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package celo

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type MerkleProof struct {
	TxRootHash [32]byte // Expected root of trie, in our case should be transactionsRoot from block
	Key        []byte   // RLP encoding of tx index, for the tx we want to prove
	Nodes      []byte   // The actual proof, all the nodes of the trie that between leaf value and root
}

type SignatureVerification struct {
	AggregatePublicKey []byte      // Aggregated public key of block validators
	BlockHash          common.Hash // Hash of block we are proving
	Signature          []byte      // Signature of block we are proving
}

type CeloMessage struct {
	Source       uint8  // Source where message was initiated
	Destination  uint8  // Destination chain of message
	Type         string // type of bridge transfer
	DepositNonce uint64 // Nonce for the deposit
	ResourceId   [32]byte
	Payload      []interface{} // data associated with event sequence
	MPParams     *MerkleProof
	SVParams     *SignatureVerification
}

func (m *CeloMessage) GetSource() uint8 {
	return m.GetSource()
}
func (m *CeloMessage) GetDestination() uint8 {
	return m.GetDestination()
}
func (m *CeloMessage) GetType() string {
	return m.GetType()
}
func (m *CeloMessage) GetDepositNonce() uint64 {
	return m.GetDepositNonce()
}
func (m *CeloMessage) GetResourceID() [32]byte {
	return m.GetResourceID()
}
func (m *CeloMessage) GetPayload() []interface{} {
	return m.GetPayload()
}
func (m *CeloMessage) CreateProposalDataHash(data []byte) common.Hash {
	return crypto.Keccak256Hash(data)
}

func (m *CeloMessage) CreateProposalData() ([]byte, error) {
	var data []byte
	var err error
	switch m.Type {
	case relayer.FungibleTransfer:
		data, err = m.createERC20ProposalData()
	case relayer.NonFungibleTransfer:
		data, err = m.createErc721ProposalData()
	case relayer.GenericTransfer:
		data, err = m.createGenericDepositProposalData()
	default:
		return nil, errors.New(fmt.Sprintf("unknown message type received %s", m.Type))
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

func sliceTo32Bytes(in []byte) [32]byte {
	var res [32]byte
	copy(res[:], in)
	return res
}
