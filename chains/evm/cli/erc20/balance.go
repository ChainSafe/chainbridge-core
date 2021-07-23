package erc20

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Query balance of an account in an ERC20 contract",
	Long:  "Query balance of an account in an ERC20 contract",
	Run:   balance,
}

func init() {
	balanceCmd.Flags().String("erc20Address", "", "ERC20 contract address")
	balanceCmd.Flags().String("accountAddress", "", "address to receive balance of")
}

func balance(cmd *cobra.Command, args []string) {
	erc20Address := cmd.Flag("erc20Address").Value
	accountAddress := cmd.Flag("account").Value
	log.Debug().Msgf(`
Account balance of ERC20
ERC20 address: %s
Account address: %s
		`, erc20Address, accountAddress)
}

/*
func balanceOf(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	erc20 := cctx.String("erc20Address")
	if !common.IsHexAddress(erc20) {
		return errors.New("invalid erc20Address address")
	}
	erc20Address := common.HexToAddress(erc20)

	address := cctx.String("address")
	if !common.IsHexAddress(address) {
		return errors.New("invalid target address")
	}
	targetAddress := common.HexToAddress(address)

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	balance, err := utils.ERC20BalanceOf(ethClient, erc20Address, targetAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("balance of %s is %s", targetAddress.String(), balance.String())
	return nil
}
*/
