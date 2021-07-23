package erc20

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var AllowanceCmd = &cobra.Command{
	Use:   "allowance",
	Short: "Get the allowance of a spender for an address",
	Long:  "Get the allowance of a spender for an address",
	Run:   allowance,
}

func init() {
	AllowanceCmd.Flags().String("erc20Address", "", "ERC20 contract address")
	AllowanceCmd.Flags().String("owner", "", "address of token owner")
	AllowanceCmd.Flags().String("spender", "", "address of spender")
}

func allowance(cmd *cobra.Command, args []string) {
	erc20Address := cmd.Flag("erc20Address").Value
	ownerAddress := cmd.Flag("owner").Value
	spenderAddress := cmd.Flag("spender").Value
	log.Debug().Msgf(`
Determing allowance
ERC20 address: %s
Owner address: %s
Spender address: %s`,
		erc20Address, ownerAddress, spenderAddress)
}

/*
func allowance(cctx *cli.Context) error {
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

	spender := cctx.String("spender")
	if !common.IsHexAddress(spender) {
		return errors.New("invalid spender address")
	}
	spenderAddress := common.HexToAddress(spender)

	owner := cctx.String("owner")
	if !common.IsHexAddress(owner) {
		return errors.New("invalid owner address")
	}
	ownerAddress := common.HexToAddress(owner)

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	balance, err := utils.ERC20Allowance(ethClient, erc20Address, spenderAddress, ownerAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("allowance of %s to spend from address %s is %s", spenderAddress.String(), ownerAddress.String(), balance.String())
	return nil
}
*/
