package bridge

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

var registerResourceCmd = &cobra.Command{
	Use:   "register-resource",
	Short: "Register a resource ID",
	Long:  "Register a resource ID with a contract address for a handler",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return RegisterResourceCmd(cmd, args, txFabric)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := validateRegisterResourceFlags(cmd, args)
		if err != nil {
			return err
		}

		err = processRegisterResourceFlags(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func BindRegisterResourceCmdFlags() {
	registerResourceCmd.Flags().StringVarP(&Handler, "handler", "h", "", "handler contract address")
	registerResourceCmd.Flags().StringVarP(&Bridge, "bridge", "b", "", "bridge contract address")
	registerResourceCmd.Flags().StringVarP(&Target, "target", "t", "", "contract address to be registered")
	registerResourceCmd.Flags().StringVarP(&ResourceID, "resourceId", "rID", "", "resource ID to be registered")
}

func init() {
	BindRegisterResourceCmdFlags()
}

func validateRegisterResourceFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Handler) {
		return fmt.Errorf("invalid handler address %s", Handler)
	}
	if !common.IsHexAddress(Target) {
		return fmt.Errorf("invalid target address %s", Target)
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func processRegisterResourceFlags(cmd *cobra.Command, args []string) error {
	var err error
	handlerAddr = common.HexToAddress(Handler)
	targetContractAddr = common.HexToAddress(Target)
	bridgeAddr = common.HexToAddress(Bridge)

	resourceIdBytesArr, err = flags.ProcessResourceID(ResourceID)
	if err != nil {
		return err
	}
	return nil
}

func RegisterResourceCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	log.Debug().Msgf(`
Registering resource
Handler address: %s
Resource ID: %s
Target address: %s
Bridge address: %s
`, Handler, ResourceID, Target, Bridge)

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	fmt.Printf("Registering contract %s with resource ID %s on handler %s", Target, ResourceID, handlerAddr)

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	registerResourceInput, err := calls.PrepareAdminSetResourceInput(handlerAddr, resourceIdBytesArr, targetContractAddr)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	_, err = calls.Transact(ethClient, txFabric, &bridgeAddr, registerResourceInput, gasLimit)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	fmt.Println("Resource registered")
	return nil
}
