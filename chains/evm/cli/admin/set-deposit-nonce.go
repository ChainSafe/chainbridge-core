package admin

import (
	"errors"
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
}

func BindSetDepositNonceFlags(cli *cobra.Command) {
	cli.Flags().Uint8("domainId", 0, "domain ID of chain")
	cli.Flags().Uint64("depositNonce", 0, "deposit nonce to set (does not decrement)")
	cli.Flags().String("bridgeAddress", "", "bridge contract address")
	cli.MarkFlagRequired("domainId")
	cli.MarkFlagRequired("depositNonce")
}

func init() {
	BindSetDepositNonceFlags(setDepositNonceCmd)
}

func SetDepositNonceEVMCMD(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	domainId, err := cmd.Flags().GetUint8("domainId")
	if err != nil {
		return err
	}

	depositNonce, err := cmd.Flags().GetUint64("depositNonce")
	if err != nil {
		return err
	}

	bridgeAddress := cmd.Flag("bridgeAddress").Value.String()

	log.Debug().Msgf(`
Set Deposit Nonce
Domain ID: %v
Deposit Nonce: %v
Bridge Address: %s`, domainId, depositNonce, bridgeAddress)

	if !common.IsHexAddress(bridgeAddress) {
		return errors.New("invalid bridge address")
	}
	bridgeAddr := common.HexToAddress(bridgeAddress)

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

	setDepositNonceInput, err := calls.PrepareSetDepositNonceInput(domainId, depositNonce)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare set deposit nonce input error: %v", err))
		return err
	}

	_, err = calls.Transact(ethClient, txFabric, &bridgeAddr, setDepositNonceInput, gasLimit)
	if err != nil {
		log.Error().Err(fmt.Errorf("transact error: %v", err))
		return err
	}
	log.Info().Msgf("[domain ID: %v] successfully set nonce: %v at address: %s", domainId, depositNonce, bridgeAddr.String())
	return nil
}
