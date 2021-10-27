// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package evm

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"

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

type E2EClient interface {
	calls.ContractCallerClient
	evmgaspricer.GasPriceClient
	calls.ClientDeployer
}

func PrepareEVME2EEnv(ethClient E2EClient, fabric calls.TxFabric, domainID uint8, threshold *big.Int, mintTo common.Address) (common.Address, common.Address, common.Address, error) {
	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(ethClient, nil)
	bridgeAddr, erc20Addr, erc20HandlerAddr, err := deployForTest(ethClient, fabric, staticGasPricer, domainID, threshold)
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
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &bridgeAddr, registerResourceInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	// Minting tokens
	tenTokens := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))
	minInput, err := calls.PrepareMintTokensInput(mintTo, tenTokens)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &erc20Addr, minInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}

	// Approving tokens
	approveInput, err := calls.PrepareErc20ApproveInput(erc20HandlerAddr, tenTokens)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &erc20Addr, approveInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}

	// Adding minter
	minterInput, err := calls.PrepareErc20AddMinterInput(ethClient, erc20Addr, erc20HandlerAddr)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &erc20Addr, minterInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}

	setBurnInput, err := calls.PrepareSetBurnableInput(erc20HandlerAddr, erc20Addr)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &bridgeAddr, setBurnInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, err
	}
	log.Debug().Msgf("All deployments and preparations are done")
	return bridgeAddr, erc20Addr, erc20HandlerAddr, nil
}

func deployForTest(c E2EClient, fabric calls.TxFabric, gasPriceClient calls.GasPricer, domainID uint8, threshold *big.Int) (common.Address, common.Address, common.Address, error) {
	erc20Addr, err := calls.DeployErc20(c, fabric, gasPriceClient, "Test", "TST")
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("ERC20 deploy failed: %w", err)
	}

	bridgeAdrr, err := calls.DeployBridge(c, fabric, gasPriceClient, domainID, DefaultRelayerAddresses, threshold)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("bridge deploy failed: %w", err)
	}

	erc20HandlerAddr, err := calls.DeployErc20Handler(c, fabric, gasPriceClient, bridgeAdrr)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("bridge deploy failed: %w", err)
	}

	log.Debug().Msgf("Smart contracts deployed.\n Bridge: %s; \n ERC20: %s;\n ERC20Handler: %s;\n", bridgeAdrr, erc20Addr, erc20HandlerAddr)
	return bridgeAdrr, erc20Addr, erc20HandlerAddr, nil
}
