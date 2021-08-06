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
	RunE:  CallAddMinter,
}

func init() {
	addMinterCmd.Flags().String("erc20Address", "", "ERC20 contract address")
	addMinterCmd.Flags().String("minter", "", "address of minter")
}

func CallAddMinter(cmd *cobra.Command, args []string) error {
	txFabric := evmtransaction.NewTransaction
	return addMinter(cmd, args, txFabric)
}
func addMinter(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	erc20Address := cmd.Flag("erc20Address").Value.String()
	minterAddress := cmd.Flag("minter").Value.String()

	// fetch global flag values
	url, _, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	if !common.IsHexAddress(erc20Address) {
		err := errors.New("invalid erc20Address address")
		log.Error().Err(err)
	}
	erc20Addr := common.HexToAddress(erc20Address)

	if !common.IsHexAddress(minterAddress) {
		err := errors.New("invalid minter address")
		log.Error().Err(err)
	}
	minterAddr := common.HexToAddress(minterAddress)

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	mintableInput, err := calls.PrepareErc20AddMinterInput(ethClient, erc20Addr, minterAddr)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	_, err = calls.SendInput(ethClient, minterAddr, mintableInput, txFabric)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("%s account granted minter roles", minterAddr.String())
	return nil
}
