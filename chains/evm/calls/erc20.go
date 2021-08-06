package calls

import (
	"context"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

// @dev
// inputs here and in bridge.go could get consolidated into something similar to txFabric in deploy.go

func PrepareMintTokensInput(destAddr common.Address, amount *big.Int) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("mint", destAddr, amount)
	if err != nil {
		return []byte{}, err
	}
	input = append(input, common.FromHex(ERC20PresetMinterPauserBin)...)
	return input, nil
}

func PrepareErc20ApproveInput(target common.Address, amount *big.Int) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("approve", target, amount)
	if err != nil {
		return []byte{}, err
	}
	input = append(input, common.FromHex(ERC20PresetMinterPauserBin)...)
	return input, nil
}

func PrepareErc20AddMinterInput(client ChainClient, erc20Contract, handler common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	role, err := mintRole(client, erc20Contract)
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("grantRole", role, handler)
	if err != nil {
		return []byte{}, err
	}
	input = append(input, common.FromHex(ERC20PresetMinterPauserBin)...)
	return input, nil
}

func PrepareErc20DepositInput(bridgeAddress, recipientAddress common.Address, amount *big.Int, rId [32]byte, destChainId uint8) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("deposit", bridgeAddress, recipientAddress, amount, rId, destChainId)
	if err != nil {
		return []byte{}, err
	}
	input = append(input, common.FromHex(ERC20PresetMinterPauserBin)...)
	return input, nil
}

func PrepareRegisterGenericResourceInput(handler common.Address, rId [32]byte, addr common.Address, depositSig, executeSig [4]byte) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err // Not sure what status to use here
	}
	input, err := a.Pack("adminSetGenericResource", handler, rId, addr, depositSig, executeSig)
	if err != nil {
		return []byte{}, err
	}
	input = append(input, common.FromHex(ERC20PresetMinterPauserBin)...)
	return input, nil
}

func PrepareERC20BalanceInput(erc20Addr, accountAddr common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("balanceOf", erc20Addr, accountAddr)
	if err != nil {
		return []byte{}, err
	}
	input = append(input, common.FromHex(ERC20PresetMinterPauserBin)...)

	return input, nil
}

// @dev
// refactor to be reusable
func SendInput(client ChainClient, dest common.Address, input []byte, txFabric TxFabric) (common.Hash, error) {
	gp, err := client.GasPrice()
	if err != nil {
		return common.Hash{}, err
	}
	client.LockNonce()
	n, err := client.UnsafeNonce()
	if err != nil {
		return common.Hash{}, err
	}
	tx := txFabric(n.Uint64(), nil, big.NewInt(0), config.DefaultGasLimit, gp, input)
	hash, err := client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return common.Hash{}, err
	}
	log.Debug().Str("hash", hash.String()).Uint64("nonce", n.Uint64()).Msg("tx success")
	err = client.UnsafeIncreaseNonce()
	if err != nil {
		return common.Hash{}, err
	}
	client.UnlockNonce()
	return tx.Hash(), nil
}
