package erc20

import (
	"errors"
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtypes"

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
}

func BindAddMinterCmdFlags(cli *cobra.Command) {
	cli.Flags().String("erc20Address", "", "ERC20 contract address")
	cli.Flags().String("minter", "", "address of minter")
}

func init() {
	BindAddMinterCmdFlags(addMinterCmd)
}

func AddMinterCmd(cmd *cobra.Command, args []string, txFabric evmtypes.TxFabric) error {
	erc20Address := cmd.Flag("erc20Address").Value.String()
	minterAddress := cmd.Flag("minter").Value.String()

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	if !common.IsHexAddress(erc20Address) {
		err := errors.New("invalid erc20Address address")
		log.Error().Err(err)
		return err
	}
	erc20Addr := common.HexToAddress(erc20Address)

	if !common.IsHexAddress(minterAddress) {
		err := errors.New("invalid minter address")
		log.Error().Err(err)
		return err
	}
	minterAddr := common.HexToAddress(minterAddress)

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
