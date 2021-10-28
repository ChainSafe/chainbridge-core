package admin

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/writer"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var removeRelayerCmd = &cobra.Command{
	Use:   "remove-relayer",
	Short: "Remove a relayer",
	Long:  "Remove a relayer",
	Run:   removeRelayer,
}

func init() {
	removeRelayerCmd.Flags().StringVar(&Relayer, "relayer", "", "address to remove")
	removeRelayerCmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	flags.MarkFlagsAsRequired(removeRelayerCmd, "relayer", "bridge")
}

func removeRelayer(cmd *cobra.Command, args []string) {
	log.Debug().Msgf(`
Removing relayer
Relayer address: %s
Bridge address: %s`, Relayer, Bridge)
	writer.WriteCliDataToFile(cmd)
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
