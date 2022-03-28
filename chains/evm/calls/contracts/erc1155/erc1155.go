package erc1155

import (
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type ERC1155Contract struct {
	contracts.Contract
}

func NewErc1155Contract(
	client calls.ContractCallerDispatcher,
	erc1155ContractAddress common.Address,
	t transactor.Transactor,
) *ERC1155Contract {
	a, _ := abi.JSON(strings.NewReader(consts.ERC1155PresetMinterPauserABI))
	b := common.FromHex(consts.ERC1155PresetMinterPauserBin)
	return &ERC1155Contract{contracts.NewContract(erc1155ContractAddress, a, b, client, t)}
}

func (c *ERC1155Contract) GetBalance(address common.Address, tokenId *big.Int) (*big.Int, error) {
	log.Debug().Msgf("Getting balance for %s %s", address.String(), tokenId.String())
	res, err := c.CallContract("balanceOf", address, tokenId)
	if err != nil {
		return nil, err
	}
	b := abi.ConvertType(res[0], new(big.Int)).(*big.Int)
	return b, nil
}

func (c *ERC1155Contract) AddMinter(
	minter common.Address, opts transactor.TransactOptions,
) (*common.Hash, error) {
	log.Debug().Msgf("Adding new minter %s", minter.String())
	role, err := c.MinterRole()
	if err != nil {
		return nil, err
	}
	return c.ExecuteTransaction("grantRole", opts, role, minter)
}

func (c *ERC1155Contract) Approve(
	recipient common.Address, approved bool, opts transactor.TransactOptions,
) (*common.Hash, error) {
	log.Debug().Msgf("Approving all token for %s", recipient.String())
	return c.ExecuteTransaction("setApprovalForAll", opts, recipient, approved)
}

func (c *ERC1155Contract) Mint(
	destination common.Address, tokenId *big.Int, amount *big.Int, data []byte, opts transactor.TransactOptions,
) (*common.Hash, error) {
	log.Debug().Msgf("Minting tokens %s to %s", tokenId.String(), destination.String())
	return c.ExecuteTransaction("mint", opts, destination, tokenId, amount, data)
}

func (c *ERC1155Contract) Owner(tokenId *big.Int) (*common.Address, error) {
	log.Debug().Msgf("Getting owner of %s", tokenId.String())
	res, err := c.CallContract("ownerOf", tokenId)
	if err != nil {
		return nil, err
	}

	ownerAddr := abi.ConvertType(res[0], new(common.Address)).(*common.Address)
	return ownerAddr, nil
}

func (c *ERC1155Contract) MinterRole() ([32]byte, error) {
	res, err := c.CallContract("MINTER_ROLE")
	if err != nil {
		return [32]byte{}, err
	}
	out := *abi.ConvertType(res[0], new([32]byte)).(*[32]byte)
	return out, nil
}
