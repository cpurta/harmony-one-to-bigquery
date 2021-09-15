// Code generated by MockGen. DO NOT EDIT.
// Source: ./bigquery.go

// Package mock_bigquery is a generated GoMock package.
package mock_bigquery

import (
	context "context"
	reflect "reflect"

	model "github.com/cpurta/harmony-one-to-bigquery/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockBigQueryClient is a mock of BigQueryClient interface.
type MockBigQueryClient struct {
	ctrl     *gomock.Controller
	recorder *MockBigQueryClientMockRecorder
}

// MockBigQueryClientMockRecorder is the mock recorder for MockBigQueryClient.
type MockBigQueryClientMockRecorder struct {
	mock *MockBigQueryClient
}

// NewMockBigQueryClient creates a new mock instance.
func NewMockBigQueryClient(ctrl *gomock.Controller) *MockBigQueryClient {
	mock := &MockBigQueryClient{ctrl: ctrl}
	mock.recorder = &MockBigQueryClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBigQueryClient) EXPECT() *MockBigQueryClientMockRecorder {
	return m.recorder
}

// BlocksTableExists mocks base method.
func (m *MockBigQueryClient) BlocksTableExists(ctx context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlocksTableExists", ctx)
	ret0, _ := ret[0].(bool)
	return ret0
}

// BlocksTableExists indicates an expected call of BlocksTableExists.
func (mr *MockBigQueryClientMockRecorder) BlocksTableExists(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlocksTableExists", reflect.TypeOf((*MockBigQueryClient)(nil).BlocksTableExists), ctx)
}

// Close mocks base method.
func (m *MockBigQueryClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockBigQueryClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBigQueryClient)(nil).Close))
}

// CreateBlocksTable mocks base method.
func (m *MockBigQueryClient) CreateBlocksTable(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBlocksTable", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBlocksTable indicates an expected call of CreateBlocksTable.
func (mr *MockBigQueryClientMockRecorder) CreateBlocksTable(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBlocksTable", reflect.TypeOf((*MockBigQueryClient)(nil).CreateBlocksTable), ctx)
}

// CreateProjectDataset mocks base method.
func (m *MockBigQueryClient) CreateProjectDataset(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProjectDataset", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProjectDataset indicates an expected call of CreateProjectDataset.
func (mr *MockBigQueryClientMockRecorder) CreateProjectDataset(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProjectDataset", reflect.TypeOf((*MockBigQueryClient)(nil).CreateProjectDataset), ctx)
}

// CreateTransactionsTable mocks base method.
func (m *MockBigQueryClient) CreateTransactionsTable(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransactionsTable", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTransactionsTable indicates an expected call of CreateTransactionsTable.
func (mr *MockBigQueryClientMockRecorder) CreateTransactionsTable(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransactionsTable", reflect.TypeOf((*MockBigQueryClient)(nil).CreateTransactionsTable), ctx)
}

// GetMostRecentBlockNumber mocks base method.
func (m *MockBigQueryClient) GetMostRecentBlockNumber(ctx context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMostRecentBlockNumber", ctx)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMostRecentBlockNumber indicates an expected call of GetMostRecentBlockNumber.
func (mr *MockBigQueryClientMockRecorder) GetMostRecentBlockNumber(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMostRecentBlockNumber", reflect.TypeOf((*MockBigQueryClient)(nil).GetMostRecentBlockNumber), ctx)
}

// InsertBlock mocks base method.
func (m *MockBigQueryClient) InsertBlock(block *model.Block, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertBlock", block, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertBlock indicates an expected call of InsertBlock.
func (mr *MockBigQueryClientMockRecorder) InsertBlock(block, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertBlock", reflect.TypeOf((*MockBigQueryClient)(nil).InsertBlock), block, ctx)
}

// InsertTransactions mocks base method.
func (m *MockBigQueryClient) InsertTransactions(transactions []*model.Transaction, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTransactions", transactions, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertTransactions indicates an expected call of InsertTransactions.
func (mr *MockBigQueryClientMockRecorder) InsertTransactions(transactions, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTransactions", reflect.TypeOf((*MockBigQueryClient)(nil).InsertTransactions), transactions, ctx)
}

// ProjectDatasetExists mocks base method.
func (m *MockBigQueryClient) ProjectDatasetExists(ctx context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProjectDatasetExists", ctx)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ProjectDatasetExists indicates an expected call of ProjectDatasetExists.
func (mr *MockBigQueryClientMockRecorder) ProjectDatasetExists(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProjectDatasetExists", reflect.TypeOf((*MockBigQueryClient)(nil).ProjectDatasetExists), ctx)
}

// TransactionsTableExists mocks base method.
func (m *MockBigQueryClient) TransactionsTableExists(ctx context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionsTableExists", ctx)
	ret0, _ := ret[0].(bool)
	return ret0
}

// TransactionsTableExists indicates an expected call of TransactionsTableExists.
func (mr *MockBigQueryClientMockRecorder) TransactionsTableExists(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionsTableExists", reflect.TypeOf((*MockBigQueryClient)(nil).TransactionsTableExists), ctx)
}
