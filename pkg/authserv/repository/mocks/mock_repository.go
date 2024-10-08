// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/repository (interfaces: AuthRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	auth "github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/auth"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockAuthRepository is a mock of AuthRepository interface
type MockAuthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAuthRepositoryMockRecorder
}

// MockAuthRepositoryMockRecorder is the mock recorder for MockAuthRepository
type MockAuthRepositoryMockRecorder struct {
	mock *MockAuthRepository
}

// NewMockAuthRepository creates a new mock instance
func NewMockAuthRepository(ctrl *gomock.Controller) *MockAuthRepository {
	mock := &MockAuthRepository{ctrl: ctrl}
	mock.recorder = &MockAuthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthRepository) EXPECT() *MockAuthRepositoryMockRecorder {
	return m.recorder
}

// AddToBlacklist mocks base method
func (m *MockAuthRepository) AddToBlacklist(arg0 string, arg1 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToBlacklist", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToBlacklist indicates an expected call of AddToBlacklist
func (mr *MockAuthRepositoryMockRecorder) AddToBlacklist(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToBlacklist", reflect.TypeOf((*MockAuthRepository)(nil).AddToBlacklist), arg0, arg1)
}

// CleanExpiredTokens mocks base method
func (m *MockAuthRepository) CleanExpiredTokens() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanExpiredTokens")
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanExpiredTokens indicates an expected call of CleanExpiredTokens
func (mr *MockAuthRepositoryMockRecorder) CleanExpiredTokens() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanExpiredTokens", reflect.TypeOf((*MockAuthRepository)(nil).CleanExpiredTokens))
}

// CreateUser mocks base method
func (m *MockAuthRepository) CreateUser(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockAuthRepositoryMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthRepository)(nil).CreateUser), arg0, arg1)
}

// GetUser mocks base method
func (m *MockAuthRepository) GetUser(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockAuthRepositoryMockRecorder) GetUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockAuthRepository)(nil).GetUser), arg0)
}

// IsInBlacklist mocks base method
func (m *MockAuthRepository) IsInBlacklist(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsInBlacklist", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsInBlacklist indicates an expected call of IsInBlacklist
func (mr *MockAuthRepositoryMockRecorder) IsInBlacklist(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsInBlacklist", reflect.TypeOf((*MockAuthRepository)(nil).IsInBlacklist), arg0)
}

// ValidateToken mocks base method
func (m *MockAuthRepository) ValidateToken(arg0 string) (*auth.Claims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", arg0)
	ret0, _ := ret[0].(*auth.Claims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken
func (mr *MockAuthRepositoryMockRecorder) ValidateToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockAuthRepository)(nil).ValidateToken), arg0)
}
