package bridge

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var RegisterGenericResourceCmd = &cobra.Command{
	Use:   "register-generic-resource",
	Short: "Register a resource ID",
	Long:  "Register a resource ID with a contract address for a generic handler",
	Run:   registerGenericResource,
}

func init() {
	RegisterGenericResourceCmd.Flags().String("handler", "", "handler contract address")
	RegisterGenericResourceCmd.Flags().String("resourceId", "", "resource ID to query")
	RegisterGenericResourceCmd.Flags().String("bridge", "", "bridge contract address")
	RegisterGenericResourceCmd.Flags().String("target", "", "contract address to be registered")
	RegisterGenericResourceCmd.Flags().String("deposit", "", "deposit function signature")
	RegisterGenericResourceCmd.Flags().String("execute", "", "execute proposal function signature")
	RegisterGenericResourceCmd.Flags().Bool("hash", false, "treat signature inputs as function prototype strings, hash and take the first 4 bytes")
}

func registerGenericResource(cmd *cobra.Command, args []string) {
	handlerAddress := cmd.Flag("handler").Value
	resourceId := cmd.Flag("resourceId").Value
	bridgeAddress := cmd.Flag("bridge").Value
	targetAddress := cmd.Flag("target").Value
	deposit := cmd.Flag("deposit").Value
	execute := cmd.Flag("execute").Value
	hash := cmd.Flag("hash").Value
	log.Debug().Msgf(`
Registering generic resource
Handler address: %s
Resource ID: %s
Bridge address: %s
Target address: %s
Deposit: %s
Execute: %s
Hash: %v
`, handlerAddress, resourceId, bridgeAddress, targetAddress, deposit, execute, hash)
}

/*
func registerGenericResource(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Int64("gasLimit")
	gasPrice := cctx.Int64("gasPrice")

	depositSig := cctx.String("deposit")
	depositSigBytes := common.Hex2Bytes(depositSig)
	depositSigBytesArr := utils.SliceTo4Bytes(depositSigBytes)

	executeSig := cctx.String("execute")
	executeSigBytes := common.Hex2Bytes(executeSig)
	executeSigBytesArr := utils.SliceTo4Bytes(executeSigBytes)

	if cctx.Bool("hash") {
		depositSigBytesArr = utils.GetSolidityFunctionSig(depositSig)
		executeSigBytesArr = utils.GetSolidityFunctionSig(executeSig)
	}

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

	log.Info().Msgf("Registering contract %s with resource ID %s on handler %s", targetContract, resourceId, handler)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(gasLimit), big.NewInt(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.RegisterGenericResource(ethClient, bridgeAddress, handlerAddress, resourceIdBytesArr, targetContractAddress, depositSigBytesArr, executeSigBytesArr)
	if err != nil {
		return err
	}
	fmt.Println("Resource registered")
	return nil
}
*/
