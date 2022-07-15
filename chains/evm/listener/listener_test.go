package listener_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/events"
	"github.com/ChainSafe/sygma-core/chains/evm/listener"
	"github.com/ChainSafe/sygma-core/relayer/message"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/suite"
)

var errIncorrectCalldataLen = errors.New("invalid calldata length: less than 84 bytes")

type ListenerTestSuite struct {
	suite.Suite
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(ListenerTestSuite))
}

func (s *ListenerTestSuite) TestErc20HandleEvent() {
	// 0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b
	recipientByteSlice := []byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}

	// construct ERC20 deposit data
	// follows behavior of solidity tests
	// https://github.com/ChainSafe/chainbridge-solidity/blob/develop/test/contractBridge/depositERC20.js#L46-L50
	var calldata []byte
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(2), 32)...)
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(len(recipientByteSlice))), 32)...)
	calldata = append(calldata, recipientByteSlice...)

	depositLog := &events.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	amountParsed := calldata[:32]
	recipientAddressParsed := calldata[64:]

	expected := &message.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			amountParsed,
			recipientAddressParsed,
		},
	}

	m, err := listener.Erc20DepositHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
	s.Nil(err)
	s.NotNil(m)
	s.Equal(m, expected)
}

func (s *ListenerTestSuite) TestErc20HandleEventIncorrectCalldataLen() {
	// 0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b
	recipientByteSlice := []byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}

	// construct ERC20 deposit data
	// follows behavior of solidity tests
	// https://github.com/ChainSafe/chainbridge-solidity/blob/develop/test/contractBridge/depositERC20.js#L46-L50
	var calldata []byte
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(2), 16)...)
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(len(recipientByteSlice))), 16)...)
	calldata = append(calldata, recipientByteSlice...)

	depositLog := &events.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)

	m, err := listener.Erc20DepositHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
	s.Nil(m)
	s.EqualError(err, errIncorrectCalldataLen.Error())
}

func (s *ListenerTestSuite) TestErc721HandleEvent_WithMetadata_Sucess() {
	// 0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b
	recipientByteSlice := []byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}

	metadataByteSlice := []byte{132, 121, 165, 51, 119, 4, 194, 218, 249, 250, 250, 212, 173, 9, 147, 218, 249, 250, 250, 4, 194, 218, 132, 121, 194, 218, 132, 121, 194, 218, 132, 121}

	var calldata []byte
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(2), 32)...)
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(len(recipientByteSlice))), 32)...)
	calldata = append(calldata, recipientByteSlice...)
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(len(metadataByteSlice))), 32)...)
	calldata = append(calldata, metadataByteSlice...)

	depositLog := &events.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	tokenIdParsed := calldata[:32]
	recipientAddressParsed := calldata[64:84]
	metadataParsed := calldata[116:]

	expected := &message.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         message.NonFungibleTransfer,
		Payload: []interface{}{
			tokenIdParsed,
			recipientAddressParsed,
			metadataParsed,
		},
	}

	m, err := listener.Erc721DepositHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
	s.Nil(err)
	s.NotNil(m)
	s.Equal(expected, m)
}

func (s *ListenerTestSuite) TestErc721HandleEvent_WithoutMetadata_Success() {
	// 0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b
	recipientByteSlice := []byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}

	var calldata []byte
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(2), 32)...)
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(len(recipientByteSlice))), 32)...)
	calldata = append(calldata, recipientByteSlice...)
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(0)), 32)...)

	depositLog := &events.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	tokenIdParsed := calldata[:32]
	recipientAddressParsed := calldata[64:84]
	var metadataParsed []byte

	expected := &message.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         message.NonFungibleTransfer,
		Payload: []interface{}{
			tokenIdParsed,
			recipientAddressParsed,
			metadataParsed,
		},
	}

	m, err := listener.Erc721DepositHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
	s.Nil(err)
	s.NotNil(m)
	s.Equal(expected, m)
}

func (s *ListenerTestSuite) TestErc721HandleEvent_IncorrectCalldataLen_Failure() {
	recipientByteSlice := []byte{241, 229, 143, 177, 119, 4, 194}

	var calldata []byte
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(2), 32)...)
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(len(recipientByteSlice))), 16)...)
	calldata = append(calldata, recipientByteSlice...)

	depositLog := &events.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)

	m, err := listener.Erc721DepositHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
	s.Nil(m)
	s.EqualError(err, errIncorrectCalldataLen.Error())
}
