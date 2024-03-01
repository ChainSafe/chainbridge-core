package centrifuge

import (
	"encoding/hex"
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Deposit a generic data hash",
	Long:  "The deposit subcommand creates a new generic data deposit on the bridge contract",
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
		return DepositCmd(cmd, args, bridge.NewBridgeContract(c, BridgeAddr, t))
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
	BindDepositFlags(depositCmd)
}

func BindDepositFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Recipient, "recipient", "", "Address of contract to receive generic data")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Address of bridge contract")
	cmd.Flags().StringVar(&Metadata, "metadata", "", "Data (hex bytes) representing params for previously registered functions. Params should be encoded as 32 bytes each")
	cmd.Flags().Uint8Var(&DomainID, "domain", 0, "Destination domain ID")
	cmd.Flags().StringVar(&ResourceID, "resource", "", "Resource ID for transfer")
	cmd.Flags().StringVar(&Priority, "priority", "none", "Transaction priority speed")
	flags.MarkFlagsAsRequired(cmd, "recipient", "bridge", "metadata", "domain", "resource")
}

func ValidateDepositFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Recipient) {
		return fmt.Errorf("invalid recipient address %s", Recipient)
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	switch Priority {
	case "none", "slow", "medium", "fast":
		return nil
	default:
		return fmt.Errorf("invalid priority value %s, supported priorities: \"slow|medium|fast\"", Priority)
	}
}

func ProcessDepositFlags(cmd *cobra.Command, args []string) error {
	var err error

	RecipientAddr = common.HexToAddress(Recipient)
	BridgeAddr = common.HexToAddress(Bridge)
	ResourceIdBytesArr, err = flags.ProcessResourceID(ResourceID)
	if err != nil {
		return err
	}

	MetadataBytes, err = hex.DecodeString(Metadata)
	if err != nil {
		return err
	}

	return err
}

func DepositCmd(cmd *cobra.Command, args []string, bridgeContract *bridge.BridgeContract) error {
	txHash, err := bridgeContract.GenericDeposit(
		MetadataBytes, ResourceIdBytesArr, uint8(DomainID), transactor.TransactOptions{GasLimit: gasLimit, Priority: transactor.TxPriorities[Priority]},
	)
	if err != nil {
		return err
	}

	log.Info().Msgf(
		`Generic deposit hash: %s
		%s metadata was transferred to %s from %s`,
		txHash.Hex(),
		Metadata,
		RecipientAddr.Hex(),
		senderKeyPair.CommonAddress().String(),
	)
	return nil
}
