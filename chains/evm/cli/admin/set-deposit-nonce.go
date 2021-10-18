package admin

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var setDepositNonceCmd = &cobra.Command{
	Use:   "set-deposit-nonce",
	Short: "Set the deposit nonce",
	Long: `Set the deposit nonce

This nonce cannot be less than what is currently stored in the contract`,
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return SetDepositNonceEVMCMD(cmd, args, txFabric)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := validateSetDepositNonceFlags(cmd, args)
		if err != nil {
			return err
		}

		processSetDepositNonceFlags(cmd, args)
		return nil
	},
}

func BindSetDepositNonceFlags() {
	setDepositNonceCmd.Flags().Uint8VarP(&DomainID, "domainId", "dID", 0, "domain ID of chain")
	setDepositNonceCmd.Flags().Uint64VarP(&DepositNonce, "depositNonce", "dn", 0, "deposit nonce to set (does not decrement)")
	setDepositNonceCmd.Flags().StringVarP(&Bridge, "bridge", "b", "", "bridge contract address")
	flags.CheckRequiredFlags(setDepositNonceCmd, "domainId", "depositNonce", "bridge")
}

func init() {
	BindSetDepositNonceFlags()
}

func validateSetDepositNonceFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Relayer) {
		return fmt.Errorf("invalid relayer address %s", Relayer)
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func processSetDepositNonceFlags(cmd *cobra.Command, args []string) {
	relayerAddr = common.HexToAddress(Relayer)
	bridgeAddr = common.HexToAddress(Bridge)
}

func SetDepositNonceEVMCMD(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {

	log.Debug().Msgf(`
Set Deposit Nonce
Domain ID: %v
Deposit Nonce: %v
Bridge Address: %s`, DomainID, DepositNonce, Bridge)

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	setDepositNonceInput, err := calls.PrepareSetDepositNonceInput(DomainID, DepositNonce)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare set deposit nonce input error: %v", err))
		return err
	}

	_, err = calls.Transact(ethClient, txFabric, &bridgeAddr, setDepositNonceInput, gasLimit)
	if err != nil {
		log.Error().Err(fmt.Errorf("transact error: %v", err))
		return err
	}
	log.Info().Msgf("[domain ID: %v] successfully set nonce: %v at address: %s", DomainID, DepositNonce, bridgeAddr.String())
	return nil
}
