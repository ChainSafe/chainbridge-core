// Code generated by MockGen. DO NOT EDIT.
// Source: ./chains/evm/listener/listener.go

// Package mock_listener is a generated GoMock package.
package mock_listener

import (
	big "math/big"
	reflect "reflect"

	message "github.com/ChainSafe/chainbridge-core/relayer/message"
	gomock "github.com/golang/mock/gomock"
)

// MockEventHandler is a mock of EventHandler interface.
type MockEventHandler struct {
	ctrl     *gomock.Controller
	recorder *MockEventHandlerMockRecorder
}

// MockEventHandlerMockRecorder is the mock recorder for MockEventHandler.
type MockEventHandlerMockRecorder struct {
	mock *MockEventHandler
}

// NewMockEventHandler creates a new mock instance.
func NewMockEventHandler(ctrl *gomock.Controller) *MockEventHandler {
	mock := &MockEventHandler{ctrl: ctrl}
	mock.recorder = &MockEventHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventHandler) EXPECT() *MockEventHandlerMockRecorder {
	return m.recorder
}

// HandleEvent mocks base method.
func (m *MockEventHandler) HandleEvent(startBlock, endBlock *big.Int, msgChan chan []*message.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleEvent", startBlock, endBlock, msgChan)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleEvent indicates an expected call of HandleEvent.
func (mr *MockEventHandlerMockRecorder) HandleEvent(startBlock, endBlock, msgChan interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleEvent", reflect.TypeOf((*MockEventHandler)(nil).HandleEvent), startBlock, endBlock, msgChan)
}

// MockChainClient is a mock of ChainClient interface.
type MockChainClient struct {
	ctrl     *gomock.Controller
	recorder *MockChainClientMockRecorder
}

// MockChainClientMockRecorder is the mock recorder for MockChainClient.
type MockChainClientMockRecorder struct {
	mock *MockChainClient
}

// NewMockChainClient creates a new mock instance.
func NewMockChainClient(ctrl *gomock.Controller) *MockChainClient {
	mock := &MockChainClient{ctrl: ctrl}
	mock.recorder = &MockChainClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChainClient) EXPECT() *MockChainClientMockRecorder {
	return m.recorder
}

// LatestBlock mocks base method.
func (m *MockChainClient) LatestBlock() (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LatestBlock")
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LatestBlock indicates an expected call of LatestBlock.
func (mr *MockChainClientMockRecorder) LatestBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LatestBlock", reflect.TypeOf((*MockChainClient)(nil).LatestBlock))
}

// MockBlockDeltaMeter is a mock of BlockDeltaMeter interface.
type MockBlockDeltaMeter struct {
	ctrl     *gomock.Controller
	recorder *MockBlockDeltaMeterMockRecorder
}

// MockBlockDeltaMeterMockRecorder is the mock recorder for MockBlockDeltaMeter.
type MockBlockDeltaMeterMockRecorder struct {
	mock *MockBlockDeltaMeter
}

// NewMockBlockDeltaMeter creates a new mock instance.
func NewMockBlockDeltaMeter(ctrl *gomock.Controller) *MockBlockDeltaMeter {
	mock := &MockBlockDeltaMeter{ctrl: ctrl}
	mock.recorder = &MockBlockDeltaMeterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockDeltaMeter) EXPECT() *MockBlockDeltaMeterMockRecorder {
	return m.recorder
}

// TrackBlockDelta mocks base method.
func (m *MockBlockDeltaMeter) TrackBlockDelta(domainID uint8, head, current *big.Int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "TrackBlockDelta", domainID, head, current)
}

// TrackBlockDelta indicates an expected call of TrackBlockDelta.
func (mr *MockBlockDeltaMeterMockRecorder) TrackBlockDelta(domainID, head, current interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TrackBlockDelta", reflect.TypeOf((*MockBlockDeltaMeter)(nil).TrackBlockDelta), domainID, head, current)
}

// MockBlockStorer is a mock of BlockStorer interface.
type MockBlockStorer struct {
	ctrl     *gomock.Controller
	recorder *MockBlockStorerMockRecorder
}

// MockBlockStorerMockRecorder is the mock recorder for MockBlockStorer.
type MockBlockStorerMockRecorder struct {
	mock *MockBlockStorer
}

// NewMockBlockStorer creates a new mock instance.
func NewMockBlockStorer(ctrl *gomock.Controller) *MockBlockStorer {
	mock := &MockBlockStorer{ctrl: ctrl}
	mock.recorder = &MockBlockStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockStorer) EXPECT() *MockBlockStorerMockRecorder {
	return m.recorder
}

// StoreBlock mocks base method.
func (m *MockBlockStorer) StoreBlock(block *big.Int, domainID uint8) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreBlock", block, domainID)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreBlock indicates an expected call of StoreBlock.
func (mr *MockBlockStorerMockRecorder) StoreBlock(block, domainID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreBlock", reflect.TypeOf((*MockBlockStorer)(nil).StoreBlock), block, domainID)
}
