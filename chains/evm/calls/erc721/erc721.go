package erc721

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type ERC721Contract struct {
	client                calls.ContractCallerDispatcherClient
	erc721ContractAddress common.Address
	abi                   abi.ABI
	transactor.Transactor
}

func NewErc721Contract(
	client calls.ContractCallerDispatcherClient,
	erc721ContractAddress common.Address,
	t transactor.Transactor,
) *ERC721Contract {
	// load ABI
	a, err := abi.JSON(strings.NewReader(consts.ERC721PresetMinterPauserABI))
	if err != nil {
		log.Fatal().Msg("Unable to load ABI") // TODO
	}
	return &ERC721Contract{
		client:                client,
		erc721ContractAddress: erc721ContractAddress,
		abi:                   a,
		Transactor:            t,
	}
}

func (c *ERC721Contract) AddMinter(minter common.Address, opts transactor.TransactOptions) (*common.Hash, error) {
	addMinterInput, err := c.prepareErc721AddMinterInput(minter)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	txHash, err := c.Transact(&c.erc721ContractAddress, addMinterInput, opts)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return txHash, err
}

func (c *ERC721Contract) Approve(tokenId *big.Int, recipient common.Address, opts transactor.TransactOptions) (*common.Hash, error) {
	approveTokenInput, err := c.prepareERC721ApproveInput(recipient, tokenId)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc721 approve input error: %v", err))
		return nil, err
	}

	txHash, err := c.Transact(&c.erc721ContractAddress, approveTokenInput, opts)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return txHash, err
}

// TODO - move to bridge - needs bridge to be refactored
func (c *ERC721Contract) Deposit(txFabric calls.TxFabric, gasPricer calls.GasPricer, tokenId *big.Int, metadata string, destinationId int, resourceId types.ResourceID, bridgeContract, recipient common.Address, opts transactor.TransactOptions) (*common.Hash, error) {
	data := calls.ConstructErc721DepositData(recipient.Bytes(), tokenId, []byte(metadata))
	txHash, err := calls.Deposit(c.client, txFabric, gasPricer, bridgeContract, resourceId, uint8(destinationId), data)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return txHash, err
}

func (c *ERC721Contract) Mint(tokenId *big.Int, metadata string, destination common.Address, opts transactor.TransactOptions) (*common.Hash, error) {
	mintTokenInput, err := c.prepareERC721MintTokensInput(destination, tokenId, metadata)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc721 mint input error: %v", err))
		return nil, err
	}

	txHash, err := c.Transact(&c.erc721ContractAddress, mintTokenInput, opts)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return txHash, err
}

func (c *ERC721Contract) Owner(tokenId *big.Int) (*common.Address, error) {
	ownerOfTokenInput, err := c.prepareERC721OwnerInput(tokenId)
	if err != nil {
		log.Error().Err(fmt.Errorf("erc721 approve input error: %v", err))
		return nil, err
	}

	msg := ethereum.CallMsg{
		From: common.Address{},
		To:   &c.erc721ContractAddress,
		Data: ownerOfTokenInput,
	}

	out, err := c.client.CallContract(context.TODO(), calls.ToCallArg(msg), nil)
	if err != nil {
		log.Error().Err(fmt.Errorf("call contract error: %v", err))
		return nil, err
	}

	res, err := c.UnpackResult("ownerOf", out)
	if err != nil {
		return nil, err
	}

	ownerAddr := abi.ConvertType(res[0], new(common.Address)).(*common.Address)

	return ownerAddr, nil
}

func (c *ERC721Contract) GetABI() abi.ABI {
	return c.abi
}

func (c *ERC721Contract) PackMethod(method string, args ...interface{}) ([]byte, error) {
	input, err := c.abi.Pack(method, args...)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func (c *ERC721Contract) UnpackResult(method string, output []byte) ([]interface{}, error) {
	res, err := c.abi.Unpack(method, output)
	if err != nil {
		log.Error().Err(fmt.Errorf("unpack output error: %v", err))
		return nil, err
	}
	return res, err
}

func (c *ERC721Contract) MinterRole() ([32]byte, error) {
	input, err := c.PackMethod("MINTER_ROLE")
	if err != nil {
		return [32]byte{}, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &c.erc721ContractAddress, Data: input}
	out, err := c.client.CallContract(context.TODO(), calls.ToCallArg(msg), nil)
	if err != nil {
		return [32]byte{}, err
	}
	res, err := c.UnpackResult("MINTER_ROLE", out)
	if err != nil {
		return [32]byte{}, err
	}
	out0 := *abi.ConvertType(res[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

func (c *ERC721Contract) prepareERC721MintTokensInput(destAddr common.Address, tokenId *big.Int, metadata string) ([]byte, error) {
	res, err := c.PackMethod(
		"mint",
		destAddr,
		tokenId,
		metadata,
	)
	return res, err
}

func (c *ERC721Contract) prepareERC721ApproveInput(recipientAddr common.Address, tokenId *big.Int) ([]byte, error) {
	res, err := c.PackMethod(
		"approve",
		recipientAddr,
		tokenId,
	)
	return res, err
}

func (c *ERC721Contract) prepareERC721OwnerInput(tokenId *big.Int) ([]byte, error) {
	res, err := c.PackMethod(
		"ownerOf",
		tokenId,
	)
	return res, err
}

func (c *ERC721Contract) prepareErc721AddMinterInput(minter common.Address) ([]byte, error) {
	role, err := calls.MinterRole(c.client, c.erc721ContractAddress)
	if err != nil {
		return []byte{}, err
	}

	res, err := c.PackMethod("grantRole", role, minter)
	return res, err
}
