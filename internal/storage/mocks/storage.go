// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/red-rocket-software/reminder-go/internal/app/model"
	storage "github.com/red-rocket-software/reminder-go/internal/storage"
	pagination "github.com/red-rocket-software/reminder-go/pkg/pagination"
)

// MockReminderRepo is a mock of ReminderRepo interface.
type MockReminderRepo struct {
	ctrl     *gomock.Controller
	recorder *MockReminderRepoMockRecorder
}

// MockReminderRepoMockRecorder is the mock recorder for MockReminderRepo.
type MockReminderRepoMockRecorder struct {
	mock *MockReminderRepo
}

// NewMockReminderRepo creates a new mock instance.
func NewMockReminderRepo(ctrl *gomock.Controller) *MockReminderRepo {
	mock := &MockReminderRepo{ctrl: ctrl}
	mock.recorder = &MockReminderRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReminderRepo) EXPECT() *MockReminderRepoMockRecorder {
	return m.recorder
}

// CreateRemind mocks base method.
func (m *MockReminderRepo) CreateRemind(ctx context.Context, todo model.Todo) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRemind", ctx, todo)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRemind indicates an expected call of CreateRemind.
func (mr *MockReminderRepoMockRecorder) CreateRemind(ctx, todo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRemind", reflect.TypeOf((*MockReminderRepo)(nil).CreateRemind), ctx, todo)
}

// CreateUser mocks base method.
func (m *MockReminderRepo) CreateUser(ctx context.Context, input model.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, input)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockReminderRepoMockRecorder) CreateUser(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockReminderRepo)(nil).CreateUser), ctx, input)
}

// DeleteRemind mocks base method.
func (m *MockReminderRepo) DeleteRemind(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRemind", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRemind indicates an expected call of DeleteRemind.
func (mr *MockReminderRepoMockRecorder) DeleteRemind(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRemind", reflect.TypeOf((*MockReminderRepo)(nil).DeleteRemind), ctx, id)
}

// GetAllReminds mocks base method.
func (m *MockReminderRepo) GetAllReminds(ctx context.Context, params pagination.Page, userID int) ([]model.Todo, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllReminds", ctx, params, userID)
	ret0, _ := ret[0].([]model.Todo)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAllReminds indicates an expected call of GetAllReminds.
func (mr *MockReminderRepoMockRecorder) GetAllReminds(ctx, params, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllReminds", reflect.TypeOf((*MockReminderRepo)(nil).GetAllReminds), ctx, params, userID)
}

// GetCompletedReminds mocks base method.
func (m *MockReminderRepo) GetCompletedReminds(ctx context.Context, params storage.Params, userID int) ([]model.Todo, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompletedReminds", ctx, params, userID)
	ret0, _ := ret[0].([]model.Todo)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCompletedReminds indicates an expected call of GetCompletedReminds.
func (mr *MockReminderRepoMockRecorder) GetCompletedReminds(ctx, params, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompletedReminds", reflect.TypeOf((*MockReminderRepo)(nil).GetCompletedReminds), ctx, params, userID)
}

// GetNewReminds mocks base method.
func (m *MockReminderRepo) GetNewReminds(ctx context.Context, params pagination.Page, userID int) ([]model.Todo, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewReminds", ctx, params, userID)
	ret0, _ := ret[0].([]model.Todo)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNewReminds indicates an expected call of GetNewReminds.
func (mr *MockReminderRepoMockRecorder) GetNewReminds(ctx, params, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewReminds", reflect.TypeOf((*MockReminderRepo)(nil).GetNewReminds), ctx, params, userID)
}

// GetRemindByID mocks base method.
func (m *MockReminderRepo) GetRemindByID(ctx context.Context, id int) (model.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRemindByID", ctx, id)
	ret0, _ := ret[0].(model.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRemindByID indicates an expected call of GetRemindByID.
func (mr *MockReminderRepoMockRecorder) GetRemindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRemindByID", reflect.TypeOf((*MockReminderRepo)(nil).GetRemindByID), ctx, id)
}

// GetRemindsForNotification mocks base method.
func (m *MockReminderRepo) GetRemindsForNotification(ctx context.Context, days int) ([]model.NotificationRemind, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRemindsForNotification", ctx, days)
	ret0, _ := ret[0].([]model.NotificationRemind)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRemindsForNotification indicates an expected call of GetRemindsForNotification.
func (mr *MockReminderRepoMockRecorder) GetRemindsForNotification(ctx, days interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRemindsForNotification", reflect.TypeOf((*MockReminderRepo)(nil).GetRemindsForNotification), ctx, days)
}

// GetUserByEmail mocks base method.
func (m *MockReminderRepo) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockReminderRepoMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockReminderRepo)(nil).GetUserByEmail), ctx, email)
}

// GetUserByID mocks base method.
func (m *MockReminderRepo) GetUserByID(ctx context.Context, id int) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, id)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockReminderRepoMockRecorder) GetUserByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockReminderRepo)(nil).GetUserByID), ctx, id)
}

// SeedTodos mocks base method.
func (m *MockReminderRepo) SeedTodos() ([]model.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SeedTodos")
	ret0, _ := ret[0].([]model.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SeedTodos indicates an expected call of SeedTodos.
func (mr *MockReminderRepoMockRecorder) SeedTodos() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SeedTodos", reflect.TypeOf((*MockReminderRepo)(nil).SeedTodos))
}

// SeedUser mocks base method.
func (m *MockReminderRepo) SeedUser() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SeedUser")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SeedUser indicates an expected call of SeedUser.
func (mr *MockReminderRepoMockRecorder) SeedUser() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SeedUser", reflect.TypeOf((*MockReminderRepo)(nil).SeedUser))
}

// Truncate mocks base method.
func (m *MockReminderRepo) Truncate() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Truncate")
	ret0, _ := ret[0].(error)
	return ret0
}

// Truncate indicates an expected call of Truncate.
func (mr *MockReminderRepoMockRecorder) Truncate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Truncate", reflect.TypeOf((*MockReminderRepo)(nil).Truncate))
}

// UpdateNotification mocks base method.
func (m *MockReminderRepo) UpdateNotification(ctx context.Context, id int, dao model.NotificationDAO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNotification", ctx, id, dao)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNotification indicates an expected call of UpdateNotification.
func (mr *MockReminderRepoMockRecorder) UpdateNotification(ctx, id, dao interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNotification", reflect.TypeOf((*MockReminderRepo)(nil).UpdateNotification), ctx, id, dao)
}

// UpdateRemind mocks base method.
func (m *MockReminderRepo) UpdateRemind(ctx context.Context, id int, input model.TodoUpdateInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRemind", ctx, id, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRemind indicates an expected call of UpdateRemind.
func (mr *MockReminderRepoMockRecorder) UpdateRemind(ctx, id, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRemind", reflect.TypeOf((*MockReminderRepo)(nil).UpdateRemind), ctx, id, input)
}

// UpdateStatus mocks base method.
func (m *MockReminderRepo) UpdateStatus(ctx context.Context, id int, updateInput model.TodoUpdateStatusInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, id, updateInput)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockReminderRepoMockRecorder) UpdateStatus(ctx, id, updateInput interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockReminderRepo)(nil).UpdateStatus), ctx, id, updateInput)
}

// UpdateUser mocks base method.
func (m *MockReminderRepo) UpdateUser(ctx context.Context, id int, input model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, id, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockReminderRepoMockRecorder) UpdateUser(ctx, id, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockReminderRepo)(nil).UpdateUser), ctx, id, input)
}

// UpdateUserNotification mocks base method.
func (m *MockReminderRepo) UpdateUserNotification(ctx context.Context, id int, status bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserNotification", ctx, id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserNotification indicates an expected call of UpdateUserNotification.
func (mr *MockReminderRepoMockRecorder) UpdateUserNotification(ctx, id, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserNotification", reflect.TypeOf((*MockReminderRepo)(nil).UpdateUserNotification), ctx, id, status)
}
