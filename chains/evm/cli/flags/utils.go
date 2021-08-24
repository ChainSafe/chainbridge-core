package flags

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/keystore"

	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func GlobalFlagValues(cmd *cobra.Command) (string, uint64, *big.Int, *secp256k1.Keypair, error) {
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		log.Error().Err(fmt.Errorf("url error: %v", err))
		return "", evmclient.DefaultGasLimit, nil, nil, err
	}

	gasLimitInt, err := cmd.Flags().GetUint64("gasLimit")
	if err != nil {
		log.Error().Err(fmt.Errorf("gas limit error: %v", err))
		return "", evmclient.DefaultGasLimit, nil, nil, err
	}

	gasPriceInt, err := cmd.Flags().GetUint64("gasPrice")
	if err != nil {
		log.Error().Err(fmt.Errorf("gas price error: %v", err))
		return "", evmclient.DefaultGasPrice, nil, nil, err
	}

	gasPrice := big.NewInt(0).SetUint64(gasPriceInt)

	senderKeyPair, err := defineSender(cmd)
	if err != nil {
		log.Error().Err(fmt.Errorf("define sender error: %v", err))
		return "", evmclient.DefaultGasLimit, nil, nil, err
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
