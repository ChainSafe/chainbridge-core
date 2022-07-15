package listener_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/deposit"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/events"
	"github.com/ChainSafe/sygma-core/chains/evm/listener"
	"github.com/ChainSafe/sygma-core/relayer/message"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/suite"
)

var errIncorrectDataLen = errors.New("invalid calldata length: less than 84 bytes")

type Erc20HandlerTestSuite struct {
	suite.Suite
}

func TestRunErc20HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(Erc20HandlerTestSuite))
}

func (s *Erc20HandlerTestSuite) SetupSuite()    {}
func (s *Erc20HandlerTestSuite) TearDownSuite() {}
func (s *Erc20HandlerTestSuite) SetupTest()     {}
func (s *Erc20HandlerTestSuite) TearDownTest()  {}

func (s *Erc20HandlerTestSuite) TestErc20HandleEvent() {
	// 0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b
	recipientByteSlice := []byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}

	calldata := deposit.ConstructErc20DepositData(recipientByteSlice, big.NewInt(2))
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

	message, err := listener.Erc20DepositHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)

	s.Nil(err)
	s.NotNil(message)
	s.Equal(message, expected)
}

func (s *Erc20HandlerTestSuite) TestErc20HandleEventWithPriority() {
	// 0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b
	recipientByteSlice := []byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}

	calldata := deposit.ConstructErc20DepositDataWithPriority(recipientByteSlice, big.NewInt(2), uint8(1))
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
	// 32-64 is recipient address length
	recipientAddressLength := big.NewInt(0).SetBytes(calldata[32:64])

	// 64 - (64 + recipient address length) is recipient address
	recipientAddressParsed := calldata[64:(64 + recipientAddressLength.Int64())]
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
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	message, err := listener.Erc20DepositHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)

	s.Nil(err)
	s.NotNil(message)
	s.Equal(message, expected)
}

func (s *Erc20HandlerTestSuite) TestErc20HandleEventIncorrectDataLen() {
	metadata := []byte("0xdeadbeef")

	var calldata []byte
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 32)...)
	calldata = append(calldata, metadata...)

	depositLog := &events.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)

	message, err := listener.Erc20DepositHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)

	s.Nil(message)
	s.EqualError(err, errIncorrectDataLen.Error())
}

type Erc721HandlerTestSuite struct {
	suite.Suite
}

func TestRunErc721HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(Erc721HandlerTestSuite))
}

func (s *Erc721HandlerTestSuite) SetupSuite()    {}
func (s *Erc721HandlerTestSuite) TearDownSuite() {}
func (s *Erc721HandlerTestSuite) SetupTest()     {}
func (s *Erc721HandlerTestSuite) TearDownTest()  {}

func (s *Erc721HandlerTestSuite) TestErc721DepositHandlerEmptyMetadata() {
	recipient := common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b")

	calldata := deposit.ConstructErc721DepositData(recipient.Bytes(), big.NewInt(2), []byte{})
	depositLog := &events.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	tokenId := calldata[:32]
	recipientAddressParsed := calldata[64:84]
	var metadata []byte

	expected := &message.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         message.NonFungibleTransfer,
		Payload: []interface{}{
			tokenId,
			recipientAddressParsed,
			metadata,
		},
	}

	m, err := listener.Erc721DepositHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
	s.Nil(err)
	s.NotNil(m)
	s.Equal(expected, m)
}

func (s *Erc721HandlerTestSuite) TestErc721DepositHandlerIncorrectDataLen() {
	metadata := []byte("0xdeadbeef")

	var calldata []byte
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 16)...)
	calldata = append(calldata, metadata...)

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
	s.EqualError(err, "invalid calldata length: less than 84 bytes")
}

func (s *Erc721HandlerTestSuite) TestErc721DepositHandler() {
	recipient := common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b")
	metadata := []byte("metadata.url")

	calldata := deposit.ConstructErc721DepositData(recipient.Bytes(), big.NewInt(2), metadata)
	depositLog := &events.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	tokenId := calldata[:32]
	recipientAddressParsed := calldata[64:84]
	parsedMetadata := calldata[116:128]

	expected := &message.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         message.NonFungibleTransfer,
		Payload: []interface{}{
			tokenId,
			recipientAddressParsed,
			parsedMetadata,
		},
	}

	m, err := listener.Erc721DepositHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
	s.Nil(err)
	s.NotNil(m)
	s.Equal(expected, m)
}
func (s *Erc721HandlerTestSuite) TestErc721DepositHandlerWithPriority() {
	recipient := common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b")
	metadata := []byte("metadata.url")

	calldata := deposit.ConstructErc721DepositDataWithPriority(recipient.Bytes(), big.NewInt(2), metadata, uint8(1))
	depositLog := &events.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	tokenId := calldata[:32]

	// 32 - 64 is recipient address length
	recipientAddressLength := big.NewInt(0).SetBytes(calldata[32:64])

	// 64 - (64 + recipient address length) is recipient address
	recipientAddressParsed := calldata[64:(64 + recipientAddressLength.Int64())]

	// (64 + recipient address length) - ((64 + recipient address length) + 32) is metadata length
	metadataLength := big.NewInt(0).SetBytes(
		calldata[(64 + recipientAddressLength.Int64()):((64 + recipientAddressLength.Int64()) + 32)],
	)

	metadataStart := (64 + recipientAddressLength.Int64()) + 32
	parsedMetadata := calldata[metadataStart : metadataStart+metadataLength.Int64()]

	expected := &message.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         message.NonFungibleTransfer,
		Payload: []interface{}{
			tokenId,
			recipientAddressParsed,
			parsedMetadata,
		},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	m, err := listener.Erc721DepositHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
	s.Nil(err)
	s.NotNil(m)
	s.Equal(expected, m)
}

type GenericHandlerTestSuite struct {
	suite.Suite
}

func TestRunGenericHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GenericHandlerTestSuite))
}

func (s *GenericHandlerTestSuite) SetupSuite()    {}
func (s *GenericHandlerTestSuite) TearDownSuite() {}
func (s *GenericHandlerTestSuite) SetupTest()     {}
func (s *GenericHandlerTestSuite) TearDownTest()  {}

func (s *GenericHandlerTestSuite) TestGenericHandleEventIncorrectDataLen() {
	metadata := []byte("0xdeadbeef")

	var calldata []byte
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 16)...)
	calldata = append(calldata, metadata...)

	depositLog := &events.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	message, err := listener.GenericDepositHandler(
		sourceID,
		depositLog.DestinationDomainID,
		depositLog.DepositNonce,
		depositLog.ResourceID,
		depositLog.Data,
		depositLog.HandlerResponse,
	)

	s.Nil(message)
	s.EqualError(err, "invalid calldata length: less than 32 bytes")
}

func (s *GenericHandlerTestSuite) TestGenericHandleEventEmptyMetadata() {
	metadata := []byte("")
	calldata := deposit.ConstructGenericDepositData(metadata)

	depositLog := &events.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	expected := &message.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         message.GenericTransfer,
		Payload: []interface{}{
			metadata,
		},
	}

	message, err := listener.GenericDepositHandler(
		sourceID,
		depositLog.DestinationDomainID,
		depositLog.DepositNonce,
		depositLog.ResourceID,
		depositLog.Data,
		depositLog.HandlerResponse,
	)

	s.Nil(err)
	s.NotNil(message)
	s.Equal(message, expected)
}

func (s *GenericHandlerTestSuite) TestGenericHandleEvent() {
	metadata := []byte("0xdeadbeef")
	calldata := deposit.ConstructGenericDepositData(metadata)

	depositLog := &events.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	expected := &message.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         message.GenericTransfer,
		Payload: []interface{}{
			metadata,
		},
	}

	message, err := listener.GenericDepositHandler(
		sourceID,
		depositLog.DestinationDomainID,
		depositLog.DepositNonce,
		depositLog.ResourceID,
		depositLog.Data,
		depositLog.HandlerResponse,
	)

	s.Nil(err)
	s.NotNil(message)
	s.Equal(message, expected)
}
