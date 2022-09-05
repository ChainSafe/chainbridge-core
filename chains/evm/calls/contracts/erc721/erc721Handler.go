package erc721

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

type ERC721HandlerContract struct {
	contracts.Contract
}

func NewERC721HandlerContract(
	client calls.ContractCallerDispatcher,
	erc721HandlerContractAddress common.Address,
	t transactor.Transactor,
) *ERC721HandlerContract {
	a, _ := abi.JSON(strings.NewReader(consts.ERC721HandlerABI))
	b := common.FromHex(consts.ERC721HandlerBin)
	return &ERC721HandlerContract{contracts.NewContract(erc721HandlerContractAddress, a, b, client, t)}
}
