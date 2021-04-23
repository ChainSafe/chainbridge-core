package relayer

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Message struct {
	Source       uint8  // Source where message was initiated
	Destination  uint8  // Destination chain of message
	DepositNonce uint64 // Nonce for the deposit
	ResourceId   [32]byte
	Payload      []interface{} // data associated with event sequence
}

type ProposalStatus uint8

const (
	ProposalStatusInactive ProposalStatus = 0
	ProposalStatusActive   ProposalStatus = 1
	ProposalStatusPassed   ProposalStatus = 2 // Ready to be executed
	ProposalStatusExecuted ProposalStatus = 3
	ProposalStatusCanceled ProposalStatus = 4
)

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

func GetIDAndNonce(p *Proposal) *big.Int {
	data := bytes.Buffer{}
	bn := big.NewInt(0).SetUint64(p.DepositNonce).Bytes()
	data.Write(bn)
	data.Write([]byte{p.Source})
	return big.NewInt(0).SetBytes(data.Bytes())
}
