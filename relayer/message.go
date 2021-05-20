// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package relayer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// XCMessage is used as a generic format cross-chain communications
// GenericMessage?
// TODO: Check change to struct bcs maybe we are not gonna need it
type XCMessager interface {
	GetSource() uint8
	GetDestination() uint8
	GetDepositNonce() uint64
	GetResourceID() [32]byte
	GetPayload() []interface{} // Maybe this should be some bytes encoding
	String() string
}

type ProposalStatus uint8

const (
	ProposalStatusInactive ProposalStatus = 0
	ProposalStatusActive   ProposalStatus = 1
	ProposalStatusPassed   ProposalStatus = 2 // Ready to be executed
	ProposalStatusExecuted ProposalStatus = 3
	ProposalStatusCanceled ProposalStatus = 4
)

// TODO: check this could be rewriten as struct
type Proposal interface {
	XCMessager
	GetProposalData() []byte
	GetProposalDataHash() common.Hash
	GetIDAndNonce() *big.Int
}
