package erc20

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var addMinterCmd = &cobra.Command{
	Use:   "add-minter",
	Short: "Add a minter to an Erc20 mintable contract",
	Long:  "Add a minter to an Erc20 mintable contract",
	Run:   addMinter,
}

func init() {
	addMinterCmd.Flags().String("erc20Address", "", "ERC20 contract address")
	addMinterCmd.Flags().String("minter", "", "address of minter")
}

func addMinter(cmd *cobra.Command, args []string) {
	erc20Address := cmd.Flag("erc20Address").Value
	minterAddress := cmd.Flag("minter").Value
	log.Debug().Msgf(`
Adding minter
Minter address: %s 
ERC20 address: %s`, minterAddress, erc20Address)
}

/*
func addMinter(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	erc20 := cctx.String("erc20Address")
	if !common.IsHexAddress(erc20) {
		return errors.New("invalid erc20Address address")
	}
	erc20Address := common.HexToAddress(erc20)

	minter := cctx.String("minter")
	if !common.IsHexAddress(minter) {
		return errors.New("invalid minter address")
	}
	minterAddress := common.HexToAddress(minter)

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.ERC20AddMinter(ethClient, erc20Address, minterAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("%s account granted minter roles", minterAddress.String())
	return nil
}
*/
