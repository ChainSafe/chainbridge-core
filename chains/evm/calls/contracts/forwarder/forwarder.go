package forwarder

import (
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type ForwarderContract struct {
	contracts.Contract
}
type ForwardRequest struct {
	From       common.Address
	To         common.Address
	Value      *big.Int
	Gas        *big.Int
	Nonce      *big.Int
	Data       []byte
	ValidUntil *big.Int
}

func NewForwarderContract(
	client calls.ContractCallerDispatcher,
	contractAddress common.Address,
	transactor transactor.Transactor,
) *ForwarderContract {
	a, _ := abi.JSON(strings.NewReader(consts.GsnForwarderABI))
	b := common.FromHex(consts.GsnForwarderBin)
	return &ForwarderContract{
		contracts.NewContract(contractAddress, a, b, client, transactor),
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

func (c *ForwarderContract) ExecuteData(
	forwardReq ForwardRequest,
	domainSeparator *[32]byte,
	typeHash *[32]byte,
	suffixData []byte,
	sig []byte,
) ([]byte, error) {
	return c.ABI.Pack("execute", forwardReq, domainSeparator, typeHash, suffixData, sig)
}
