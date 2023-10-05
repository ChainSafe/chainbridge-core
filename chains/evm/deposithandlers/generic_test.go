package deposithandlers_test

import (
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

	genericDepositHandler := deposithandlers.GenericDepositHandler{
		ArbitraryFunction: testFunc,
		Config:            conf,
	}
	message, err := genericDepositHandler.HandleDeposit(
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

	depositLog := &eventhandlers.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	expected := &types.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         types.GenericTransfer,
		Payload: []interface{}{
			metadata,
		},
	}
	conf := &Config{}

	genericDepositHandler := deposithandlers.GenericDepositHandler{
		ArbitraryFunction: testFunc,
		Config:            conf,
	}
	message, err := genericDepositHandler.HandleDeposit(
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

	depositLog := &eventhandlers.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	expected := &types.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         types.GenericTransfer,
		Payload: []interface{}{
			metadata,
		},
	}

	conf := &Config{}
	genericDepositHandler := deposithandlers.GenericDepositHandler{
		ArbitraryFunction: testFunc,
		Config:            conf,
	}
	message, err := genericDepositHandler.HandleDeposit(
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

func (s *GenericHandlerTestSuite) TestGenericHandleEventArbitraryFunctionError() {
	metadata := []byte("0xdeadbeef")
	calldata := deposit.ConstructGenericDepositData(metadata)

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
	genericDepositHandler := deposithandlers.GenericDepositHandler{
		ArbitraryFunction: testErrFunc,
		Config:            conf,
	}
	_, err := genericDepositHandler.HandleDeposit(
		sourceID,
		depositLog.DestinationDomainID,
		depositLog.DepositNonce,
		depositLog.ResourceID,
		depositLog.Data,
		depositLog.HandlerResponse,
	)

	s.NotNil(err)
}
