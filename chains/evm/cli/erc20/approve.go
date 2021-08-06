package erc20

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/cliutils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
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
	erc20Address := cmd.Flag("erc20Address").Value.String()
	recipientAddress := cmd.Flag("recipient").Value.String()
	amount := cmd.Flag("amount").Value.String()
	decimals, err := cmd.Flags().GetUint64("decimals")
	if err != nil {
		return err
	}
	log.Debug().Msgf(`
Approving ERC20
ERC20 address: %s
Recipient address: %s
Amount: %s
Decimals: %s`,
		erc20Address, recipientAddress, amount, decimals)

	// fetch global flag values
	url, _, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	decimalsBigInt := big.NewInt(0).SetUint64(decimals)

	if !common.IsHexAddress(erc20Address) {
		return errors.New("invalid erc20Address address")
	}
	erc20Addr := common.HexToAddress(erc20Address)

	if !common.IsHexAddress(recipientAddress) {
		return errors.New("invalid minter address")
	}
	recipientAddr := common.HexToAddress(recipientAddress)

	realAmount, err := cliutils.UserAmountToWei(amount, decimalsBigInt)
	if err != nil {
		log.Fatal().Err(err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Fatal().Err(err)
	}

	i, err := calls.PrepareErc20ApproveInput(erc20Addr, realAmount)
	if err != nil {
		log.Fatal().Err(err)
	}
	_, err = calls.SendInput(ethClient, erc20Addr, i, txFabric)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("%s account granted allowance on %v tokens of %s", recipientAddr.String(), amount, erc20Addr.String())
	return nil
}
