package erc20

import (
	"errors"
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
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
		return AddMinterCmd(cmd, args, txFabric)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := validateAddMinterFlags(cmd, args)
		if err != nil {
			return err
		}
		processAddMinterFlags(cmd, args)

		return nil
	},
}

func BindAddMinterCmdFlags() {
	addMinterCmd.Flags().StringVarP(&Erc20Address, "erc20Address", "erc20add", "", "ERC20 contract address")
	addMinterCmd.Flags().StringVarP(&Minter, "minter", "m", "", "address of minter")
}

func init() {
	BindAddMinterCmdFlags()
}

func validateAddMinterFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc20Address) {
		return errors.New("invalid erc20Address address")
	}
	if !common.IsHexAddress(Minter) {
		return errors.New("invalid minter address")
	}
	return nil
}

var minterAddr common.Address

func processAddMinterFlags(cmd *cobra.Command, args []string) {
	erc20Addr = common.HexToAddress(Erc20Address)
	minterAddr = common.HexToAddress(Minter)
}

func AddMinterCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}
	mintableInput, err := calls.PrepareErc20AddMinterInput(ethClient, erc20Addr, minterAddr)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	_, err = calls.Transact(ethClient, txFabric, &erc20Addr, mintableInput, gasLimit)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("%s account granted minter roles", minterAddr.String())
	return nil
}
