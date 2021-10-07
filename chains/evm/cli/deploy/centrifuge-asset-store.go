package deploy

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var deployCentrifugeAssetStoreCmd = &cobra.Command{
	Use:   "centrifuge-asset-store",
	Short: "Deploy centrifuge asset store contract",
	Long:  "This command can be used to deploy Centrifuge asset store contract that represents bridged Centrifuge assets.",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return DeployCentrifugeAssetStoreCmd(cmd, args, txFabric)
	},
}

func DeployCentrifugeAssetStoreCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return err
	}

	log.Debug().Msgf("url: %s gas limit: %v gas price: %v", url, gasLimit, gasPrice)

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("ethereum client error: %v", err)).Msg("error initializing new EVM client")
		return err
	}

	assetStoreAddr, err := calls.DeployCentrifugeAssetStore(ethClient, txFabric)
	if err != nil {
		log.Error().Err(fmt.Errorf("Centrifuge asset store deploy failed: %w", err))
	}

	log.Info().Msgf("Deployed Centrifuge asset store to address: %s", assetStoreAddr.String())
	return nil
}
