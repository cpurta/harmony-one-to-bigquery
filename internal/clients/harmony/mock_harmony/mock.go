// Code generated by MockGen. DO NOT EDIT.
// Source: ./harmony.go

// Package mock_harmony is a generated GoMock package.
package mock_harmony

import (
	reflect "reflect"

	model "github.com/cpurta/harmony-one-to-bigquery/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockHarmonyClient is a mock of HarmonyClient interface.
type MockHarmonyClient struct {
	ctrl     *gomock.Controller
	recorder *MockHarmonyClientMockRecorder
}

// MockHarmonyClientMockRecorder is the mock recorder for MockHarmonyClient.
type MockHarmonyClientMockRecorder struct {
	mock *MockHarmonyClient
}

// NewMockHarmonyClient creates a new mock instance.
func NewMockHarmonyClient(ctrl *gomock.Controller) *MockHarmonyClient {
	mock := &MockHarmonyClient{ctrl: ctrl}
	mock.recorder = &MockHarmonyClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHarmonyClient) EXPECT() *MockHarmonyClientMockRecorder {
	return m.recorder
}

// GetBlockByNumber mocks base method.
func (m *MockHarmonyClient) GetBlockByNumber(blockNumber int64) (*model.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockByNumber", blockNumber)
	ret0, _ := ret[0].(*model.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockByNumber indicates an expected call of GetBlockByNumber.
func (mr *MockHarmonyClientMockRecorder) GetBlockByNumber(blockNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockByNumber", reflect.TypeOf((*MockHarmonyClient)(nil).GetBlockByNumber), blockNumber)
}

// GetLatestHeader mocks base method.
func (m *MockHarmonyClient) GetLatestHeader() (*model.Header, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestHeader")
	ret0, _ := ret[0].(*model.Header)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestHeader indicates an expected call of GetLatestHeader.
func (mr *MockHarmonyClientMockRecorder) GetLatestHeader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestHeader", reflect.TypeOf((*MockHarmonyClient)(nil).GetLatestHeader))
}