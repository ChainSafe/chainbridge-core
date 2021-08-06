package calls

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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
