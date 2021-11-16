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

func ERC721AddMinter(client ContractCallerDispatcherClient, txFabric TxFabric, gasPricer GasPricer, gasLimit uint64, erc721Contract, minter common.Address) (*common.Hash, error) {
	addMinterInput, err := prepareErc721AddMinterInput(client, erc721Contract, minter)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	txHash, err := Transact(client, txFabric, gasPricer, &erc721Contract, addMinterInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return &txHash, err
}

func ERC721Approve(client ClientDispatcher, txFabric TxFabric, gasPricer GasPricer, gasLimit uint64, tokenId *big.Int, erc721Contract, recipient common.Address) (*common.Hash, error) {
	approveTokenInput, err := prepareERC721ApproveInput(recipient, tokenId)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc721 approve input error: %v", err))
		return nil, err
	}

	txHash, err := Transact(client, txFabric, gasPricer, &erc721Contract, approveTokenInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return &txHash, err
}

func ERC721Deposit(client ClientDispatcher, txFabric TxFabric, gasPricer GasPricer, gasLimit uint64, tokenId *big.Int, metadata string, destinationId int, resourceId types.ResourceID, bridgeContract, recipient common.Address) (*common.Hash, error) {
	data := ConstructErc721DepositData(recipient.Bytes(), tokenId, []byte(metadata))
	txHash, err := Deposit(client, txFabric, gasPricer, bridgeContract, resourceId, uint8(destinationId), data)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return txHash, err
}

func ERC721Mint(client ClientDispatcher, txFabric TxFabric, gasPricer GasPricer, gasLimit uint64, tokenId *big.Int, metadata string, erc721Contract, destination common.Address) (*common.Hash, error) {
	mintTokenInput, err := prepareERC721MintTokensInput(destination, tokenId, metadata)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc721 mint input error: %v", err))
		return nil, err
	}

	txHash, err := Transact(client, txFabric, gasPricer, &erc721Contract, mintTokenInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return &txHash, err
}

func ERC721Owner(client ContractCallerDispatcherClient, tokenId *big.Int, erc721Contract common.Address) (*common.Address, error) {
	ownerOfTokenInput, err := prepareERC721OwnerInput(tokenId)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc721 approve input error: %v", err))
		return nil, err
	}

	msg := ethereum.CallMsg{
		From: common.Address{},
		To:   &erc721Contract,
		Data: ownerOfTokenInput,
	}

	out, err := client.CallContract(context.TODO(), ToCallArg(msg), nil)
	if err != nil {
		log.Error().Err(fmt.Errorf("call contract error: %v", err))
		return nil, err
	}

	res, err := UnpackErc721Result("ownerOf", out)
	if err != nil {
		return nil, err
	}

	ownerAddr := abi.ConvertType(res[0], new(common.Address)).(*common.Address)

	return ownerAddr, nil
}

func PackERC721Method(method string, args ...interface{}) (abi.ABI, []byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.ERC721PresetMinterPauserABI))
	if err != nil {
		return a, []byte{}, err
	}
	input, err := a.Pack(method, args...)
	if err != nil {
		return a, []byte{}, err
	}
	return a, input, nil
}

func UnpackErc721Result(method string, output []byte) ([]interface{}, error) {
	a, err := abi.JSON(strings.NewReader(consts.ERC721PresetMinterPauserABI))
	if err != nil {
		return nil, err
	}

	res, err := a.Unpack(method, output)
	if err != nil {
		log.Error().Err(fmt.Errorf("unpack output error: %v", err))
		return nil, err
	}
	return res, err
}

func prepareERC721MintTokensInput(destAddr common.Address, tokenId *big.Int, metadata string) ([]byte, error) {
	_, res, err := PackERC721Method(
		"mint",
		destAddr,
		tokenId,
		metadata,
	)
	return res, err
}

func prepareERC721ApproveInput(recipientAddr common.Address, tokenId *big.Int) ([]byte, error) {
	_, res, err := PackERC721Method(
		"approve",
		recipientAddr,
		tokenId,
	)
	return res, err
}

func prepareERC721OwnerInput(tokenId *big.Int) ([]byte, error) {
	_, res, err := PackERC721Method(
		"ownerOf",
		tokenId,
	)
	return res, err
}

func prepareErc721AddMinterInput(client ContractCallerClient, erc721Contract, minter common.Address) ([]byte, error) {
	role, err := MinterRole(client, erc721Contract)
	if err != nil {
		return []byte{}, err
	}

	_, res, err := PackERC721Method("grantRole", role, minter)
	return res, err
}
