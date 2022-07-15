package admin

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

var setDepositNonceCmd = &cobra.Command{
	Use:   "set-deposit-nonce",
	Short: "Set the deposit nonce",
	Long:  "The set-deposit-nonce subcommand sets the deposit nonce. This nonce cannot be less than what is currently stored in the contract",
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
		return SetDepositNonceEVMCMD(cmd, args, bridge.NewBridgeContract(c, BridgeAddr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateSetDepositNonceFlags(cmd, args)
		if err != nil {
			return err
		}

		ProcessSetDepositNonceFlags(cmd, args)
		return nil
	},
}

func BindSetDepositNonceFlags(cmd *cobra.Command) {
	cmd.Flags().Uint8Var(&DomainID, "domain", 0, "Domain ID of chain")
	cmd.Flags().Uint64Var(&DepositNonce, "deposit-nonce", 0, "Deposit nonce to set (does not decrement)")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	flags.MarkFlagsAsRequired(cmd, "domain", "deposit-nonce", "bridge")
}

func init() {
	BindSetDepositNonceFlags(setDepositNonceCmd)
}

func ValidateSetDepositNonceFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func ProcessSetDepositNonceFlags(cmd *cobra.Command, args []string) {
	BridgeAddr = common.HexToAddress(Bridge)
}

func SetDepositNonceEVMCMD(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	log.Debug().Msgf(`
Set Deposit Nonce
Domain ID: %v
Deposit Nonce: %v
Bridge Address: %s`, DomainID, DepositNonce, Bridge)
	_, err := contract.SetDepositNonce(DomainID, DepositNonce, transactor.TransactOptions{GasLimit: gasLimit})
	if err != nil {
		return err
	}
	log.Info().Msgf("[domain ID: %v] successfully set nonce: %v at address: %s", DomainID, DepositNonce, BridgeAddr.String())
	return nil
}
