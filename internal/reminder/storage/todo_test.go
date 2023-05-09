package storage

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"

	model "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestStorageTodo_CreateRemind(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	userID, err := testStorage.SeedUserConfig()
	require.NoError(t, err)

	date := time.Date(2023, time.April, 1, 1, 0, 0, 0, time.UTC)

	type args struct {
		ctx  context.Context
		todo model.Todo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{
			context.Background(),
			model.Todo{
				Description: "test text",
				Title:       "test text",
				UserID:      userID,
				DeadlineAt:  date,
				CreatedAt:   time.Now(),
			},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := testStorage.CreateRemind(tt.args.ctx, tt.args.todo)
			require.NoError(t, err)
			require.NotZero(t, id)
		})
	}
}

func TestStorageTodo_GetRemindByID(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	date := time.Date(2023, time.April, 1, 1, 0, 0, 0, time.UTC)

	userID, err := testStorage.SeedUserConfig()
	require.NoError(t, err)

	insertTodo := model.Todo{
		Description: "test",
		Title:       "test",
		UserID:      userID,
		DeadlineAt:  date,
		CreatedAt:   time.Now(),
	}

	remind, _ := testStorage.CreateRemind(context.Background(), insertTodo)

	got, err := testStorage.GetRemindByID(context.Background(), remind.ID)

	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.Equal(t, insertTodo.Description, got.Description)
	require.Equal(t, insertTodo.DeadlineAt, got.DeadlineAt)

}

func TestStorageTodo_GetReminds(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	expectedTodo, err := testStorage.SeedTodos()
	if err != nil {
		log.Fatal("error seed reminds")
	}

	type args struct {
		ctx         context.Context
		fetchParams FetchParams
		userID      string
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Todo
		want1   int
		want2   int
		wantErr bool
	}{
		{name: "success get all reminds", args: args{context.Background(), FetchParams{
			Page: utils.Page{
				Cursor: 0,
				Limit:  10,
			},
			FilterByDate:  "CreatedAt",
			FilterBySort:  "DESC",
			FilterByQuery: "all",
		}, expectedTodo[0].UserID},
			want:    expectedTodo,
			want1:   5,
			want2:   expectedTodo[len(expectedTodo)-1].ID,
			wantErr: false},
		{name: "success get completed", args: args{context.Background(), FetchParams{
			Page: utils.Page{
				Cursor: 0,
				Limit:  10,
			},
			FilterByDate:  "CreatedAt",
			FilterBySort:  "DESC",
			FilterByQuery: "completed",
		}, expectedTodo[0].UserID},
			want:    []model.Todo{expectedTodo[1]},
			want1:   1,
			want2:   expectedTodo[len(expectedTodo)-4].ID,
			wantErr: false},
		{name: "success get current", args: args{context.Background(), FetchParams{
			Page: utils.Page{
				Cursor: 0,
				Limit:  10,
			},
			FilterByDate:  "CreatedAt",
			FilterBySort:  "DESC",
			FilterByQuery: "current",
		}, expectedTodo[0].UserID},
			want:    []model.Todo{expectedTodo[0], expectedTodo[2], expectedTodo[3], expectedTodo[4]},
			want1:   4,
			want2:   expectedTodo[len(expectedTodo)-1].ID,
			wantErr: false},
		{name: "empty filterParams value", args: args{context.Background(), FetchParams{
			Page: utils.Page{
				Cursor: 0,
				Limit:  10,
			},
			FilterByDate:  "CreatedAt",
			FilterBySort:  "DESC",
			FilterByQuery: "",
		}, expectedTodo[0].UserID},
			want:    nil,
			want1:   0,
			wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := testStorage.GetReminds(tt.args.ctx, tt.args.fetchParams, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllReminds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAReminds() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetAReminds() got1 = %v, want1 %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("GetAReminds() got2 = %v, want2 %v", got2, tt.want2)
			}
		})
	}
}

func TestStorageTodo_DeleteRemind(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	expectedTodo, err := testStorage.SeedTodos()
	if err != nil {
		log.Fatal("error seed reminds")
	}
	err = testStorage.DeleteRemind(context.Background(), expectedTodo[1].ID)
	require.NoError(t, err)

	todo, _ := testStorage.GetRemindByID(context.Background(), expectedTodo[1].ID)
	require.Empty(t, todo)
}

func TestStorageTodo_UpdateRemind(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	expectedTodo, err := testStorage.SeedTodos()
	if err != nil {
		log.Fatal("error seed reminds")
	}

	tn := time.Now()
	updateInput := model.TodoUpdateInput{
		Description: "New text",
		FinishedAt:  &tn,
		Completed:   true,
		DeadlineAt:  "2023-01-26T17:05:00Z",
	}

	_, err = testStorage.UpdateRemind(context.Background(), expectedTodo[1].ID, updateInput)
	require.NoError(t, err)

	newTodo, _ := testStorage.GetRemindByID(context.Background(), expectedTodo[1].ID)
	require.Equal(t, updateInput.Description, newTodo.Description)
	require.Equal(t, updateInput.Completed, newTodo.Completed)
}

func TestStorageTodo_SeedTodos(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	todos, err := testStorage.SeedTodos()

	require.NoError(t, err)
	require.Equal(t, len(todos), 5)
}

func TestStorage_GetRemindsForNotification(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	_, err := testStorage.SeedTodos()
	if err != nil {
		log.Fatal("error seed reminds")
	}
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		reminds, err := testStorage.GetRemindsForNotification(context.Background())
		require.NoError(t, err)
		require.Equal(t, 4, len(reminds))

		tn := time.Now().Truncate(1 * time.Second).UTC()

		for _, remind := range reminds {
			userID := remind.UserID
			userConfig, err := testStorage.GetUserConfigs(context.Background(), userID)
			if err != nil {
				log.Fatal("error to get userConfig")
			}
			tfromPeriod := time.Now().AddDate(0, 0, userConfig.Period)
			expr := remind.DeadlineAt.After(tn) && remind.DeadlineAt.Before(tfromPeriod)
			require.Equal(t, true, expr)
		}

		err = testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
		require.NoError(t, err)
	})

	t.Run("no reminds for Notification at this moment", func(t *testing.T) {
		reminds, err := testStorage.GetRemindsForNotification(context.Background())
		require.NoError(t, err)
		require.Empty(t, reminds)
	})
}

func TestStorage_GetRemindsForDeadlineNotification(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	expectedReminds, err := testStorage.SeedTodos()
	if err != nil {
		log.Fatal("error seed reminds")
	}
	require.NoError(t, err)

	tn := time.Now().Truncate(time.Minute).Format(time.RFC3339)

	t.Run("success", func(t *testing.T) {
		reminds, timeNow, err := testStorage.GetRemindsForDeadlineNotification(context.Background())
		require.NoError(t, err)
		require.Equal(t, 1, len(reminds))
		require.Equal(t, tn, timeNow)
		require.Equal(t, expectedReminds[0].UserID, reminds[0].UserID)

		err = testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
		require.NoError(t, err)
	})

	t.Run("no reminds for DeadlineNotification at this moment", func(t *testing.T) {
		reminds, timeNow, err := testStorage.GetRemindsForDeadlineNotification(context.Background())
		require.Empty(t, reminds)
		require.Equal(t, tn, timeNow)
		require.Empty(t, err)
	})
}

func TestStorage_UpdateNotifyPeriod(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	expectedTodos, err := testStorage.SeedTodos()
	if err != nil {
		log.Fatal("error seed reminds")
	}
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		err = testStorage.UpdateNotifyPeriod(context.Background(), expectedTodos[0].ID, (expectedTodos[0].NotifyPeriod[0]).Format("2006-01-02 15:04:05"))
		require.NoError(t, err)
	})

	t.Run("remind not found", func(t *testing.T) {
		err = testStorage.UpdateNotifyPeriod(context.Background(), 0, (expectedTodos[0].NotifyPeriod[0]).Format("2006-01-02 15:04:05"))
		require.Error(t, err)
	})
}

func TestStorageTodo_UpdateUserConfig(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	expectedUserID, err := testStorage.SeedUserConfig()
	if err != nil {
		log.Fatal("error truncate config")
	}
	require.NoError(t, err)
	tn := time.Now()

	updateConfigInput := model.UserConfigs{
		ID:           expectedUserID,
		Notification: true,
		Period:       1,
		UpdatedAt:    &tn,
	}

	t.Run("success", func(t *testing.T) {
		err = testStorage.UpdateUserConfig(context.Background(), expectedUserID, updateConfigInput)
		require.NoError(t, err)
	})
	t.Run("user configs not found", func(t *testing.T) {
		err = testStorage.UpdateUserConfig(context.Background(), "0", updateConfigInput)
		require.Error(t, err)
	})

}

func TestStorage_CreateUserConfigs(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	expectedUserConfig := model.UserConfigs{
		ID:           "1",
		Notification: false,
		Period:       2,
	}

	t.Run("success", func(t *testing.T) {
		got, err := testStorage.CreateUserConfigs(context.Background(), expectedUserConfig.ID)
		require.NoError(t, err)
		require.Equal(t, got.ID, expectedUserConfig.ID)
		require.Equal(t, got.Notification, expectedUserConfig.Notification)
		require.Equal(t, got.Period, expectedUserConfig.Period)
	})
}

func TestStorage_GetUserConfigs(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	expectedUserID, err := testStorage.SeedUserConfig()
	if err != nil {
		log.Fatal("error truncate config")
	}
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		config, err := testStorage.GetUserConfigs(context.Background(), expectedUserID)
		require.NoError(t, err)
		require.Equal(t, config.ID, expectedUserID)
		require.Equal(t, config.Notification, true)
		require.Equal(t, config.Period, 2)
	})
	t.Run("no rows in result set", func(t *testing.T) {
		config, err := testStorage.GetUserConfigs(context.Background(), "0")
		require.Empty(t, err)
		require.Empty(t, config)
	})
}
