package erc721

import (
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
	Short: "Add a minter to an ERC721 mintable contract",
	Long:  "Add a minter to an ERC721 mintable contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return AddMinterCmd(cmd, args, txFabric)
	},
}

func init() {
	addMinterCmd.Flags().String("erc721Address", "", "ERC721 contract address")
	addMinterCmd.Flags().String("minter", "", "address of minter")
}

func AddMinterCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(
		url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	erc721Address := cmd.Flag("erc721Address").Value.String()
	if !common.IsHexAddress(erc721Address) {
		return fmt.Errorf("invalid erc20Address address")
	}
	erc721Addr := common.HexToAddress(erc721Address)

	minterAddress := cmd.Flag("minter").Value.String()
	if !common.IsHexAddress(minterAddress) {
		return fmt.Errorf("invalid erc20Address address")
	}
	minterAddr := common.HexToAddress(minterAddress)

	addMinterInput, err := calls.PrepareErc721AddMinterInput(ethClient, erc721Addr, minterAddr)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	_, err = calls.Transact(ethClient, txFabric, &erc721Addr, addMinterInput, gasLimit)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Debug().Msgf(`
	Adding minter
	Minter address: %s
	ERC721 address: %s`,
		minterAddress, erc721Address)

	return nil
}
