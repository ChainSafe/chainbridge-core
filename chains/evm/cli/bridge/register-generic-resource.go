package bridge

import (
	"encoding/hex"
	"fmt"
	"math/big"

	callsUtil "github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/util"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/initialize"
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
		t, err := initialize.InitializeTransactor(gasPrice, evmtransaction.NewTransaction, c, prepare)
		if err != nil {
			return err
		}
		return RegisterGenericResource(cmd, args, bridge.NewBridgeContract(c, BridgeAddr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateRegisterGenericResourceFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessRegisterGenericResourceFlags(cmd, args)
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
	cmd.Flags().StringVar(&Deposit, "deposit", "00000000", "Deposit function signature")
	cmd.Flags().Uint64Var(&DepositerOffset, "depositerOffset", 0, "Offset to find the bridge tx depositer address inside the metadata sent on a deposit")
	cmd.Flags().StringVar(&Execute, "execute", "00000000", "Execute proposal function signature")
	cmd.Flags().BoolVar(&Hash, "hash", false, "Treat signature inputs as function prototype strings, hash and take the first 4 bytes")
	flags.MarkFlagsAsRequired(cmd, "handler", "resource", "bridge", "target")
}

func init() {
	BindRegisterGenericResourceFlags(registerGenericResourceCmd)
}

func ValidateRegisterGenericResourceFlags(cmd *cobra.Command, args []string) error {
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

func ProcessRegisterGenericResourceFlags(cmd *cobra.Command, args []string) error {
	HandlerAddr = common.HexToAddress(Handler)
	TargetContractAddr = common.HexToAddress(Target)
	BridgeAddr = common.HexToAddress(Bridge)

	if ResourceID[0:2] == "0x" {
		ResourceID = ResourceID[2:]
	}

	resourceIdBytes, err := hex.DecodeString(ResourceID)
	if err != nil {
		return err
	}

	ResourceIdBytesArr = callsUtil.SliceTo32Bytes(resourceIdBytes)

	if Hash {
		// We must check whether both a deposit and execute function signature is provide or else
		// an invalid hash of 0x00000000 will be taken and set as a function selector in the handler
		DepositSigBytes = callsUtil.GetSolidityFunctionSig([]byte(Deposit))
		ExecuteSigBytes = callsUtil.GetSolidityFunctionSig([]byte(Execute))
		if Deposit == "00000000" {
			DepositSigBytes = [4]byte{}
		}
		if Execute == "00000000" {
			ExecuteSigBytes = [4]byte{}
		}
	} else {
		depositBytes, err := hex.DecodeString(Deposit)
		if err != nil {
			return err
		}
		copy(DepositSigBytes[:], depositBytes[:])

		executeBytes, err := hex.DecodeString(Execute)
		if err != nil {
			return err
		}
		copy(ExecuteSigBytes[:], executeBytes[:])
	}

	log.Debug().Msgf("DepositSigBytes: %x\n", DepositSigBytes[:])
	log.Debug().Msgf("ExecuteSigBytes: %x\n", ExecuteSigBytes[:])

	return nil
}

func RegisterGenericResource(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	log.Info().Msgf("Registering contract %s with resource ID %s on handler %s", TargetContractAddr, ResourceID, HandlerAddr)

	depositerOffsetBigInt := new(big.Int).SetUint64(DepositerOffset)
	log.Info().Msgf("handlerAddr: %s, resourceId: %s, targetcontract: %s, depositsigbytes: %s, depositerOffset: %s, executeSigBytes: %s",
		HandlerAddr,
		ResourceIdBytesArr,
		TargetContractAddr,
		string(DepositSigBytes[:]),
		depositerOffsetBigInt,
		string(DepositSigBytes[:]))
	h, err := contract.AdminSetGenericResource(
		HandlerAddr,
		ResourceIdBytesArr,
		TargetContractAddr,
		DepositSigBytes,
		depositerOffsetBigInt,
		ExecuteSigBytes,
		transactor.TransactOptions{GasLimit: gasLimit},
	)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("Generic resource registered with transaction: %s", h.Hex())
	return nil
}
