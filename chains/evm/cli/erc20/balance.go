package erc20

import (
	"errors"
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Query balance of an account in an ERC20 contract",
	Long:  "Query balance of an account in an ERC20 contract",
	RunE:  CallBalance,
}

func init() {
	balanceCmd.Flags().String("erc20Address", "", "ERC20 contract address")
	balanceCmd.Flags().String("accountAddress", "", "address to receive balance of")
}

func CallBalance(cmd *cobra.Command, args []string) error {
	txFabric := evmtransaction.NewTransaction
	return balance(cmd, args, txFabric)
}

func balance(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	erc20Address := cmd.Flag("erc20Address").Value.String()
	accountAddress := cmd.Flag("accountAddress").Value.String()

	// fetch global flag values
	url, _, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	if !common.IsHexAddress(erc20Address) {
		return errors.New("invalid erc20Address address")
	}
	erc20Addr := common.HexToAddress(erc20Address)

	if !common.IsHexAddress(accountAddress) {
		return errors.New("invalid account address")
	}
	accountAddr := common.HexToAddress(accountAddress)

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	// erc20Addr, accountAddr
	input, err := calls.PrepareERC20BalanceInput(erc20Addr, accountAddr)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare input error: %v", err))
		return err
	}

	txHash, err := calls.SendInput(ethClient, erc20Addr, input, txFabric)
	if err != nil {
		log.Error().Err(fmt.Errorf("send input error: %v", err))
		return err
	}
	log.Debug().Msgf("tx hash: %v", txHash.Hex())

	// determine balance

	// log.Info().Msgf("balance of %s is %s", accountAddr.String(), balance.String())
	return nil
}
