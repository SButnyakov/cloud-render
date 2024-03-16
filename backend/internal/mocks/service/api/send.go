// Code generated by MockGen. DO NOT EDIT.
// Source: http/api/send.go

// Package mock_api is a generated GoMock package.
package mock_api

import (
	dto "cloud-render/internal/dto"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOrderCreator is a mock of OrderCreator interface.
type MockOrderCreator struct {
	ctrl     *gomock.Controller
	recorder *MockOrderCreatorMockRecorder
}

// MockOrderCreatorMockRecorder is the mock recorder for MockOrderCreator.
type MockOrderCreatorMockRecorder struct {
	mock *MockOrderCreator
}

// NewMockOrderCreator creates a new mock instance.
func NewMockOrderCreator(ctrl *gomock.Controller) *MockOrderCreator {
	mock := &MockOrderCreator{ctrl: ctrl}
	mock.recorder = &MockOrderCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderCreator) EXPECT() *MockOrderCreatorMockRecorder {
	return m.recorder
}

// CreateOrder mocks base method.
func (m *MockOrderCreator) CreateOrder(dto dto.CreateOrderDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrderCreatorMockRecorder) CreateOrder(dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrderCreator)(nil).CreateOrder), dto)
}
