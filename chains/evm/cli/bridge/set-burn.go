package bridge

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var setBurnCmd = &cobra.Command{
	Use:   "set-burn",
	Short: "Set a token contract as mintable/burnable",
	Long:  "The set-burn subcommand sets a token contract as mintable/burnable in a handler",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return SetBurnCmd(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateSetBurnFlags(cmd, args)
		if err != nil {
			return err
		}

		ProcessSetBurnFlags(cmd, args)
		return nil
	},
}

func BindSetBurnCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Handler, "handler", "", "ERC20 handler contract address")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	cmd.Flags().StringVar(&TokenContract, "token", "", "Token contract to be registered")
	flags.MarkFlagsAsRequired(cmd, "handler", "bridge", "token")
}

func init() {
	BindSetBurnCmdFlags(setBurnCmd)
}
func ValidateSetBurnFlags(cmd *cobra.Command, args []string) error {

	if !common.IsHexAddress(Handler) {
		return fmt.Errorf("invalid handler address %s", Handler)
	}
	if !common.IsHexAddress(TokenContract) {
		return fmt.Errorf("invalid token contract address %s", TokenContract)
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func ProcessSetBurnFlags(cmd *cobra.Command, args []string) {
	handlerAddr = common.HexToAddress(Handler)
	bridgeAddr = common.HexToAddress(Bridge)
	tokenContractAddr = common.HexToAddress(TokenContract)
}
func SetBurnCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}
	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})

	log.Info().Msgf("Setting contract %s as burnable on handler %s", tokenContractAddr.String(), handlerAddr.String())
	setBurnableInput, err := calls.PrepareSetBurnableInput(handlerAddr, tokenContractAddr)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	_, err = calls.Transact(ethClient, txFabric, gasPricer, &bridgeAddr, setBurnableInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(err)
		return err
	}
	log.Info().Msg("Burnable set")
	return nil
}
