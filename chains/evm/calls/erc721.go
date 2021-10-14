package calls

import (
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

func PrepareERC721MintTokensInput(destAddr common.Address, tokenId *big.Int, metadata []byte) ([]byte, error) {
	log.Debug().Msgf("")
	_, res, err := packERC721Method(
		"mint",
		destAddr,
		tokenId,
		metadata,
	)
	return res, err
}

func PrepareERC721ApproveInput(recipientAddr common.Address, tokenId *big.Int) ([]byte, error) {
	_, res, err := packERC721Method(
		"approve",
		recipientAddr,
		tokenId,
	)
	return res, err
}

func PrepareERC721OwnerInput(tokenId *big.Int) ([]byte, error) {
	_, res, err := packERC721Method(
		"ownerOf",
		tokenId,
	)
	return res, err
}

func PrepareErc721AddMinterInput(client ChainClient, erc721Contract, handler common.Address) ([]byte, error) {
	role, err := MinterRole(client, erc721Contract)
	if err != nil {
		return []byte{}, err
	}

	_, res, err := packERC721Method("grantRole", role, handler)
	return res, err
}

func packERC721Method(method string, args ...interface{}) (abi.ABI, []byte, error) {
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
