package utils

import (
	"errors"
	gomath "math"
	"math/big"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

var UtilsCmd = &cobra.Command{
	Use:   "utils",
	Short: "Utils-related instructions",
	Long:  "Utils-related instructions",
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

func ConstructErc20DepositData(destRecipient []byte, amount *big.Int) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(amount, 32)...)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...)
	data = append(data, destRecipient...)
	return data
}

// constructErc20Data constructs the data field to be passed into an erc721 deposit call
func ConstructErc721DepositData(tokenId *big.Int, destRecipient []byte) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(tokenId, 32)...)                               // Resource Id + Token Id
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...) // Length of recipient
	data = append(data, destRecipient...)                                                  // Recipient

	return data
}

func ConstructGenericDepositData(metadata []byte) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 32)...)
	data = append(data, metadata...)
	return data
}
