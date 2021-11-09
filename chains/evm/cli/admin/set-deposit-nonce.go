package admin

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
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
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return SetDepositNonceEVMCMD(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
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

func SetDepositNonceEVMCMD(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {

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

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}
	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})
	setDepositNonceInput, err := calls.PrepareSetDepositNonceInput(DomainID, DepositNonce)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare set deposit nonce input error: %v", err))
		return err
	}

	_, err = calls.Transact(ethClient, txFabric, gasPricer, &bridgeAddr, setDepositNonceInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(fmt.Errorf("transact error: %v", err))
		return err
	}
	log.Info().Msgf("[domain ID: %v] successfully set nonce: %v at address: %s", DomainID, DepositNonce, bridgeAddr.String())
	return nil
}
