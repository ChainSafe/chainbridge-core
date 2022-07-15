package bridge

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ChainSafe/sygma-core/util"
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
		return SetBurnCmd(cmd, args, bridge.NewBridgeContract(c, BridgeAddr, t))
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

func BindSetBurnFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Handler, "handler", "", "ERC20 handler contract address")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	cmd.Flags().StringVar(&TokenContract, "token-contract", "", "Token contract to be registered")
	flags.MarkFlagsAsRequired(cmd, "handler", "bridge", "token-contract")
}

func init() {
	BindSetBurnFlags(setBurnCmd)
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
	HandlerAddr = common.HexToAddress(Handler)
	BridgeAddr = common.HexToAddress(Bridge)
	TokenContractAddr = common.HexToAddress(TokenContract)
}
func SetBurnCmd(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	log.Info().Msgf(
		"Setting contract %s as burnable on handler %s",
		TokenContractAddr.String(), HandlerAddr.String(),
	)
	_, err := contract.SetBurnableInput(
		HandlerAddr, TokenContractAddr, transactor.TransactOptions{GasLimit: gasLimit},
	)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	log.Info().Msg("Burnable set")
	return nil
}
