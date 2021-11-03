package admin

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var setThresholdCmd = &cobra.Command{
	Use:   "set-threshold",
	Short: "Set a new relayer vote threshold",
	Long:  "Set a new relayer vote threshold",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	Run: setThreshold,
}

func init() {
	setThresholdCmd.Flags().Uint64Var(&RelayerThreshold, "threshold", 0, "new relayer threshold")
	setThresholdCmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	flags.MarkFlagsAsRequired(setThresholdCmd, "threshold", "bridge")
}

func setThreshold(cmd *cobra.Command, args []string) {
	log.Debug().Msgf(`
Setting new threshold
Threshold: %d
Bridge address: %s`, RelayerThreshold, Bridge)
}

/*

func setThreshold(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	bridgeAddress, err := cliutils.DefineBridgeAddress(cctx)
	if err != nil {
		return err
	}
	threshold := cctx.Uint64("threshold")
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.AdminSetThreshHold(ethClient, bridgeAddress, big.NewInt(0).SetUint64(threshold))
	if err != nil {
		return err
	}
	log.Info().Msgf("New threshold set for %v", threshold)
	return nil
}
*/
