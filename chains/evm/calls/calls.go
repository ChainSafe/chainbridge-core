package calls

import (
	"context"
	"encoding/hex"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"math/big"
)

type TxFabric func(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrices []*big.Int, data []byte) (evmclient.CommonTransaction, error)

type ContractChecker interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
}

type ContractCaller interface {
	CallContract(ctx context.Context, callArgs map[string]interface{}, blockNumber *big.Int) ([]byte, error)
}

type GasPricer interface {
	GasPrice() ([]*big.Int, error)
}

type ClientDispatcher interface {
	WaitAndReturnTxReceipt(h common.Hash) (*types.Receipt, error)
	SignAndSendTransaction(ctx context.Context, tx evmclient.CommonTransaction) (common.Hash, error)
	GetTransactionByHash(h common.Hash) (tx *types.Transaction, isPending bool, err error)
	UnsafeNonce() (*big.Int, error)
	LockNonce()
	UnlockNonce()
	UnsafeIncreaseNonce() error
	From() common.Address
}

type ContractCallerDispatcher interface {
	ContractCaller
	ClientDispatcher
	ContractChecker
}

type SimulateCaller interface {
	ContractCaller
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
}

// Simulate function gets transaction info by hash and then executes a message call transaction, which is directly executed in the VM
// of the node, but never mined into the blockchain. Execution happens against provided block.
func Simulate(c SimulateCaller, block *big.Int, txHash common.Hash, from common.Address) ([]byte, error) {
	tx, _, err := c.TransactionByHash(context.TODO(), txHash)
	if err != nil {
		log.Debug().Msgf("[client] tx by hash error: %v", err)
		return nil, err
	}

	log.Debug().Msgf("from: %v to: %v gas: %v gasPrice: %v value: %v data: %v", from, tx.To(), tx.Gas(), tx.GasPrice(), tx.Value(), tx.Data())

	msg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
	res, err := c.CallContract(context.TODO(), ToCallArg(msg), block)
	if err != nil {
		log.Debug().Msgf("[client] call contract error: %v", err)
		return nil, err
	}
	bs, err := hex.DecodeString(common.Bytes2Hex(res))
	if err != nil {
		log.Debug().Msgf("[client] decode string error: %v", err)
		return nil, err
	}
	log.Debug().Msg(string(bs))
	return bs, nil
}
