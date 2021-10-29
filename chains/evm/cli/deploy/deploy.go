package deploy

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var ErrNoDeploymentFlagsProvided = errors.New("provide at least one deployment flag. For help use --help")
var ErrErc20TokenAndSymbolNotProvided = errors.New("erc20Name and erc20Symbol flags should be provided")

var DeployEVM = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy smart contracts",
	Long:  "This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags",
	RunE:  CallDeployCLI,
	Args: func(cmd *cobra.Command, args []string) error {
		err := ProcessDeployFlags(cmd, args)
		return err
	},
}

var (
	// Flags for all EVM Deploy CLI commands
	Bridge           bool
	Erc20Handler     bool
	Erc20            bool
	Erc721           bool
	DeployAll        bool
	RelayerThreshold uint64
	DomainId         uint8
	Relayers         []string
	Fee              string
	BridgeAddress    string
	Erc20Symbol      string
	Erc20Name        string
)

func BindDeployEVMFlags(deployCmd *cobra.Command) {
	deployCmd.Flags().BoolVar(&Bridge, "bridge", false, "deploy bridge")
	deployCmd.Flags().BoolVar(&Erc20Handler, "erc20Handler", false, "deploy ERC20 handler")
	//deployCmd.Flags().Bool("erc721Handler", false, "deploy ERC721 handler")
	//deployCmd.Flags().Bool("genericHandler", false, "deploy generic handler")
	deployCmd.Flags().BoolVar(&Erc20, "erc20", false, "deploy ERC20")
	deployCmd.Flags().BoolVar(&Erc721, "erc721", false, "deploy ERC721")
	deployCmd.Flags().BoolVar(&DeployAll, "all", false, "deploy all")
	deployCmd.Flags().Uint64Var(&RelayerThreshold, "relayerThreshold", 1, "number of votes required for a proposal to pass")
	deployCmd.Flags().Uint8Var(&DomainId, "domainId", 1, "domain ID for the instance")
	deployCmd.Flags().StringSliceVar(&Relayers, "relayers", []string{}, "list of initial relayers")
	deployCmd.Flags().StringVar(&Fee, "fee", "0", "fee to be taken when making a deposit (in ETH, decimas are allowed)")
	deployCmd.Flags().StringVar(&BridgeAddress, "bridgeAddress", "", "bridge contract address. Should be provided if handlers are deployed separately")
	deployCmd.Flags().StringVar(&Erc20Symbol, "erc20Symbol", "", "ERC20 contract symbol")
	deployCmd.Flags().StringVar(&Erc20Name, "erc20Name", "", "ERC20 contract name")
	flags.MarkFlagsAsRequired(deployCmd, "relayerThreshold", "domainId", "fee", "erc20Symbol", "erc20Name")
}

func init() {
	BindDeployEVMFlags(DeployEVM)
}

var bridgeAddr common.Address
var relayerAddresses []common.Address

func ProcessDeployFlags(cmd *cobra.Command, args []string) error {

	if common.IsHexAddress(BridgeAddress) {
		bridgeAddr = common.HexToAddress(BridgeAddress)
	}
	for _, addr := range Relayers {
		if !common.IsHexAddress(addr) {
			return fmt.Errorf("invalid relayer address %s", addr)
		}
		relayerAddresses = append(relayerAddresses, common.HexToAddress(addr))
	}
	return nil
}

func CallDeployCLI(cmd *cobra.Command, args []string) error {
	txFabric := evmtransaction.NewTransaction
	return DeployCLI(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
}

func DeployCLI(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return err
	}
	log.Debug().Msgf("url: %s gas limit: %v gas price: %v", url, gasLimit, gasPrice)
	log.Debug().Msgf("SENDER Private key 0x%s", hex.EncodeToString(crypto.FromECDSA(senderKeyPair.PrivateKey())))
	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("ethereum client error: %v", err)).Msg("error initializing new EVM client")
		return err
	}
	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})
	log.Debug().Msgf("Relaysers for deploy %+v", Relayers)
	deployments := make([]string, 0)
	log.Debug().Msgf("all bool: %v", DeployAll)

	if DeployAll {
		deployments = append(deployments, []string{"bridge", "erc20Handler", "erc721Handler", "genericHandler", "erc20", "erc721"}...)
	} else {
		if Bridge {
			deployments = append(deployments, "bridge")
		}
		if Erc20Handler {
			deployments = append(deployments, "erc20Handler")
		}
		if Erc20 {
			deployments = append(deployments, "erc20")
		}
	}
	if len(deployments) == 0 {
		log.Error().Err(ErrNoDeploymentFlagsProvided)
		return err
	}
	deployedContracts := make(map[string]string)
	for _, v := range deployments {
		switch v {
		case "bridge":
			log.Debug().Msgf("deploying bridge..")

			bridgeAddr, err = calls.DeployBridge(ethClient, txFabric, gasPricer, DomainId, relayerAddresses, big.NewInt(0).SetUint64(RelayerThreshold))
			if err != nil {
				log.Error().Err(fmt.Errorf("bridge deploy failed: %w", err))
				return err
			}
			deployedContracts["bridge"] = bridgeAddr.String()

			log.Debug().Msgf("bridge address; %v", bridgeAddr.String())
		case "erc20Handler":
			log.Debug().Msgf("deploying ERC20 handler..")
			emptyAddr := common.Address{}
			if bridgeAddr == emptyAddr {
				log.Error().Err(errors.New("bridge flag or bridgeAddress param should be set for contracts deployments"))
				return err
			}

			erc20HandlerAddr, err := calls.DeployErc20Handler(ethClient, txFabric, gasPricer, bridgeAddr)
			if err != nil {
				log.Error().Err(fmt.Errorf("ERC20 handler deploy failed: %w", err))
				return err
			}
			deployedContracts["erc20Handler"] = erc20HandlerAddr.String()
		case "erc20":
			log.Debug().Msgf("deploying ERC20..")
			if Erc20Name == "" || Erc20Symbol == "" {
				log.Error().Err(ErrErc20TokenAndSymbolNotProvided)
				return ErrErc20TokenAndSymbolNotProvided
			}

			erc20Addr, err := calls.DeployErc20(ethClient, txFabric, gasPricer, Erc20Name, Erc20Symbol)
			if err != nil {
				log.Error().Err(fmt.Errorf("erc 20 deploy failed: %w", err))
				return err
			}
			deployedContracts["erc20Token"] = erc20Addr.String()
			if err != nil {
				log.Error().Err(err)
				return err
			}
			if Erc20Name == "" || Erc20Symbol == "" {
				log.Error().Err(ErrErc20TokenAndSymbolNotProvided)
				return ErrErc20TokenAndSymbolNotProvided
			}
		}
	}
	fmt.Printf("%+v", deployedContracts)
	return nil
}
