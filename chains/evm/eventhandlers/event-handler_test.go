package eventhandlers_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ChainSafe/sygma-core/chains/evm/eventhandlers"
	mock_listener "github.com/ChainSafe/sygma-core/chains/evm/eventhandlers/mock"
	"github.com/ChainSafe/sygma-core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type DepositHandlerTestSuite struct {
	suite.Suite
	depositEventHandler *eventhandlers.DepositEventHandler
	mockDepositHandler  *mock_listener.MockDepositHandler
	mockEventListener   *mock_listener.MockEventListener
	domainID            uint8
}

func TestRunDepositHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(DepositHandlerTestSuite))
}

func (s *DepositHandlerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.domainID = 1
	s.mockEventListener = mock_listener.NewMockEventListener(ctrl)
	s.mockDepositHandler = mock_listener.NewMockDepositHandler(ctrl)
	s.depositEventHandler = eventhandlers.NewDepositEventHandler(s.mockEventListener, s.mockDepositHandler, common.Address{}, s.domainID)
}

func (s *DepositHandlerTestSuite) Test_FetchDepositFails() {
	s.mockEventListener.EXPECT().FetchDeposits(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*eventhandlers.Deposit{}, fmt.Errorf("error"))

	msgChan := make(chan []*types.Message, 1)
	err := s.depositEventHandler.HandleEvent(big.NewInt(0), big.NewInt(5), msgChan)

	s.NotNil(err)
	s.Equal(len(msgChan), 0)
}

func (s *DepositHandlerTestSuite) Test_HandleDepositFails_ExecutionContinue() {
	d1 := &eventhandlers.Deposit{
		DepositNonce:        1,
		DestinationDomainID: 2,
		ResourceID:          types.ResourceID{},
		HandlerResponse:     []byte{},
		Data:                []byte{},
	}
	d2 := &eventhandlers.Deposit{
		DepositNonce:        2,
		DestinationDomainID: 2,
		ResourceID:          types.ResourceID{},
		HandlerResponse:     []byte{},
		Data:                []byte{},
	}
	deposits := []*eventhandlers.Deposit{d1, d2}
	s.mockEventListener.EXPECT().FetchDeposits(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(deposits, nil)
	s.mockDepositHandler.EXPECT().HandleDeposit(
		s.domainID,
		d1.DestinationDomainID,
		d1.DepositNonce,
		d1.ResourceID,
		d1.Data,
		d1.HandlerResponse,
	).Return(&types.Message{}, fmt.Errorf("error"))
	s.mockDepositHandler.EXPECT().HandleDeposit(
		s.domainID,
		d2.DestinationDomainID,
		d2.DepositNonce,
		d2.ResourceID,
		d2.Data,
		d2.HandlerResponse,
	).Return(
		&types.Message{DepositNonce: 2},
		nil,
	)

	msgChan := make(chan []*types.Message, 2)
	err := s.depositEventHandler.HandleEvent(big.NewInt(0), big.NewInt(5), msgChan)
	msgs := <-msgChan

	s.Nil(err)
	s.Equal(msgs, []*types.Message{{DepositNonce: 2}})
}

func (s *DepositHandlerTestSuite) Test_HandleDepositPanis_ExecutionContinues() {
	d1 := &eventhandlers.Deposit{
		DepositNonce:        1,
		DestinationDomainID: 2,
		ResourceID:          types.ResourceID{},
		HandlerResponse:     []byte{},
		Data:                []byte{},
	}
	d2 := &eventhandlers.Deposit{
		DepositNonce:        2,
		DestinationDomainID: 2,
		ResourceID:          types.ResourceID{},
		HandlerResponse:     []byte{},
		Data:                []byte{},
	}
	deposits := []*eventhandlers.Deposit{d1, d2}
	s.mockEventListener.EXPECT().FetchDeposits(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(deposits, nil)
	s.mockDepositHandler.EXPECT().HandleDeposit(
		s.domainID,
		d1.DestinationDomainID,
		d1.DepositNonce,
		d1.ResourceID,
		d1.Data,
		d1.HandlerResponse,
	).Do(func(sourceID, destID, nonce, resourceID, calldata, handlerResponse interface{}) {
		panic("error")
	})
	s.mockDepositHandler.EXPECT().HandleDeposit(
		s.domainID,
		d2.DestinationDomainID,
		d2.DepositNonce,
		d2.ResourceID,
		d2.Data,
		d2.HandlerResponse,
	).Return(
		&types.Message{DepositNonce: 2},
		nil,
	)

	msgChan := make(chan []*types.Message, 2)
	err := s.depositEventHandler.HandleEvent(big.NewInt(0), big.NewInt(5), msgChan)
	msgs := <-msgChan

	s.Nil(err)
	s.Equal(msgs, []*types.Message{{DepositNonce: 2}})
}

func (s *DepositHandlerTestSuite) Test_SuccessfulHandleDeposit() {
	d1 := &eventhandlers.Deposit{
		DepositNonce:        1,
		DestinationDomainID: 2,
		ResourceID:          types.ResourceID{},
		HandlerResponse:     []byte{},
		Data:                []byte{},
	}
	d2 := &eventhandlers.Deposit{
		DepositNonce:        2,
		DestinationDomainID: 2,
		ResourceID:          types.ResourceID{},
		HandlerResponse:     []byte{},
		Data:                []byte{},
	}
	deposits := []*eventhandlers.Deposit{d1, d2}
	s.mockEventListener.EXPECT().FetchDeposits(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(deposits, nil)
	s.mockDepositHandler.EXPECT().HandleDeposit(
		s.domainID,
		d1.DestinationDomainID,
		d1.DepositNonce,
		d1.ResourceID,
		d1.Data,
		d1.HandlerResponse,
	).Return(
		&types.Message{DepositNonce: 1},
		nil,
	)
	s.mockDepositHandler.EXPECT().HandleDeposit(
		s.domainID,
		d2.DestinationDomainID,
		d2.DepositNonce,
		d2.ResourceID,
		d2.Data,
		d2.HandlerResponse,
	).Return(
		&types.Message{DepositNonce: 2},
		nil,
	)

	msgChan := make(chan []*types.Message, 2)
	err := s.depositEventHandler.HandleEvent(big.NewInt(0), big.NewInt(5), msgChan)
	msgs := <-msgChan

	s.Nil(err)
	s.Equal(msgs, []*types.Message{{DepositNonce: 1}, {DepositNonce: 2}})
}
