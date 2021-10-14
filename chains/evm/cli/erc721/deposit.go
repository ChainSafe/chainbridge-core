package erc721

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

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
	Short: "Deposit ERC721 token",
	Long:  "Deposit ERC721 token",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return DepositCmd(cmd, args, txFabric)
	},
}

func BindDepositCmdFlags(cli *cobra.Command) {
	cli.Flags().String("recipient", "", "address of recipient")
	cli.Flags().String("bridge", "", "address of bridge contract")
	cli.Flags().String("destId", "", "destination domain ID")
	cli.Flags().String("resourceId", "", "resource ID for transfer")
	cli.Flags().Uint64("tokenId", 0, "ERC721 token id")
}

func init() {
	BindDepositCmdFlags(approveCmd)
}

func DepositCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(
		url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	recipientAddress := cmd.Flag("recipient").Value.String()
	if !common.IsHexAddress(recipientAddress) {
		return fmt.Errorf("invalid recipient address")
	}
	recipientAddr := common.HexToAddress(recipientAddress)

	bridgeAddress := cmd.Flag("bridge").Value.String()
	if !common.IsHexAddress(bridgeAddress) {
		return fmt.Errorf("invalid bridge address")
	}
	bridgeAddr := common.HexToAddress(bridgeAddress)

	destinationId := cmd.Flag("destId").Value.String()
	destinationIdInt, err := strconv.Atoi(destinationId)
	if err != nil {
		log.Error().Err(fmt.Errorf("destination ID conversion error: %v", err))
		return err
	}

	resourceId := cmd.Flag("resourceId").Value.String()
	if resourceId[0:2] == "0x" {
		resourceId = resourceId[2:]
	}
	resourceIdBytes, err := hex.DecodeString(resourceId)
	if err != nil {
		return err
	}
	resourceIdBytesArr := calls.SliceTo32Bytes(resourceIdBytes)

	tokenIdAsString := cmd.Flag("tokenId").Value.String()
	tokenId, ok := big.NewInt(0).SetString(tokenIdAsString, 10)
	if !ok {
		return fmt.Errorf("invalid token id value")
	}

	data := calls.ConstructErc721DepositData(tokenId, recipientAddr.Bytes())

	input, err := calls.PrepareErc20DepositInput(uint8(destinationIdInt), resourceIdBytesArr, data)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc20 deposit input error: %v", err))
		return err
	}

	blockNum, err := ethClient.BlockNumber(context.Background())
	if err != nil {
		log.Error().Err(fmt.Errorf("block fetch error: %v", err))
		return err
	}

	log.Debug().Msgf("blockNum: %v", blockNum)

	txHash, err := calls.Transact(ethClient, txFabric, &bridgeAddr, input, gasLimit)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Debug().Msgf("erc721 deposit hash: %s", txHash.Hex())

	log.Info().Msgf("%s token were transferred to %s from %s", tokenId.String(), recipientAddr.Hex(), senderKeyPair.CommonAddress().String())
	return nil
}
