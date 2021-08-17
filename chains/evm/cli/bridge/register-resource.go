package bridge

import (
	"encoding/hex"
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
}

func BindRegisterResourceCmdFlags(cli *cobra.Command) {
	cli.Flags().String("handler", "", "handler contract address")
	cli.Flags().String("bridge", "", "bridge contract address")
	cli.Flags().String("target", "", "contract address to be registered")
	cli.Flags().String("resourceId", "", "resource ID to be registered")
}

func init() {
	BindRegisterResourceCmdFlags(registerResourceCmd)
}

func RegisterResourceCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	handlerAddressString := cmd.Flag("handler").Value.String()
	resourceId := cmd.Flag("resourceId").Value.String()
	targetAddress := cmd.Flag("target").Value.String()
	bridgeAddressStr := cmd.Flag("bridge").Value.String()
	log.Debug().Msgf(`
Registering resource
Handler address: %s
Resource ID: %s
Target address: %s
Bridge address: %s
`, handlerAddressString, resourceId, targetAddress, bridgeAddressStr)

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
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
	bridgeAddress := common.HexToAddress(bridgeAddressStr)
	if resourceId[0:2] == "0x" {
		resourceId = resourceId[2:]
	}
	resourceIdBytes, err := hex.DecodeString(resourceId)
	if err != nil {
		return err
	}
	resourceIdBytesArr := calls.SliceTo32Bytes(resourceIdBytes)

	fmt.Printf("Registering contract %s with resource ID %s on handler %s", targetAddress, resourceId, handlerAddr)

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

	_, err = calls.Transact(ethClient, txFabric, &bridgeAddress, registerResourceInput, gasLimit)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	fmt.Println("Resource registered")
	return nil
}
