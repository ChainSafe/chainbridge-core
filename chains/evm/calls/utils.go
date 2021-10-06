package calls

import (
	"context"
	"errors"
	gomath "math"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)


type TxFabric func(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrices []*big.Int, data []byte) (evmclient.CommonTransaction, error)

type ChainClient interface {
	SignAndSendTransaction(ctx context.Context, tx evmclient.CommonTransaction) (common.Hash, error)
	CallContract(ctx context.Context, callArgs map[string]interface{}, blockNumber *big.Int) ([]byte, error)
	WaitAndReturnTxReceipt(h common.Hash) (*types.Receipt, error)
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	UnsafeNonce() (*big.Int, error)
	LockNonce()
	UnlockNonce()
	UnsafeIncreaseNonce() error
	// GasPrices method returns array of gasPRices to be compatibale with pre- and post- London fork
	// if array size is bigger than 1, then it must contians maxTipFee and maxCapFee for post London Fork chains,
	// otherwise it contains regular gasPrice as first element
	GasPrices() []*big.Int
	From() common.Address
	ChainID(ctx context.Context) (*big.Int, error)
	Simulate(block *big.Int, txHash common.Hash, fromAddress common.Address) ([]byte, error)
}

func SliceTo32Bytes(in []byte) [32]byte {
	var res [32]byte
	copy(res[:], in)
	return res
}

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

func Transact(client ChainClient, txFabric TxFabric, to *common.Address, data []byte, gasLimit uint64) (common.Hash, error) {
	client.LockNonce()
	n, err := client.UnsafeNonce()
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := txFabric(n.Uint64(), to, big.NewInt(0), gasLimit, client.GasPrices(), data)
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
	client.UnlockNonce()
	return tx.Hash(), nil
}

func ConstructErc20DepositData(destRecipient []byte, amount *big.Int) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(amount, 32)...)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...)
	data = append(data, destRecipient...)
	return data
}
