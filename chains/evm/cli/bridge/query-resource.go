package bridge

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ChainSafe/sygma-core/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var queryResourceCmd = &cobra.Command{
	Use:   "query-resource",
	Short: "Query the resource ID for a handler contract",
	Long:  "The query-resource subcommand queries the contract address with the provided resource ID for a specific handler contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: queryResource,
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateQueryResourceFlags(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func BindQueryResourceFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Handler, "handler", "", "Handler contract address")
	cmd.Flags().StringVar(&ResourceID, "resource", "", "Resource ID to query")
	flags.MarkFlagsAsRequired(cmd, "handler", "resource")
}

func init() {
	BindQueryResourceFlags(queryResourceCmd)
}

func ValidateQueryResourceFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Handler) {
		return fmt.Errorf("invalid handler address: %s", Handler)
	}
	return nil
}

func queryResource(cmd *cobra.Command, args []string) error {
	log.Debug().Msgf(`
Querying resource
Handler address: %s
Resource ID: %s`, Handler, ResourceID)
	return nil
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
