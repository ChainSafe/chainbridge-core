package erc721

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var addMinterCmd = &cobra.Command{
	Use:   "add-minter",
	Short: "Add a minter to an ERC721 mintable contract",
	Long:  "Add a minter to an ERC721 mintable contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return AddMinterCmd(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateAddMinterFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessAddMinterFlags(cmd, args)
		return err
	},
}

func BindAddMinterCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Erc721Address, "contract-address", "", "address of contract")
	cmd.Flags().StringVar(&Minter, "minter", "", "minter address")
}

func init() {
	BindAddMinterCmdFlags(addMinterCmd)
}

func ValidateAddMinterFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc721Address) {
		return fmt.Errorf("invalid ERC721 contract address %s", Erc721Address)
	}
	if !common.IsHexAddress(Minter) {
		return fmt.Errorf("invalid minter address %s", Minter)
	}
	return nil
}

func ProcessAddMinterFlags(cmd *cobra.Command, args []string) error {
	erc721Addr = common.HexToAddress(Erc721Address)
	minterAddr = common.HexToAddress(Minter)
	return nil
}

func AddMinterCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	ethClient, err := evmclient.NewEVMClientFromParams(
		url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})

	_, err = calls.ERC721AddMinter(ethClient, txFabric, gasPricer.(calls.GasPricer), gasLimit, erc721Addr, minterAddr)
	if err != nil {
		return err
	}
	log.Debug().Msgf(`
	Adding minter
	Minter address: %s
	ERC721 address: %s`,
		minterAddr, erc721Addr)
	return err
}
