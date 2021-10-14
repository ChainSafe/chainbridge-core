package erc721

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var approveCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint ERC721 token",
	Long:  "Mint ERC721 token",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return ApproveCmd(cmd, args, txFabric)
	},
}

func BindApproveCmdFlags(cli *cobra.Command) {
	cli.Flags().String("erc721Address", "", "ERC721 contract address")
	cli.Flags().String("recipientAddress", "", "todo")
	cli.Flags().Uint64("tokenId", 0, "ERC721 token id")
}

func init() {
	BindApproveCmdFlags(approveCmd)
}

func ApproveCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	erc721Address := cmd.Flag("erc721Address").Value.String()
	if !common.IsHexAddress(erc721Address) {
		return fmt.Errorf("invalid erc20Address address")
	}
	erc721Addr := common.HexToAddress(erc721Address)

	recipientAddress := cmd.Flag("recipient").Value.String()
	if !common.IsHexAddress(recipientAddress) {
		return fmt.Errorf("invalid recipient address")
	}
	recipientAddr := common.HexToAddress(recipientAddress)

	tokenIdAsString := cmd.Flag("tokenId").Value.String()
	tokenId, ok := big.NewInt(0).SetString(tokenIdAsString, 10)
	if !ok {
		return fmt.Errorf("invalid token id value")
	}

	ethclient, err := evmclient.NewEVMClientFromParams(
		url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	approveTokenInput, err := calls.PrepareERC721ApproveInput(recipientAddr, tokenId)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc721 approve input error: %v", err))
		return err
	}

	_, err = calls.Transact(ethclient, txFabric, &erc721Addr, approveTokenInput, gasLimit)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("%v token approved", tokenId)
	return nil
}
