package admin

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var setFeeCmd = &cobra.Command{
	Use:   "set-fee",
	Short: "Set a new fee for deposits",
	Long:  "Set a new fee for deposits",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	Run: setFee,
}

func BindSetFeeFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Fee, "fee", "", "New fee (in ether)")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	flags.MarkFlagsAsRequired(cmd, "fee", "bridge")
}
func init() {
	BindSetFeeFlags(setFeeCmd)
}

func setFee(cmd *cobra.Command, args []string) {
	log.Debug().Msgf(`
Setting new fee
Fee amount: %s
Bridge address: %s`, Fee, Bridge)
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
