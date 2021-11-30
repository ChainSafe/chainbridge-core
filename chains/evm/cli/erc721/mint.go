package erc721

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/contracts"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/erc721"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint token on an ERC721 mintable contract",
	Long:  "Mint token on an ERC721 mintable contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		erc721Contract, err := contracts.InitializeErc721Contract(
			url, gasLimit, gasPrice, senderKeyPair, erc721Addr,
		)
		if err != nil {
			return err
		}
		return MintCmd(cmd, args, erc721Contract)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateMintFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessMintFlags(cmd, args)
		return err
	},
}

func init() {
	BindMintFlags(mintCmd)
}

func BindMintFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Erc721Address, "contract-address", "", "address of contract")
	cmd.Flags().StringVar(&DstAddress, "destination-address", "", "address of recipient")
	cmd.Flags().StringVar(&TokenId, "tokenId", "", "ERC721 token ID")
	cmd.Flags().StringVar(&Metadata, "metadata", "", "ERC721 token metadata")
	flags.MarkFlagsAsRequired(cmd, "contract-address", "destination-address", "tokenId", "metadata", "contract-address")
}

func ValidateMintFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc721Address) {
		return fmt.Errorf("invalid ERC721 contract address %s", Erc721Address)
	}
	if !common.IsHexAddress(DstAddress) {
		return fmt.Errorf("invalid recipient address %s", DstAddress)
	}
	return nil
}

func ProcessMintFlags(cmd *cobra.Command, args []string) error {
	var err error
	erc721Addr = common.HexToAddress(Erc721Address)

	if !common.IsHexAddress(DstAddress) {
		dstAddress = senderKeyPair.CommonAddress()
	} else {
		dstAddress = common.HexToAddress(DstAddress)
	}

	var ok bool
	if tokenId, ok = big.NewInt(0).SetString(TokenId, 10); !ok {
		return fmt.Errorf("invalid token id value")
	}

	return err
}

func MintCmd(cmd *cobra.Command, args []string, erc721Contract *erc721.ERC721Contract) error {
	_, err = erc721Contract.Mint(tokenId, Metadata, dstAddress, transactor.TransactOptions{})
	if err != nil {
		return err
	}

	log.Info().Msgf("%v token minted", tokenId)
	return err
}
