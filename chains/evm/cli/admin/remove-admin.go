package admin

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/writer"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var removeAdminCmd = &cobra.Command{
	Use:   "remove-admin",
	Short: "Remove an existing admin",
	Long:  "Remove an existing admin",
	Run:   removeAdmin,
}

func init() {
	removeAdminCmd.Flags().StringVar(&Admin, "admin", "", "address to remove")
	removeAdminCmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	flags.MarkFlagsAsRequired(removeAdminCmd, "admin", "bridge")
}

func removeAdmin(cmd *cobra.Command, args []string) {

	log.Debug().Msgf(`
Removing admin
Admin address: %s
Bridge address: %s`, Admin, Bridge)
	writer.WriteCliDataToFile(cmd)
}

/*
func removeAdmin(cctx *cli.Context) error {
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
	err = utils.AdminRemoveAdmin(ethClient, bridgeAddress, adminAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("Address %s is removed from admins", adminAddress.String())
	return nil
}
*/
