package forwarder

import (
	"math/big"
	"strings"

	"github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/consts"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// ForwarderContract matches an instance of https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/metatx/MinimalForwarder.sol
type ForwarderContract struct {
	contracts.Contract
}

type ForwardRequest struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Gas   *big.Int
	Nonce *big.Int
	Data  []byte
}

func NewForwarderContract(
	client calls.ContractCallerDispatcher,
	contractAddress common.Address,
) *ForwarderContract {
	a, _ := abi.JSON(strings.NewReader(consts.MinimalForwarderABI))
	b := common.FromHex(consts.MinimalForwarderBin)
	return &ForwarderContract{
		contracts.NewContract(contractAddress, a, b, client, nil),
	}
}

func (c *ForwarderContract) GetNonce(from common.Address) (*big.Int, error) {
	res, err := c.CallContract("getNonce", from)
	if err != nil {
		return nil, err
	}

	nonce := abi.ConvertType(res[0], new(big.Int)).(*big.Int)
	return nonce, nil
}

func (c *ForwarderContract) PrepareExecute(
	forwardReq ForwardRequest,
	sig []byte,
) ([]byte, error) {
	return c.ABI.Pack("execute", forwardReq, sig)
}
