package evmtypes

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// TxFabric is the function type represents abstract fabric that produce Transactions to implement CommonTransaction interface for different evm based chains
type TxFabric func(chainId *big.Int, nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPricer GasPricer, data []byte) (CommonTransaction, error)

// CommonTransaction is the abstract representation of transaction that could be processed by the relayer
type CommonTransaction interface {
	// Hash returns the transaction hash.
	Hash() common.Hash

	// RawWithSignature Returns signed transaction by provided private key
	RawWithSignature(key *ecdsa.PrivateKey, chainID *big.Int) ([]byte, error)
}
