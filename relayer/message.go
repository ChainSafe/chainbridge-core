package relayer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// XCMessage is used as a generic format cross-chain communications
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

type Proposal interface {
	XCMessager
	GetProposalData() []byte
	GetProposalDataHash() common.Hash
	GetIDAndNonce() *big.Int
}
