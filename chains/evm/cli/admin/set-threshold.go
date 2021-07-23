package admin

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var SetThresholdCmd = &cobra.Command{
	Use:   "set-threshold",
	Short: "Set a new relayer vote threshold",
	Long:  "Set a new relayer vote threshold",
	Run:   setThreshold,
}

func init() {
	SetThresholdCmd.Flags().Uint64("threshold", 0, "new relayer threshold")
	SetThresholdCmd.Flags().String("bridge", "", "bridge contract address")
}

func setThreshold(cmd *cobra.Command, args []string) {
	threshold := cmd.Flag("threshold").Value
	bridgeAddress := cmd.Flag("bridge").Value
	log.Debug().Msgf("Setting new threshold: %v for bridge address: %s", threshold, bridgeAddress)
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
