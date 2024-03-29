// Code generated by MockGen. DO NOT EDIT.
// Source: todo.go

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
)

// MockTodoRepository is a mock of TodoRepository interface.
type MockTodoRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTodoRepositoryMockRecorder
}

// MockTodoRepositoryMockRecorder is the mock recorder for MockTodoRepository.
type MockTodoRepositoryMockRecorder struct {
	mock *MockTodoRepository
}

// NewMockTodoRepository creates a new mock instance.
func NewMockTodoRepository(ctrl *gomock.Controller) *MockTodoRepository {
	mock := &MockTodoRepository{ctrl: ctrl}
	mock.recorder = &MockTodoRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodoRepository) EXPECT() *MockTodoRepositoryMockRecorder {
	return m.recorder
}

// CreateRemind mocks base method.
func (m *MockTodoRepository) CreateRemind(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRemind", ctx, todo)
	ret0, _ := ret[0].(domain.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRemind indicates an expected call of CreateRemind.
func (mr *MockTodoRepositoryMockRecorder) CreateRemind(ctx, todo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRemind", reflect.TypeOf((*MockTodoRepository)(nil).CreateRemind), ctx, todo)
}

// DeleteRemind mocks base method.
func (m *MockTodoRepository) DeleteRemind(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRemind", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRemind indicates an expected call of DeleteRemind.
func (mr *MockTodoRepositoryMockRecorder) DeleteRemind(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRemind", reflect.TypeOf((*MockTodoRepository)(nil).DeleteRemind), ctx, id)
}

// GetRemindByID mocks base method.
func (m *MockTodoRepository) GetRemindByID(ctx context.Context, id int) (domain.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRemindByID", ctx, id)
	ret0, _ := ret[0].(domain.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRemindByID indicates an expected call of GetRemindByID.
func (mr *MockTodoRepositoryMockRecorder) GetRemindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRemindByID", reflect.TypeOf((*MockTodoRepository)(nil).GetRemindByID), ctx, id)
}

// GetReminds mocks base method.
func (m *MockTodoRepository) GetReminds(ctx context.Context, params domain.FetchParams, userID string) ([]domain.Todo, int, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReminds", ctx, params, userID)
	ret0, _ := ret[0].([]domain.Todo)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(int)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetReminds indicates an expected call of GetReminds.
func (mr *MockTodoRepositoryMockRecorder) GetReminds(ctx, params, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReminds", reflect.TypeOf((*MockTodoRepository)(nil).GetReminds), ctx, params, userID)
}

// GetRemindsForDeadlineNotification mocks base method.
func (m *MockTodoRepository) GetRemindsForDeadlineNotification(ctx context.Context) ([]domain.NotificationRemind, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRemindsForDeadlineNotification", ctx)
	ret0, _ := ret[0].([]domain.NotificationRemind)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetRemindsForDeadlineNotification indicates an expected call of GetRemindsForDeadlineNotification.
func (mr *MockTodoRepositoryMockRecorder) GetRemindsForDeadlineNotification(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRemindsForDeadlineNotification", reflect.TypeOf((*MockTodoRepository)(nil).GetRemindsForDeadlineNotification), ctx)
}

// GetRemindsForNotification mocks base method.
func (m *MockTodoRepository) GetRemindsForNotification(ctx context.Context) ([]domain.NotificationRemind, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRemindsForNotification", ctx)
	ret0, _ := ret[0].([]domain.NotificationRemind)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRemindsForNotification indicates an expected call of GetRemindsForNotification.
func (mr *MockTodoRepositoryMockRecorder) GetRemindsForNotification(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRemindsForNotification", reflect.TypeOf((*MockTodoRepository)(nil).GetRemindsForNotification), ctx)
}

// UpdateNotification mocks base method.
func (m *MockTodoRepository) UpdateNotification(ctx context.Context, id int, dao domain.NotificationDAO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNotification", ctx, id, dao)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNotification indicates an expected call of UpdateNotification.
func (mr *MockTodoRepositoryMockRecorder) UpdateNotification(ctx, id, dao interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNotification", reflect.TypeOf((*MockTodoRepository)(nil).UpdateNotification), ctx, id, dao)
}

// UpdateNotifyPeriod mocks base method.
func (m *MockTodoRepository) UpdateNotifyPeriod(ctx context.Context, id int, timeToDelete string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNotifyPeriod", ctx, id, timeToDelete)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNotifyPeriod indicates an expected call of UpdateNotifyPeriod.
func (mr *MockTodoRepositoryMockRecorder) UpdateNotifyPeriod(ctx, id, timeToDelete interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNotifyPeriod", reflect.TypeOf((*MockTodoRepository)(nil).UpdateNotifyPeriod), ctx, id, timeToDelete)
}

// UpdateRemind mocks base method.
func (m *MockTodoRepository) UpdateRemind(ctx context.Context, id int, input domain.TodoUpdateInput) (domain.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRemind", ctx, id, input)
	ret0, _ := ret[0].(domain.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateRemind indicates an expected call of UpdateRemind.
func (mr *MockTodoRepositoryMockRecorder) UpdateRemind(ctx, id, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRemind", reflect.TypeOf((*MockTodoRepository)(nil).UpdateRemind), ctx, id, input)
}

// UpdateStatus mocks base method.
func (m *MockTodoRepository) UpdateStatus(ctx context.Context, id int, updateInput domain.TodoUpdateStatusInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, id, updateInput)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockTodoRepositoryMockRecorder) UpdateStatus(ctx, id, updateInput interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockTodoRepository)(nil).UpdateStatus), ctx, id, updateInput)
}
