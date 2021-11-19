// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package message

import (
	"errors"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/types"
)

type TransferType string

const (
	FungibleTransfer    TransferType = "FungibleTransfer"
	NonFungibleTransfer TransferType = "NonFungibleTransfer"
	GenericTransfer     TransferType = "GenericTransfer"
)

type ProposalStatus struct {
	Status        uint8
	YesVotes      *big.Int
	YesVotesTotal uint8
	ProposedBlock *big.Int
}

const (
	ProposalStatusInactive uint8 = iota
	ProposalStatusActive
	ProposalStatusPassed // Ready to be executed
	ProposalStatusExecuted
	ProposalStatusCanceled
)

var (
	StatusMap = map[uint8]string{ProposalStatusInactive: "inactive", ProposalStatusActive: "active", ProposalStatusPassed: "passed", ProposalStatusExecuted: "executed", ProposalStatusCanceled: "canceled"}
)

type Message struct {
	Source       uint8  // Source where message was initiated
	Destination  uint8  // Destination chain of message
	DepositNonce uint64 // Nonce for the deposit
	ResourceId   types.ResourceID
	Payload      []interface{} // data associated with event sequence
	Type         TransferType
}

// extractAmountTransferred is a private method to extract and transform the transfer amount
// from the Payload field within the Message struct
func (m *Message) extractAmountTransferred() (float64, error) {
	// parse payload field from event log message to obtain transfer amount
	// payload slice of interfaces includes..
	// index 0: amount ([]byte)
	// index 1: destination recipient address ([]byte)

	// declare new float64 as return value
	var payloadAmountFloat float64

	// cast interface to byte slice
	amountByteSlice, ok := m.Payload[0].([]byte)
	if !ok {
		err := errors.New("could not cast interface to byte slice")
		return payloadAmountFloat, err
	}

	// convert big int => float64
	// ignore accuracy (rounding)
	payloadAmountFloat, _ = new(big.Float).SetInt(big.NewInt(0).SetBytes(amountByteSlice)).Float64()

	return payloadAmountFloat, nil
}
