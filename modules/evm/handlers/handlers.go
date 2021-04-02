package handlers

import (
	"math/big"

	erc20Handler "github.com/ChainSafe/chainbridgev2/bindings/eth/bindings/ERC20Handler"
	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func newErc721Message(source, dest uint8, nonce uint64, resourceId [32]byte, tokenId *big.Int, recipient, metadata []byte) *relayer.XCMessage {
	return &relayer.XCMessage{
		Source:       source,
		Destination:  dest,
		Type:         relayer.NonFungibleTransfer,
		DepositNonce: nonce,
		ResourceId:   resourceId,
		Payload: []interface{}{
			tokenId.Bytes(),
			recipient,
			metadata,
		},
	}
}

func newGenericHandledMessage(source, dest uint8, nonce uint64, resourceId [32]byte, metadata []byte) *relayer.XCMessage {
	return &relayer.XCMessage{
		Source:       source,
		Destination:  dest,
		Type:         relayer.GenericTransfer,
		DepositNonce: nonce,
		ResourceId:   resourceId,
		Payload: []interface{}{
			metadata,
		},
	}
}

func HandleErc20DepositedEvent(sourceID, destId uint8, nonce uint64, handlerContractAddress common.Address, backend bind.ContractBackend) (*relayer.XCMessage, error) {
	contract, err := erc20Handler.NewERC20Handler(handlerContractAddress, backend)
	if err != nil {
		return nil, err
	}
	record, err := contract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
	if err != nil {
		return nil, err
	}

	return &relayer.XCMessage{
		Source:       sourceID,
		Destination:  destId,
		Type:         relayer.FungibleTransfer,
		DepositNonce: nonce,
		ResourceId:   record.ResourceID,
		Payload: []interface{}{
			record.Amount.Bytes(),
			record.DestinationRecipientAddress,
		},
	}, nil
}

//
//func HandleErc721DepositedEvent(destId uint8, nonce uint64) (*relayer.XCMessage, error) {
//	//TODO no call opts. should have From in original chainbridge.
//	record, err := l.erc721HandlerContract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
//	if err != nil {
//		return nil, err
//	}
//	return newErc721Message(
//		l.cfg.ID,
//		destId,
//		nonce,
//		record.ResourceID,
//		nil,
//		nil,
//		record.TokenID,
//		record.DestinationRecipientAddress,
//		record.MetaData,
//	), nil
//}
//
//func HandleGenericDepositedEvent(destId uint8, nonce uint64) (*relayer.XCMessage, error) {
//	record, err := l.genericHandlerContract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
//	if err != nil {
//		log.Error().Err(err).Msg("Error Unpacking Generic Deposit Record")
//		return nil, err
//	}
//	log.Info().Interface("dest", destId).Interface("nonce", nonce).Str("resourceID", common.Bytes2Hex(record.ResourceID[:])).Msg("Handling generic deposit event")
//	return newGenericHandledMessage(
//		l.cfg.ID,
//		destId,
//		nonce,
//		record.ResourceID,
//		nil,
//		nil,
//		record.MetaData[:],
//	), nil
//}
