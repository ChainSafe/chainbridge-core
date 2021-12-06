package account

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/client"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/init"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"math/big"
)

var transferBaseCurrencyCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer base currency",
	Long:  "The generate subcommand is used to transfer the base currency",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := init.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := init.InitializeTransactor(gasPrice, evmtransaction.NewTransaction, c)
		if err != nil {
			return err
		}
		return TransferBaseCurrency(cmd, args, t)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateTransferBaseCurrencyFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessTransferBaseCurrencyFlags(cmd, args)
		return err
	},
}

func BindTransferCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Recipient, "recipient", "", "recipient address")
	cmd.Flags().StringVar(&Amount, "amount", "", "transfer amount")
	cmd.Flags().Uint64Var(&Decimals, "decimals", 0, "base token decimals")
	flags.MarkFlagsAsRequired(cmd, "recipient", "amount", "decimals")
}

func init() {
	BindTransferCmdFlags(transferBaseCurrencyCmd)
}
func ValidateTransferBaseCurrencyFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Recipient) {
		return fmt.Errorf("invalid recipient address %s", Recipient)
	}
	return nil
}

func ProcessTransferBaseCurrencyFlags(cmd *cobra.Command, args []string) error {
	var err error
	recipientAddress = common.HexToAddress(Recipient)
	decimals := big.NewInt(int64(Decimals))
	weiAmount, err = client.UserAmountToWei(Amount, decimals)
	return err
}
func TransferBaseCurrency(cmd *cobra.Command, args []string, t transactor.Transactor) error {
	hash, err := t.Transact(&recipientAddress, nil, transactor.TransactOptions{Value: weiAmount})
	if err != nil {
		log.Error().Err(fmt.Errorf("base currency deposit error: %v", err))
		return err
	}
	log.Debug().Msgf("base currency transaction hash: %s", hash.Hex())

	log.Info().Msgf("%s tokens were transferred to %s from %s", Amount, recipientAddress.Hex(), senderKeyPair.CommonAddress().String())
	return nil
}
