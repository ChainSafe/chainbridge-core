package bridge

import (
	"encoding/hex"
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

var registerGenericResourceCmd = &cobra.Command{
	Use:   "register-generic-resource",
	Short: "Register a generic resource ID",
	Long:  "Register a resource ID with a contract address for a generic handler",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateGenericResourceFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessGenericResourceFlags(cmd, args)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return RegisterGenericResource(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
	},
}

func BindRegisterGenericResourceFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Handler, "handler", "", "handler contract address")
	cmd.Flags().StringVar(&ResourceID, "resourceId", "", "resource ID to query")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	cmd.Flags().StringVar(&Target, "target", "", "contract address to be registered") // TODO change the description (target is not necessary a contract address, could be hash storage)
	cmd.Flags().StringVar(&Deposit, "deposit", "0x00000000", "deposit function signature")
	cmd.Flags().StringVar(&Execute, "execute", "0x00000000", "execute proposal function signature")
	cmd.Flags().BoolVar(&Hash, "hash", false, "treat signature inputs as function prototype strings, hash and take the first 4 bytes")
	flags.MarkFlagsAsRequired(cmd, "handler", "resourceId", "bridge", "target")
}

func init() {
	BindRegisterGenericResourceFlags(registerGenericResourceCmd)
}

func ValidateGenericResourceFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Handler) {
		return fmt.Errorf("invalid handler address %s", Handler)
	}

	if !common.IsHexAddress(Target) {
		return fmt.Errorf("invalid target address %s", Target)
	}

	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Target)
	}

	return nil
}

func ProcessGenericResourceFlags(cmd *cobra.Command, args []string) error {
	handlerAddr = common.HexToAddress(Handler)
	targetContractAddr = common.HexToAddress(Target)
	bridgeAddr = common.HexToAddress(Bridge)

	if ResourceID[0:2] == "0x" {
		ResourceID = ResourceID[2:]
	}

	resourceIdBytes, err := hex.DecodeString(ResourceID)
	if err != nil {
		return err
	}

	resourceIdBytesArr = calls.SliceTo32Bytes(resourceIdBytes)

	if Hash {
		depositSigBytes = calls.GetSolidityFunctionSig([]byte(Deposit))
		executeSigBytes = calls.GetSolidityFunctionSig([]byte(Execute))
	} else {
		copy(depositSigBytes[:], []byte(Deposit)[:])
		copy(executeSigBytes[:], []byte(Execute)[:])
	}

	return nil
}

func RegisterGenericResource(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	url, _, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	log.Info().Msgf("Registering contract %s with resource ID %s on handler %s", targetContractAddr, ResourceID, handlerAddr)

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}
	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})

	h, err := calls.AdminSetGenericResource(
		ethClient,
		txFabric,
		gasPricer,
		handlerAddr,
		resourceIdBytesArr,
		targetContractAddr,
		depositSigBytes,
		big.NewInt(int64(DepositerOffset)),
		executeSigBytes,
	)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("Generic resource registered with hash: %s", h.Hex())
	return nil
}
