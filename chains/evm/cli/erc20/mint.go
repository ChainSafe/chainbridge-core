package erc20

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint tokens on an ERC20 mintable contract",
	Long:  "Mint tokens on an ERC20 mintable contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return MintCmd(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateMintFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessMintFlags(cmd, args)
		return err
	},
}

func init() {
	BindMintCmdFlags(mintCmd)
}
func BindMintCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Amount, "amount", "", "amount to deposit")
	cmd.Flags().Uint64Var(&Decimals, "decimals", 0, "ERC20 token decimals")
	cmd.Flags().StringVar(&DstAddress, "dstAddress", "", "Where tokens should be minted. Defaults to TX sender")
	cmd.Flags().StringVar(&Erc20Address, "erc20Address", "", "ERC20 contract address")
	flags.MarkFlagsAsRequired(cmd, "amount", "decimals", "dstAddress", "erc20Address")
}

func ValidateMintFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc20Address) {
		return fmt.Errorf("invalid erc20address %s", Erc20Address)
	}
	return nil
}

var (
	dstAddress    common.Address
	url           string
	gasLimit      uint64
	gasPrice      *big.Int
	senderKeyPair *secp256k1.Keypair
)

func ProcessMintFlags(cmd *cobra.Command, args []string) error {
	var err error
	decimals := big.NewInt(int64(Decimals))
	erc20Addr = common.HexToAddress(Erc20Address)
	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err = flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	if !common.IsHexAddress(DstAddress) {
		dstAddress = senderKeyPair.CommonAddress()
	} else {
		dstAddress = common.HexToAddress(DstAddress)
	}

	realAmount, err = utils.UserAmountToWei(Amount, decimals)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	return nil
}

func MintCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})

	mintTokensInput, err := calls.PrepareMintTokensInput(dstAddress, realAmount)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc20 mint input error: %v", err))
		return err
	}

	_, err = calls.Transact(ethClient, txFabric, gasPricer, &erc20Addr, mintTokensInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(err)
		return err
	}
	log.Info().Msgf("%v tokens minted", Amount)
	return nil
}
