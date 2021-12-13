package bridge

import (
	"encoding/hex"
	"fmt"
	callsUtil "github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/chainbridge-core/util"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var registerGenericResourceCmd = &cobra.Command{
	Use:   "register-generic-resource",
	Short: "Register a generic resource ID",
	Long:  "The register-generic-resource subcommand registers a resource ID with a contract address for a generic handler",
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
		t, err := initialize.InitializeTransactor(gasPrice, evmtransaction.NewTransaction, c)
		if err != nil {
			return err
		}
		return RegisterGenericResource(cmd, args, bridge.NewBridgeContract(c, bridgeAddr, t))
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
}

func BindRegisterGenericResourceFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Handler, "handler", "", "Handler contract address")
	cmd.Flags().StringVar(&ResourceID, "resource", "", "Resource ID to query")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	cmd.Flags().StringVar(&Target, "target", "", "Contract address or hash storage to be registered")
	cmd.Flags().StringVar(&Deposit, "deposit", "0x00000000", "Deposit function signature")
	cmd.Flags().StringVar(&Execute, "execute", "0x00000000", "Execute proposal function signature")
	cmd.Flags().BoolVar(&Hash, "hash", false, "Treat signature inputs as function prototype strings, hash and take the first 4 bytes")
	flags.MarkFlagsAsRequired(cmd, "handler", "resource", "bridge", "target")
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

	resourceIdBytesArr = callsUtil.SliceTo32Bytes(resourceIdBytes)

	if Hash {
		depositSigBytes = callsUtil.GetSolidityFunctionSig([]byte(Deposit))
		executeSigBytes = callsUtil.GetSolidityFunctionSig([]byte(Execute))
	} else {
		copy(depositSigBytes[:], []byte(Deposit)[:])
		copy(executeSigBytes[:], []byte(Execute)[:])
	}

	return nil
}

func RegisterGenericResource(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	log.Info().Msgf("Registering contract %s with resource ID %s on handler %s", targetContractAddr, ResourceID, handlerAddr)

	h, err := contract.AdminSetGenericResource(
		handlerAddr,
		resourceIdBytesArr,
		targetContractAddr,
		depositSigBytes,
		big.NewInt(int64(DepositerOffset)),
		executeSigBytes,
		transactor.TransactOptions{GasLimit: gasLimit},
	)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("Generic resource registered with hash: %s", h.Hex())
	return nil
}
