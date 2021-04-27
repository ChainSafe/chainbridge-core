package writer

import (
	"math/big"

	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

func CreateFungibleProposal(m *relayer.Message) (*SubstrateProposal, error) {
	bigAmt := big.NewInt(0).SetBytes(m.Payload[0].([]byte))
	amount := types.NewU128(*bigAmt)
	recipient := types.NewAccountID(m.Payload[1].([]byte))
	depositNonce := types.U64(m.DepositNonce)

	meta := w.conn.getMetadata()
	method, err := w.resolveResourceId(m.ResourceId)
	if err != nil {
		return nil, err
	}
	call, err := types.NewCall(
		&meta,
		method,
		recipient,
		amount,
	)
	if err != nil {
		return nil, err
	}
	if w.extendCall {
		eRID, err := types.EncodeToBytes(m.ResourceId)
		if err != nil {
			return nil, err
		}
		call.Args = append(call.Args, eRID...)
	}

	return &SubstrateProposal{
		DepositNonce: depositNonce,
		Call:         call,
		SourceId:     types.U8(m.Source),
		ResourceId:   types.NewBytes32(m.ResourceId),
		Method:       method,
	}, nil
}

func CreateNonFungibleProposal(m *relayer.Message) (*SubstrateProposal, error) {
	tokenId := types.NewU256(*big.NewInt(0).SetBytes(m.Payload[0].([]byte)))
	recipient := types.NewAccountID(m.Payload[1].([]byte))
	metadata := types.Bytes(m.Payload[2].([]byte))
	depositNonce := types.U64(m.DepositNonce)

	meta := w.conn.getMetadata()
	method, err := w.resolveResourceId(m.ResourceId)
	if err != nil {
		return nil, err
	}

	call, err := types.NewCall(
		&meta,
		method,
		recipient,
		tokenId,
		metadata,
	)
	if err != nil {
		return nil, err
	}
	if w.extendCall {
		eRID, err := types.EncodeToBytes(m.ResourceId)
		if err != nil {
			return nil, err
		}
		call.Args = append(call.Args, eRID...)
	}

	return &SubstrateProposal{
		DepositNonce: depositNonce,
		Call:         call,
		SourceId:     types.U8(m.Source),
		ResourceId:   types.NewBytes32(m.ResourceId),
		Method:       method,
	}, nil
}

func CreateGenericProposal(m *relayer.Message) (*SubstrateProposal, error) {
	meta := w.conn.getMetadata()
	method, err := w.resolveResourceId(m.ResourceId)
	if err != nil {
		return nil, err
	}

	call, err := types.NewCall(
		&meta,
		method,
		types.NewHash(m.Payload[0].([]byte)),
	)
	if err != nil {
		return nil, err
	}
	if w.extendCall {
		eRID, err := types.EncodeToBytes(m.ResourceId)
		if err != nil {
			return nil, err
		}

		call.Args = append(call.Args, eRID...)
	}
	return &SubstrateProposal{
		DepositNonce: types.U64(m.DepositNonce),
		Call:         call,
		SourceId:     types.U8(m.Source),
		ResourceId:   types.NewBytes32(m.ResourceId),
		Method:       method,
	}, nil
}
