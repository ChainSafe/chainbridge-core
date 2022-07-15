package erc20

import (
	"errors"
	"math/big"

	callsUtil "github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/erc20"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/util"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var approveCmd = &cobra.Command{
	Use:   "approve",
	Short: "Approve an ERC20 tokens",
	Long:  "The approve subcommand approves tokens in an ERC20 contract for transfer",
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
		return ApproveCmd(cmd, args, erc20.NewERC20Contract(c, Erc20Addr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateApproveFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessApproveFlags(cmd, args)
		return err
	},
}

func BindApproveFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Erc20Address, "contract", "", "ERC20 contract address")
	cmd.Flags().StringVar(&Amount, "amount", "", "Amount to grant allowance")
	cmd.Flags().StringVar(&Recipient, "recipient", "", "Recipient address")
	cmd.Flags().Uint64Var(&Decimals, "decimals", 0, "ERC20 token decimals")
	flags.MarkFlagsAsRequired(cmd, "contract", "amount", "recipient", "decimals")
}

func init() {
	BindApproveFlags(approveCmd)
}

func ValidateApproveFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc20Address) {
		return errors.New("invalid erc20Address address")
	}
	if !common.IsHexAddress(Recipient) {
		return errors.New("invalid minter address")
	}
	return nil
}

func ProcessApproveFlags(cmd *cobra.Command, args []string) error {
	var err error

	decimals := big.NewInt(int64(Decimals))
	Erc20Addr = common.HexToAddress(Erc20Address)
	RecipientAddress = common.HexToAddress(Recipient)
	RealAmount, err = callsUtil.UserAmountToWei(Amount, decimals)
	if err != nil {
		return err
	}

	return nil
}

func ApproveCmd(cmd *cobra.Command, args []string, contract *erc20.ERC20Contract) error {
	log.Debug().Msgf(`
Approving ERC20
ERC20 address: %s
Recipient address: %s
Amount: %s
Decimals: %v`,
		Erc20Address, Recipient, Amount, Decimals)

	_, err := contract.ApproveTokens(RecipientAddress, RealAmount, transactor.TransactOptions{GasLimit: gasLimit})
	if err != nil {
		log.Fatal().Err(err)
		return err
	}
	log.Info().Msgf(
		"%s account granted allowance on %v tokens of %s",
		RecipientAddress.String(), Amount, RecipientAddress.String(),
	)
	return nil
}
