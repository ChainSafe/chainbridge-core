package erc20

import (
	"context"
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Query balance of an account in an ERC20 contract",
	Long:  "Query balance of an account in an ERC20 contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return BalanceCmd(cmd, args, txFabric)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := validateBalanceFlags(cmd, args)
		if err != nil {
			return err
		}

		processBalanceFlags(cmd, args)
		return nil
	},
}

func BindBalanceCmdFlags() {
	balanceCmd.Flags().StringVar(&Erc20Address, "erc20Address", "", "ERC20 contract address")
	balanceCmd.Flags().StringVar(&AccountAddress, "accountAddress", "", "address to receive balance of")
	flags.MarkFlagsAsRequired(balanceCmd, "erc20Address", "accountAddress")
}

func init() {
	BindBalanceCmdFlags()
}

var accountAddr common.Address

func validateBalanceFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc20Address) {
		return fmt.Errorf("invalid recipient address %s", Recipient)
	}
	if !common.IsHexAddress(AccountAddress) {
		return fmt.Errorf("invalid account address %s", AccountAddress)
	}
	return nil
}

func processBalanceFlags(cmd *cobra.Command, args []string) {
	erc20Addr = common.HexToAddress(Erc20Address)
	accountAddr = common.HexToAddress(AccountAddress)
}

func BalanceCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	// fetch global flag values
	url, _, _, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	// erc20Addr, accountAddr
	input, err := calls.PrepareERC20BalanceInput(accountAddr)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare input error: %v", err))
		return err
	}

	msg := ethereum.CallMsg{
		From: common.Address{},
		To:   &erc20Addr,
		Data: input,
	}

	out, err := ethClient.CallContract(context.TODO(), calls.ToCallArg(msg), nil)
	if err != nil {
		log.Error().Err(fmt.Errorf("call contract error: %v", err))
		return err
	}

	if len(out) == 0 {
		// Make sure we have a contract to operate on, and bail out otherwise.
		if code, err := ethClient.CodeAt(context.Background(), erc20Addr, nil); err != nil {
			return err
		} else if len(code) == 0 {
			return fmt.Errorf("no code at provided address %s", erc20Addr.String())
		}
	}

	balance, err := calls.ParseERC20BalanceOutput(out)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare output error: %v", err))
		return err
	}

	log.Info().Msgf("balance of %s is %s", accountAddr.String(), balance.String())
	return nil
}
