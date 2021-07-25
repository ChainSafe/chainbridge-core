package cliutils

import (
	"bytes"
	"errors"
	"fmt"
	gomath "math"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-core/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
	BobKp   = keystore.TestKeyRing.EthereumKeys[keystore.BobKey]
	EveKp   = keystore.TestKeyRing.EthereumKeys[keystore.EveKey]

	DefaultRelayerAddresses = []common.Address{
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.AliceKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.BobKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.CharlieKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.DaveKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.EveKey].Address()),
	}
)

type EventSig string

func (es EventSig) GetTopic() common.Hash {
	return crypto.Keccak256Hash([]byte(es))
}

func IsActive(status uint8) bool {
	return ProposalStatus(status) == Active
}

func IsPassed(status uint8) bool {
	return ProposalStatus(status) == Passed
}

func IsExecuted(status uint8) bool {
	return ProposalStatus(status) == Executed
}

func SliceTo32Bytes(in []byte) [32]byte {
	var res [32]byte
	copy(res[:], in)
	return res
}

func SliceTo4Bytes(in []byte) [4]byte {
	var res [4]byte
	copy(res[:], in)
	return res
}

func GetSolidityFunctionSig(in string) [4]byte {
	res := crypto.Keccak256(bytes.NewBufferString(in).Bytes())
	return SliceTo4Bytes(res)
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

// RlpEncodeHeader is method to RLP encode data stored in a block header
func RlpEncodeHeader(header *types.Header) ([]byte, error) {
	// deep copy of header
	newHeader := types.CopyHeader(header)

	// encode copied header into byte slice
	rlpEncodedHeader, err := rlp.EncodeToBytes(newHeader)
	if err != nil {
		// return empty byte slice, error
		return []byte{}, fmt.Errorf("error encoding header: %w", err)
	}

	return rlpEncodedHeader, nil
}

func DefineSender(privateKey string) (*secp256k1.Keypair, error) {
	if privateKey != "" {
		kp, err := secp256k1.NewKeypairFromString(privateKey)
		if err != nil {
			return nil, err
		}
		return kp, nil
	}
	return AliceKp, nil
}

func DefineBridgeAddress(bridgeAddress string) (common.Address, error) {
	if !common.IsHexAddress(bridgeAddress) {
		return common.Address{}, fmt.Errorf("invalid bridge address %s", bridgeAddress)
	}
	return common.HexToAddress(bridgeAddress), nil
}
