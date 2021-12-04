package flags

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"

	"github.com/ChainSafe/chainbridge-core/keystore"
	"github.com/ChainSafe/chainbridge-core/types"

	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func GlobalFlagValues(cmd *cobra.Command) (string, uint64, *big.Int, *secp256k1.Keypair, error) {
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		log.Error().Err(fmt.Errorf("url error: %v", err))
		return "", consts.DefaultGasLimit, nil, nil, err
	}

	gasLimitInt, err := cmd.Flags().GetUint64("gasLimit")
	if err != nil {
		log.Error().Err(fmt.Errorf("gas limit error: %v", err))
		return "", consts.DefaultGasLimit, nil, nil, err
	}

	gasPriceInt, err := cmd.Flags().GetUint64("gasPrice")
	if err != nil {
		log.Error().Err(fmt.Errorf("gas price error: %v", err))
		return "", consts.DefaultGasLimit, nil, nil, err
	}
	var gasPrice *big.Int = nil
	if gasPriceInt != 0 {
		gasPrice = big.NewInt(0).SetUint64(gasPriceInt)
	}

	senderKeyPair, err := defineSender(cmd)
	if err != nil {
		log.Error().Err(fmt.Errorf("define sender error: %v", err))
		return "", consts.DefaultGasLimit, nil, nil, err
	}
	return url, gasLimitInt, gasPrice, senderKeyPair, nil
}

func defineSender(cmd *cobra.Command) (*secp256k1.Keypair, error) {
	privateKey, err := cmd.Flags().GetString("privateKey")
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
		return [32]byte{}, err
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
