package calls

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

func PrepareMintTokensInput(destAddr common.Address, amount *big.Int) ([]byte, error) {
	log.Debug().Msgf("Minting tokens %s %s", destAddr.String(), amount.String())
	a, err := abi.JSON(strings.NewReader(consts.ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("mint", destAddr, amount)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareErc20ApproveInput(target common.Address, amount *big.Int) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("approve", target, amount)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareErc20AddMinterInput(client ContractCallerClient, erc20Contract, handler common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	role, err := MinterRole(client, erc20Contract)
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("grantRole", role, handler)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareRegisterGenericResourceInput(handler common.Address, resourceID types.ResourceID, addr common.Address, depositSig, executeSig [4]byte) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err // Not sure what status to use here
	}
	input, err := a.Pack("adminSetGenericResource", handler, resourceID, addr, depositSig, executeSig)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareERC20BalanceInput(accountAddr common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("balanceOf", accountAddr)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func ParseERC20BalanceOutput(output []byte) (*big.Int, error) {
	a, err := abi.JSON(strings.NewReader(consts.ERC20PresetMinterPauserABI))
	if err != nil {
		return new(big.Int), err
	}

	res, err := a.Unpack("balanceOf", output)
	if err != nil {
		log.Error().Err(fmt.Errorf("unpack output error: %v", err))
		return new(big.Int), err
	}

	balance := abi.ConvertType(res[0], new(big.Int)).(*big.Int)

	return balance, nil
}

func MinterRole(chainClient ContractCallerClient, erc20Contract common.Address) ([32]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.ERC20PresetMinterPauserABI))
	if err != nil {
		return [32]byte{}, err
	}
	input, err := a.Pack("MINTER_ROLE")
	if err != nil {
		return [32]byte{}, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &erc20Contract, Data: input}
	out, err := chainClient.CallContract(context.TODO(), ToCallArg(msg), nil)
	if err != nil {
		return [32]byte{}, err
	}
	res, err := a.Unpack("MINTER_ROLE", out)
	if err != nil {
		return [32]byte{}, err
	}
	out0 := *abi.ConvertType(res[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

func GetERC20Balance(ethClient ContractCheckerCallerClient, erc20Addr, address common.Address) (*big.Int, error) {
	input, err := PrepareERC20BalanceInput(address)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare input error: %v", err))
		return nil, err
	}

	msg := ethereum.CallMsg{
		From: common.Address{},
		To:   &erc20Addr,
		Data: input,
	}

	out, err := ethClient.CallContract(context.TODO(), ToCallArg(msg), nil)
	if err != nil {
		log.Error().Err(fmt.Errorf("call contract error: %v", err))
		return nil, err
	}

	if len(out) == 0 {
		// Make sure we have a contract to operate on, and bail out otherwise.
		if code, err := ethClient.CodeAt(context.Background(), erc20Addr, nil); err != nil {
			return nil, err
		} else if len(code) == 0 {
			return nil, fmt.Errorf("no code at provided address %s", erc20Addr.String())
		}
	}

	balance, err := ParseERC20BalanceOutput(out)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare output error: %v", err))
		return nil, err
	}
	return balance, nil
}
