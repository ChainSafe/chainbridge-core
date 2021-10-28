package erc721

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var approveCmd = &cobra.Command{
	Use:   "mint",
	Short: "Approve token in an ERC721 contract for transfer.",
	Long:  "Approve token in an ERC721 contract for transfer.",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return ApproveCmd(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
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

func BindApproveCmdFlags(cli *cobra.Command) {
	mintCmd.Flags().StringVar(&Erc721Address, "contract-address", "", "address of contract")
	mintCmd.Flags().StringVar(&Recipient, "recipient", "", "address of recipient")
	mintCmd.Flags().StringVar(&TokenId, "tokenId", "", "ERC721 token ID")
	flags.MarkFlagsAsRequired(mintCmd, "contract-address", "recipient", "tokenId")
}

func init() {
	BindApproveCmdFlags(approveCmd)
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
	recipientAddr = common.HexToAddress(Recipient)
	erc721Addr = common.HexToAddress(Erc721Address)

	var ok bool
	if tokenId, ok = big.NewInt(0).SetString(TokenId, 10); !ok {
		return fmt.Errorf("invalid token id value")
	}
	return nil
}

func ApproveCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	ethclient, err := evmclient.NewEVMClientFromParams(
		url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	approveTokenInput, err := calls.PrepareERC721ApproveInput(recipientAddr, tokenId)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc721 approve input error: %v", err))
		return err
	}

	_, err = calls.Transact(ethclient, txFabric, gasPricer, &erc721Addr, approveTokenInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("%v token approved", tokenId)
	return nil
}
