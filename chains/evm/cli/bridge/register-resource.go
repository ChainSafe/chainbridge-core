package bridge

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var RegisterResourceCmd = &cobra.Command{
	Use:   "register-resource",
	Short: "Register a resource ID",
	Long:  "Register a resource ID with a contract address for a handler",
	Run:   registerResource,
}

func init() {
	RegisterResourceCmd.Flags().String("handler", "", "handler contract address")
	RegisterResourceCmd.Flags().String("bridge", "", "bridge contract address")
	RegisterResourceCmd.Flags().String("target", "", "contract address to be registered")
	RegisterResourceCmd.Flags().String("resourceId", "", "resource ID to be registered")
}

func registerResource(cmd *cobra.Command, args []string) {
	handlerAddress := cmd.Flag("handler").Value
	resourceId := cmd.Flag("resourceId").Value
	targetAddress := cmd.Flag("target").Value
	bridgeAddress := cmd.Flag("bridge").Value
	log.Debug().Msgf(`
Registering resource
Handler address: %s
Resource ID: %s
Target address: %s
Bridge address: %s
`, handlerAddress, resourceId, targetAddress, bridgeAddress)
}

/*
func registerResource(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Int64("gasLimit")
	gasPrice := cctx.Int64("gasPrice")

	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}

	bridgeAddress, err := cliutils.DefineBridgeAddress(cctx)
	if err != nil {
		return err
	}

	handler := cctx.String("handler")
	if !common.IsHexAddress(handler) {
		return fmt.Errorf("invalid handler address %s", handler)
	}
	handlerAddress := common.HexToAddress(handler)
	targetContract := cctx.String("targetContract")
	if !common.IsHexAddress(targetContract) {
		return fmt.Errorf("invalid targetContract address %s", targetContract)
	}
	targetContractAddress := common.HexToAddress(targetContract)
	resourceId := cctx.String("resourceId")
	resourceIdBytes := common.Hex2Bytes(resourceId)
	resourceIdBytesArr := utils.SliceTo32Bytes(resourceIdBytes)

	fmt.Printf("Registering contract %s with resource ID %s on handler %s", targetContract, resourceId, handler)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(gasLimit), big.NewInt(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.RegisterResource(ethClient, bridgeAddress, handlerAddress, resourceIdBytesArr, targetContractAddress)
	if err != nil {
		return err
	}
	fmt.Println("Resource registered")

	return nil
}
*/
