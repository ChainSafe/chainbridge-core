package admin

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ChainSafe/sygma-core/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var removeAdminCmd = &cobra.Command{
	Use:   "remove-admin",
	Short: "Remove an existing admin",
	Long:  "The remove-admin subcommand removes an existing admin",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: removeAdmin,
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateRemoveAdminFlags(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func BindRemoveAdminFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Admin, "admin", "", "Address to remove")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	flags.MarkFlagsAsRequired(cmd, "admin", "bridge")
}

func init() {
	BindRemoveAdminFlags(removeAdminCmd)
}
func ValidateRemoveAdminFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Admin) {
		return fmt.Errorf("invalid admin address %s", Admin)
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func removeAdmin(cmd *cobra.Command, args []string) error {

	log.Debug().Msgf(`
Removing admin
Admin address: %s
Bridge address: %s`, Admin, Bridge)
	return nil
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
