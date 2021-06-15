package evmsender

import (
	"context"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type CommonTransaction interface {
	// Hash returns the transaction hash.
	Hash() common.Hash
	// RawWithSignature mostly copies WithSignature interface of type.Transaction from go-ethereum,
	// but return raw rlp encoded signed transaction to be compatible and interchangeable between different go-ethereum implementations
	RawWithSignature(signer types.Signer, sig []byte) ([]byte, error)
}

type EVMSender struct {
	client evmclient.EVMClient
}

func (s *EVMSender) From() common.Address {
	return s.sender.From()
}

func (s *EVMSender) SignAndSendTransaction(tx CommonTransaction) (common.Hash, error) {
	h := tx.Hash()
	sig, err := crypto.Sign(h[:], prvKey)
	if err != nil {
		return common.Hash{}, err
	}

	rawTX, err := tx.RawWithSignature(types.HomesteadSigner{}, sig)
	if err != nil {
		return common.Hash{}, err
	}

	err = s.client.SendRawTransaction(context.TODO(), rawTX)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}
