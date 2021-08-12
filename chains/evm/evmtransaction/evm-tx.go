package evmtransaction

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

type TX struct {
	tx *types.Transaction
}

// RawWithSignature mostly copies WithSignature interface of type.Transaction from go-ethereum,
// but return raw byte representation of transaction to be compatible and interchangeable between different go-ethereum forks
// WithSignature returns a new transaction with the given signature.
// This signature needs to be in the [R || S || V] format where V is 0 or 1.
func (a *TX) RawWithSignature(opts EVMTransactor, key *ecdsa.PrivateKey) ([]byte, error) {
	signer := opts.Signer()
	tx, err := signer(crypto.PubkeyToAddress(key.PublicKey), a.tx)
	if err != nil {
		return nil, err
	}
	a.tx = tx

	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// TODO: change *bind.TransactOpts to wrapper, implement getters for GasFeeCap, GasTipCap, and GasPrice
func NewTransaction(opts EVMTransactor, to common.Address, amount *big.Int, chainId *big.Int, data []byte) *TX {
	var tx *types.Transaction
	log.Info().Msgf("gas fee cap: %v", opts.GasFeeCap())
	log.Info().Msgf("opts: %v", opts)
	if opts.GasFeeCap() != nil {
		log.Debug().Msgf("Using DynamicFeeTx, nonce: %v", opts.Nonce())
		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainId,
			Nonce:     opts.Nonce().Uint64(),
			To:        &to,
			GasFeeCap: opts.GasFeeCap(),
			GasTipCap: opts.GasTipCap(),
			Gas:       opts.GasLimit(),
			Value:     amount,
			Data:      data,
		})
	} else {
		log.Debug().Msgf("Using LegacyTx, nonce: %v", opts.Nonce())
		tx = types.NewTransaction(opts.Nonce().Uint64(), to, amount, opts.GasLimit(), opts.GasPrice(), data)
	}
	return &TX{tx: tx}
}

func (a *TX) Hash() common.Hash {
	return a.tx.Hash()
}
