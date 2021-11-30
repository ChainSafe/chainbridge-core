package erc721

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/contracts"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/erc721"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
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
		erc721Contract, err := contracts.InitializeErc721Contract(
			url, gasLimit, gasPrice, senderKeyPair, erc721Addr,
		)
		if err != nil {
			return err
		}
		return AddMinterCmd(cmd, args, erc721Contract)
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

func AddMinterCmd(cmd *cobra.Command, args []string, erc721Contract *erc721.ERC721Contract) error {
	_, err = erc721Contract.AddMinter(minterAddr, transactor.TransactOptions{})
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
