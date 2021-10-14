package bridge

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var registerGenericResourceCmd = &cobra.Command{
	Use:   "register-generic-resource",
	Short: "Register a generic resource ID",
	Long:  "Register a resource ID with a contract address for a generic handler",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return RegisterGenericResourceCmd(cmd, args, txFabric)
	},
}

func BindRegisterGenericResourceCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String("handler", "", "handler contract address")
	cmd.Flags().String("resourceId", "", "resource ID to query")
	cmd.Flags().String("bridge", "", "bridge contract address")
	cmd.Flags().String("target", "", "contract address to be registered")
	cmd.Flags().String("deposit", "0x00000000", "deposit function signature")
	cmd.Flags().String("execute", "0x00000000", "execute proposal function signature")
	cmd.Flags().Int("depositerOffset", 0, "depositer address position offset in the metadata, in bytes")
	cmd.Flags().Bool("hash", false, "treat signature inputs as function prototype strings, hash and take the first 4 bytes")

	err := cmd.MarkFlagRequired("handler")
	if err != nil {
		panic(err)
	}
	err = cmd.MarkFlagRequired("resourceId")
	if err != nil {
		panic(err)
	}
	err = cmd.MarkFlagRequired("bridge")
	if err != nil {
		panic(err)
	}
	err = cmd.MarkFlagRequired("target")
	if err != nil {
		panic(err)
	}
}

func init() {
	BindRegisterGenericResourceCmdFlags(registerGenericResourceCmd)
}

func RegisterGenericResourceCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	handlerAddressStr := cmd.Flag("handler").Value.String()
	resourceId := cmd.Flag("resourceId").Value.String()
	bridgeAddressStr := cmd.Flag("bridge").Value.String()
	targetAddressStr := cmd.Flag("target").Value.String()
	depositSig := cmd.Flag("deposit").Value.String()
	executeSig := cmd.Flag("execute").Value.String()
	depositerOffset, err := cmd.Flags().GetInt("depositerOffset")
	if err != nil {
		log.Error().Err(err)
		return fmt.Errorf("could not get depositer offset value: %v", err)
	}
	hash, err := cmd.Flags().GetBool("hash")
	if err != nil {
		log.Error().Err(err)
		return fmt.Errorf("could not get hash value: %v", err)
	}

	log.Debug().Msgf(`
Registering generic resource
Handler address: %s
Resource ID: %s
Bridge address: %s
Target address: %s
Deposit: %s
Execute: %s
Hash: %v
`, handlerAddressStr, resourceId, bridgeAddressStr, targetAddressStr, depositSig, executeSig, hash)

	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	var depositSigBytes [4]byte
	var executeSigBytes [4]byte
	if hash {
		depositSigBytes = calls.GetSolidityFunctionSig([]byte(depositSig))
		executeSigBytes = calls.GetSolidityFunctionSig([]byte(executeSig))
	} else {
		copy(depositSigBytes[:], []byte(depositSig)[:])
		copy(executeSigBytes[:], []byte(executeSig)[:])
	}

	if !common.IsHexAddress(handlerAddressStr) {
		err := fmt.Errorf("invalid handler address %s", handlerAddressStr)
		log.Error().Err(err)
		return err
	}
	handlerAddr := common.HexToAddress(handlerAddressStr)

	if !common.IsHexAddress(targetAddressStr) {
		err := fmt.Errorf("invalid target address %s", targetAddressStr)
		log.Error().Err(err)
		return err
	}
	targetAddr := common.HexToAddress(targetAddressStr)

	bridgeAddress := common.HexToAddress(bridgeAddressStr)
	if resourceId[0:2] == "0x" {
		resourceId = resourceId[2:]
	}
	resourceIdBytes, err := hex.DecodeString(resourceId)
	if err != nil {
		return err
	}
	resourceIdBytesArr := calls.SliceTo32Bytes(resourceIdBytes)

	log.Info().Msgf("Registering contract %s with resource ID %s on handler %s", targetAddressStr, resourceId, handlerAddr)

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	registerGenericResourceInput, err := calls.PrepareAdminSetGenericResourceInput(
		handlerAddr,
		resourceIdBytesArr,
		targetAddr,
		depositSigBytes,
		big.NewInt(int64(depositerOffset)),
		executeSigBytes,
	)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	_, err = calls.Transact(ethClient, txFabric, &bridgeAddress, registerGenericResourceInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msg("Generic resource registered")
	return nil
}
