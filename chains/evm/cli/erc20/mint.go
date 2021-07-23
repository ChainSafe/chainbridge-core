package erc20

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var MintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint tokens on an ERC20 mintable contract",
	Long:  "Mint tokens on an ERC20 mintable contract",
	Run:   mint,
}

func init() {
	MintCmd.Flags().String("amount", "", "amount to mint fee (in wei)")
	MintCmd.Flags().String("bridge", "", "bridge contract address")

}

func mint(cmd *cobra.Command, args []string) {
	amount := cmd.Flag("amount").Value
	bridgeAddress := cmd.Flag("bridge").Value
	log.Debug().Msgf(`
Minting token
Amount: %s
Bridge address: %s`, amount, bridgeAddress)
}

/*
func mint(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	decimals := big.NewInt(0).SetUint64(cctx.Uint64("decimals"))
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	erc20 := cctx.String("erc20Address")
	if !common.IsHexAddress(erc20) {
		return errors.New("invalid erc20Address address")
	}
	erc20Address := common.HexToAddress(erc20)

	amount := cctx.String("amount")

	realAmount, err := utils.UserAmountToWei(amount, decimals)
	if err != nil {
		return err
	}

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.ERC20Mint(ethClient, realAmount, erc20Address, sender.CommonAddress())
	if err != nil {
		return err
	}
	log.Info().Msgf("%v tokens minted", amount)
	return nil
}
*/
