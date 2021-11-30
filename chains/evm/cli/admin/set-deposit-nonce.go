package admin

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/contracts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var setDepositNonceCmd = &cobra.Command{
	Use:   "set-deposit-nonce",
	Short: "Set the deposit nonce",
	Long: `Set the deposit nonce

This nonce cannot be less than what is currently stored in the contract`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		bridgeContract, err := contracts.InitializeBridgeContract(
			url, gasLimit, gasPrice, senderKeyPair, bridgeAddr,
		)
		if err != nil {
			return err
		}
		return SetDepositNonceEVMCMD(cmd, args, bridgeContract)
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
	cmd.Flags().Uint8Var(&DomainID, "domainId", 0, "domain ID of chain")
	cmd.Flags().Uint64Var(&DepositNonce, "depositNonce", 0, "deposit nonce to set (does not decrement)")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	flags.MarkFlagsAsRequired(cmd, "domainId", "depositNonce", "bridge")
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
	bridgeAddr = common.HexToAddress(Bridge)
}

func SetDepositNonceEVMCMD(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	log.Debug().Msgf(`
Set Deposit Nonce
Domain ID: %v
Deposit Nonce: %v
Bridge Address: %s`, DomainID, DepositNonce, Bridge)
	_, err := contract.SetDepositNonce(DomainID, DepositNonce, transactor.TransactOptions{})
	if err != nil {
		return err
	}
	log.Info().Msgf("[domain ID: %v] successfully set nonce: %v at address: %s", DomainID, DepositNonce, bridgeAddr.String())
	return nil
}
