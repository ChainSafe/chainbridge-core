package deposithandlers_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/deposit"
	"github.com/ChainSafe/sygma-core/chains/evm/deposithandlers"
	"github.com/ChainSafe/sygma-core/chains/evm/eventhandlers"
	"github.com/ChainSafe/sygma-core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/suite"
)

var errIncorrectDataLen = errors.New("invalid calldata length: less than 84 bytes")

type Erc20HandlerTestSuite struct {
	suite.Suite
}

func testFunc(Config interface{}) error {
	return nil
}

func testErrFunc(Config interface{}) error {
	return errors.New("Error")
}

type Config struct{}

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
	depositLog := &eventhandlers.Deposit{
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

	expected := &types.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         types.FungibleTransfer,
		Payload: []interface{}{
			amountParsed,
			recipientAddressParsed,
		},
	}
	conf := &Config{}
	erc20DepositHandler := deposithandlers.Erc20DepositHandler{
		ArbitraryFunction: testFunc,
		Config:            conf,
	}
	message, err := erc20DepositHandler.HandleDeposit(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)

	s.Nil(err)
	s.NotNil(message)
	s.Equal(message, expected)
}

func (s *Erc20HandlerTestSuite) TestErc20HandleEventArbitraryFunctionError() {
	// 0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b
	recipientByteSlice := []byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}

	calldata := deposit.ConstructErc20DepositData(recipientByteSlice, big.NewInt(2))
	depositLog := &eventhandlers.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)

	conf := &Config{}
	erc20DepositHandler := deposithandlers.Erc20DepositHandler{
		ArbitraryFunction: testErrFunc,
		Config:            conf,
	}
	_, err := erc20DepositHandler.HandleDeposit(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)

	s.NotNil(err)
}

func (s *Erc20HandlerTestSuite) TestErc20HandleEventWithPriority() {
	// 0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b
	recipientByteSlice := []byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}

	calldata := deposit.ConstructErc20DepositDataWithPriority(recipientByteSlice, big.NewInt(2), uint8(1))
	depositLog := &eventhandlers.Deposit{
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
	expected := &types.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         types.FungibleTransfer,
		Payload: []interface{}{
			amountParsed,
			recipientAddressParsed,
		},
		Metadata: types.Metadata{
			Data: map[string]interface{}{
				"Priority": uint8(1),
			},
		},
	}

	conf := &Config{}
	erc20DepositHandler := deposithandlers.Erc20DepositHandler{
		ArbitraryFunction: testFunc,
		Config:            conf,
	}
	message, err := erc20DepositHandler.HandleDeposit(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)

	s.Nil(err)
	s.NotNil(message)
	s.Equal(message, expected)
}

func (s *Erc20HandlerTestSuite) TestErc20HandleEventIncorrectDataLen() {
	metadata := []byte("0xdeadbeef")

	var calldata []byte
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 32)...)
	calldata = append(calldata, metadata...)

	depositLog := &eventhandlers.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	conf := &Config{}

	erc20DepositHandler := deposithandlers.Erc20DepositHandler{
		ArbitraryFunction: testFunc,
		Config:            conf,
	}
	message, err := erc20DepositHandler.HandleDeposit(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)

	s.Nil(message)
	s.EqualError(err, errIncorrectDataLen.Error())
}
