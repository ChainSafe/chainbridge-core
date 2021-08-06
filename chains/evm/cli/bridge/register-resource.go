package bridge

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/cliutils"
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
	url := cmd.Flag("url").Value.String()
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

	if !common.IsHexAddress(handlerAddressString) {
		log.Fatal().Err(fmt.Errorf("invalid handler address %s", handlerAddressString))
	}
	handlerAddress := common.HexToAddress(handlerAddressString)

	if !common.IsHexAddress(targetAddress) {
		log.Fatal().Err(fmt.Errorf("invalid target address %s", targetAddress))
	}
	targetContractAddress := common.HexToAddress(targetAddress)
	resourceIdBytes := common.Hex2Bytes(resourceId)
	resourceIdBytesArr := calls.SliceTo32Bytes(resourceIdBytes)

	fmt.Printf("Registering contract %s with resource ID %s on handler %s", targetAddress, resourceId, handlerAddress)

	senderKeyPair, err := cliutils.DefineSender(cmd)
	if err != nil {
		log.Fatal().Err(err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Fatal().Err(err)
	}

	registerResourceInput, err := calls.PrepareAdminSetResourceInput(handlerAddress, resourceIdBytesArr, targetContractAddress)
	if err != nil {
		log.Fatal().Err(err)
	}

	_, err = calls.SendInput(ethClient, targetContractAddress, registerResourceInput, txFabric)
	if err != nil {
		log.Fatal().Err(err)
	}

	fmt.Println("Resource registered")
	return nil
}
