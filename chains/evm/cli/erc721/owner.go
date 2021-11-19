package erc721

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return OwnerCmd(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
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

func OwnerCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	ethClient, err := evmclient.NewEVMClientFromParams(
		url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})

	owner, err := calls.ERC721Owner(ethClient, tokenId, erc721Addr)
	if err != nil {
		return err
	}

	log.Info().Msgf("%v token owner: %v", tokenId, owner)
	return err
}
