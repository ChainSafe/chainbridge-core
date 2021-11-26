package erc20

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var getAllowanceCmd = &cobra.Command{
	Use:   "get-allowance",
	Short: "Get the allowance of a spender for an address",
	Long:  "The get-allowance subcommand returns the allowance of a spender for an address",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return GetAllowanceCmd(cmd, args, txFabric)
	},
}

func BindGetAllowanceCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Erc20Address, "contract", "", "ERC20 contract address")
	cmd.Flags().StringVar(&OwnerAddress, "owner", "", "Address of token owner")
	cmd.Flags().StringVar(&SpenderAddress, "spender", "", "Address of spender")
	flags.MarkFlagsAsRequired(cmd, "contract", "owner", "spender")
}

func init() {
	BindGetAllowanceCmdFlags(getAllowanceCmd)
}

func GetAllowanceCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	log.Debug().Msgf(`
Determing allowance
ERC20 address: %s
Owner address: %s
Spender address: %s`,
		Erc20Address, OwnerAddress, SpenderAddress)
	return nil

	/*
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
			log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
			return err
		}
		balance, err := utils.ERC20Allowance(ethClient, erc20Address, spenderAddress, ownerAddress)
		if err != nil {
			return err
		}
		log.Info().Msgf("allowance of %s to spend from address %s is %s", spenderAddress.String(), ownerAddress.String(), balance.String())
		return nil
	*/
}
