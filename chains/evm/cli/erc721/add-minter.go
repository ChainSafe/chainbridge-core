package erc721

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/erc721"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/util"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var addMinterCmd = &cobra.Command{
	Use:   "add-minter",
	Short: "Add a new ERC721 minter",
	Long:  "The add-minter subcommand adds a new minter address to an ERC721 mintable contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, evmtransaction.NewTransaction, c, prepare)
		if err != nil {
			return err
		}
		return AddMinterCmd(cmd, args, erc721.NewErc721Contract(c, Erc721Addr, t))
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

func BindAddMinterFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Erc721Address, "contract", "", "ERC721 contract address")
	cmd.Flags().StringVar(&Minter, "minter", "", "Minter address")
}

func init() {
	BindAddMinterFlags(addMinterCmd)
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
	Erc721Addr = common.HexToAddress(Erc721Address)
	MinterAddr = common.HexToAddress(Minter)
	return nil
}

func AddMinterCmd(cmd *cobra.Command, args []string, erc721Contract *erc721.ERC721Contract) error {
	_, err = erc721Contract.AddMinter(
		MinterAddr, transactor.TransactOptions{GasLimit: gasLimit},
	)
	if err != nil {
		return err
	}
	log.Debug().Msgf(`
	Adding minter
	Minter address: %s
	ERC721 address: %s`,
		MinterAddr, Erc721Addr)
	return err
}
