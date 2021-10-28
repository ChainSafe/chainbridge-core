package erc721

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
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
		return DepositCmd(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateDepositFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessDepositFlags(cmd, args)
		return err
	},
}

func BindDepositCmdFlags() {
	mintCmd.Flags().StringVar(&Recipient, "recipient", "", "address of recipient")
	mintCmd.Flags().StringVar(&Bridge, "bridge", "", "address of bridge contract")
	mintCmd.Flags().StringVar(&DestionationID, "destId", "", "destination domain ID")
	mintCmd.Flags().StringVar(&ResourceID, "resourceId", "", "resource ID for transfer")
	mintCmd.Flags().StringVar(&TokenId, "tokenId", "", "ERC721 token ID")
	flags.MarkFlagsAsRequired(mintCmd, "recipient", "bridge", "destId", "resourceId", "tokenId")
}

func init() {
	BindDepositCmdFlags()
}

func ValidateDepositFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Recipient) {
		return fmt.Errorf("invalid recipient address")
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address")
	}
	return nil
}

func ProcessDepositFlags(cmd *cobra.Command, args []string) error {
	var err error

	recipientAddr = common.HexToAddress(Recipient)
	bridgeAddr = common.HexToAddress(Bridge)

	destinationID, err = strconv.Atoi(DestionationID)
	if err != nil {
		log.Error().Err(fmt.Errorf("destination ID conversion error: %v", err))
		return err
	}

	var ok bool
	tokenId, ok = big.NewInt(0).SetString(TokenId, 10)
	if !ok {
		return fmt.Errorf("invalid token id value")
	}

	resourceId, err = flags.ProcessResourceID(ResourceID)
	return err
}

func DepositCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	ethClient, err := evmclient.NewEVMClientFromParams(
		url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	data := calls.ConstructErc721DepositData(tokenId, recipientAddr.Bytes())

	depositInput, err := calls.PrepareErc20DepositInput(uint8(destinationID), resourceId, data)
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

	txHash, err := calls.Transact(ethClient, txFabric, gasPricer, &bridgeAddr, depositInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Debug().Msgf("erc721 deposit hash: %s", txHash.Hex())

	log.Info().Msgf("%s token were transferred to %s from %s", tokenId.String(), recipientAddr.Hex(), senderKeyPair.CommonAddress().String())
	return nil
}
