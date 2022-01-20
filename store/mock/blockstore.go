// Code generated by MockGen. DO NOT EDIT.
// Source: blockstore/blockstore.go

// Package mock_blockstore is a generated GoMock package.
package mock_blockstore

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockKeyValueReaderWriter is a mock of KeyValueReaderWriter interface.
type MockKeyValueReaderWriter struct {
	ctrl     *gomock.Controller
	recorder *MockKeyValueReaderWriterMockRecorder
}

// MockKeyValueReaderWriterMockRecorder is the mock recorder for MockKeyValueReaderWriter.
type MockKeyValueReaderWriterMockRecorder struct {
	mock *MockKeyValueReaderWriter
}

// NewMockKeyValueReaderWriter creates a new mock instance.
func NewMockKeyValueReaderWriter(ctrl *gomock.Controller) *MockKeyValueReaderWriter {
	mock := &MockKeyValueReaderWriter{ctrl: ctrl}
	mock.recorder = &MockKeyValueReaderWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyValueReaderWriter) EXPECT() *MockKeyValueReaderWriterMockRecorder {
	return m.recorder
}

// GetByKey mocks base method.
func (m *MockKeyValueReaderWriter) GetByKey(key []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByKey", key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByKey indicates an expected call of GetByKey.
func (mr *MockKeyValueReaderWriterMockRecorder) GetByKey(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKey", reflect.TypeOf((*MockKeyValueReaderWriter)(nil).GetByKey), key)
}

// SetByKey mocks base method.
func (m *MockKeyValueReaderWriter) SetByKey(key, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetByKey", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetByKey indicates an expected call of SetByKey.
func (mr *MockKeyValueReaderWriterMockRecorder) SetByKey(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetByKey", reflect.TypeOf((*MockKeyValueReaderWriter)(nil).SetByKey), key, value)
}

// MockKeyValueReader is a mock of KeyValueReader interface.
type MockKeyValueReader struct {
	ctrl     *gomock.Controller
	recorder *MockKeyValueReaderMockRecorder
}

// MockKeyValueReaderMockRecorder is the mock recorder for MockKeyValueReader.
type MockKeyValueReaderMockRecorder struct {
	mock *MockKeyValueReader
}

// NewMockKeyValueReader creates a new mock instance.
func NewMockKeyValueReader(ctrl *gomock.Controller) *MockKeyValueReader {
	mock := &MockKeyValueReader{ctrl: ctrl}
	mock.recorder = &MockKeyValueReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyValueReader) EXPECT() *MockKeyValueReaderMockRecorder {
	return m.recorder
}

// GetByKey mocks base method.
func (m *MockKeyValueReader) GetByKey(key []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByKey", key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByKey indicates an expected call of GetByKey.
func (mr *MockKeyValueReaderMockRecorder) GetByKey(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKey", reflect.TypeOf((*MockKeyValueReader)(nil).GetByKey), key)
}

// MockKeyValueWriter is a mock of KeyValueWriter interface.
type MockKeyValueWriter struct {
	ctrl     *gomock.Controller
	recorder *MockKeyValueWriterMockRecorder
}

// MockKeyValueWriterMockRecorder is the mock recorder for MockKeyValueWriter.
type MockKeyValueWriterMockRecorder struct {
	mock *MockKeyValueWriter
}

// NewMockKeyValueWriter creates a new mock instance.
func NewMockKeyValueWriter(ctrl *gomock.Controller) *MockKeyValueWriter {
	mock := &MockKeyValueWriter{ctrl: ctrl}
	mock.recorder = &MockKeyValueWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyValueWriter) EXPECT() *MockKeyValueWriterMockRecorder {
	return m.recorder
}

// SetByKey mocks base method.
func (m *MockKeyValueWriter) SetByKey(key, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetByKey", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetByKey indicates an expected call of SetByKey.
func (mr *MockKeyValueWriterMockRecorder) SetByKey(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetByKey", reflect.TypeOf((*MockKeyValueWriter)(nil).SetByKey), key, value)
}