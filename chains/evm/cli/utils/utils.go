package utils

import (
	"github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmgaspricer"

	"github.com/spf13/cobra"
)

var UtilsCmd = &cobra.Command{
	Use:   "utils",
	Short: "Set of utility commands",
	Long:  "Set of utility commands",
}

func init() {
	UtilsCmd.AddCommand(simulateCmd)
	UtilsCmd.AddCommand(hashListCmd)
}

type GasPricerWithPostConfig interface {
	calls.GasPricer
	SetClient(client evmgaspricer.LondonGasClient)
	SetOpts(opts *evmgaspricer.GasPricerOpts)
}
