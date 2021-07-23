package erc20

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var DepositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Initiate a transfer of ERC20 tokens",
	Long:  "Initiate a transfer of ERC20 tokens",
	Run:   deposit,
}

func init() {
	DepositCmd.Flags().String("recipient", "", "address of recipient")
	DepositCmd.Flags().String("bridge", "", "address of bridge contract")
	DepositCmd.Flags().String("amount", "", "amount to deposit")
	DepositCmd.Flags().String("value", "", "value of ETH that should be sent along with deposit to cover possible fees. In ETH (decimals are allowed)")
	DepositCmd.Flags().String("destId", "", "destination chain ID")
	DepositCmd.Flags().String("resourceId", "", "resource ID for transfer")
	DepositCmd.Flags().Uint64("decimals", 0, "ERC20 token decimals")
}

func deposit(cmd *cobra.Command, args []string) {
	recipientAddress := cmd.Flag("recipient").Value
	bridgeAddress := cmd.Flag("bridge").Value
	amount := cmd.Flag("amount").Value
	value := cmd.Flag("value").Value
	destinationId := cmd.Flag("destId").Value
	resourceId := cmd.Flag("resourceId").Value
	decimals := cmd.Flag("decimals").Value
	log.Debug().Msgf(`
Initiating deposit of ERC20
Recipient address: %s
Bridge address: %s
Amount: %s
Value: %s
Destination chain ID: %s
Resource ID: %s
Decimals: %d
		`, recipientAddress, bridgeAddress, amount, value, destinationId, resourceId, decimals)
}

/*
func deposit(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	decimals := big.NewInt(0).SetUint64(cctx.Uint64("decimals"))

	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	bridgeAddress, err := cliutils.DefineBridgeAddress(cctx)
	if err != nil {
		return err
	}

	recipient := cctx.String("recipient")
	if !common.IsHexAddress(recipient) {
		return fmt.Errorf("invalid recipient address %s", recipient)
	}
	recipientAddress := common.HexToAddress(recipient)

	amount := cctx.String("amount")

	realAmount, err := utils.UserAmountToWei(amount, decimals)
	if err != nil {
		return err
	}

	value := cctx.String("value")

	realValue, err := utils.UserAmountToWei(value, big.NewInt(18))
	if err != nil {
		return err
	}
	dest := cctx.Uint64("dest")

	resourceId := cctx.String("resourceId")
	resourceIDBytes := utils.SliceTo32Bytes(common.Hex2Bytes(resourceId))

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}

	ethClient.ClientWithArgs(client.ClientWithValue(realValue))

	err = utils.MakeAndSendERC20Deposit(ethClient, bridgeAddress, recipientAddress, realAmount, resourceIDBytes, uint8(dest))
	if err != nil {
		return err
	}
	log.Info().Msgf("%s tokens were transferred to %s from %s", amount, recipientAddress.String(), sender.CommonAddress().String())
	return nil
}
*/
