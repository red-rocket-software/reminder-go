package server

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	mockdb "github.com/red-rocket-software/reminder-go/internal/reminder/domain/mocks"
	"github.com/red-rocket-software/reminder-go/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestControllers_AddRemind(t *testing.T) {
	dTime, err := time.Parse(time.RFC3339, "2023-04-15T16:27:00+02:00")
	require.NoError(t, err)
	now, err := time.Parse("02.01.2006, 15:04:05", "14.04.2023, 15:30:35")
	require.NoError(t, err)
	b := false

	testCases := []struct {
		name                 string
		body                 string
		inputTodo            domain.Todo
		mockBehavior         func(store *mockdb.MockTodoRepository, input domain.Todo)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			body: `{"description": "Test", "title": "Title", "user_id": "GxRlwVXMF0UAc15VwtkYJGWdKmj2", "deadline_at": "2023-04-15T16:27:00+02:00", "created_at": "14.04.2023, 15:30:35", "deadline_notify": false, "notify_period": []}`,
			inputTodo: domain.Todo{
				Description:    "Test",
				Title:          "Title",
				UserID:         "GxRlwVXMF0UAc15VwtkYJGWdKmj2",
				DeadlineAt:     dTime,
				CreatedAt:      now,
				DeadlineNotify: &b,
				NotifyPeriod:   []time.Time{},
			},
			mockBehavior: func(store *mockdb.MockTodoRepository, input domain.Todo) {
				store.EXPECT().CreateRemind(gomock.Any(), input).Return(domain.Todo{}, nil)
			},
			expectedStatusCode: 201,
		},
		{
			name:                 "Error - wrong input",
			body:                 `{"description":"", "user_id": "1", "deadline_at": "2023-02-02"}`,
			inputTodo:            domain.Todo{},
			mockBehavior:         func(store *mockdb.MockTodoRepository, input domain.Todo) {},
			expectedStatusCode:   422,
			expectedResponseBody: "nothing to save",
		},
		{
			name:                 "Error - notify period after deadline",
			body:                 `{"description": "Test", "title": "Title", "user_id": "GxRlwVXMF0UAc15VwtkYJGWdKmj2", "deadline_at": "2023-04-15T16:27:00+02:00", "created_at": "14.04.2023, 15:30:35", "deadline_notify": false, "notify_period": ["2023-05-15T16:27:00+02:00"]}`,
			inputTodo:            domain.Todo{},
			mockBehavior:         func(store *mockdb.MockTodoRepository, input domain.Todo) {},
			expectedStatusCode:   400,
			expectedResponseBody: "time to deadline notification can't be more than deadline time",
		},
		{
			name:                 "Error - notify period more than 2 days before deadline",
			body:                 `{"description": "Test", "title": "Title", "user_id": "GxRlwVXMF0UAc15VwtkYJGWdKmj2", "deadline_at": "2023-04-15T16:27:00+02:00", "created_at": "14.04.2023, 15:30:35", "deadline_notify": false, "notify_period": ["2023-04-12T16:27:00+02:00"]}`,
			inputTodo:            domain.Todo{},
			mockBehavior:         func(store *mockdb.MockTodoRepository, input domain.Todo) {},
			expectedStatusCode:   400,
			expectedResponseBody: "time to deadline notification can't be less than 2 days to deadline time",
		},
		{
			name:                 "Error - wrong deadline time format",
			body:                 `{"description": "Test", "title": "Title", "user_id": "GxRlwVXMF0UAc15VwtkYJGWdKmj2", "deadline_at": "2023-04-15", "created_at": "14.04.2023, 15:30:35", "deadline_notify": false, "notify_period": []}`,
			inputTodo:            domain.Todo{},
			mockBehavior:         func(store *mockdb.MockTodoRepository, input domain.Todo) {},
			expectedStatusCode:   400,
			expectedResponseBody: "cannot parse",
		},
		{
			name: "Error - Service error",
			body: `{"description": "Test", "title": "Title", "user_id": "GxRlwVXMF0UAc15VwtkYJGWdKmj2", "deadline_at": "2023-04-15T16:27:00+02:00", "created_at": "14.04.2023, 15:30:35", "deadline_notify": false, "notify_period": []}`,
			inputTodo: domain.Todo{
				Description:    "Test",
				Title:          "Title",
				UserID:         "GxRlwVXMF0UAc15VwtkYJGWdKmj2",
				DeadlineAt:     dTime,
				CreatedAt:      now,
				DeadlineNotify: &b,
				NotifyPeriod:   []time.Time{},
			},
			mockBehavior: func(store *mockdb.MockTodoRepository, input domain.Todo) {
				store.EXPECT().CreateRemind(gomock.Any(), input).Return(domain.Todo{}, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `"something went wrong"`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			configStore := mockdb.NewMockConfigRepository(c)
			todoStore := mockdb.NewMockTodoRepository(c)
			test.mockBehavior(todoStore, test.inputTodo)

			server := newTestServer(todoStore, configStore)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/remind", bytes.NewBufferString(test.body))
			ctx := req.Context()
			ctx = context.WithValue(ctx, "userID", "GxRlwVXMF0UAc15VwtkYJGWdKmj2")
			req = req.WithContext(ctx)

			handler := http.HandlerFunc(server.AddRemind)
			handler.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)
			if test.name != "OK" {
				require.Contains(t, w.Body.String(), test.expectedResponseBody)
			}
		})
	}
}

func TestControllers_GetRemindByID(t *testing.T) {
	testCases := []struct {
		name               string
		id                 int
		resTodo            domain.Todo
		mockBehavior       func(store *mockdb.MockTodoRepository, id int)
		expectedStatusCode int
	}{
		{
			name: "OK",
			id:   1,
			mockBehavior: func(store *mockdb.MockTodoRepository, id int) {
				store.EXPECT().GetRemindByID(gomock.Any(), gomock.Eq(1)).Return(domain.Todo{
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
			mockBehavior: func(store *mockdb.MockTodoRepository, id int) {
				store.EXPECT().GetRemindByID(gomock.Any(), gomock.Eq(id)).Return(domain.Todo{}, sql.ErrNoRows).Times(1)
			},
			expectedStatusCode: 404,
		},
		{
			name: "InternalError",
			id:   1,
			mockBehavior: func(store *mockdb.MockTodoRepository, id int) {
				store.EXPECT().GetRemindByID(gomock.Any(), gomock.Eq(id)).Return(domain.Todo{}, sql.ErrConnDone).Times(1)
			},
			expectedStatusCode: 500,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			configStore := mockdb.NewMockConfigRepository(c)
			todoStore := mockdb.NewMockTodoRepository(c)
			test.mockBehavior(todoStore, test.id)

			server := newTestServer(todoStore, configStore)

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
		mockBehavior       func(store *mockdb.MockTodoRepository, id int)
		expectedStatusCode int
	}{
		{
			name: "OK",
			id:   1,
			body: `{"description":"new test", "title":"new test"}`,
			mockBehavior: func(store *mockdb.MockTodoRepository, id int) {
				store.EXPECT().UpdateRemind(gomock.Any(), id, domain.TodoUpdateInput{
					Description: "new test",
					Title:       "new test",
				}).Return(domain.Todo{Description: "new test", Title: "new test"}, nil).Times(1)
			},
			expectedStatusCode: 200,
		},
		{
			name:               "Error - no description",
			body:               `{"description":"", "title":"title"}`,
			mockBehavior:       func(store *mockdb.MockTodoRepository, id int) {},
			expectedStatusCode: 422,
		},
		{
			name:               "Error - no title",
			body:               `{"description":"test", "title":""}`,
			mockBehavior:       func(store *mockdb.MockTodoRepository, id int) {},
			expectedStatusCode: 422,
		},
		{
			name: "Error - Internal error",
			id:   1,
			body: `{"description":"new test", "title":"new test"}`,
			mockBehavior: func(store *mockdb.MockTodoRepository, id int) {
				store.EXPECT().UpdateRemind(gomock.Any(), gomock.Eq(id), domain.TodoUpdateInput{
					Description: "new test",
					Title:       "new test",
				}).Return(domain.Todo{}, errors.New("something went wrong")).Times(1)
			},
			expectedStatusCode: 500,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			configStore := mockdb.NewMockConfigRepository(c)
			todoStore := mockdb.NewMockTodoRepository(c)
			test.mockBehavior(todoStore, test.id)

			server := newTestServer(todoStore, configStore)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/remind", bytes.NewBufferString(test.body))
			req = mux.SetURLVars(req, map[string]string{"id": "1"})

			handler := http.HandlerFunc(server.UpdateRemind)
			handler.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

func TestControllers_GetReminds(t *testing.T) {
	testCases := []struct {
		name               string
		params             domain.FetchParams
		userID             string
		mockBehavior       func(store *mockdb.MockTodoRepository, params domain.FetchParams, userID string)
		expectedStatusCode int
	}{
		{
			name: "OK",
			params: domain.FetchParams{
				Page: utils.Page{
					Cursor: 0,
					Limit:  10,
				},
				FilterByDate:  "createdAt",
				FilterBySort:  "ASC",
				FilterByQuery: "all",
			},
			userID: "rrdZH9ERxueDxj2m1e1T2vIQKBP2",
			mockBehavior: func(store *mockdb.MockTodoRepository, params domain.FetchParams, userID string) {
				store.EXPECT().GetReminds(context.Background(), params, userID).Return([]domain.Todo{{
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
		{
			name: "Error wrong filter",
			params: domain.FetchParams{
				Page: utils.Page{
					Cursor: 0,
					Limit:  10,
				},
				FilterByDate:  "",
				FilterBySort:  "ASC",
				FilterByQuery: "all",
			},
			userID:             "rrdZH9ERxueDxj2m1e1T2vIQKBP2",
			mockBehavior:       func(store *mockdb.MockTodoRepository, params domain.FetchParams, userID string) {},
			expectedStatusCode: 400,
		},
		{
			name: "Error wrong filter params",
			params: domain.FetchParams{
				Page: utils.Page{
					Cursor: 0,
					Limit:  10,
				},
				FilterByDate:  "CratedAt",
				FilterBySort:  "ASC",
				FilterByQuery: "",
			},
			userID:             "rrdZH9ERxueDxj2m1e1T2vIQKBP2",
			mockBehavior:       func(store *mockdb.MockTodoRepository, params domain.FetchParams, userID string) {},
			expectedStatusCode: 400,
		},
		{
			name: "Error wrong filter options",
			params: domain.FetchParams{
				Page: utils.Page{
					Cursor: 0,
					Limit:  10,
				},
				FilterByDate:  "CratedAt",
				FilterBySort:  "",
				FilterByQuery: "all",
			},
			userID:             "rrdZH9ERxueDxj2m1e1T2vIQKBP2",
			mockBehavior:       func(store *mockdb.MockTodoRepository, params domain.FetchParams, userID string) {},
			expectedStatusCode: 400,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			configStore := mockdb.NewMockConfigRepository(c)
			todoStore := mockdb.NewMockTodoRepository(c)
			test.mockBehavior(todoStore, test.params, test.userID)

			server := newTestServer(todoStore, configStore)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/reminds", http.NoBody)

			ctx := req.Context()
			ctx = context.WithValue(ctx, "userID", "rrdZH9ERxueDxj2m1e1T2vIQKBP2")
			req = req.WithContext(ctx)

			// Add query parameters to request URL
			q := req.URL.Query()
			q.Add("cursor", fmt.Sprintf("%d", test.params.Cursor))
			q.Add("limit", fmt.Sprintf("%d", test.params.Limit))
			q.Add("filter", fmt.Sprintf("%s", test.params.FilterByDate))
			q.Add("filterOption", fmt.Sprintf("%s", test.params.FilterBySort))
			q.Add("filterParams", fmt.Sprintf("%s", test.params.FilterByQuery))
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
		mockBehavior   func(store *mockdb.MockTodoRepository, id int)
		expectedStatus int
	}{
		{
			name: "OK",
			id:   1,
			mockBehavior: func(store *mockdb.MockTodoRepository, id int) {
				store.EXPECT().DeleteRemind(gomock.Any(), gomock.Eq(id)).Return(nil).Times(1)
			},
			expectedStatus: 204,
		},
		{
			name: "InternalError",
			id:   1,
			mockBehavior: func(store *mockdb.MockTodoRepository, id int) {
				store.EXPECT().DeleteRemind(gomock.Any(), gomock.Eq(id)).Return(sql.ErrConnDone).Times(1)
			},
			expectedStatus: 500,
		},
		{
			name: "Error remind doesn't exist",
			id:   1,
			mockBehavior: func(store *mockdb.MockTodoRepository, id int) {
				store.EXPECT().DeleteRemind(gomock.Any(), gomock.Eq(id)).Return(domain.ErrCantFindRemindWithID).Times(1)
			},
			expectedStatus: 404,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			configStore := mockdb.NewMockConfigRepository(c)
			todoStore := mockdb.NewMockTodoRepository(c)
			test.mockBehavior(todoStore, test.id)

			server := newTestServer(todoStore, configStore)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/remind", http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": "1"})

			handler := http.HandlerFunc(server.DeleteRemind)
			handler.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatus, w.Code)

		})
	}

}

func TestServer_UpdateUserConfig(t *testing.T) {
	testCases := []struct {
		name               string
		id                 string
		body               string
		mockBehavior       func(store *mockdb.MockConfigRepository, id string)
		expectedStatusCode int
	}{
		{
			name: "OK",
			id:   "rrdZH9ERxueDxj2m1e1T2vIQKBP2",
			body: `{"notification": true, "period": 1}`,
			mockBehavior: func(store *mockdb.MockConfigRepository, id string) {
				store.EXPECT().UpdateUserConfig(gomock.Any(), gomock.Eq(id), domain.UserConfigs{
					Notification: true,
					Period:       1,
				}).Return(nil).Times(1)
			},
			expectedStatusCode: 200,
		},
		{
			name: "Error - internal error",
			id:   "rrdZH9ERxueDxj2m1e1T2vIQKBP2",
			body: `{"notification": true, "period": 1}`,
			mockBehavior: func(store *mockdb.MockConfigRepository, id string) {
				store.EXPECT().UpdateUserConfig(gomock.Any(), gomock.Eq(id), domain.UserConfigs{
					Notification: true,
					Period:       1,
				}).Return(errors.New("something went wrong")).Times(1)
			},
			expectedStatusCode: 500,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			configStore := mockdb.NewMockConfigRepository(c)
			todoStore := mockdb.NewMockTodoRepository(c)
			test.mockBehavior(configStore, test.id)

			server := newTestServer(todoStore, configStore)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/configs", bytes.NewBufferString(test.body))
			req = mux.SetURLVars(req, map[string]string{"id": "rrdZH9ERxueDxj2m1e1T2vIQKBP2"})

			handler := http.HandlerFunc(server.UpdateUserConfig)
			handler.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

// func TestServer_AuthMiddleware(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		token          string
// 		mockBehavior   func(store *mock_firestore.MockClient, token string)
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{
// 			name:  "valid token",
// 			token: "Bearer valid_token",
// 			mockBehavior: func(store *mock_firestore.MockClient, token string) {
// 				store.EXPECT().VerifyIDToken("valid_token").Return(&auth.Token{
// 					UID: "user123",
// 					Claims: map[string]interface{}{
// 						"user_id": "user123",
// 					},
// 				}, nil)
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   "OK",
// 		},
// 		{
// 			name:           "no authorization header",
// 			token:          "",
// 			mockBehavior:   func(store *mock_firestore.MockClient, token string) {},
// 			expectedStatus: http.StatusUnauthorized,
// 			expectedBody:   "you are not logged in",
// 		},
// 		{
// 			name:           "invalid authorization header",
// 			token:          "invalid_header",
// 			mockBehavior:   func(store *mock_firestore.MockClient, token string) {},
// 			expectedStatus: http.StatusUnauthorized,
// 			expectedBody:   "you are not logged in",
// 		},
// 		{
// 			name:  "invalid token",
// 			token: "Bearer invalid_token",
// 			mockBehavior: func(store *mock_firestore.MockClient, token string) {
// 				store.EXPECT().VerifyIDToken("invalid_token").Return(nil, errors.New("invalid token"))
// 			},
// 			expectedStatus: http.StatusUnauthorized,
// 			expectedBody:   "error verify token",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			req, err := http.NewRequest(http.MethodGet, "/", nil)
// 			if err != nil {
// 				t.Fatalf("failed to create request: %v", err)
// 			}

// 			if tt.token != "" {
// 				req.Header.Set("Authorization", tt.token)
// 			}

// 			rec := httptest.NewRecorder()

// 			client := mock_firestore.NewMockClient(c)
// 			tt.mockBehavior(client, tt.token)

// 			server := &Server{
// 				FireClient: client,
// 			}

// 			handler := server.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 				w.WriteHeader(http.StatusOK)
// 				w.Write([]byte("OK"))
// 			}))

// 			handler.ServeHTTP(rec, req)

// 			if rec.Code != tt.expectedStatus {
// 				t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, tt.expectedStatus)
// 			}

// 			if !strings.Contains(rec.Body.String(), tt.expectedBody) {
// 				t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), tt.expectedBody)
// 			}
// 		})
// 	}
// }

func Test_UpdateCompleteStatus(t *testing.T) {
	tn := time.Now().Truncate(1 * time.Second)

	testCases := []struct {
		name               string
		id                 int
		body               string
		mockBehavior       func(store *mockdb.MockTodoRepository, id int)
		expectedStatusCode int
	}{
		{
			name: "OK",
			id:   1,
			body: `{"completed": true}`,
			mockBehavior: func(store *mockdb.MockTodoRepository, id int) {
				store.EXPECT().UpdateStatus(gomock.Any(), gomock.Eq(id), domain.TodoUpdateStatusInput{
					Completed:  true,
					FinishedAt: &tn,
				}).Return(nil).Times(1)
			},
			expectedStatusCode: 200,
		},
		{
			name:               "Error - Missed Body",
			id:                 1,
			body:               "",
			mockBehavior:       func(store *mockdb.MockTodoRepository, id int) {},
			expectedStatusCode: 422,
		},
		{
			name: "Error - Internal error",
			id:   1,
			body: `{"completed": true}`,
			mockBehavior: func(store *mockdb.MockTodoRepository, id int) {
				store.EXPECT().UpdateStatus(gomock.Any(), gomock.Eq(id), domain.TodoUpdateStatusInput{
					Completed:  true,
					FinishedAt: &tn,
				}).Return(errors.New("remind not found")).Times(1)
			},
			expectedStatusCode: 500,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			configStore := mockdb.NewMockConfigRepository(c)
			todoStore := mockdb.NewMockTodoRepository(c)
			test.mockBehavior(todoStore, test.id)

			server := newTestServer(todoStore, configStore)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/status", bytes.NewBufferString(test.body))

			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(test.id)})

			handler := http.HandlerFunc(server.UpdateCompleteStatus)
			handler.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

func TestServer_GetOrCreateUserConfig(t *testing.T) {
	tn := time.Now()

	testCases := []struct {
		name               string
		id                 string
		mockBehavior       func(store *mockdb.MockConfigRepository, id string)
		expectedStatusCode int
	}{
		{
			name: "OK",
			id:   "rrdZH9ERxueDxj2m1e1T2vIQKBP2",
			mockBehavior: func(store *mockdb.MockConfigRepository, id string) {
				store.EXPECT().GetUserConfigs(gomock.Any(), gomock.Eq(id)).Return(domain.UserConfigs{}, nil)
				store.EXPECT().CreateUserConfigs(gomock.Any(), gomock.Eq(id)).Return(domain.UserConfigs{
					ID:           id,
					Notification: false,
					Period:       2,
					CreatedAt:    tn,
				}, nil)
			},
			expectedStatusCode: 200,
		},
		{
			name: "Error - Wrong uID",
			id:   "",
			mockBehavior: func(store *mockdb.MockConfigRepository, id string) {
			},
			expectedStatusCode: 500,
		},
		{
			name: "Error - GetUserConfigs Error",
			id:   "rrdZH9ERxueDxj2m1e1T2vIQKBP2",
			mockBehavior: func(store *mockdb.MockConfigRepository, id string) {
				store.EXPECT().GetUserConfigs(gomock.Any(), gomock.Eq(id)).Return(domain.UserConfigs{}, errors.New("cannot get user-configs from database"))
			},
			expectedStatusCode: 500,
		},
		{
			name: "Error - CreateUserConfigs Error",
			id:   "rrdZH9ERxueDxj2m1e1T2vIQKBP2",
			mockBehavior: func(store *mockdb.MockConfigRepository, id string) {
				store.EXPECT().GetUserConfigs(gomock.Any(), gomock.Eq(id)).Return(domain.UserConfigs{}, nil)
				store.EXPECT().CreateUserConfigs(gomock.Any(), gomock.Eq(id)).Return(domain.UserConfigs{}, errors.New("Error create userConfigs"))
			},
			expectedStatusCode: 500,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(*testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			configStore := mockdb.NewMockConfigRepository(c)
			todoStore := mockdb.NewMockTodoRepository(c)
			test.mockBehavior(configStore, test.id)

			server := newTestServer(todoStore, configStore)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/configs", http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": test.id})

			handler := http.HandlerFunc(server.GetOrCreateUserConfig)
			handler.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}
