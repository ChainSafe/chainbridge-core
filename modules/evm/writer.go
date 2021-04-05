package evm

import (
	"context"
	"math/big"

	"github.com/ChainSafe/chainbridgev2/bindings/eth/bindings/Bridge"
	goeth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ContractCaller interface {
	FilterLogs(ctx context.Context, q goeth.FilterQuery) ([]types.Log, error)
	LatestBlock() (*big.Int, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	CallOpts() *bind.CallOpts
	Opts() *bind.TransactOpts
	LockAndUpdateOpts() error
	UnlockOpts()
	WaitForBlock(block *big.Int) error
}

type Bridger interface {
	GetProposal(opts *bind.CallOpts, originChainID uint8, depositNonce uint64, dataHash [32]byte) (Bridge.BridgeProposal, error)
	HasVotedOnProposal(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte, arg2 common.Address) (bool, error)
	VoteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, resourceID [32]byte, dataHash [32]byte) (*types.Transaction, error)
	ExecuteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, data []byte, resourceID [32]byte, signatureHeader []byte, aggregatePublicKey []byte, hashedMessage [32]byte, rootHash [32]byte, key []byte, nodes []byte) (*types.Transaction, error)
}

type EVMWriter struct {
	chainID uint8
	client  ContractCaller
	stop    <-chan struct{}
	sysErr  chan<- error
	bridge  Bridger
}

// NewWriter creates and returns writer
func NewWriter(client ContractCaller, stop <-chan struct{}, sysErr chan<- error) *EVMWriter {
	return &EVMWriter{
		client: client,
		stop:   stop,
		sysErr: sysErr,
	}
}
