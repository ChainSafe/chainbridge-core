package account

import (
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ChainSafe/sygma-core/crypto/secp256k1"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var generateKeyPairCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a bridge keystore (Secp256k1)",
	Long:  "The generate subcommand is used to generate the bridge keystore. If no options are specified, a Secp256k1 key will be made",
	RunE:  generateKeyPair,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
}

func generateKeyPair(cmd *cobra.Command, args []string) error {
	kp, err := secp256k1.GenerateKeypair()
	if err != nil {
		return err
	}
	log.Debug().Msgf("Address: %s,  Private key: %x", kp.CommonAddress().String(), kp.Encode())
	return nil
}
