// Code generated by MockGen. DO NOT EDIT.
// Source: user-configs.go

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
)

// MockConfigRepository is a mock of ConfigRepository interface.
type MockConfigRepository struct {
	ctrl     *gomock.Controller
	recorder *MockConfigRepositoryMockRecorder
}

// MockConfigRepositoryMockRecorder is the mock recorder for MockConfigRepository.
type MockConfigRepositoryMockRecorder struct {
	mock *MockConfigRepository
}

// NewMockConfigRepository creates a new mock instance.
func NewMockConfigRepository(ctrl *gomock.Controller) *MockConfigRepository {
	mock := &MockConfigRepository{ctrl: ctrl}
	mock.recorder = &MockConfigRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigRepository) EXPECT() *MockConfigRepositoryMockRecorder {
	return m.recorder
}

// CreateUserConfigs mocks base method.
func (m *MockConfigRepository) CreateUserConfigs(ctx context.Context, userID string) (domain.UserConfigs, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserConfigs", ctx, userID)
	ret0, _ := ret[0].(domain.UserConfigs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserConfigs indicates an expected call of CreateUserConfigs.
func (mr *MockConfigRepositoryMockRecorder) CreateUserConfigs(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserConfigs", reflect.TypeOf((*MockConfigRepository)(nil).CreateUserConfigs), ctx, userID)
}

// GetUserConfigs mocks base method.
func (m *MockConfigRepository) GetUserConfigs(ctx context.Context, userID string) (domain.UserConfigs, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserConfigs", ctx, userID)
	ret0, _ := ret[0].(domain.UserConfigs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserConfigs indicates an expected call of GetUserConfigs.
func (mr *MockConfigRepositoryMockRecorder) GetUserConfigs(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserConfigs", reflect.TypeOf((*MockConfigRepository)(nil).GetUserConfigs), ctx, userID)
}

// UpdateUserConfig mocks base method.
func (m *MockConfigRepository) UpdateUserConfig(ctx context.Context, id string, input domain.UserConfigs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserConfig", ctx, id, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserConfig indicates an expected call of UpdateUserConfig.
func (mr *MockConfigRepositoryMockRecorder) UpdateUserConfig(ctx, id, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserConfig", reflect.TypeOf((*MockConfigRepository)(nil).UpdateUserConfig), ctx, id, input)
}
