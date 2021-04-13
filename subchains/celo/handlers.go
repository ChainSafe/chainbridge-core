package celo

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	erc20Handler "github.com/ChainSafe/chainbridgev2/bindings/celo/bindings/ERC20Handler"
	"github.com/ChainSafe/chainbridgev2/chains/evmd"
	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/ChainSafe/chainbridgev2/subchains/celo/txtrie"
	goeth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func handleCeloProofs(m *CeloMessage, cr goeth.ChainReader, txBlock *big.Int, txIndex uint) error {
	blockData, err := cr.BlockByNumber(context.Background(), txBlock)
	if err != nil {
		return err
	}
	trie, err := txtrie.CreateNewTrie(blockData.TxHash(), blockData.Transactions())
	if err != nil {
		return err
	}
	apk, err := l.valsAggr.GetAPKForBlock(txBlock, uint8(l.cfg.ID), l.cfg.EpochSize)
	if err != nil {
		return err

	}
	keyRlp, err := rlp.EncodeToBytes(txIndex)
	if err != nil {
		return fmt.Errorf("encoding TxIndex to rlp: %w", err)
	}
	proof, key, err := txtrie.RetrieveProof(trie, keyRlp)
	if err != nil {
		return err
	}
	m.SVParams = &SignatureVerification{AggregatePublicKey: apk, BlockHash: blockData.Header().Hash(), Signature: blockData.EpochSnarkData().Signature}
	m.MPParams = &MerkleProof{TxRootHash: sliceTo32Bytes(blockData.TxHash().Bytes()), Nodes: proof, Key: key}
	return nil
}

func HandleErc20DepositedEvent(sourceID, destId uint8, nonce uint64, handlerContractAddress common.Address, backend bind.ContractBackend) (relayer.XCMessager, error) {
	contract, err := erc20Handler.NewERC20Handler(handlerContractAddress, backend)
	if err != nil {
		return nil, err
	}
	record, err := contract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
	if err != nil {
		return nil, err
	}

	return &evmd.DefaultEVMMessage{
		Source:       sourceID,
		Destination:  destId,
		Type:         FungibleTransfer,
		DepositNonce: nonce,
		ResourceId:   record.ResourceID,
		Payload: []interface{}{
			record.Amount.Bytes(),
			record.DestinationRecipientAddress,
		},
	}, nil
}

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
	Source       uint8                // Source where message was initiated
	Destination  uint8                // Destination chain of message
	Type         relayer.TransferType // type of bridge transfer
	DepositNonce uint64               // Nonce for the deposit
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
