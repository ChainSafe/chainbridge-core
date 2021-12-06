package erc721

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/contracts"
	"github.com/ChainSafe/chainbridge-core/util"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/erc721"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var ownerCmd = &cobra.Command{
	Use:   "owner",
	Short: "Get token owner from an ERC721 mintable contract",
	Long:  "Get token owner from an ERC721 mintable contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		erc721Contract, err := contracts.InitializeErc721Contract(
			url, gasLimit, gasPrice, senderKeyPair, erc721Addr,
		)
		if err != nil {
			return err
		}
		return OwnerCmd(cmd, args, erc721Contract)
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

func BindOwnerCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Erc721Address, "contract-address", "", "address of contract")
	cmd.Flags().StringVar(&TokenId, "tokenId", "", "ERC721 token ID")
	flags.MarkFlagsAsRequired(cmd, "contract-address", "tokenId")
}

func init() {
	BindOwnerCmdFlags(ownerCmd)
}

func ValidateOwnerFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc721Address) {
		return fmt.Errorf("invalid ERC721 contract address %s", Erc721Address)
	}
	return nil
}

func ProcessOwnerFlags(cmd *cobra.Command, args []string) error {
	erc721Addr = common.HexToAddress(Erc721Address)

	var ok bool
	if tokenId, ok = big.NewInt(0).SetString(TokenId, 10); !ok {
		return fmt.Errorf("invalid token id value")
	}

	return nil
}

func OwnerCmd(cmd *cobra.Command, args []string, erc721Contract *erc721.ERC721Contract) error {
	owner, err := erc721Contract.Owner(tokenId)
	if err != nil {
		return err
	}

	log.Info().Msgf("%v token owner: %v", tokenId, owner)
	return err
}
