package erc721

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/erc721"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/util"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var ownerCmd = &cobra.Command{
	Use:   "owner",
	Short: "Get an ERC721 token owner",
	Long:  "The owner subcommand gets a token owner from an ERC721 mintable contract",
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
		return OwnerCmd(cmd, args, erc721.NewErc721Contract(c, Erc721Addr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateOwnerFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessOwnerFlags(cmd, args)
		return err
	},
}

func BindOwnerFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Erc721Address, "contract", "", "ERC721 contract address")
	cmd.Flags().StringVar(&Token, "token", "", "ERC721 token ID")
	flags.MarkFlagsAsRequired(cmd, "contract", "token")
}

func init() {
	BindOwnerFlags(ownerCmd)
}

func ValidateOwnerFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc721Address) {
		return fmt.Errorf("invalid ERC721 contract address %s", Erc721Address)
	}
	return nil
}

func ProcessOwnerFlags(cmd *cobra.Command, args []string) error {
	Erc721Addr = common.HexToAddress(Erc721Address)

	var ok bool
	if TokenId, ok = big.NewInt(0).SetString(Token, 10); !ok {
		return fmt.Errorf("invalid token id value")
	}

	return nil
}

func OwnerCmd(cmd *cobra.Command, args []string, erc721Contract *erc721.ERC721Contract) error {
	owner, err := erc721Contract.Owner(TokenId)
	if err != nil {
		return err
	}

	log.Info().Msgf("%v token owner: %v", TokenId, owner)
	return err
}
