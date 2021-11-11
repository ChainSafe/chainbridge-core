package bridge

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var queryResourceCmd = &cobra.Command{
	Use:   "query-resource",
	Short: "Query the contract address",
	Long:  "Query the contract address with the provided resource ID for a specific handler contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	Run: queryResource,
}

func BindQueryResourceFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Handler, "handler", "", "handler contract address")
	cmd.Flags().StringVar(&ResourceID, "resourceId", "", "resource ID to query")
	flags.MarkFlagsAsRequired(cmd, "handler", "resourceId")
}

func init() {
	BindQueryResourceFlags(queryResourceCmd)
}

func queryResource(cmd *cobra.Command, args []string) {
	log.Debug().Msgf(`
Querying resource
Handler address: %s
Resource ID: %s`, Handler, ResourceID)
}

/*
func queryResource(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	handlerS := cctx.String("handler")
	if !common.IsHexAddress(handlerS) {
		return errors.New("provided handler address is not valid")
	}
	handlerAddr := common.HexToAddress(handlerS)
	resourceIDs := cctx.String("resourceId")
	resourceID := utils.SliceTo32Bytes(common.Hex2Bytes(resourceIDs))
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	res, err := utils.QueryResource(ethClient, handlerAddr, resourceID)
	if err != nil {
		return err
	}
	log.Info().Msgf("Resource address that associated with ID %s is %s", common.Bytes2Hex(resourceID[:]), res.String())
	return nil
}
*/
