package celo

import (
	"context"
	"fmt"

	erc20Handler "github.com/ChainSafe/chainbridgev2/bindings/celo/bindings/ERC20Handler"
	"github.com/ChainSafe/chainbridgev2/chains/evm"
	"github.com/ChainSafe/chainbridgev2/chains/evm/listener"
	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/ChainSafe/chainbridgev2/subchains/celo/txtrie"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

func HandleErc20DepositedEventCelo(sourceID, destId uint8, nonce uint64, handlerContractAddress string, backend listener.ChainClient) (relayer.XCMessager, error) {
	contract, err := erc20Handler.NewERC20HandlerCaller(common.HexToAddress(handlerContractAddress), backend)
	if err != nil {
		return nil, err
	}
	record, err := contract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
	if err != nil {
		return nil, err
	}
	m := &CeloMessage{
		Source:       sourceID,
		Destination:  destId,
		Type:         evm.FungibleTransfer,
		DepositNonce: nonce,
		ResourceId:   record.ResourceID,
		Payload: []interface{}{
			record.Amount.Bytes(),
			record.DestinationRecipientAddress,
		},
	}
	blockData, err := backend.BlockByNumber(context.Background(), txBlock)
	if err != nil {
		return nil, err
	}
	trie, err := txtrie.CreateNewTrie(blockData.TxHash(), blockData.Transactions())
	if err != nil {
		return nil, err
	}
	apk, err := l.valsAggr.GetAPKForBlock(txBlock, uint8(l.cfg.ID), l.cfg.EpochSize)
	if err != nil {
		return nil, err

	}
	keyRlp, err := rlp.EncodeToBytes(txIndex)
	if err != nil {
		return nil, fmt.Errorf("encoding TxIndex to rlp: %w", err)
	}
	proof, key, err := txtrie.RetrieveProof(trie, keyRlp)
	if err != nil {
		return nil, err
	}
	m.SVParams = &SignatureVerification{AggregatePublicKey: apk, BlockHash: blockData.Header().Hash(), Signature: blockData.EpochSnarkData().Signature}
	m.MPParams = &MerkleProof{TxRootHash: sliceTo32Bytes(blockData.TxHash().Bytes()), Nodes: proof, Key: key}
	return m, nil
}
