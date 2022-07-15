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

var removeRelayerCmd = &cobra.Command{
	Use:   "remove-relayer",
	Short: "Remove an existing relayer",
	Long:  "The remove-relayer subcommand removes an existing relayer",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: removeRelayer,
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateRemoveRelayerFlags(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func BindRemoveRelayerFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Relayer, "relayer", "", "Address to remove")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	flags.MarkFlagsAsRequired(cmd, "relayer", "bridge")
}
func init() {
	BindRemoveRelayerFlags(removeRelayerCmd)
}
func ValidateRemoveRelayerFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Relayer) {
		return fmt.Errorf("invalid relayer address %s", Relayer)
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func removeRelayer(cmd *cobra.Command, args []string) error {
	log.Debug().Msgf(`
Removing relayer
Relayer address: %s
Bridge address: %s`, Relayer, Bridge)
	return nil
}

/*
func removeRelayer(cctx *cli.Context) error {
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
	relayer := cctx.String("relayer")
	if !common.IsHexAddress(relayer) {
		return fmt.Errorf("invalid bridge address %s", relayer)
	}
	relayerAddress := common.HexToAddress(relayer)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.AdminRemoveRelayer(ethClient, bridgeAddress, relayerAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("Address %s is relayer now", relayerAddress.String())
	return nil
}
*/
