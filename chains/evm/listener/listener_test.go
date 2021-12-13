package listener_test

import (
	"errors"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/relayer/message"

	mock_listener "github.com/ChainSafe/chainbridge-core/chains/evm/listener/mock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

var errIncorrectCalldataLen = errors.New("invalid calldata length: less than 84 bytes")

type ListenerTestSuite struct {
	suite.Suite
	mockEventHandler *mock_listener.MockEventHandler
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(ListenerTestSuite))
}

func (s *ListenerTestSuite) SetupSuite()    {}
func (s *ListenerTestSuite) TearDownSuite() {}
func (s *ListenerTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.mockEventHandler = mock_listener.NewMockEventHandler(gomockController)
}
func (s *ListenerTestSuite) TearDownTest() {}

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

	depositLog := &evmclient.DepositLogs{
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

	m, err := listener.Erc20EventHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
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

	depositLog := &evmclient.DepositLogs{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)

	m, err := listener.Erc20EventHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
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

	depositLog := &evmclient.DepositLogs{
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

	m, err := listener.Erc721EventHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
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

	depositLog := &evmclient.DepositLogs{
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

	m, err := listener.Erc721EventHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
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

	depositLog := &evmclient.DepositLogs{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)

	m, err := listener.Erc721EventHandler(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
	s.Nil(m)
	s.EqualError(err, errIncorrectCalldataLen.Error())
}
