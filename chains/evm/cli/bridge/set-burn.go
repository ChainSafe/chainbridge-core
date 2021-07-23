package bridge

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var SetBurnCmd = &cobra.Command{
	Use:   "set-burn",
	Short: "Set a token contract as mintable/burnable",
	Long:  "Set a token contract as mintable/burnable in a handler",
	Run:   setBurn,
}

func init() {
	SetBurnCmd.Flags().String("handler", "", "ERC20 handler contract address")
	SetBurnCmd.Flags().String("bridge", "", "bridge contract address")
	SetBurnCmd.Flags().String("tokenContract", "", "token contract to be registered")
}

func setBurn(cmd *cobra.Command, args []string) {
	handlerAddress := cmd.Flag("handler").Value
	bridgeAddress := cmd.Flag("bridge").Value
	tokenAddress := cmd.Flag("tokenContract").Value
	log.Debug().Msgf(`
Setting contract as mintable/burnable
Handler address: %s
Bridge address: %s
Token contract address: %s`, handlerAddress, bridgeAddress, tokenAddress)
}

/*
func setBurn(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Int64("gasLimit")
	gasPrice := cctx.Int64("gasPrice")
	bridgeAddress, err := cliutils.DefineBridgeAddress(cctx)
	if err != nil {
		return err
	}
	handler := cctx.String("handler")
	if !common.IsHexAddress(handler) {
		return errors.New("handler address is incorrect format")
	}
	tokenContract := cctx.String("tokenContract")
	if !common.IsHexAddress(tokenContract) {
		return errors.New("tokenContract address is incorrect format")
	}
	handlerAddress := common.HexToAddress(handler)
	tokenContractAddress := common.HexToAddress(tokenContract)

	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(gasLimit), big.NewInt(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	log.Info().Msgf("Setting contract %s as burnable on handler %s", tokenContractAddress.String(), handlerAddress.String())
	err = utils.SetBurnable(ethClient, bridgeAddress, handlerAddress, tokenContractAddress)
	if err != nil {
		return err
	}
	log.Info().Msg("Burnable set")
	return nil
}
*/
