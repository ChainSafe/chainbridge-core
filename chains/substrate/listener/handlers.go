package listener

import (
	"fmt"

	"github.com/ChainSafe/chainbridgev2/chains/substrate"
	"github.com/ChainSafe/chainbridgev2/relayer"
)

func FungibleTransferHandler(evtI interface{}) (*relayer.Message, error) {
	evt, ok := evtI.(substrate.EventFungibleTransfer)
	if !ok {
		return nil, fmt.Errorf("failed to cast EventFungibleTransfer type")
	}
	return &relayer.Message{
		Source:       0, // Unset
		Destination:  uint8(evt.Destination),
		DepositNonce: uint64(evt.DepositNonce),
		ResourceId:   evt.ResourceId,
		Payload: []interface{}{
			evt.Amount.Int,
			evt.Recipient,
		},
	}, nil
}

func NonFungibleTransferHandler(evtI interface{}) (*relayer.Message, error) {
	evt, ok := evtI.(substrate.EventNonFungibleTransfer)
	if !ok {
		return nil, fmt.Errorf("failed to cast EventNonFungibleTransfer type")
	}

	return &relayer.Message{
		Source:       0, // Unset
		Destination:  uint8(evt.Destination),
		DepositNonce: uint64(evt.DepositNonce),
		ResourceId:   evt.ResourceId,
		Payload: []interface{}{
			evt.Recipient,
			evt.Metadata,
		},
	}, nil
}

func GenericTransferHandler(evtI interface{}) (*relayer.Message, error) {
	evt, ok := evtI.(substrate.EventGenericTransfer)
	if !ok {
		return nil, fmt.Errorf("failed to cast EventGenericTransfer type")
	}
	return &relayer.Message{
		Source:       0, // Unset
		Destination:  uint8(evt.Destination),
		DepositNonce: uint64(evt.DepositNonce),
		ResourceId:   evt.ResourceId,
		Payload: []interface{}{
			evt.Metadata,
		},
	}, nil
}
