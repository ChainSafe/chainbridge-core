package admin

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var IsRelayerCmd = &cobra.Command{
	Use:   "is-relayer",
	Short: "Checks if an address is registered as a relayer",
	Long:  "Checks if an address is registered as a relayer",
	Run:   isRelayer,
}

func init() {
	IsRelayerCmd.Flags().String("relayer", "", "address to check")
	IsRelayerCmd.Flags().String("bridge", "", "bridge contract address")
}

func isRelayer(cmd *cobra.Command, args []string) {
	relayerAddress := cmd.Flag("relayer").Value
	bridgeAddress := cmd.Flag("bridge").Value
	log.Debug().Msgf("Adding relayer: %s to bridge address: %s", relayerAddress, bridgeAddress)
}

/*
func isRelayer(cctx *cli.Context) error {
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
		return fmt.Errorf("invalid relayer address %s", relayer)
	}
	relayerAddress := common.HexToAddress(relayer)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	isRelayer, err := utils.AdminIsRelayer(ethClient, bridgeAddress, relayerAddress)
	if err != nil {
		return err
	}
	if isRelayer {
		log.Info().Msgf("Requested address %s is relayer", relayerAddress.String())
	} else {
		log.Info().Msgf("Requested address %s is not a relayer", relayerAddress.String())
	}
	return nil
}
*/
