package admin

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var addAdminCmd = &cobra.Command{
	Use:   "add-admin",
	Short: "Add a new admin",
	Long:  "Add a new admin",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	Run: addAdmin,
}

func BindAddAdminFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Admin, "admin", "", "address to add")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	flags.MarkFlagsAsRequired(cmd, "admin", "bridge")
}

func init() {
	BindAddAdminFlags(addAdminCmd)
}

func addAdmin(cmd *cobra.Command, args []string) {
	log.Debug().Msgf(`
Adding admin
Admin address: %s
Bridge address: %s`, Admin, Bridge)
}

/*
func addAdmin(cctx *cli.Context) error {
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
	admin := cctx.String("admin")
	if !common.IsHexAddress(admin) {
		return fmt.Errorf("invalid admin address %s", admin)
	}
	adminAddress := common.HexToAddress(admin)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.AdminAddAdmin(ethClient, bridgeAddress, adminAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("Address %s is set to admin", adminAddress.String())
	return nil
}
*/
