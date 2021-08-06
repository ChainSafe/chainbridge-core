package deploy

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ErrNoDeploymentFalgsProvided = errors.New("provide at least one deployment flag. For help use --help.")

var DeployEVM = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy smart contracts",
	Long:  "This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags",
	RunE:  CallDeployCLI,
}

func init() {
	flags.BindDeployEVMFlags(DeployEVM)
}

func CallDeployCLI(cmd *cobra.Command, args []string) error {
	txFabric := evmtransaction.NewTransaction
	return DeployCLI(cmd, args, txFabric)
}

func DeployCLI(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return err
	}
	log.Debug().Msgf("url: %s gas limit: %v gas price: %v", url, gasLimit, gasPrice)
	log.Debug().Msgf("SENDER Private key 0x%s", hex.EncodeToString(crypto.FromECDSA(senderKeyPair.PrivateKey())))
	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		return err
	}
	relayerThreshold, err := cmd.Flags().GetUint64("relayerThreshold")
	if err != nil {
		log.Error().Err(fmt.Errorf("relayer threshold error: %v", err)).Msg("error parsing relayersTreshold")
		return err
	}
	relayerAddressesStringSlice := viper.GetStringSlice(flags.RelayersFlagName)
	if err != nil {
		log.Error().Err(fmt.Errorf("relayer threshold error: %v", err))
		return err
	}
	var relayerAddresses []common.Address
	for _, addr := range relayerAddressesStringSlice {
		relayerAddresses = append(relayerAddresses, common.HexToAddress(addr))
	}

	var bridgeAddress common.Address
	bridgeAddressString := cmd.Flag("bridgeAddress").Value.String()
	if common.IsHexAddress(bridgeAddressString) {
		bridgeAddress = common.HexToAddress(bridgeAddressString)
	}

	deployments := make([]string, 0)

	// flag bools
	log.Debug().Msgf("all bool: %v", viper.GetBool("all"))
	allBool, err := cmd.Flags().GetBool("all")
	if err != nil {
		log.Error().Err(fmt.Errorf("all flag error: %v", err))
		return err
	}
	bridgeBool, err := cmd.Flags().GetBool("bridge")
	if err != nil {
		log.Error().Err(fmt.Errorf("bridge flag error: %v", err))
		return err
	}
	erc20HandlerBool, err := cmd.Flags().GetBool("erc20Handler")
	if err != nil {
		log.Error().Err(fmt.Errorf("erc20 handler flag error: %v", err))
		return err
	}
	erc20Bool, err := cmd.Flags().GetBool("erc20")
	if err != nil {
		log.Error().Err(fmt.Errorf("erc20 flag error: %v", err))
		return err
	}

	if allBool {
		deployments = append(deployments, []string{"bridge", "erc20Handler", "erc721Handler", "genericHandler", "erc20", "erc721"}...)
	} else {
		if bridgeBool {
			deployments = append(deployments, "bridge")
		}
		if erc20HandlerBool {
			deployments = append(deployments, "erc20Handler")
		}
		if erc20Bool {
			deployments = append(deployments, "erc20")
		}
	}
	if len(deployments) == 0 {
		log.Error().Err(ErrNoDeploymentFalgsProvided)
		return err
	}
	chainId := cmd.Flag("chainId").Value.String()
	deployedContracts := make(map[string]string)
	for _, v := range deployments {
		switch v {
		case "bridge":
			log.Debug().Msgf("deploying bridge..")
			// convert chain ID to uint
			chainIdInt, err := strconv.Atoi(chainId)
			if err != nil {
				log.Error().Err(fmt.Errorf("chain ID flag error: %v", err))
				return err
			}
			bridgeAddr, err := calls.DeployBridge(ethClient, txFabric, uint8(chainIdInt), relayerAddresses, big.NewInt(0).SetUint64(relayerThreshold))
			if err != nil {
				log.Error().Err(fmt.Errorf("bridge deploy failed: %w", err))
				return err
			}
			deployedContracts["bridge"] = bridgeAddr.String()

			log.Debug().Msgf("bridge address; %v", bridgeAddr.String())
		case "erc20Handler":
			log.Debug().Msgf("deploying ERC20 handler..")
			if bridgeAddress.String() == "" {
				log.Error().Err(errors.New("bridge flag or bridgeAddress param should be set for contracts deployments"))
				return err
			}
			erc20HandlerAddr, err := calls.DeployErc20Handler(ethClient, txFabric, bridgeAddress)
			if err != nil {
				log.Error().Err(fmt.Errorf("ERC20 handler deploy failed: %w", err))
				return err
			}
			deployedContracts["erc20Handler"] = erc20HandlerAddr.String()
		case "erc20":
			log.Debug().Msgf("deploying ERC20..")
			name := cmd.Flag("erc20Name").Value.String()
			symbol := cmd.Flag("erc20Symbol").Value.String()
			if name == "" || symbol == "" {
				log.Error().Err(errors.New("erc20Name and erc20Symbol flags should be provided"))
				return err
			}

			erc20Addr, err := calls.DeployErc20(ethClient, txFabric, name, symbol)
			if err != nil {
				log.Error().Err(fmt.Errorf("erc 20 deploy failed: %w", err))
				return err
			}
			deployedContracts["erc20Token"] = erc20Addr.String()
			if err != nil {
				log.Error().Err(err)
				return err
			}
			if name == "" || symbol == "" {
				log.Error().Err(errors.New("erc20Name and erc20Symbol flags should be provided"))
				return err
			}
		}
	}
	fmt.Printf("%+v", deployedContracts)
	return nil
}
