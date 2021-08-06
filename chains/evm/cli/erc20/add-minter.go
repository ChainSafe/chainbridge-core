package erc20

import (
	"errors"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/cliutils"
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
	erc20Address := cmd.Flag("erc20Address").Value
	minterAddress := cmd.Flag("minter").Value
	log.Debug().Msgf(`
Adding minter
Minter address: %s 
ERC20 address: %s`, minterAddress, erc20Address)

	url := cmd.Flag("url").Value.String()

	erc20 := cmd.Flag("erc20Address").Value.String()
	if !common.IsHexAddress(erc20) {
		log.Fatal().Err(errors.New("invalid erc20Address address"))
	}
	erc20Addr := common.HexToAddress(erc20)

	minter := cmd.Flag("minter").Value.String()
	if !common.IsHexAddress(minter) {
		log.Fatal().Err(errors.New("invalid minter address"))
	}
	minterAddr := common.HexToAddress(minter)

	senderKeyPair, err := cliutils.DefineSender(cmd)
	if err != nil {
		log.Fatal().Err(err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Fatal().Err(err)
	}
	mintableInput, err := calls.PrepareErc20AddMinterInput(ethClient, erc20Addr, minterAddr)
	if err != nil {
		log.Fatal().Err(err)
	}
	_, err = calls.SendInput(ethClient, minterAddr, mintableInput, txFabric)
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msgf("%s account granted minter roles", minterAddress.String())
	return nil
}
