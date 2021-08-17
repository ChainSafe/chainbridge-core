package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var isRelayerCmd = &cobra.Command{
	Use:   "is-relayer",
	Short: "Check if an address is registered as a relayer",
	Long:  "Check if an address is registered as a relayer",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return IsRelayer(cmd, args, txFabric)
	},
}

func BindIsRelayerFlags(cli *cobra.Command) {
	cli.Flags().String("relayer", "", "address to check")
	cli.Flags().String("bridge", "", "bridge contract address")
}

func init() {
	BindIsRelayerFlags(isRelayerCmd)
}

func IsRelayer(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	relayerAddress := cmd.Flag("relayer").Value.String()
	bridgeAddress := cmd.Flag("bridge").Value.String()
	log.Debug().Msgf(`
	Checking relayer
	Relayer address: %s
	Bridge address: %s`, relayerAddress, bridgeAddress)

	// fetch global flag values
	url, _, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	if !common.IsHexAddress(relayerAddress) {
		err := errors.New("handler address is incorrect format")
		log.Error().Err(err)
		return err
	}

	if !common.IsHexAddress(bridgeAddress) {
		err := errors.New("tokenContract address is incorrect format")
		log.Error().Err(err)
		return err
	}
	relayer := common.HexToAddress(relayerAddress)
	bridge := common.HexToAddress(bridgeAddress)
	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	// erc20Addr, accountAddr
	input, err := calls.PrepareIsRelayerInput(relayer)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare input error: %v", err))
		return err
	}

	msg := ethereum.CallMsg{
		From: common.Address{},
		To:   &bridge,
		Data: input,
	}

	out, err := ethClient.CallContract(context.TODO(), calls.ToCallArg(msg), nil)
	if err != nil {
		log.Error().Err(fmt.Errorf("call contract error: %v", err))
		return err
	}

	if len(out) == 0 {
		// Make sure we have a contract to operate on, and bail out otherwise.
		if code, err := ethClient.CodeAt(context.Background(), bridge, nil); err != nil {
			return err
		} else if len(code) == 0 {
			return fmt.Errorf("no code at provided address %s", bridge.String())
		}
	}
	b, err := calls.ParseIsRelayerOutput(out)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare output error: %v", err))
		return err
	}
	if !b {
		log.Info().Msgf("Address %s is NOT relayer", relayer.String())
	} else {
		log.Info().Msgf("Address %s is relayer", relayer.String())
	}
	return nil
}
