package admin

import (
	"context"
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
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
		return IsRelayer(cmd, args)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := validateIsRelayerFlags(cmd, args)
		if err != nil {
			return err
		}

		processIsRelayerFlags(cmd, args)
		return nil
	},
}

func BindIsRelayerFlags() {
	isRelayerCmd.Flags().StringVarP(&Relayer, "relayer", "r", "", "address to check")
	isRelayerCmd.Flags().StringVarP(&Bridge, "bridge", "b", "", "bridge contract address")
	flags.MarkFlagsAsRequired(isRelayerCmd, "relayer", "bridge")
}

func init() {
	BindIsRelayerFlags()
}

func validateIsRelayerFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Relayer) {
		return fmt.Errorf("invalid relayer address %s", Relayer)
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func processIsRelayerFlags(cmd *cobra.Command, args []string) {
	relayerAddr = common.HexToAddress(Relayer)
	bridgeAddr = common.HexToAddress(Bridge)
}

func IsRelayer(cmd *cobra.Command, args []string) error {
	log.Debug().Msgf(`
	Checking relayer
	Relayer address: %s
	Bridge address: %s`, Relayer, Bridge)

	// fetch global flag values
	url, _, _, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(err)
		return err
	}
	// erc20Addr, accountAddr
	input, err := calls.PrepareIsRelayerInput(relayerAddr)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare input error: %v", err))
		return err
	}

	msg := ethereum.CallMsg{
		From: common.Address{},
		To:   &bridgeAddr,
		Data: input,
	}

	out, err := ethClient.CallContract(context.TODO(), calls.ToCallArg(msg), nil)
	if err != nil {
		log.Error().Err(fmt.Errorf("call contract error: %v", err))
		return err
	}

	if len(out) == 0 {
		// Make sure we have a contract to operate on, and bail out otherwise.
		if code, err := ethClient.CodeAt(context.Background(), bridgeAddr, nil); err != nil {
			return err
		} else if len(code) == 0 {
			return fmt.Errorf("no code at provided address %s", bridgeAddr.String())
		}
	}
	b, err := calls.ParseIsRelayerOutput(out)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare output error: %v", err))
		return err
	}
	if !b {
		log.Info().Msgf("Address %s is NOT relayer", relayerAddr.String())
	} else {
		log.Info().Msgf("Address %s is relayer", relayerAddr.String())
	}
	return nil
}
