// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ruslanjo/url_shortener/internal/app/storage (interfaces: Storage)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/ruslanjo/url_shortener/internal/app/storage/models"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddShortURL mocks base method.
func (m *MockStorage) AddShortURL(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddShortURL", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddShortURL indicates an expected call of AddShortURL.
func (mr *MockStorageMockRecorder) AddShortURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddShortURL", reflect.TypeOf((*MockStorage)(nil).AddShortURL), arg0, arg1)
}

// GetURLByShortLink mocks base method.
func (m *MockStorage) GetURLByShortLink(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLByShortLink", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLByShortLink indicates an expected call of GetURLByShortLink.
func (mr *MockStorageMockRecorder) GetURLByShortLink(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLByShortLink", reflect.TypeOf((*MockStorage)(nil).GetURLByShortLink), arg0)
}

// PingContext mocks base method.
func (m *MockStorage) PingContext(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingContext", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PingContext indicates an expected call of PingContext.
func (mr *MockStorageMockRecorder) PingContext(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingContext", reflect.TypeOf((*MockStorage)(nil).PingContext), arg0)
}

// SaveURLBatched mocks base method.
func (m *MockStorage) SaveURLBatched(arg0 context.Context, arg1 []models.URLBatch) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveURLBatched", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveURLBatched indicates an expected call of SaveURLBatched.
func (mr *MockStorageMockRecorder) SaveURLBatched(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveURLBatched", reflect.TypeOf((*MockStorage)(nil).SaveURLBatched), arg0, arg1)
}
