package erc20

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/erc20"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/util"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint an ERC20 token",
	Long:  "The mint subcommand mints a token on an ERC20 mintable contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, evmtransaction.NewTransaction, c, prepare)
		if err != nil {
			return err
		}
		return MintCmd(cmd, args, erc20.NewERC20Contract(c, Erc20Addr, t))
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
	BindMintFlags(mintCmd)
}
func BindMintFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Erc20Address, "contract", "", "ERC20 contract address")
	cmd.Flags().StringVar(&Amount, "amount", "", "Deposit amount")
	cmd.Flags().Uint64Var(&Decimals, "decimals", 0, "ERC20 token decimals")
	cmd.Flags().StringVar(&DstAddress, "recipient", "", "Recipient address")
	flags.MarkFlagsAsRequired(cmd, "contract", "amount", "decimals", "recipient")
}

func ValidateMintFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc20Address) {
		return fmt.Errorf("invalid ERC20 address %s", Erc20Address)
	}
	if !common.IsHexAddress(DstAddress) {
		return fmt.Errorf("invalid recipient address %s", DstAddress)
	}
	return nil
}

func ProcessMintFlags(cmd *cobra.Command, args []string) error {
	var err error
	decimals := big.NewInt(int64(Decimals))
	Erc20Addr = common.HexToAddress(Erc20Address)

	if !common.IsHexAddress(DstAddress) {
		dstAddress = senderKeyPair.CommonAddress()
	} else {
		dstAddress = common.HexToAddress(DstAddress)
	}

	RealAmount, err = calls.UserAmountToWei(Amount, decimals)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	return nil
}

func MintCmd(cmd *cobra.Command, args []string, contract *erc20.ERC20Contract) error {
	_, err := contract.MintTokens(
		dstAddress, RealAmount, transactor.TransactOptions{GasLimit: gasLimit},
	)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	log.Info().Msgf("%v tokens minted", Amount)
	return nil
}
