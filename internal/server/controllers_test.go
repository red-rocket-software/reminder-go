package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/red-rocket-software/reminder-go/internal/app/model"
	mockdb "github.com/red-rocket-software/reminder-go/internal/storage/mocks"
	"github.com/stretchr/testify/require"
)

func TestControllers_AddRemind(t *testing.T) {
	dTime, err := time.Parse("2006-01-02", "2023-02-02")
	require.NoError(t, err)
	now, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
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
			body: `{"description":"test", "deadline_at": "2023-02-02"}`,
			inputTodo: model.Todo{
				CreatedAt:   now,
				Description: "test",
				DeadlineAt:  dTime,
			},
			mockBehavior: func(store *mockdb.MockReminderRepo, input model.Todo) {
				store.EXPECT().CreateRemind(gomock.Any(), input).Times(1).Return(0, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: "Remind is successfully created",
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

			handler := http.HandlerFunc(server.AddRemind)
			handler.ServeHTTP(w, req)

			require.Equal(t, w.Code, test.expectedStatusCode)
		})
	}
}
