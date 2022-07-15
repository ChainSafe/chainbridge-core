package flags

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ChainSafe/sygma-core/chains/evm/calls"

	"github.com/ChainSafe/sygma-core/keystore"
	"github.com/ChainSafe/sygma-core/types"

	"github.com/ChainSafe/sygma-core/crypto/secp256k1"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const DefaultGasLimit = 2000000

func GlobalFlagValues(cmd *cobra.Command) (string, uint64, *big.Int, *secp256k1.Keypair, bool, error) {
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		log.Error().Err(fmt.Errorf("url error: %v", err))
		return "", DefaultGasLimit, nil, nil, false, err
	}

	gasLimitInt, err := cmd.Flags().GetUint64("gas-limit")
	if err != nil {
		log.Error().Err(fmt.Errorf("gas limit error: %v", err))
		return "", DefaultGasLimit, nil, nil, false, err
	}

	gasPriceInt, err := cmd.Flags().GetUint64("gas-price")
	if err != nil {
		log.Error().Err(fmt.Errorf("gas price error: %v", err))
		return "", DefaultGasLimit, nil, nil, false, err
	}
	var gasPrice *big.Int = nil
	if gasPriceInt != 0 {
		gasPrice = big.NewInt(0).SetUint64(gasPriceInt)
	}

	senderKeyPair, err := defineSender(cmd)
	if err != nil {
		log.Error().Err(fmt.Errorf("define sender error: %v", err))
		return "", DefaultGasLimit, nil, nil, false, err
	}

	prepare, err := cmd.Flags().GetBool("prepare")
	if err != nil {
		log.Error().Err(fmt.Errorf("generate calldata error: %v", err))
		return "", DefaultGasLimit, nil, nil, false, err
	}
	return url, gasLimitInt, gasPrice, senderKeyPair, prepare, nil
}

func defineSender(cmd *cobra.Command) (*secp256k1.Keypair, error) {
	privateKey, err := cmd.Flags().GetString("private-key")
	if err != nil {
		return nil, err
	}
	if privateKey != "" {
		kp, err := secp256k1.NewKeypairFromString(privateKey)
		if err != nil {
			return nil, err
		}
		return kp, nil
	}
	var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
	return AliceKp, nil
}

func ProcessResourceID(resourceID string) (types.ResourceID, error) {
	if resourceID[0:2] == "0x" {
		resourceID = resourceID[2:]
	}
	resourceIdBytes, err := hex.DecodeString(resourceID)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed decoding resourceID hex string: %s", err)
	}
	return calls.SliceTo32Bytes(resourceIdBytes), nil
}

func MarkFlagsAsRequired(cmd *cobra.Command, flags ...string) {
	for _, flag := range flags {
		err := cmd.MarkFlagRequired(flag)
		if err != nil {
			panic(err)
		}
	}
}
