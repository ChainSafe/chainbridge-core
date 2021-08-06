package admin

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var withdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "Withdraw tokens from the handler contract",
	Long:  "Withdraw tokens from the handler contract",
	Run:   withdraw,
}

func init() {
	withdrawCmd.Flags().String("amount", "", "token amount to withdraw. Should be set or ID or amount if both set error will occur")
	withdrawCmd.Flags().String("id", "", "token ID to withdraw. Should be set or ID or amount if both set error will occur")
	withdrawCmd.Flags().String("bridge", "", "bridge contract address")
	withdrawCmd.Flags().String("handler", "", "handler contract address")
	withdrawCmd.Flags().String("token", "", "ERC20 or ERC721 token contract address")
	withdrawCmd.Flags().String("recipient", "", "address to withdraw to")
	withdrawCmd.Flags().Uint64("decimals", 0, "ERC20 token decimals")
}

func withdraw(cmd *cobra.Command, args []string) {
	amount := cmd.Flag("amount").Value
	id := cmd.Flag("id").Value
	bridgeAddress := cmd.Flag("bridge").Value
	handler := cmd.Flag("handler").Value
	token := cmd.Flag("token").Value
	recipient := cmd.Flag("recipient").Value
	decimals := cmd.Flag("decimals").Value
	log.Debug().Msgf(`
Withdrawing
Amount: %s
ID: %s
Bridge address: %s
Handler: %s
Token: %s
Recipient: %s
Decimals: %v`, amount, id, bridgeAddress, handler, token, recipient, decimals)
}

/*

func withdraw(cctx *cli.Context) error {
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

	handler := cctx.String("handler")
	if !common.IsHexAddress(handler) {
		return fmt.Errorf("invalid handler address %s", handler)
	}
	handlerAddress := common.HexToAddress(handler)

	token := cctx.String("token")
	if !common.IsHexAddress(token) {
		return fmt.Errorf("invalid token address %s", token)
	}
	tokenAddress := common.HexToAddress(token)

	recipient := cctx.String("recipient")
	if !common.IsHexAddress(recipient) {
		return fmt.Errorf("invalid recipient address %s", recipient)
	}
	recipientAddress := common.HexToAddress(recipient)

	amount := cctx.String("amount")
	id := cctx.String("id")

	if id != "" && amount != "" {
		return errors.New("Only id or amount should be set.")
	}
	if id == "" && amount == "" {
		return errors.New("id or amount flag should be set")
	}
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	idOrAmountToWithdraw := new(big.Int)
	if amount != "" {
		decimals := big.NewInt(0).SetUint64(cctx.Uint64("decimals"))
		idOrAmountToWithdraw, err = utils.UserAmountToWei(amount, decimals)
		if err != nil {
			return err
		}
	} else {
		idOrAmountToWithdraw.SetString(id, 10)
	}

	err = utils.AdminWithdraw(ethClient, bridgeAddress, handlerAddress, tokenAddress, recipientAddress, idOrAmountToWithdraw)
	if err != nil {
		return err
	}

	log.Info().Msgf("Withdrawn %s to %s", idOrAmountToWithdraw.String(), recipient)
	return nil
}
*/
