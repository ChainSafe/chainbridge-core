package calls

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

// func DeployAndPrepareEnv(ethClient *evmclient.EVMClient, chainID uint8) (common.Address, common.Address, common.Address, error) {
// 	bridgeAddr, erc20Addr, erc20HandlerAddr, err := Deploy(ethClient, chainID)
// 	if err != nil {
// 		return common.Address{}, common.Address{}, common.Address{}, err
// 	}
// 	resourceId := sliceTo32Bytes(append(common.LeftPadBytes(erc20Addr.Bytes(), 31), chainID))
// 	resourceIdBytesArr := utils.SliceTo32Bytes(resourceId)

// 	registerResourceInput, err := PrepareRegisterResourceInput(erc20HandlerAddr, resourceID, erc20Addr, resourceIdBytesArr)
// 	if err != nil {
// 		return common.Address{}, common.Address{}, common.Address{}, err
// 	}

// 	tenTokens := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))
// 	mintTokensInput, err := PrepareMintTokensInput(erc20Addr, tenTokens)
// 	if err != nil {
// 		return common.Address{}, common.Address{}, common.Address{}, err
// 	}
// 	erc20ApproveInput, err := PrepareErc20ApproveInput(erc20Addr, tenTokens)
// 	if err != nil {
// 		return common.Address{}, common.Address{}, common.Address{}, err
// 	}

// 	erc20AddMinterInput, err := PrepareErc20AddMinterInput(ethClient, erc20Addr, erc20HandlerAddr)
// 	if err != nil {
// 		return common.Address{}, common.Address{}, common.Address{}, err
// 	}

// 	setBurnableInput, err := PrepareSetBurnableInput(ethClient, bridgeAddr, erc20HandlerAddr, erc20Addr)
// 	if err != nil {
// 		return common.Address{}, common.Address{}, common.Address{}, err
// 	}
// 	txHash, err := SendInput(c, erc20HandlerAddr, setBurnableInput)

// 	log.Debug().Msgf("All deployments and preparations are done")
// 	return bridgeAddr, erc20Addr, erc20HandlerAddr, nil
// }

// Deploy is a public function to deploy all ChainBridge contracts
// WIP
// TODO: add additional contracts (ERC721, ERC721 handler, generic handler)
func Deploy(c *evmclient.EVMClient, chainID uint8) (common.Address, common.Address, common.Address, error) {
	erc20Addr, err := DeployContract(c, ERC20PresetMinterPauserABI, ERC20PresetMinterPauserBin, "Test", "TST")
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("erc 20 deploy failed: %w", err)
	}

	bridgeAddr, err := DeployContract(c, BridgeABI, BridgeBin, chainID, DefaultRelayerAddresses, big.NewInt(1), big.NewInt(0), big.NewInt(100))
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("bridge deploy failed: %w", err)
	}

	erc20HandlerAddr, err := DeployContract(c, ERC20HandlerABI, ERC20HandlerBin, bridgeAddr, [][32]byte{}, []common.Address{}, []common.Address{})
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("erc 20 Handler deploy failed: %w", err)
	}

	log.Debug().Msgf("Smart contracts deployed.\n Bridge: %s; \n ERC20: %s;\n ERC20 Handler: %s;\n", bridgeAddr, erc20Addr, erc20HandlerAddr)

	return bridgeAddr, erc20Addr, erc20HandlerAddr, nil
}

// COMMANDS

// DeployContract is a public function to manage parsing the contract ABI, constructing a transaction payload and submitting this transaction through the client
func DeployContract(c *evmclient.EVMClient, abiString, binString string, params ...interface{}) (common.Address, error) {
	parsed, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return common.Address{}, err
	}

	input, err := PrepareDeployContractInput(parsed, common.FromHex(binString), params...)
	if err != nil {
		return common.Address{}, err
	}
	address, err := SendInputContract(c, input)
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}
