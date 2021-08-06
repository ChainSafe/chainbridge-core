package erc20

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/cliutils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var approveCmd = &cobra.Command{
	Use:   "approve",
	Short: "Approve tokens in an ERC20 contract for transfer",
	Long:  "Approve tokens in an ERC20 contract for transfer",
	RunE:  CallApprove,
}

func init() {
	approveCmd.Flags().String("erc20Address", "", "ERC20 contract address")
	approveCmd.Flags().String("amount", "", "amount to grant allowance")
	approveCmd.Flags().String("recipient", "", "address of recipient")
	approveCmd.Flags().Uint64("decimals", 0, "ERC20 token decimals")
	approveCmd.MarkFlagRequired("decimals")
}

func CallApprove(cmd *cobra.Command, args []string) error {
	txFabric := evmtransaction.NewTransaction
	return approve(cmd, args, txFabric)
}

func approve(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	erc20Address := common.HexToAddress(cmd.Flag("erc20Address").Value.String())
	recipientAddress := common.HexToAddress(cmd.Flag("recipient").Value.String())
	amount := cmd.Flag("amount").Value.String()
	decimals := cmd.Flag("decimals").Value.String()
	log.Debug().Msgf(`
Approving ERC20
ERC20 address: %s
Recipient address: %s
Amount: %s
Decimals: %s`,
		erc20Address, recipientAddress, amount, decimals)

	url := cmd.Flag("url").Value.String()
	decimalsBigInt, _ := big.NewInt(0).SetString(decimals, 10)

	// erc20 := cctx.String("erc20Address")
	// if !common.IsHexAddress(erc20) {
	// 	return errors.New("invalid erc20Address address")
	// }
	// erc20Address := common.HexToAddress(erc20)

	// recipient := cctx.String("recipient")
	// if !common.IsHexAddress(recipient) {
	// 	return errors.New("invalid minter address")
	// }
	// recipientAddress := common.HexToAddress(recipient)

	realAmount, err := cliutils.UserAmountToWei(amount, decimalsBigInt)
	if err != nil {
		log.Fatal().Err(err)
	}

	senderKeyPair, err := cliutils.DefineSender(cmd)
	if err != nil {
		log.Fatal().Err(err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Fatal().Err(err)
	}

	i, err := calls.PrepareErc20ApproveInput(erc20Address, realAmount)
	if err != nil {
		log.Fatal().Err(err)
	}
	_, err = calls.SendInput(ethClient, erc20Address, i, txFabric)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("%s account granted allowance on %v tokens of %s", recipientAddress.String(), amount, erc20Address.String())
	return nil
}
