package erc721

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/client"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contract"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type ERC721Contract struct {
	contract.Contract
}

func NewErc721Contract(
	client client.ContractCallerDispatcherClient,
	erc721ContractAddress common.Address,
	t transactor.Transactor,
) *ERC721Contract {
	// load ABI
	a, err := abi.JSON(strings.NewReader(consts.PresetMinterPauserABI))
	if err != nil {
		log.Fatal().Msg("Unable to load ABI") // TODO
	}
	b := common.FromHex(consts.PresetMinterPauserBin)
	return &ERC721Contract{contract.NewContract(erc721ContractAddress, a, b, client, t)}
}

// Add new minter for ERC721 contract
func (c *ERC721Contract) AddMinter(minter common.Address, opts transactor.TransactOptions) (*common.Hash, error) {
	role, err := c.MinterRole()
	if err != nil {
		return nil, err
	}
	return c.ExecuteTransaction("grantRole", opts, role, minter)
}

func (c *ERC721Contract) Approve(tokenId *big.Int, recipient common.Address, opts transactor.TransactOptions) (*common.Hash, error) {
	return c.ExecuteTransaction("approve", opts, recipient, tokenId)
}

func (c *ERC721Contract) Mint(tokenId *big.Int, metadata string, destination common.Address, opts transactor.TransactOptions) (*common.Hash, error) {
	return c.ExecuteTransaction("mint", opts, destination, tokenId, metadata)
}

func (c *ERC721Contract) Owner(tokenId *big.Int) (*common.Address, error) {
	res, err := c.CallContract("ownerOf", tokenId)
	if err != nil {
		return nil, err
	}

	ownerAddr := abi.ConvertType(res[0], new(common.Address)).(*common.Address)
	return ownerAddr, nil
}

func (c *ERC721Contract) MinterRole() ([32]byte, error) {
	res, err := c.CallContract("MINTER_ROLE")
	if err != nil {
		return [32]byte{}, err
	}
	out := *abi.ConvertType(res[0], new([32]byte)).(*[32]byte)
	return out, nil
}
