package erc721

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint ERC721 token",
	Long:  "Mint ERC721 token",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return MintCmd(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := validateFlags(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	mintCmd.Flags().StringVarP(&Erc721Address, "contract-address", "con", "", "address of contract")
	mintCmd.Flags().StringVarP(&DstAddress, "destination-address", "dest", "", "address of recipient")
	mintCmd.Flags().StringVarP(&TokenId, "token-id", "tid", "", "token id")
	mintCmd.Flags().StringVarP(&Metadata, "metadata", "met", "", "token metadata")
}

func validateFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc721Address) {
		return fmt.Errorf("invalid ERC721 contract address %s", Erc721Address)
	}
	if !common.IsHexAddress(DstAddress) {
		return fmt.Errorf("invalid recipient address %s", DstAddress)
	}
	return nil
}

func processMintFlags(cmd *cobra.Command, args []string) error {
	var err error

	erc721Addr = common.HexToAddress(DstAddress)

	if !common.IsHexAddress(DstAddress) {
		dstAddress = senderKeyPair.CommonAddress()
	} else {
		dstAddress = common.HexToAddress(DstAddress)
	}

	var ok bool
	if tokenId, ok = big.NewInt(0).SetString(TokenId, 10); !ok {
		return fmt.Errorf("invalid token id value")
	}

	if Metadata[0:2] == "0x" {
		Metadata = Metadata[2:]
	}
	metadata, err = hex.DecodeString(Metadata)
	return err
}

func MintCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {

	ethClient, err := evmclient.NewEVMClientFromParams(
		url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})

	mintTokenInput, err := calls.PrepareERC721MintTokensInput(dstAddress, tokenId, metadata)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc721 mint input error: %v", err))
		return err
	}

	_, err = calls.Transact(ethClient, txFabric, gasPricer, &erc721Addr, mintTokenInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("%v token minted", tokenId)
	return nil
}
