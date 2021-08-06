package erc721

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var addMinterCmd = &cobra.Command{
	Use:   "add-minter",
	Short: "Add a minter to an ERC721 mintable contract",
	Long:  "Add a minter to an ERC721 mintable contract",
	Run:   addMinter,
}

func init() {
	addMinterCmd.Flags().String("erc721Address", "", "ERC721 contract address")
	addMinterCmd.Flags().String("minter", "", "address of minter")
}

func addMinter(cmd *cobra.Command, args []string) {
	erc721Address := cmd.Flag("erc721Address").Value
	minterAddress := cmd.Flag("minter").Value
	log.Debug().Msgf(`
Adding minter
Minter address: %s
ERC721 address: %s`, minterAddress, erc721Address)
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
	erc721 := cctx.String("erc721Address")
	if !common.IsHexAddress(erc721) {
		return errors.New("invalid erc20Address address")
	}
	erc721Address := common.HexToAddress(erc721)

	minter := cctx.String("minter")
	if !common.IsHexAddress(minter) {
		return errors.New("invalid minter address")
	}
	minterAddress := common.HexToAddress(minter)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.ERC721AddMinter(ethClient, erc721Address, minterAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("Minter with address %s added", minterAddress.String())
	return nil
}
*/
