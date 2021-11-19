package calls

import (
	"context"
	"encoding/hex"
	"errors"
	gomath "math"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

type TxFabric func(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrices []*big.Int, data []byte) (evmclient.CommonTransaction, error)

type ClientContractChecker interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
}

type ContractCallerClient interface {
	CallContract(ctx context.Context, callArgs map[string]interface{}, blockNumber *big.Int) ([]byte, error)
}

type ContractCheckerCallerClient interface {
	ContractCallerClient
	ClientContractChecker
}

type ClientDeployer interface {
	ClientDispatcher
	ClientContractChecker
}

type ContractCallerDispatcherClient interface {
	ContractCallerClient
	ClientDispatcher
}

type GasPricer interface {
	GasPrice() ([]*big.Int, error)
}

type ClientDispatcher interface {
	WaitAndReturnTxReceipt(h common.Hash) (*types.Receipt, error)
	SignAndSendTransaction(ctx context.Context, tx evmclient.CommonTransaction) (common.Hash, error)
	UnsafeNonce() (*big.Int, error)
	LockNonce()
	UnlockNonce()
	UnsafeIncreaseNonce() error
	From() common.Address
}

type SimulateCallerClient interface {
	ContractCallerClient
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
}

func GetSolidityFunctionSig(in []byte) [4]byte {
	var res [4]byte
	hash := crypto.Keccak256Hash(in)
	copy(res[:], hash[:])
	return res
}

func SliceTo32Bytes(in []byte) [32]byte {
	var res [32]byte
	copy(res[:], in)
	return res
}

// ToCallArg is the function that converts ethereum.CallMsg into more abstract map
// This is done for matter of  making EVMClient more abstract since some go-ethereum forks uses different messages types
func ToCallArg(msg ethereum.CallMsg) map[string]interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}

// UserAmountToWei converts decimal user friendly representation of token amount to 'Wei' representation with provided amount of decimal places
// eg UserAmountToWei(1, 5) => 100000
func UserAmountToWei(amount string, decimal *big.Int) (*big.Int, error) {
	amountFloat, ok := big.NewFloat(0).SetString(amount)
	if !ok {
		return nil, errors.New("wrong amount format")
	}
	ethValueFloat := new(big.Float).Mul(amountFloat, big.NewFloat(gomath.Pow10(int(decimal.Int64()))))
	ethValueFloatString := strings.Split(ethValueFloat.Text('f', int(decimal.Int64())), ".")

	i, ok := big.NewInt(0).SetString(ethValueFloatString[0], 10)
	if !ok {
		return nil, errors.New(ethValueFloat.Text('f', int(decimal.Int64())))
	}

	return i, nil
}

func Transact(client ClientDispatcher, txFabric TxFabric, gasPriceClient GasPricer, to *common.Address, data []byte, gasLimit uint64, value *big.Int) (common.Hash, error) {
	defer client.UnlockNonce()
	client.LockNonce()
	n, err := client.UnsafeNonce()
	if err != nil {
		return common.Hash{}, nil
	}
	gp, err := gasPriceClient.GasPrice()
	if err != nil {
		return common.Hash{}, err
	}
	tx, err := txFabric(n.Uint64(), to, value, gasLimit, gp, data)
	if err != nil {
		return common.Hash{}, err
	}
	_, err = client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return common.Hash{}, err
	}
	log.Debug().Msgf("hash: %v from: %s", tx.Hash(), client.From())
	_, err = client.WaitAndReturnTxReceipt(tx.Hash())
	if err != nil {
		return common.Hash{}, err
	}
	err = client.UnsafeIncreaseNonce()
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

// Simulate function gets transaction info by hash and then executes a message call transaction, which is directly executed in the VM
// of the node, but never mined into the blockchain. Execution happens against provided block.
func Simulate(c SimulateCallerClient, block *big.Int, txHash common.Hash, from common.Address) ([]byte, error) {
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
