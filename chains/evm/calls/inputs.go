package calls

import (
	"context"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

// INPUT FACTORIES

func SendInputContract(client ChainClient, input []byte) (common.Address, error) {
	gasLimit := uint64(2000000)
	gp, err := client.GasPrice()
	if err != nil {
		return common.Address{}, err
	}
	client.LockNonce()
	n, err := client.UnsafeNonce()
	if err != nil {
		return common.Address{}, err
	}
	tx := evmtransaction.NewContractTransaction(n.Uint64(), big.NewInt(0), gasLimit, gp, input)
	hash, err := client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return common.Address{}, err
	}
	log.Debug().Str("hash", hash.String()).Uint64("nonce", n.Uint64()).Msg("contract tx success")
	err = client.UnsafeIncreaseNonce()
	if err != nil {
		return common.Address{}, err
	}
	client.UnlockNonce()
	address := crypto.CreateAddress(client.From(), n.Uint64())

	return address, nil
}

func SendInput(client ChainClient, dest common.Address, input []byte) (common.Hash, error) {
	gasLimit := uint64(2000000)
	gp, err := client.GasPrice()
	if err != nil {
		return common.Hash{}, err
	}
	client.LockNonce()
	n, err := client.UnsafeNonce()
	if err != nil {
		return common.Hash{}, err
	}
	tx := evmtransaction.NewTransaction(n.Uint64(), dest, big.NewInt(0), gasLimit, gp, input)
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

// PREPARED INPUTS

func PrepareDeployContractInput(abi abi.ABI, bytecode []byte, params ...interface{}) ([]byte, error) {
	input, err := abi.Pack("", params...)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareRegisterGenericResourceInput(handler common.Address, rId msg.ResourceId, addr common.Address, depositSig, executeSig [4]byte) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err // Not sure what status to use here
	}
	input, err := a.Pack("adminSetGenericResource", handler, rId, addr, depositSig, executeSig)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareAdminSetResourceInput(handler common.Address, rId [32]byte, addr common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminSetResource", handler, rId, addr)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareMintTokensInput(erc20Addr common.Address, amount *big.Int) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("mint", AliceKp.CommonAddress(), amount)
	if err != nil {
		return []byte{}, err
	}
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
	return input, nil
}

func PrepareSetBurnableInput(client ChainClient, bridge, handler, tokenAddress common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminSetBurnable", handler, tokenAddress)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}
