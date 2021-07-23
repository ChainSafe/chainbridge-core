package admin

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var SetFeeCmd = &cobra.Command{
	Use:   "set-fee",
	Short: "Set a new fee for deposits",
	Long:  "Set a new fee for deposits",
	Run:   setFee,
}

func init() {
	SetFeeCmd.Flags().String("fee", "", "New fee (in ether)")
	SetFeeCmd.Flags().String("bridge", "", "bridge contract address")
}

func setFee(cmd *cobra.Command, args []string) {
	feeAmount := cmd.Flag("fee").Value
	bridgeAddress := cmd.Flag("bridge").Value
	log.Debug().Msgf(`
Setting new fee
Fee amount: %s
Bridge address: %s`, feeAmount, bridgeAddress)
}

/*
func setFee(cctx *cli.Context) error {
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
	fee := cctx.String("fee")

	realFeeAmount, err := utils.UserAmountToWei(fee, big.NewInt(18))
	if err != nil {
		return err
	}

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.AdminSetFee(ethClient, bridgeAddress, realFeeAmount)
	if err != nil {
		return err
	}
	log.Info().Msgf("Fee set to %s", realFeeAmount.String())
	return nil
}
*/
