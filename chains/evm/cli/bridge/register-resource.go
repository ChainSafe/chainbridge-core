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
	RunE:  CallRegisterResource,
}

func init() {
	registerResourceCmd.Flags().String("handler", "", "handler contract address")
	registerResourceCmd.Flags().String("bridge", "", "bridge contract address")
	registerResourceCmd.Flags().String("target", "", "contract address to be registered")
	registerResourceCmd.Flags().String("resourceId", "", "resource ID to be registered")
}

func CallRegisterResource(cmd *cobra.Command, args []string) error {
	txFabric := evmtransaction.NewTransaction
	return registerResource(cmd, args, txFabric)
}

func registerResource(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	handlerAddressString := cmd.Flag("handler").Value.String()
	resourceId := cmd.Flag("resourceId").Value.String()
	targetAddress := cmd.Flag("target").Value.String()
	bridgeAddress := cmd.Flag("bridge").Value.String()
	log.Debug().Msgf(`
Registering resource
Handler address: %s
Resource ID: %s
Target address: %s
Bridge address: %s
`, handlerAddressString, resourceId, targetAddress, bridgeAddress)
	// bridgeAddress := cmd.Flag("bridge").Value.String()

	// fetch global flag values
	url, _, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	if !common.IsHexAddress(handlerAddressString) {
		err := fmt.Errorf("invalid handler address %s", handlerAddressString)
		log.Error().Err(err)
		return err
	}
	handlerAddr := common.HexToAddress(handlerAddressString)

	if !common.IsHexAddress(targetAddress) {
		err := fmt.Errorf("invalid target address %s", targetAddress)
		log.Error().Err(err)
		return err
	}
	targetContractAddr := common.HexToAddress(targetAddress)
	resourceIdBytes := common.Hex2Bytes(resourceId)
	resourceIdBytesArr := calls.SliceTo32Bytes(resourceIdBytes)

	fmt.Printf("Registering contract %s with resource ID %s on handler %s", targetAddress, resourceId, handlerAddr)

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(),gasPrice)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	registerResourceInput, err := calls.PrepareAdminSetResourceInput(handlerAddr, resourceIdBytesArr, targetContractAddr)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	_, err = calls.SendInput(ethClient, targetContractAddr, registerResourceInput, txFabric)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	fmt.Println("Resource registered")
	return nil
}
