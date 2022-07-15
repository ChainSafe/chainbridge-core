package erc721

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/erc721"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/util"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var approveCmd = &cobra.Command{
	Use:   "approve",
	Short: "Approve an ERC721 token",
	Long:  "The approve subcommand approves a token in an ERC721 contract for transfer",
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
		return ApproveCmd(cmd, args, erc721.NewErc721Contract(c, Erc721Addr, t))
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
	cmd.Flags().StringVar(&Erc721Address, "contract", "", "ERC721 contract address")
	cmd.Flags().StringVar(&Recipient, "recipient", "", "Recipient address")
	cmd.Flags().StringVar(&Token, "token", "", "ERC721 token ID")
	flags.MarkFlagsAsRequired(cmd, "contract", "recipient", "token")
}

func init() {
	BindApproveFlags(approveCmd)
}

func ValidateApproveFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc721Address) {
		return fmt.Errorf("invalid ERC721 contract address %s", Erc721Address)
	}
	if !common.IsHexAddress(Recipient) {
		return fmt.Errorf("invalid recipient address")
	}
	return nil
}

func ProcessApproveFlags(cmd *cobra.Command, args []string) error {
	RecipientAddr = common.HexToAddress(Recipient)
	Erc721Addr = common.HexToAddress(Erc721Address)

	var ok bool
	if TokenId, ok = big.NewInt(0).SetString(Token, 10); !ok {
		return fmt.Errorf("invalid token id value")
	}
	return nil
}

func ApproveCmd(cmd *cobra.Command, args []string, erc721Contract *erc721.ERC721Contract) error {
	_, err = erc721Contract.Approve(
		TokenId, RecipientAddr, transactor.TransactOptions{GasLimit: gasLimit},
	)
	if err != nil {
		return err
	}

	log.Info().Msgf("%v token approved", TokenId)
	return err
}
