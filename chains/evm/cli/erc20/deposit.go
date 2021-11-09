package erc20

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

var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Initiate a transfer of ERC20 tokens",
	Long:  "Initiate a transfer of ERC20 tokens",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return DepositCmd(cmd, args, evmtransaction.NewTransaction, &evmgaspricer.LondonGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateDepositFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessDepositFlags(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	BindDepositCmdFlags(depositCmd)
}

func BindDepositCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Recipient, "recipient", "", "address of recipient")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "address of bridge contract")
	cmd.Flags().StringVar(&Amount, "amount", "", "amount to deposit")
	cmd.Flags().Uint64Var(&DomainID, "domainId", 0, "destination domain ID")
	cmd.Flags().StringVar(&ResourceID, "resourceId", "", "resource ID for transfer")
	cmd.Flags().Uint64Var(&Decimals, "decimals", 0, "ERC20 token decimals")
	flags.MarkFlagsAsRequired(cmd, "recipient", "bridge", "amount", "domainId", "resourceId", "decimals")
}

func ValidateDepositFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Recipient) {
		return fmt.Errorf("invalid recipient address %s", Recipient)
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func ProcessDepositFlags(cmd *cobra.Command, args []string) error {
	var err error

	recipientAddress = common.HexToAddress(Recipient)
	decimals := big.NewInt(int64(Decimals))
	bridgeAddr = common.HexToAddress(Bridge)
	realAmount, err = calls.UserAmountToWei(Amount, decimals)
	if err != nil {
		return err
	}
	resourceIdBytesArr, err = flags.ProcessResourceID(ResourceID)
	return err
}

func DepositCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {

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

	data := calls.ConstructErc20DepositData(recipientAddress.Bytes(), realAmount)
	input, err := calls.PrepareDepositInput(uint8(DomainID), resourceIdBytesArr, data)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc20 deposit input error: %v", err))
		return err
	}

	// destinationId
	txHash, err := calls.Transact(ethClient, txFabric, gasPricer, &bridgeAddr, input, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(fmt.Errorf("erc20 deposit error: %v", err))
		return err
	}

	log.Info().Msgf("%s tokens were transferred to %s from %s with hash %s", Amount, recipientAddress.Hex(), senderKeyPair.CommonAddress().String(), txHash.Hex())
	return nil
}
