package admin

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause deposits and proposals",
	Long:  "Pause deposits and proposals",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	Run: pause,
}

func BindPauseFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	flags.MarkFlagsAsRequired(cmd, "bridge")
}
func init() {
	BindPauseFlags(pauseCmd)
}

func pause(cmd *cobra.Command, args []string) {
	log.Debug().Msgf(`
Pausing
Bridge address: %s`, Bridge)
}

/*
func pause(cctx *cli.Context) error {
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
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.AdminPause(ethClient, bridgeAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("Deposits and proposals are paused")
	return nil
}
*/
