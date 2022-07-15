package erc20

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/erc20"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/util"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Query an ERC20 token balance",
	Long:  "The balance subcommand queries the balance of an account in an ERC20 contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, evmtransaction.NewTransaction, c, prepare)
		if err != nil {
			return err
		}
		return BalanceCmd(cmd, args, erc20.NewERC20Contract(c, Erc20Addr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateBalanceFlags(cmd, args)
		if err != nil {
			return err
		}

		ProcessBalanceFlags(cmd, args)
		return nil
	},
}

func BindBalanceFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Erc20Address, "contract", "", "ERC20 contract address")
	cmd.Flags().StringVar(&AccountAddress, "address", "", "Address to receive balance of")
	flags.MarkFlagsAsRequired(cmd, "contract", "address")
}

func init() {
	BindBalanceFlags(balanceCmd)
}

var accountAddr common.Address

func ValidateBalanceFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc20Address) {
		return fmt.Errorf("invalid recipient address %s", Recipient)
	}
	if !common.IsHexAddress(AccountAddress) {
		return fmt.Errorf("invalid account address %s", AccountAddress)
	}
	return nil
}

func ProcessBalanceFlags(cmd *cobra.Command, args []string) {
	Erc20Addr = common.HexToAddress(Erc20Address)
	accountAddr = common.HexToAddress(AccountAddress)
}

func BalanceCmd(cmd *cobra.Command, args []string, contract *erc20.ERC20Contract) error {
	balance, err := contract.GetBalance(accountAddr)
	if err != nil {
		log.Error().Err(fmt.Errorf("failed contract call error: %v", err))
		return err
	}

	log.Info().Msgf("balance of %s is %s", accountAddr.String(), balance.String())
	return nil
}
