package erc20

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/cliutils"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Initiate a transfer of ERC20 tokens",
	Long:  "Initiate a transfer of ERC20 tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return DepositCmd(cmd, args, txFabric)
	},
}

func BindDepositCmdFlags(cli *cobra.Command) {
	cli.Flags().String("recipient", "", "address of recipient")
	cli.Flags().String("bridge", "", "address of bridge contract")
	cli.Flags().String("amount", "", "amount to deposit")
	cli.Flags().String("destId", "", "destination chain ID")
	cli.Flags().String("resourceId", "", "resource ID for transfer")
	cli.Flags().Uint64("decimals", 0, "ERC20 token decimals")
	cli.Flags().Bool("simulate", false, "simulate transaction invocation")
	cli.MarkFlagRequired("decimals")
}

func init() {
	BindDepositCmdFlags(depositCmd)
}

func DepositCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	recipient := cmd.Flag("recipient").Value.String()
	bridgeAddress := cmd.Flag("bridge").Value.String()
	amount := cmd.Flag("amount").Value.String()
	destinationId := cmd.Flag("destId").Value.String()
	resourceId := cmd.Flag("resourceId").Value.String()

	simulateBool, err := cmd.Flags().GetBool("simulate")
	if err != nil {
		return fmt.Errorf("could not get simulate bool flag: %v", err)
	}

	if !common.IsHexAddress(recipient) {
		return fmt.Errorf("invalid recipient address %s", recipient)
	}
	recipientAddress := common.HexToAddress(recipient)
	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	// ignore success bool
	decimals, _ := big.NewInt(0).SetString(cmd.Flag("decimals").Value.String(), 10)

	if !common.IsHexAddress(bridgeAddress) {
		return fmt.Errorf("invalid bridge address %s", bridgeAddress)
	}

	bridgeAddr := common.HexToAddress(bridgeAddress)

	if !common.IsHexAddress(recipient) {
		return fmt.Errorf("invalid recipient address %s", recipientAddress)
	}

	realAmount, err := calls.UserAmountToWei(amount, decimals)
	if err != nil {
		return err
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	if resourceId[0:2] == "0x" {
		resourceId = resourceId[2:]
	}
	resourceIdBytes, err := hex.DecodeString(resourceId)
	if err != nil {
		return err
	}
	resourceIdBytesArr := calls.SliceTo32Bytes(resourceIdBytes)

	destinationIdInt, err := strconv.Atoi(destinationId)
	if err != nil {
		log.Error().Err(fmt.Errorf("destination ID conversion error: %v", err))
		return err
	}
	data := cliutils.ConstructErc20DepositData(recipientAddress.Bytes(), realAmount)
	// TODO: confirm correct arguments
	input, err := calls.PrepareErc20DepositInput(uint8(destinationIdInt), resourceIdBytesArr, data)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc20 deposit input error: %v", err))
		return err
	}

	if simulateBool {
		block, err := ethClient.BlockNumber(context.Background())
		if err != nil {
			log.Error().Err(fmt.Errorf("block fetch error: %v", err))
			return err
		}

		blockNumBigInt := new(big.Int).SetUint64(block)

		simulationData, err := calls.SimulateTransact(ethClient, txFabric, &bridgeAddr, input, gasLimit, blockNumBigInt)
		if err != nil {
			log.Error().Err(fmt.Errorf("simulate transact error: %v", err))
			return err
		}
		log.Debug().Msgf("simulate transaction data: %v", string(simulationData))

		return nil
	}

	// destinationId
	txHash, err := calls.Transact(ethClient, txFabric, &bridgeAddr, input, gasLimit)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc20 deposit error: %v", err))
		return err
	}

	log.Debug().Msgf("erc20 deposit hash: %s", txHash.Hex())

	log.Info().Msgf("%s tokens were transferred to %s from %s", amount, recipientAddress.Hex(), senderKeyPair.CommonAddress().String())
	return nil
}
