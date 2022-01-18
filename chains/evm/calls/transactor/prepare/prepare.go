package prepare

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/common"
)

type Transactor interface {
	Transact(to *common.Address, data []byte, opts transactor.TransactOptions) (*common.Hash, error)
}
type prepareTransactor struct{}

// Initializes PrepareTransactor which is used when --prepare flag value is set as true from CLI
// PrepareTransactor outputs calldata to stdout for multisig calls (it doesn't make any contract calls)
func NewPrepareTransactor() Transactor {
	return &prepareTransactor{}
}

// Outputs calldata to stdout (called when --prepare flag value is set as true from CLI)
func (t *prepareTransactor) Transact(to *common.Address, data []byte, opts transactor.TransactOptions) (*common.Hash, error) {
	log.Info().Msgf(`
To:
%s
Calldata:
%+v`, to, common.Bytes2Hex(data))
	return &common.Hash{}, nil
}
