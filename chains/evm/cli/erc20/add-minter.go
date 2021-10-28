package erc20

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/writer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var addMinterCmd = &cobra.Command{
	Use:   "add-minter",
	Short: "Add a minter to an Erc20 mintable contract",
	Long:  "Add a minter to an Erc20 mintable contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return AddMinterCmd(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateAddMinterFlags(cmd, args)
		if err != nil {
			return err
		}
		ProcessAddMinterFlags(cmd, args)
		return nil
	},
}

func BindAddMinterCmdFlags() {
	addMinterCmd.Flags().StringVar(&Erc20Address, "erc20Address", "", "ERC20 contract address")
	addMinterCmd.Flags().StringVar(&Minter, "minter", "", "handler contract address")
	flags.MarkFlagsAsRequired(addMinterCmd, "erc20Address", "minter")
}

func init() {
	BindAddMinterCmdFlags()
}

func ValidateAddMinterFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc20Address) {
		return errors.New("invalid erc20Address address")
	}
	if !common.IsHexAddress(Minter) {
		return errors.New("invalid minter address")
	}
	return nil
}

func ProcessAddMinterFlags(cmd *cobra.Command, args []string) {
	erc20Addr = common.HexToAddress(Erc20Address)
	minterAddr = common.HexToAddress(Minter)
}

func AddMinterCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}
	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})
	mintableInput, err := calls.PrepareErc20AddMinterInput(ethClient, erc20Addr, minterAddr)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	_, err = calls.Transact(ethClient, txFabric, gasPricer, &erc20Addr, mintableInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("%s account granted minter roles", minterAddr.String())
	writer.WriteCliDataToFile(cmd)
	return nil
}
