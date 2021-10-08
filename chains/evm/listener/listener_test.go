package listener_test

import (
	"errors"
	"math/big"
	"testing"
	"time"

	mock_blockstore "github.com/ChainSafe/chainbridge-core/blockstore/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	mock_listener "github.com/ChainSafe/chainbridge-core/chains/evm/listener/mock"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type EVMListenerTestSuite struct {
	suite.Suite
	chainReaderMock    *mock_listener.MockChainClient
	eventHandlerMock   *mock_listener.MockEventHandler
	keyValueWriterMock *mock_blockstore.MockKeyValueWriter
	gomockController   *gomock.Controller
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(EVMListenerTestSuite))
}

func (s *EVMListenerTestSuite) SetupSuite()    {}
func (s *EVMListenerTestSuite) TearDownSuite() {}
func (s *EVMListenerTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.chainReaderMock = mock_listener.NewMockChainClient(s.gomockController)
	s.eventHandlerMock = mock_listener.NewMockEventHandler(s.gomockController)
	s.keyValueWriterMock = mock_blockstore.NewMockKeyValueWriter(s.gomockController)
}
func (s *EVMListenerTestSuite) TearDownTest() {}

var (
	testBridgeAddress     = common.HexToAddress("")
	testDomainID          = uint8(0)
	testLogsForStartBlock = &listener.DepositLogs{
		DestinationID: 0,
		ResourceID:    [32]byte{},
		DepositNonce:  5,
	}
	testMessage = &relayer.Message{
		Source:       0,
		Destination:  1,
		DepositNonce: 5,
		ResourceId:   [32]byte{},
		Payload:      []interface{}{},
		Type:         relayer.NonFungibleTransfer,
	}
)

func (s *EVMListenerTestSuite) TestValidMessageIsReturned() {
	latestBlock := big.NewInt(100)
	startBlock := big.NewInt(0).Sub(latestBlock, listener.BlockDelay)

	// Define mocks

	s.chainReaderMock.EXPECT().LatestBlock().Return(latestBlock, nil).AnyTimes()

	s.chainReaderMock.EXPECT().FetchDepositLogs(
		gomock.Any(),
		gomock.Eq(testBridgeAddress),
		gomock.Eq(startBlock),
		gomock.Eq(startBlock),
	).Return([]*listener.DepositLogs{testLogsForStartBlock}, nil).AnyTimes()

	s.eventHandlerMock.EXPECT().HandleEvent(
		gomock.Eq(testDomainID),
		gomock.Eq(testLogsForStartBlock.DestinationID),
		gomock.Eq(testLogsForStartBlock.DepositNonce),
		gomock.Eq(testLogsForStartBlock.ResourceID),
	).Return(testMessage, nil).AnyTimes()

	s.keyValueWriterMock.EXPECT().SetByKey(
		gomock.Any(),
		gomock.Any(),
	).Return(nil).AnyTimes()

	// Start EVMListener

	l := listener.NewEVMListener(
		s.chainReaderMock, s.eventHandlerMock, testBridgeAddress,
	)

	stopCh := make(chan struct{})
	errCh := make(chan error)

	messageCh := l.ListenToEvents(
		startBlock,
		testDomainID,
		s.keyValueWriterMock,
		stopCh,
		errCh,
	)

	// Check that message is sent to message channel
	s.Equal(testMessage, <-messageCh)
	stopCh <- struct{}{}
}

func (s *EVMListenerTestSuite) TestBlockDelay() {
	latestBlock := big.NewInt(100)
	delta := big.NewInt(0).Sub(listener.BlockDelay, big.NewInt(1))
	startBlock := big.NewInt(0).Sub(latestBlock, delta)

	// Define mocks

	s.chainReaderMock.EXPECT().LatestBlock().Return(latestBlock, nil).AnyTimes()

	s.chainReaderMock.EXPECT().FetchDepositLogs(
		gomock.Any(),
		gomock.Eq(testBridgeAddress),
		gomock.Eq(startBlock),
		gomock.Eq(startBlock),
	).Return([]*listener.DepositLogs{testLogsForStartBlock}, nil).AnyTimes()

	s.eventHandlerMock.EXPECT().HandleEvent(
		gomock.Eq(testDomainID),
		gomock.Eq(testLogsForStartBlock.DestinationID),
		gomock.Eq(testLogsForStartBlock.DepositNonce),
		gomock.Eq(testLogsForStartBlock.ResourceID),
	).Return(testMessage, nil).AnyTimes()

	s.keyValueWriterMock.EXPECT().SetByKey(
		gomock.Any(),
		gomock.Any(),
	).Return(nil).AnyTimes()

	// Start EVMListener

	l := listener.NewEVMListener(
		s.chainReaderMock, s.eventHandlerMock, testBridgeAddress,
	)

	stopCh := make(chan struct{})
	errCh := make(chan error)

	messageCh := l.ListenToEvents(
		startBlock,
		testDomainID,
		s.keyValueWriterMock,
		stopCh,
		errCh,
	)

	time.Sleep(listener.BlockRetryInterval * 2)

	// Check that message and error channels didn't recieve anything
	select {
	case m := <-messageCh:
		s.Fail("", m)
	case e := <-errCh:
		s.Fail("", e)
	default:
		stopCh <- struct{}{}
	}
}

func (s *EVMListenerTestSuite) TestErrorOnHandleEventReturnsError() {
	latestBlock := big.NewInt(100)
	startBlock := big.NewInt(0).Sub(latestBlock, listener.BlockDelay)

	// Define mocks

	s.chainReaderMock.EXPECT().LatestBlock().Return(latestBlock, nil)

	s.chainReaderMock.EXPECT().FetchDepositLogs(
		gomock.Any(),
		gomock.Eq(testBridgeAddress),
		gomock.Eq(startBlock),
		gomock.Eq(startBlock),
	).Return([]*listener.DepositLogs{testLogsForStartBlock}, nil)

	s.eventHandlerMock.EXPECT().HandleEvent(
		gomock.Eq(testDomainID),
		gomock.Eq(testLogsForStartBlock.DestinationID),
		gomock.Eq(testLogsForStartBlock.DepositNonce),
		gomock.Eq(testLogsForStartBlock.ResourceID),
	).Return(nil, errors.New(""))

	// Start EVMListener

	l := listener.NewEVMListener(
		s.chainReaderMock, s.eventHandlerMock, testBridgeAddress,
	)

	stopCh := make(chan struct{})
	errCh := make(chan error)

	_ = l.ListenToEvents(
		startBlock,
		testDomainID,
		s.keyValueWriterMock,
		stopCh,
		errCh,
	)

	// Check that error from event handler is propagated
	s.Error(<-errCh)
}
