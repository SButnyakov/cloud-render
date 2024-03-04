// Code generated by MockGen. DO NOT EDIT.
// Source: http/api/orderdelete.go

// Package mock_api is a generated GoMock package.
package mock_api

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOneOrderSoftDelter is a mock of OneOrderSoftDelter interface.
type MockOneOrderSoftDelter struct {
	ctrl     *gomock.Controller
	recorder *MockOneOrderSoftDelterMockRecorder
}

// MockOneOrderSoftDelterMockRecorder is the mock recorder for MockOneOrderSoftDelter.
type MockOneOrderSoftDelterMockRecorder struct {
	mock *MockOneOrderSoftDelter
}

// NewMockOneOrderSoftDelter creates a new mock instance.
func NewMockOneOrderSoftDelter(ctrl *gomock.Controller) *MockOneOrderSoftDelter {
	mock := &MockOneOrderSoftDelter{ctrl: ctrl}
	mock.recorder = &MockOneOrderSoftDelterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOneOrderSoftDelter) EXPECT() *MockOneOrderSoftDelterMockRecorder {
	return m.recorder
}

// SoftDeleteOneOrder mocks base method.
func (m *MockOneOrderSoftDelter) SoftDeleteOneOrder(id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SoftDeleteOneOrder", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// SoftDeleteOneOrder indicates an expected call of SoftDeleteOneOrder.
func (mr *MockOneOrderSoftDelterMockRecorder) SoftDeleteOneOrder(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SoftDeleteOneOrder", reflect.TypeOf((*MockOneOrderSoftDelter)(nil).SoftDeleteOneOrder), id)
}
