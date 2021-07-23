package erc20

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var approveCmd = &cobra.Command{
	Use:   "approve",
	Short: "Approve tokens in an ERC20 contract for transfer",
	Long:  "Approve tokens in an ERC20 contract for transfer",
	Run:   approve,
}

func init() {
	approveCmd.Flags().String("erc20Address", "", "ERC20 contract address")
	approveCmd.Flags().String("amount", "", "amount to grant allowance")
	approveCmd.Flags().String("recipient", "", "address of recipient")
	approveCmd.Flags().Uint64("decimals", 0, "ERC20 token decimals")
	approveCmd.MarkFlagRequired("decimals")
}

func approve(cmd *cobra.Command, args []string) {
	erc20Address := cmd.Flag("erc20Address").Value
	recipientAddress := cmd.Flag("recipient").Value
	amount := cmd.Flag("amount").Value
	decimals := cmd.Flag("decimals").Value
	log.Debug().Msgf(`
Approving ERC20
ERC20 address: %s
Recipient address: %s
Amount: %s
Decimals: %d`,
		erc20Address, recipientAddress, amount, decimals)
}

/*
func approve(cctx *cli.Context) error {
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

	recipient := cctx.String("recipient")
	if !common.IsHexAddress(recipient) {
		return errors.New("invalid minter address")
	}
	recipientAddress := common.HexToAddress(recipient)

	amount := cctx.String("amount")

	realAmount, err := utils.UserAmountToWei(amount, decimals)
	if err != nil {
		return err
	}

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.Erc20Approve(ethClient, erc20Address, recipientAddress, realAmount)
	if err != nil {
		return err
	}
	log.Info().Msgf("%s account granted allowance on %v tokens of %s", recipientAddress.String(), amount, sender.CommonAddress().String())
	return nil
}
*/
