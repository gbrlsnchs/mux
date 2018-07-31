// Code generated by MockGen. DO NOT EDIT.
// Source: net/http (interfaces: Handler)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// MockHandler is a mock of Handler interface
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// ServeHTTP mocks base method
func (m *MockHandler) ServeHTTP(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.Call(m, "ServeHTTP", arg0, arg1)
}

// ServeHTTP indicates an expected call of ServeHTTP
func (mr *MockHandlerMockRecorder) ServeHTTP(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServeHTTP", reflect.TypeOf((*MockHandler)(nil).ServeHTTP), arg0, arg1)
}
