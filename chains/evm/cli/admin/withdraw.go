package admin

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var withdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "Withdraw tokens from the handler contract",
	Long:  "Withdraw tokens from the handler contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	Run: withdraw,
}

func BindWithdrawFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Amount, "amount", "", "token amount to withdraw. Should be set or ID or amount if both set error will occur")
	cmd.Flags().StringVar(&TokenID, "tokenId", "", "token ID to withdraw. Should be set or ID or amount if both set error will occur")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	cmd.Flags().StringVar(&Handler, "handler", "", "handler contract address")
	cmd.Flags().StringVar(&Token, "token", "", "ERC20 or ERC721 token contract address")
	cmd.Flags().StringVar(&Recipient, "recipient", "", "address to withdraw to")
	cmd.Flags().Uint64Var(&Decimals, "decimals", 0, "ERC20 token decimals")
	flags.MarkFlagsAsRequired(cmd, "amount", "tokenId", "bridge", "handler", "token", "recipient", "decimals")
}

func init() {
	BindWithdrawFlags(withdrawCmd)
}

func withdraw(cmd *cobra.Command, args []string) {
	log.Debug().Msgf(`
Withdrawing
Amount: %s
TokenID: %s
Bridge address: %s
Handler: %s
Token: %s
Recipient: %s
Decimals: %v`, Amount, TokenID, Bridge, Handler, Token, Recipient, Decimals)
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
