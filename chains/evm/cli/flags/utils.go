package flags

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/cliutils"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Globals(cmd *cobra.Command) (string, *big.Int, *big.Int, *secp256k1.Keypair, error) {
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		log.Error().Err(fmt.Errorf("url error: %v", err))
		return "", nil, nil, nil, err
	}

	gasLimitInt, err := cmd.Flags().GetUint64("gasLimit")
	if err != nil {
		log.Error().Err(fmt.Errorf("gas limit error: %v", err))
		return "", nil, nil, nil, err
	}

	gasLimit := big.NewInt(0).SetUint64(gasLimitInt)

	gasPriceInt, err := cmd.Flags().GetUint64("gasPrice")
	if err != nil {
		log.Error().Err(fmt.Errorf("gas price error: %v", err))
		return "", nil, nil, nil, err
	}

	gasPrice := big.NewInt(0).SetUint64(gasPriceInt)

	senderKeyPair, err := cliutils.DefineSender(cmd)
	if err != nil {
		log.Error().Err(fmt.Errorf("define sender error: %v", err))
		return "", nil, nil, nil, err
	}

	return url, gasLimit, gasPrice, senderKeyPair, nil
}
