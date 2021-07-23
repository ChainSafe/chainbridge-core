package admin

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var AddAdminCmd = &cobra.Command{
	Use:   "add-admin",
	Short: "Add a new admin",
	Long:  "Add a new admin",
	Run:   addAdmin,
}

func init() {
	AddAdminCmd.Flags().String("admin", "", "address to add")
	AddAdminCmd.Flags().String("bridge", "", "bridge contract address")
}

func addAdmin(cmd *cobra.Command, args []string) {
	adminAddress := cmd.Flag("admin").Value
	bridgeAddress := cmd.Flag("bridge").Value
	log.Debug().Msgf(`
Adding admin
Admin address: %s
Bridge address: %s`, adminAddress, bridgeAddress)
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
