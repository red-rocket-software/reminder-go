package server

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/red-rocket-software/reminder-go/internal/storage"
	mockdb "github.com/red-rocket-software/reminder-go/internal/storage/mocks"
	"github.com/red-rocket-software/reminder-go/pkg/pagination"
	"github.com/stretchr/testify/require"
)

func TestControllers_AddRemind(t *testing.T) {
	dTime, err := time.Parse("2006-01-02T15:04", "2023-01-19T22:02")
	require.NoError(t, err)
	now, err := time.Parse("02.01.2006, 15:04:05", "19.01.2023, 22:15:30")
	require.NoError(t, err)

	testCases := []struct {
		name                 string
		body                 string
		inputTodo            model.Todo
		mockBehavior         func(store *mockdb.MockReminderRepo, input model.Todo)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			body: `{"description":"test", "user_id": "1", "deadline_at": "2023-01-19T22:02", "created_at": "19.01.2023, 22:15:30"}`,
			inputTodo: model.Todo{
				CreatedAt:   now,
				UserID:      1,
				Description: "test",
				DeadlineAt:  dTime,
			},
			mockBehavior: func(store *mockdb.MockReminderRepo, input model.Todo) {
				store.EXPECT().CreateRemind(gomock.Any(), input).Times(1).Return(0, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: "Remind is successfully created",
		},
		{
			name:                 "Error - wrong input",
			body:                 `{"description":"", "user_id": "1", "deadline_at": "2023-02-02"}`,
			inputTodo:            model.Todo{},
			mockBehavior:         func(store *mockdb.MockReminderRepo, input model.Todo) {},
			expectedStatusCode:   422,
			expectedResponseBody: "nothing to save",
		},
		{
			name: "Error - Service error",
			body: `{"description":"test", "user_id": "1", "deadline_at": "2023-01-19T22:02", "created_at": "19.01.2023, 22:15:30"}`,
			inputTodo: model.Todo{
				CreatedAt:   now,
				UserID:      1,
				Description: "test",
				DeadlineAt:  dTime,
			},
			mockBehavior: func(store *mockdb.MockReminderRepo, input model.Todo) {
				store.EXPECT().CreateRemind(gomock.Any(), input).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `"error":"something went wrong"`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := mockdb.NewMockReminderRepo(c)
			test.mockBehavior(store, test.inputTodo)

			server := newTestServer(store)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/remind", bytes.NewBufferString(test.body))
			ctx := req.Context()
			ctx = context.WithValue(ctx, "currentUser", model.User{
				ID: 1,
			})
			req = req.WithContext(ctx)

			handler := http.HandlerFunc(server.AddRemind)
			handler.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)
			require.Contains(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestControllers_GetRemindByID(t *testing.T) {
	testCases := []struct {
		name               string
		id                 int
		resTodo            model.Todo
		mockBehavior       func(store *mockdb.MockReminderRepo, id int)
		expectedStatusCode int
	}{
		{
			name: "OK",
			id:   1,
			mockBehavior: func(store *mockdb.MockReminderRepo, id int) {
				store.EXPECT().GetRemindByID(gomock.Any(), gomock.Eq(1)).Return(model.Todo{
					ID:          1,
					Description: "test",
					CreatedAt:   time.Now(),
					DeadlineAt:  time.Now(),
					Completed:   false,
				}, nil).Times(1)
			},
			expectedStatusCode: 200,
		},
		{
			name: "Not found",
			id:   1,
			mockBehavior: func(store *mockdb.MockReminderRepo, id int) {
				store.EXPECT().GetRemindByID(gomock.Any(), gomock.Eq(id)).Return(model.Todo{}, sql.ErrNoRows).Times(1)
			},
			expectedStatusCode: 404,
		},
		{
			name: "InternalError",
			id:   1,
			mockBehavior: func(store *mockdb.MockReminderRepo, id int) {
				store.EXPECT().GetRemindByID(gomock.Any(), gomock.Eq(id)).Return(model.Todo{}, sql.ErrConnDone).Times(1)
			},
			expectedStatusCode: 500,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := mockdb.NewMockReminderRepo(c)
			test.mockBehavior(store, test.id)

			server := newTestServer(store)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/remind", http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": "1"})

			handler := http.HandlerFunc(server.GetRemindByID)
			handler.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

func TestServer_UpdateRemind(t *testing.T) {
	testCases := []struct {
		name               string
		id                 int
		body               string
		mockBehavior       func(store *mockdb.MockReminderRepo, id int)
		expectedStatusCode int
	}{
		{
			name: "OK",
			id:   1,
			body: `{"description":"new test"}`,
			mockBehavior: func(store *mockdb.MockReminderRepo, id int) {
				store.EXPECT().UpdateRemind(gomock.Any(), gomock.Eq(id), model.TodoUpdateInput{
					Description: "new test",
				}).Return(nil).Times(1)
			},
			expectedStatusCode: 200,
		},
		{
			name:               "Error - wrong input",
			body:               `{"description":""}`,
			mockBehavior:       func(store *mockdb.MockReminderRepo, id int) {},
			expectedStatusCode: 422,
		},
		{
			name: "Error - Internal error",
			id:   1,
			body: `{"description":"new test"}`,
			mockBehavior: func(store *mockdb.MockReminderRepo, id int) {
				store.EXPECT().UpdateRemind(gomock.Any(), gomock.Eq(id), model.TodoUpdateInput{
					Description: "new test",
				}).Return(errors.New("something went wrong")).Times(1)
			},
			expectedStatusCode: 500,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := mockdb.NewMockReminderRepo(c)
			test.mockBehavior(store, test.id)

			server := newTestServer(store)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/remind", bytes.NewBufferString(test.body))
			req = mux.SetURLVars(req, map[string]string{"id": "1"})

			handler := http.HandlerFunc(server.UpdateRemind)
			handler.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

func TestServer_GetCurrentReminds(t *testing.T) {
	testCases := []struct {
		name               string
		params             pagination.Page
		userID             int
		mockBehavior       func(store *mockdb.MockReminderRepo, params pagination.Page, userID int)
		expectedStatusCode int
	}{
		{
			name:   "OK",
			params: pagination.Page{Limit: 5},
			userID: 1,
			mockBehavior: func(store *mockdb.MockReminderRepo, params pagination.Page, userID int) {
				store.EXPECT().GetNewReminds(gomock.Any(), params, userID).Return([]model.Todo{{
					ID:          1,
					Description: "test",
					CreatedAt:   time.Now(),
					DeadlineAt:  time.Now(),
					Completed:   false,
				}}, 1, nil).Times(1)
			},
			expectedStatusCode: 200,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := mockdb.NewMockReminderRepo(c)
			test.mockBehavior(store, test.params, test.userID)

			server := newTestServer(store)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/current", http.NoBody)

			ctx := req.Context()
			ctx = context.WithValue(ctx, "currentUser", model.User{
				ID: 1,
			})
			req = req.WithContext(ctx)

			handler := http.HandlerFunc(server.GetCurrentReminds)
			handler.ServeHTTP(w, req)

			// Add query parameters to request URL
			q := req.URL.Query()
			q.Add("limit", fmt.Sprintf("%d", test.params.Limit))
			q.Add("cursor", fmt.Sprintf("%d", test.params.Cursor))
			req.URL.RawQuery = q.Encode()

			require.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

func TestControllers_GetAllReminds(t *testing.T) {
	testCases := []struct {
		name               string
		params             pagination.Page
		userID             int
		mockBehavior       func(store *mockdb.MockReminderRepo, params pagination.Page, userID int)
		expectedStatusCode int
	}{
		{
			name:   "OK",
			params: pagination.Page{Limit: 5},
			userID: 1,
			mockBehavior: func(store *mockdb.MockReminderRepo, params pagination.Page, userID int) {
				store.EXPECT().GetAllReminds(gomock.Any(), params, userID).Return([]model.Todo{{
					ID:          1,
					Description: "test",
					CreatedAt:   time.Now(),
					DeadlineAt:  time.Now(),
					Completed:   false,
				}}, 1, nil).Times(1)
			},
			expectedStatusCode: 200,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := mockdb.NewMockReminderRepo(c)
			test.mockBehavior(store, test.params, test.userID)

			server := newTestServer(store)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/remind", http.NoBody)

			ctx := req.Context()
			ctx = context.WithValue(ctx, "currentUser", model.User{
				ID: 1,
			})
			req = req.WithContext(ctx)

			handler := http.HandlerFunc(server.GetAllReminds)
			handler.ServeHTTP(w, req)

			// Add query parameters to request URL
			q := req.URL.Query()
			q.Add("cursor", fmt.Sprintf("%d", test.params.Cursor))
			q.Add("limit", fmt.Sprintf("%d", test.params.Limit))
			req.URL.RawQuery = q.Encode()

			require.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

func Test_DeleteRemind(t *testing.T) {
	testCases := []struct {
		name           string
		id             int
		mockBehavior   func(store *mockdb.MockReminderRepo, id int)
		expectedStatus int
	}{
		{
			name: "OK",
			id:   1,
			mockBehavior: func(store *mockdb.MockReminderRepo, id int) {
				store.EXPECT().DeleteRemind(gomock.Any(), gomock.Eq(id)).Return(nil).Times(1)
			},
			expectedStatus: 201,
		},
		{
			name: "InternalError",
			id:   1,
			mockBehavior: func(store *mockdb.MockReminderRepo, id int) {
				store.EXPECT().DeleteRemind(gomock.Any(), gomock.Eq(id)).Return(sql.ErrConnDone).Times(1)
			},
			expectedStatus: 500,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := mockdb.NewMockReminderRepo(c)
			test.mockBehavior(store, test.id)

			server := newTestServer(store)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/remind", http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": "1"})

			handler := http.HandlerFunc(server.DeleteRemind)
			handler.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatus, w.Code)

		})
	}

}

func TestControllers_GetCompletedReminds(t *testing.T) {
	testCases := []struct {
		name               string
		params             storage.Params
		userID             int
		mockBehavior       func(store *mockdb.MockReminderRepo, params storage.Params, userID int)
		expectedStatusCode int
	}{
		{
			name: "OK",
			params: storage.Params{
				Page: pagination.Page{
					Limit: 5,
				},
			},
			userID: 1,
			mockBehavior: func(store *mockdb.MockReminderRepo, params storage.Params, userID int) {
				store.EXPECT().GetCompletedReminds(gomock.Any(), params, userID).Return([]model.Todo{{
					ID:          1,
					Description: "test",
					CreatedAt:   time.Now(),
					DeadlineAt:  time.Now(),
					Completed:   true,
				}}, 1, nil).Times(1)
			},
			expectedStatusCode: 200,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := mockdb.NewMockReminderRepo(c)
			test.mockBehavior(store, test.params, test.userID)

			server := newTestServer(store)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/completed", http.NoBody)

			ctx := req.Context()
			ctx = context.WithValue(ctx, "currentUser", model.User{
				ID: 1,
			})
			req = req.WithContext(ctx)
			handler := http.HandlerFunc(server.GetCompletedReminds)
			handler.ServeHTTP(w, req)

			// Add query parameters to request URL
			q := req.URL.Query()
			q.Add("limit", fmt.Sprintf("%d", test.params.Limit))
			q.Add("cursor", fmt.Sprintf("%d", test.params.Cursor))
			req.URL.RawQuery = q.Encode()

			require.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}
