package erc20

import (
	"fmt"
	"math/big"

	callsUtil "github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/util"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Deposit an ERC20 token",
	Long:  "The deposit subcommand creates a new ERC20 token deposit on the bridge contract",
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
	cmd.Flags().StringVar(&Recipient, "recipient", "", "Address of recipient")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Address of bridge contract")
	cmd.Flags().StringVar(&Amount, "amount", "", "Amount to deposit")
	cmd.Flags().Uint8Var(&DomainID, "domain", 0, "Destination domain ID")
	cmd.Flags().StringVar(&ResourceID, "resource", "", "Resource ID for transfer")
	cmd.Flags().Uint64Var(&Decimals, "decimals", 0, "ERC20 token decimals")
	cmd.Flags().StringVar(&Priority, "priority", "none", "Transaction priority speed")
	flags.MarkFlagsAsRequired(cmd, "recipient", "bridge", "amount", "domain", "resource", "decimals")
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

	RecipientAddress = common.HexToAddress(Recipient)
	decimals := big.NewInt(int64(Decimals))
	BridgeAddr = common.HexToAddress(Bridge)
	RealAmount, err = callsUtil.UserAmountToWei(Amount, decimals)
	if err != nil {
		return err
	}
	ResourceIdBytesArr, err = flags.ProcessResourceID(ResourceID)
	return err
}

func DepositCmd(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	hash, err := contract.Erc20Deposit(
		RecipientAddress, RealAmount, ResourceIdBytesArr,
		uint8(DomainID), transactor.TransactOptions{GasLimit: gasLimit, Priority: transactor.TxPriorities[Priority]},
	)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc20 deposit error: %v", err))
		return err
	}

	log.Info().Msgf(
		"%s tokens were transferred to %s from %s with hash %s",
		Amount, RecipientAddress.Hex(), senderKeyPair.CommonAddress().String(), hash.Hex(),
	)
	return nil
}
