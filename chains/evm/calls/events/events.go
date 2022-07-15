package events

import (
	"github.com/ChainSafe/sygma-core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type EventSig string

func (es EventSig) GetTopic() common.Hash {
	return crypto.Keccak256Hash([]byte(es))
}

const (
	DepositSig          EventSig = "Deposit(uint8,bytes32,uint64,address,bytes,bytes)"
	ThresholdChangedSig EventSig = "RelayerThresholdChanged(uint256)"
	ProposalEventSig    EventSig = "ProposalEvent(uint8,uint64,uint8,bytes32)"
	ProposalVoteSig     EventSig = "ProposalVote(uint8,uint64,uint8,bytes32)"
)

// Deposit struct holds event data with all necessary parameters and a handler response
// https://github.com/ChainSafe/chainbridge-solidity/blob/develop/contracts/Bridge.sol#L47
type Deposit struct {
	// ID of chain deposit will be bridged to
	DestinationDomainID uint8
	// ResourceID used to find address of handler to be used for deposit
	ResourceID types.ResourceID
	// Nonce of deposit
	DepositNonce uint64
	// Address of sender (msg.sender: user)
	SenderAddress common.Address
	// Additional data to be passed to specified handler
	Data []byte
	// ERC20Handler: responds with empty data
	// ERC721Handler: responds with deposited token metadata acquired by calling a tokenURI method in the token contract
	// GenericHandler: responds with the raw bytes returned from the call to the target contract
	HandlerResponse []byte
}
