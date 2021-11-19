package writer

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

func CreateFungibleProposal(m *message.Message) []interface{} {
	bigAmt := big.NewInt(0).SetBytes(m.Payload[0].([]byte))
	amount := types.NewU128(*bigAmt)
	recipient := types.NewAccountID(m.Payload[1].([]byte))

	t := make([]interface{}, 2)
	t[0] = recipient
	t[1] = amount
	return t
}

func CreateNonFungibleProposal(m *message.Message) []interface{} {
	tokenId := types.NewU256(*big.NewInt(0).SetBytes(m.Payload[0].([]byte)))
	recipient := types.NewAccountID(m.Payload[1].([]byte))
	metadata := types.Bytes(m.Payload[2].([]byte))
	t := make([]interface{}, 3)
	t[0] = recipient
	t[1] = tokenId
	t[2] = metadata
	return t
}

func CreateGenericProposal(m *message.Message) []interface{} {
	t := make([]interface{}, 1)
	t[0] = types.NewHash(m.Payload[0].([]byte))
	return t
}
