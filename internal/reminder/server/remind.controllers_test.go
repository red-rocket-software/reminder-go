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
	model "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/internal/reminder/storage"
	mockdb "github.com/red-rocket-software/reminder-go/internal/reminder/storage/mocks"
	"github.com/red-rocket-software/reminder-go/pkg/pagination"
	"github.com/stretchr/testify/require"
)

func TestControllers_AddRemind(t *testing.T) {
	//dTime, err := time.Parse(time.RFC3339, "2023-03-21T16:22:00+02:00")
	//require.NoError(t, err)
	//now, err := time.Parse("02.01.2006, 15:04:05", "19.01.2023, 22:15:30")
	//require.NoError(t, err)

	testCases := []struct {
		name                 string
		body                 string
		inputTodo            model.Todo
		mockBehavior         func(store *mockdb.MockReminderRepo, input model.Todo)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		//{
		//	name: "OK",
		//	body: `{"description":"test", "user_id": "1", "deadline_at": "2023-03-21T16:22:00+02:00", "created_at": "19.01.2023, 22:15:30"}`,
		//	inputTodo: domain.Todo{
		//		CreatedAt:   now,
		//		UserID:      1,
		//		Description: "test",
		//		DeadlineAt:  dTime,
		//	},
		//	mockBehavior: func(store *mockdb.MockReminderRepo, input domain.Todo) {
		//		store.EXPECT().CreateRemind(gomock.Any(), input).Return(0, nil)
		//	},
		//	expectedStatusCode:   201,
		//	expectedResponseBody: "Remind is successfully created",
		//},
		{
			name:                 "Error - wrong input",
			body:                 `{"description":"", "user_id": "1", "deadline_at": "2023-02-02"}`,
			inputTodo:            model.Todo{},
			mockBehavior:         func(store *mockdb.MockReminderRepo, input model.Todo) {},
			expectedStatusCode:   422,
			expectedResponseBody: "nothing to save",
		},
		//{
		//	name: "Error - Service error",
		//	body: `{"description":"test", "user_id": "1", "deadline_at": "2023-03-21T16:22:00+02:00", "created_at": "19.01.2023, 22:15:30"}`,
		//	inputTodo: domain.Todo{
		//		CreatedAt:   now,
		//		UserID:      1,
		//		Description: "test",
		//		DeadlineAt:  dTime,
		//	},
		//	mockBehavior: func(store *mockdb.MockReminderRepo, input domain.Todo) {
		//		store.EXPECT().CreateRemind(gomock.Any(), input).Return(0, errors.New("something went wrong"))
		//	},
		//	expectedStatusCode:   500,
		//	expectedResponseBody: `"something went wrong"`,
		//},
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
			ctx = context.WithValue(ctx, "uerID", "rrdZH9ERxueDxj2m1e1T2vIQKBP2")
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
			body: `{"description":"new test", "title":"new test"}`,
			mockBehavior: func(store *mockdb.MockReminderRepo, id int) {
				store.EXPECT().UpdateRemind(gomock.Any(), gomock.Eq(id), model.TodoUpdateInput{
					Description: "new test",
					Title:       "new test",
				}).Return(model.Todo{Description: "new test", Title: "new test"}, nil).Times(1)
			},
			expectedStatusCode: 200,
		},
		{
			name:               "Error - wrong input",
			body:               `{"description":"", "title":""}`,
			mockBehavior:       func(store *mockdb.MockReminderRepo, id int) {},
			expectedStatusCode: 422,
		},
		{
			name: "Error - Internal error",
			id:   1,
			body: `{"description":"new test", "title":"new test"}`,
			mockBehavior: func(store *mockdb.MockReminderRepo, id int) {
				store.EXPECT().UpdateRemind(gomock.Any(), gomock.Eq(id), model.TodoUpdateInput{
					Description: "new test",
					Title:       "new test",
				}).Return(model.Todo{}, errors.New("something went wrong")).Times(1)
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

func TestControllers_GetAllReminds(t *testing.T) {
	testCases := []struct {
		name               string
		params             storage.FetchParams
		userID             string
		mockBehavior       func(store *mockdb.MockReminderRepo, params storage.FetchParams, userID string)
		expectedStatusCode int
	}{
		{
			name: "OK",
			params: storage.FetchParams{
				Page: pagination.Page{
					Cursor: 0,
					Limit:  10,
				},
				Filter:       "createdAt",
				FilterOption: "ASC",
				FilterParam:  "all",
			},
			userID: "rrdZH9ERxueDxj2m1e1T2vIQKBP2",
			mockBehavior: func(store *mockdb.MockReminderRepo, params storage.FetchParams, userID string) {
				store.EXPECT().GetReminds(context.Background(), params, userID).Return([]model.Todo{{
					ID:          1,
					Title:       "test",
					Description: "test",
					CreatedAt:   time.Now(),
					DeadlineAt:  time.Now(),
					Completed:   false,
				}}, 1, 1, nil).Times(1)
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
			req, _ := http.NewRequest(http.MethodGet, "/reminds", http.NoBody)

			ctx := req.Context()
			ctx = context.WithValue(ctx, "userID", "rrdZH9ERxueDxj2m1e1T2vIQKBP2")
			req = req.WithContext(ctx)

			// Add query parameters to request URL
			q := req.URL.Query()
			q.Add("cursor", fmt.Sprintf("%d", test.params.Cursor))
			q.Add("limit", fmt.Sprintf("%d", test.params.Limit))
			q.Add("filter", fmt.Sprintf("%s", test.params.Filter))
			q.Add("filterOption", fmt.Sprintf("%s", test.params.FilterOption))
			q.Add("filterParams", fmt.Sprintf("%s", test.params.FilterParam))
			req.URL.RawQuery = q.Encode()

			handler := http.HandlerFunc(server.GetReminds)
			handler.ServeHTTP(w, req)

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
			expectedStatus: 204,
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
