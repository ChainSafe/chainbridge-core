package listener_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	mock_listener "github.com/ChainSafe/chainbridge-core/chains/evm/listener/mock"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ListenerTestSuite struct {
	suite.Suite
	listener            *listener.EVMListener
	mockClient          *mock_listener.MockChainClient
	mockEventHandler    *mock_listener.MockEventHandler
	mockBlockStorer     *mock_listener.MockBlockStorer
	mockBlockDeltaMeter *mock_listener.MockBlockDeltaMeter
	domainID            uint8
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(ListenerTestSuite))
}

func (s *ListenerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.domainID = 1
	s.mockClient = mock_listener.NewMockChainClient(ctrl)
	s.mockEventHandler = mock_listener.NewMockEventHandler(ctrl)
	s.mockBlockStorer = mock_listener.NewMockBlockStorer(ctrl)
	s.mockBlockDeltaMeter = mock_listener.NewMockBlockDeltaMeter(ctrl)
	s.listener = listener.NewEVMListener(
		s.mockClient,
		[]listener.EventHandler{s.mockEventHandler, s.mockEventHandler},
		s.mockBlockStorer,
		s.mockBlockDeltaMeter,
		s.domainID,
		time.Millisecond*75,
		big.NewInt(5),
		big.NewInt(5))
}

func (s *ListenerTestSuite) Test_ListenToEvents_RetriesIfBlockUnavailable() {
	s.mockClient.EXPECT().LatestBlock().Return(big.NewInt(0), fmt.Errorf("error"))

	ctx, cancel := context.WithCancel(context.Background())
	msgChan := make(chan []*message.Message, 2)
	go s.listener.ListenToEvents(ctx, nil, msgChan, nil)

	time.Sleep(time.Millisecond * 50)
	cancel()
}

func (s *ListenerTestSuite) Test_ListenToEvents_SleepsIfBlockTooNew() {
	s.mockClient.EXPECT().LatestBlock().Return(big.NewInt(109), nil)

	ctx, cancel := context.WithCancel(context.Background())
	msgChan := make(chan []*message.Message, 2)

	go s.listener.ListenToEvents(ctx, big.NewInt(100), msgChan, nil)

	time.Sleep(time.Millisecond * 50)
	cancel()
}

func (s *ListenerTestSuite) Test_ListenToEvents_RetriesInCaseOfHandlerFailure() {
	startBlock := big.NewInt(100)
	endBlock := big.NewInt(105)
	head := big.NewInt(110)
	msgChan := make(chan []*message.Message, 2)

	// First pass
	s.mockClient.EXPECT().LatestBlock().Return(head, nil)
	s.mockBlockDeltaMeter.EXPECT().TrackBlockDelta(uint8(1), head, endBlock)
	s.mockEventHandler.EXPECT().HandleEvent(startBlock, new(big.Int).Sub(endBlock, big.NewInt(1)), msgChan).Return(fmt.Errorf("error"))
	// Second pass
	s.mockClient.EXPECT().LatestBlock().Return(head, nil)
	s.mockBlockDeltaMeter.EXPECT().TrackBlockDelta(uint8(1), head, endBlock)
	s.mockEventHandler.EXPECT().HandleEvent(startBlock, new(big.Int).Sub(endBlock, big.NewInt(1)), msgChan).Return(nil)
	s.mockEventHandler.EXPECT().HandleEvent(startBlock, new(big.Int).Sub(endBlock, big.NewInt(1)), msgChan).Return(nil)
	s.mockBlockStorer.EXPECT().StoreBlock(endBlock, s.domainID).Return(nil)
	// third pass
	s.mockClient.EXPECT().LatestBlock().Return(head, nil)

	ctx, cancel := context.WithCancel(context.Background())

	go s.listener.ListenToEvents(ctx, big.NewInt(100), msgChan, nil)

	time.Sleep(time.Millisecond * 50)
	cancel()
}

func (s *ListenerTestSuite) Test_ListenToEvents_StoresBlockIfEventHandlingSuccessful() {
	startBlock := big.NewInt(100)
	endBlock := big.NewInt(105)
	head := big.NewInt(110)
	msgChan := make(chan []*message.Message, 2)

	s.mockClient.EXPECT().LatestBlock().Return(head, nil)
	// prevent infinite runs
	s.mockClient.EXPECT().LatestBlock().Return(big.NewInt(95), nil)
	s.mockBlockDeltaMeter.EXPECT().TrackBlockDelta(uint8(1), head, endBlock)
	s.mockEventHandler.EXPECT().HandleEvent(startBlock, new(big.Int).Sub(endBlock, big.NewInt(1)), msgChan).Return(nil)
	s.mockEventHandler.EXPECT().HandleEvent(startBlock, new(big.Int).Sub(endBlock, big.NewInt(1)), msgChan).Return(nil)
	s.mockBlockStorer.EXPECT().StoreBlock(endBlock, s.domainID).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())

	go s.listener.ListenToEvents(ctx, big.NewInt(100), msgChan, nil)

	time.Sleep(time.Millisecond * 50)
	cancel()
}

func (s *ListenerTestSuite) Test_ListenToEvents_IgnoresBlocStorerError() {
	startBlock := big.NewInt(100)
	endBlock := big.NewInt(105)
	head := big.NewInt(110)
	msgChan := make(chan []*message.Message, 2)

	s.mockClient.EXPECT().LatestBlock().Return(head, nil)
	// prevent infinite runs
	s.mockClient.EXPECT().LatestBlock().Return(big.NewInt(95), nil)
	s.mockBlockDeltaMeter.EXPECT().TrackBlockDelta(uint8(1), head, endBlock)
	s.mockEventHandler.EXPECT().HandleEvent(startBlock, new(big.Int).Sub(endBlock, big.NewInt(1)), msgChan).Return(nil)
	s.mockEventHandler.EXPECT().HandleEvent(startBlock, new(big.Int).Sub(endBlock, big.NewInt(1)), msgChan).Return(nil)
	s.mockBlockStorer.EXPECT().StoreBlock(endBlock, s.domainID).Return(fmt.Errorf("error"))

	ctx, cancel := context.WithCancel(context.Background())

	go s.listener.ListenToEvents(ctx, big.NewInt(100), msgChan, nil)

	time.Sleep(time.Millisecond * 50)
	cancel()
}

func (s *ListenerTestSuite) Test_ListenToEvents_UsesHeadAsStartBlockIfNilPassed() {
	startBlock := big.NewInt(110)
	endBlock := big.NewInt(115)
	oldHead := big.NewInt(110)
	newHead := big.NewInt(120)
	msgChan := make(chan []*message.Message, 2)

	s.mockClient.EXPECT().LatestBlock().Return(oldHead, nil)
	s.mockClient.EXPECT().LatestBlock().Return(newHead, nil)
	s.mockClient.EXPECT().LatestBlock().Return(big.NewInt(65), nil)

	s.mockBlockDeltaMeter.EXPECT().TrackBlockDelta(uint8(1), big.NewInt(120), endBlock)

	s.mockEventHandler.EXPECT().HandleEvent(startBlock, new(big.Int).Sub(endBlock, big.NewInt(1)), msgChan).Return(nil)
	s.mockEventHandler.EXPECT().HandleEvent(startBlock, new(big.Int).Sub(endBlock, big.NewInt(1)), msgChan).Return(nil)
	s.mockBlockStorer.EXPECT().StoreBlock(endBlock, s.domainID).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())

	go s.listener.ListenToEvents(ctx, nil, msgChan, nil)

	time.Sleep(time.Millisecond * 100)
	cancel()
}
