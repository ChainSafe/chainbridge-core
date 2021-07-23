package admin

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var UnpauseCmd = &cobra.Command{
	Use:   "unpause",
	Short: "Unpause deposits and proposals",
	Long:  "Unpause deposits and proposals",
	Run:   unpause,
}

func init() {
	UnpauseCmd.Flags().String("bridge", "", "bridge contract address")
}

func unpause(cmd *cobra.Command, args []string) {
	bridgeAddress := cmd.Flag("bridge").Value
	log.Debug().Msgf(`
Unpausing
Bridge address: %s`, bridgeAddress)
}

/*
func unpause(cctx *cli.Context) error {
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
	err = utils.AdminUnpause(ethClient, bridgeAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("Deposits and proposals are Unpaused")
	return nil
}
*/
