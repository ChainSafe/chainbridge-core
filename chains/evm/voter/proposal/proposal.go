package proposal

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewProposal(source uint8, depositNonce uint64, resourceId [32]byte, data []byte, handlerAddress, bridgeAddress common.Address) *Proposal {
	return &Proposal{
		Source:         source,
		DepositNonce:   depositNonce,
		ResourceId:     resourceId,
		Data:           data,
		HandlerAddress: handlerAddress,
		BridgeAddress:  bridgeAddress,
	}
}

type Proposal struct {
	Source         uint8  // Source domainID where message was initiated
	DepositNonce   uint64 // Nonce for the deposit
	ResourceId     [32]byte
	Payload        []interface{} // data associated with event sequence
	Data           []byte
	HandlerAddress common.Address
	BridgeAddress  common.Address
}

// GetDataHash constructs and returns proposal data hash
func (p *Proposal) GetDataHash() common.Hash {
	return crypto.Keccak256Hash(append(p.HandlerAddress.Bytes(), p.Data...))
}