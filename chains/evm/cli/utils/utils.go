package utils

import (
	"errors"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmgaspricer"
	gomath "math"
	"math/big"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var UtilsCmd = &cobra.Command{
	Use:   "utils",
	Short: "Set of utility commands",
	Long:  "Set of utility commands",
}

func init() {
	UtilsCmd.AddCommand(simulateCmd)
	UtilsCmd.AddCommand(hashListCmd)
}

type EventSig string

func (es EventSig) GetTopic() common.Hash {
	return crypto.Keccak256Hash([]byte(es))
}

//
//func IsActive(status uint8) bool {
//	return ProposalStatus(status) == Active
//}
//
//func IsPassed(status uint8) bool {
//	return ProposalStatus(status) == Passed
//}
//
//func IsExecuted(status uint8) bool {
//	return ProposalStatus(status) == Executed
//}

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

func WeiAmountToUser(amount *big.Int, decimals *big.Int) (*big.Float, error) {
	amountFloat, ok := big.NewFloat(0).SetString(amount.String())
	if !ok {
		return nil, errors.New("wrong amount format")
	}
	return new(big.Float).Quo(amountFloat, big.NewFloat(gomath.Pow10(int(decimals.Int64())))), nil
}

type GasPricerWithPostConfig interface {
	calls.GasPricer
	SetClient(client evmgaspricer.LondonGasClient)
	SetOpts(opts *evmgaspricer.GasPricerOpts)
}
