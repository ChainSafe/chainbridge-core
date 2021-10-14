package erc721

import (
	"encoding/hex"
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

var (
	Erc721ContractAddressForMint string
	DestinationAddressForMint    string
	TokenIdForMint               string
	MetadataForMint              string
)

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint ERC721 token",
	Long:  "Mint ERC721 token",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return MintCmd(cmd, args, txFabric)
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
	mintCmd.Flags().StringVarP(&Erc721ContractAddressForMint, "contract-address", "con", "", "address of contract")
	mintCmd.Flags().StringVarP(&DestinationAddressForMint, "destination-address", "dest", "", "address of recipient")
	mintCmd.Flags().StringVarP(&TokenIdForMint, "token-id", "tid", "", "token id")
	mintCmd.Flags().StringVarP(&MetadataForMint, "metadata", "met", "", "token metadata")
}

func validateFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc721ContractAddressForMint) {
		return fmt.Errorf("invalid ERC721 contract address %s", Erc721ContractAddressForMint)
	}
	if !common.IsHexAddress(DestinationAddressForMint) {
		return fmt.Errorf("invalid recipient address %s", DestinationAddressForMint)
	}
	return nil
}

type MintFlags struct {
	erc721ContractAddress common.Address
	destinationAddress    common.Address
	tokenId               *big.Int
	metadata              []byte
}

var (
	mintFlags MintFlags
)

func MintCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	mintFlags.erc721ContractAddress = common.HexToAddress(Erc721ContractAddressForMint)

	if !common.IsHexAddress(DestinationAddressForMint) {
		mintFlags.destinationAddress = senderKeyPair.CommonAddress()
	} else {
		mintFlags.destinationAddress = common.HexToAddress(DestinationAddressForMint)
	}

	var ok bool
	if mintFlags.tokenId, ok = big.NewInt(0).SetString(TokenIdForMint, 10); !ok {
		return fmt.Errorf("invalid token id value")
	}

	if MetadataForMint[0:2] == "0x" {
		MetadataForMint = MetadataForMint[2:]
	}
	mintFlags.metadata, err = hex.DecodeString(MetadataForMint)
	if err != nil {
		return err
	}

	ethclient, err := evmclient.NewEVMClientFromParams(
		url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	mintTokenInput, err := calls.PrepareERC721MintTokensInput(mintFlags.destinationAddress, mintFlags.tokenId, mintFlags.metadata)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc721 mint input error: %v", err))
		return err
	}

	_, err = calls.Transact(ethclient, txFabric, &mintFlags.erc721ContractAddress, mintTokenInput, gasLimit)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("%v token minted", mintFlags.tokenId)
	return nil
}
