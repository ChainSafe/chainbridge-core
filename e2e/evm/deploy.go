package evm

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var BobKp = keystore.TestKeyRing.EthereumKeys[keystore.BobKey]
var EveKp = keystore.TestKeyRing.EthereumKeys[keystore.EveKey]

var (
	DefaultRelayerAddresses = []common.Address{
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.AliceKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.BobKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.CharlieKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.DaveKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.EveKey].Address()),
	}
)

func PrepareEVME2EEnv(ethClient calls.ChainClient, fabric calls.TxFabric, domainID uint8, treshHold *big.Int, mintTo common.Address) (common.Address, common.Address, common.Address, error) {
	bridgeAddr, erc20Addr, erc20HandlerAddr, err := deployForTest(ethClient, fabric, domainID, treshHold)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	gasLimit := uint64(2000000)
	// Registering resource
	resourceID := calls.SliceTo32Bytes(append(common.LeftPadBytes(erc20Addr.Bytes(), 31), 0))
	registerResourceInput, err := calls.PrepareAdminSetResourceInput(erc20HandlerAddr, resourceID, erc20Addr)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, &bridgeAddr, registerResourceInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	// Minting tokens
	tenTokens := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))
	minInput, err := calls.PrepareMintTokensInput(mintTo, tenTokens)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, &erc20Addr, minInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}

	// Approving tokens
	approveInput, err := calls.PrepareErc20ApproveInput(erc20HandlerAddr, tenTokens)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, &erc20Addr, approveInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}

	// Adding minter
	minterInput, err := calls.PrepareErc20AddMinterInput(ethClient, erc20Addr, erc20HandlerAddr)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, &erc20Addr, minterInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}

	setBurnInput, err := calls.PrepareSetBurnableInput(erc20HandlerAddr, erc20Addr)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, &bridgeAddr, setBurnInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	log.Debug().Msgf("All deployments and preparations are done")
	return bridgeAddr, erc20Addr, erc20HandlerAddr, nil
}

func deployForTest(c calls.ChainClient, fabric calls.TxFabric, domainID uint8, treshHold *big.Int) (common.Address, common.Address, common.Address, error) {
	erc20Addr, err := calls.DeployErc20(c, fabric, "Test", "TST")
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("ERC20 deploy failed: %w", err)
	}

	bridgeAdrr, err := calls.DeployBridge(c, fabric, domainID, DefaultRelayerAddresses, treshHold)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("Bridge deploy failed: %w", err)
	}

	erc20HandlerAddr, err := calls.DeployErc20Handler(c, fabric, bridgeAdrr)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("ERC20 handler deploy failed: %w", err)
	}

	genericHandlerAddr, err := calls.DeployGenericHandler(c, fabric, bridgeAdrr)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("Generic handler deploy failed: %w", err)
	}

	log.Debug().Msgf("Smart contracts deployed.\n Bridge: %s; \n ERC20: %s;\n ERC20Handler: %s;\n GenericHandler: %s; \n", bridgeAdrr, erc20Addr, erc20HandlerAddr, genericHandlerAddr)
	return bridgeAdrr, erc20Addr, erc20HandlerAddr, nil
}
