package account

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"

	callsUtil "github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ChainSafe/sygma-core/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var transferBaseCurrencyCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer base currency",
	Long:  "The transfer subcommand is used to transfer the base currency",
	PreRun: func(cmd *cobra.Command, args []string) {
		confirmTransfer(cmd, args)
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
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

func BindTransferBaseCurrencyFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Recipient, "recipient", "", "Recipient address")
	cmd.Flags().StringVar(&Amount, "amount", "", "Transfer amount")
	cmd.Flags().Uint64Var(&Decimals, "decimals", 0, "Base token decimals")
	flags.MarkFlagsAsRequired(cmd, "recipient", "amount", "decimals")
}

func init() {
	BindTransferBaseCurrencyFlags(transferBaseCurrencyCmd)
}
func ValidateTransferBaseCurrencyFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Recipient) {
		return fmt.Errorf("invalid recipient address %s", Recipient)
	}
	return nil
}

func ProcessTransferBaseCurrencyFlags(cmd *cobra.Command, args []string) error {
	var err error
	RecipientAddress = common.HexToAddress(Recipient)
	decimals := big.NewInt(int64(Decimals))
	WeiAmount, err = callsUtil.UserAmountToWei(Amount, decimals)
	return err
}
func TransferBaseCurrency(cmd *cobra.Command, args []string, t transactor.Transactor) error {
	hash, err := t.Transact(
		&RecipientAddress, nil, transactor.TransactOptions{Value: WeiAmount, GasLimit: gasLimit},
	)
	if err != nil {
		log.Error().Err(fmt.Errorf("base currency deposit error: %v", err))
		return err
	}
	log.Debug().Msgf("base currency transaction hash: %s", hash.Hex())

	log.Info().Msgf("%s tokens were transferred to %s from %s", Amount, RecipientAddress.Hex(), senderKeyPair.CommonAddress().String())
	return nil
}

func confirmTransfer(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Send transaction %s(%d) to %s (Y/N)?", Amount, Decimals, Recipient)
		s, _ := reader.ReadString('\n')

		s = strings.ToLower(strings.TrimSuffix(s, "\n"))

		if strings.Compare(s, "n") == 0 {
			os.Exit(0)
		} else if strings.Compare(s, "y") == 0 {
			break
		} else {
			continue
		}
	}
}
