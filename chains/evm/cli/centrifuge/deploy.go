package centrifuge

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy centrifuge asset store contract",
	Long:  "This command can be used to deploy Centrifuge asset store contract that represents bridged Centrifuge assets.",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return DeployCentrifugeAssetStoreCmd(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
	},
}

func BindDeployCmdFlags(cmd *cobra.Command) {}

func init() {
	BindDeployCmdFlags(deployCmd)
}

func DeployCentrifugeAssetStoreCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	url, _, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return err
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("ethereum client error: %v", err)).Msg("error initializing new EVM client")
		return err
	}

	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})

	assetStoreAddr, err := calls.DeployCentrifugeAssetStore(ethClient, txFabric, gasPricer)
	if err != nil {
		log.Error().Err(fmt.Errorf("Centrifuge asset store deploy failed: %w", err))
		return err
	}

	log.Info().Msgf("Deployed Centrifuge asset store to address: %s", assetStoreAddr.String())
	return nil
}
