package utils

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate transaction invocation",
	Long:  "Replay a failed transaction by simulating invocation; not state-altering",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return SimulateCmd(cmd, args, txFabric)
	},
}

func BindSimulateCmdFlags(cli *cobra.Command) {
	cli.Flags().String("txHash", "", "transaction hash")
	cli.Flags().String("blockNumber", "", "block number")
	cli.Flags().String("fromAddress", "", "address of sender")
}

func init() {
	BindSimulateCmdFlags(simulateCmd)
}

func SimulateCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	txHash := cmd.Flag("txHash").Value.String()
	blockNumber := cmd.Flag("blockNumber").Value.String()
	fromAddress := cmd.Flag("erc20Address").Value.String()

	// convert string block number to big.Int
	// ignore success bool
	blockNumberBigInt, _ := new(big.Int).SetString(blockNumber, 10)

	log.Debug().Msgf(`
Simulating transaction
Tx hash: %s
Block number: %v
From address: %s`,
		txHash, blockNumberBigInt, fromAddress)

	return nil
}
