package evm

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type EVMProposal struct {
	Source         uint8  // Source where message was initiated
	Destination    uint8  // Destination chain of message
	DepositNonce   uint64 // Nonce for the deposit
	ResourceId     [32]byte
	Payload        []interface{} // data associated with event sequence
	data           []byte
	dataHash       common.Hash
	handlerAddress common.Address
}

func (p EVMProposal) String() string {
	return fmt.Sprintf("evmProposal from %v to %v nonce %v rID %v", p.Source, p.Destination, p.DepositNonce, hexutils.BytesToHex(p.ResourceId[:]))
}

func (p *EVMProposal) GetSource() uint8 {
	return p.GetSource()
}
func (p *EVMProposal) GetDestination() uint8 {
	return p.GetDestination()
}

func (p *EVMProposal) GetDepositNonce() uint64 {
	return p.GetDepositNonce()
}
func (p *EVMProposal) GetResourceID() [32]byte {
	return p.GetResourceID()
}
func (p *EVMProposal) GetPayload() []interface{} {
	return p.GetPayload()
}
func (p *EVMProposal) CreateProposalDataHash(data []byte) common.Hash {
	return crypto.Keccak256Hash(data)
}

func (p *EVMProposal) GetProposalData() []byte {
	return p.data
}
func (p *EVMProposal) GetProposalDataHash() common.Hash {
	return crypto.Keccak256Hash(append(p.handlerAddress.Bytes(), p.data...))
}

func (p *EVMProposal) ShouldBeVotedFor() bool {
	return false
}
func (p *EVMProposal) ProposalIsComplete() bool {
	return false
}
func (p *EVMProposal) ProposalIsReadyForExecute() bool {
	return false
}

func ERC20ProposalHandler(m relayer.XCMessager, handlerAddr string) (relayer.Proposal, error) {
	if len(m.GetPayload()) != 2 {
		return nil, errors.New("malformed payload. Len  of payload should be 2")
	}
	amount, ok := m.GetPayload()[0].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads amount format")
	}

	recipient, ok := m.GetPayload()[1].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads recipient format")

	}
	b := bytes.Buffer{}
	b.Write(common.LeftPadBytes(amount, 32)) // amount (uint256)
	recipientLen := big.NewInt(int64(len(recipient))).Bytes()
	b.Write(common.LeftPadBytes(recipientLen, 32))
	b.Write(recipient)
	return &EVMProposal{
		data:           b.Bytes(),
		handlerAddress: common.HexToAddress(handlerAddr),
	}, nil
}

func ERC721ProposalHandler(msg relayer.XCMessager, handlerAddr string) (relayer.Proposal, error) {
	if len(msg.GetPayload()) != 3 {
		return nil, errors.New("malformed payload. Len  of payload should be 3")
	}
	tokenID, ok := msg.GetPayload()[0].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads tokenID format")
	}
	recipient, ok := msg.GetPayload()[1].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads recipient format")
	}
	metadata, ok := msg.GetPayload()[2].([]byte)
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
	return &EVMProposal{
		data:           data.Bytes(),
		handlerAddress: common.HexToAddress(handlerAddr),
	}, nil
}

func GenericProposalHandler(msg relayer.XCMessager, handlerAddr string) (relayer.Proposal, error) {
	if len(msg.GetPayload()) != 1 {
		return nil, errors.New("malformed payload. Len  of payload should be 1")
	}
	metadata, ok := msg.GetPayload()[0].([]byte)
	if !ok {
		return nil, errors.New("unable to convert metadata to []byte")
	}
	data := bytes.Buffer{}
	metadataLen := big.NewInt(int64(len(metadata))).Bytes()
	data.Write(common.LeftPadBytes(metadataLen, 32)) // length of metadata (uint256)
	data.Write(metadata)
	return &EVMProposal{
		data:           data.Bytes(),
		handlerAddress: common.HexToAddress(handlerAddr),
	}, nil
}
