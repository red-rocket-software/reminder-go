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
		{name: "error invalid user ID", args: args{
			context.Background(),
			model.Todo{
				Description: "test text",
				Title:       "test text",
				UserID:      "o", // an invalid user ID
				DeadlineAt:  date,
				CreatedAt:   time.Now(),
			},
		}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := testStorage.CreateRemind(tt.args.ctx, tt.args.todo)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
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
		{name: "success get all reminds with cursor", args: args{context.Background(), FetchParams{
			Page: utils.Page{
				Cursor: expectedTodo[3].ID,
				Limit:  2,
			},
			FilterByDate:  "CreatedAt",
			FilterBySort:  "DESC",
			FilterByQuery: "all",
		}, expectedTodo[0].UserID},
			want:    []model.Todo{expectedTodo[1], expectedTodo[0]},
			want1:   5,
			want2:   expectedTodo[len(expectedTodo)-5].ID,
			wantErr: false},
		{name: "success get completed with time range ", args: args{ctx: context.Background(), fetchParams: FetchParams{
			Page: utils.Page{
				Cursor: 0,
				Limit:  10,
			},
			TimeRangeFilter: TimeRangeFilter{
				StartRange: "2023-03-01 11:50:34",
				EndRange:   "2023-04-01 11:50:34",
			},
			FilterByDate:  "CreatedAt",
			FilterBySort:  "DESC",
			FilterByQuery: "completed",
		}, userID: expectedTodo[0].UserID},
			want:    []model.Todo{expectedTodo[1]},
			want1:   1,
			want2:   expectedTodo[len(expectedTodo)-4].ID,
			wantErr: false},
		{name: "success get completed not existing in time range", args: args{ctx: context.Background(), fetchParams: FetchParams{
			Page: utils.Page{
				Cursor: 0,
				Limit:  10,
			},
			TimeRangeFilter: TimeRangeFilter{
				StartRange: "2023-05-28 11:50:34",
				EndRange:   "2023-06-12 11:50:34",
			},
			FilterByDate:  "CreatedAt",
			FilterBySort:  "DESC",
			FilterByQuery: "completed",
		}, userID: expectedTodo[0].UserID},
			want:    []model.Todo{},
			want1:   0,
			want2:   0,
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
		require.NoError(t, err)
	}()

	expectedTodo, err := testStorage.SeedTodos()
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		err = testStorage.DeleteRemind(context.Background(), expectedTodo[1].ID)
		require.NoError(t, err)

		todo, _ := testStorage.GetRemindByID(context.Background(), expectedTodo[1].ID)
		require.Empty(t, todo)
	})
	t.Run("error remind doesn't existing", func(t *testing.T) {
		err = testStorage.DeleteRemind(context.Background(), 99)
		require.Error(t, err)
	})
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

	t.Run("success", func(t *testing.T) {
		updateInput := model.TodoUpdateInput{
			Description:  "New text",
			FinishedAt:   &tn,
			Completed:    true,
			DeadlineAt:   "2023-01-26T17:05:00Z",
			NotifyPeriod: []string{"2023-01-26T16:05:00Z"},
		}

		_, err = testStorage.UpdateRemind(context.Background(), expectedTodo[1].ID, updateInput)
		require.NoError(t, err)

		newTodo, _ := testStorage.GetRemindByID(context.Background(), expectedTodo[1].ID)
		require.Equal(t, updateInput.Description, newTodo.Description)
		require.Equal(t, updateInput.Completed, newTodo.Completed)
	})
	t.Run("error wrong notify period", func(t *testing.T) {
		updateInput := model.TodoUpdateInput{
			Description:  "New text",
			FinishedAt:   &tn,
			Completed:    true,
			DeadlineAt:   "2023-01-26T17:05:00Z",
			NotifyPeriod: []string{"2023-01-26"},
		}

		_, err = testStorage.UpdateRemind(context.Background(), expectedTodo[1].ID, updateInput)
		require.Error(t, err)
	})
	t.Run("error wrong deadlineAt ", func(t *testing.T) {
		updateInput := model.TodoUpdateInput{
			Description: "New text",
			FinishedAt:  &tn,
			Completed:   true,
			DeadlineAt:  "2023-01-26",
		}

		_, err = testStorage.UpdateRemind(context.Background(), expectedTodo[1].ID, updateInput)
		require.Error(t, err)
	})
	t.Run("error not existing remind ", func(t *testing.T) {
		updateInput := model.TodoUpdateInput{
			Description: "New text",
			FinishedAt:  &tn,
			Completed:   true,
			DeadlineAt:  "2023-01-26T17:05:00Z",
		}

		_, err = testStorage.UpdateRemind(context.Background(), 9999, updateInput)
		require.Error(t, err)
	})
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

func TestTodoStorage_UpdateNotification(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		require.NoError(t, err)
	}()

	expectedTodo, err := testStorage.SeedTodos()
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		id  int
		dao model.NotificationDAO
	}
	ctx := context.Background()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{
			ctx: ctx,
			id:  expectedTodo[0].ID,
			dao: model.NotificationDAO{Notificated: true},
		}, wantErr: false,
		},
		{name: "error doesn't existing remind", args: args{
			ctx: ctx,
			id:  9999,
			dao: model.NotificationDAO{Notificated: true},
		}, wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err = testStorage.UpdateNotification(tt.args.ctx, tt.args.id, tt.args.dao)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			remind, err := testStorage.GetRemindByID(ctx, expectedTodo[0].ID)
			require.NoError(t, err)

			require.Equal(t, tt.args.dao.Notificated, remind.Notificated)
		})
	}
}

func TestTodoStorage_UpdateStatus(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		require.NoError(t, err)
	}()

	expectedTodo, err := testStorage.SeedTodos()
	require.NoError(t, err)

	tn := time.Now().UTC()

	type args struct {
		ctx context.Context
		id  int
		dao model.TodoUpdateStatusInput
	}
	ctx := context.Background()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{
			ctx: ctx,
			id:  expectedTodo[0].ID,
			dao: model.TodoUpdateStatusInput{
				Completed:  true,
				FinishedAt: &tn,
			},
		}, wantErr: false,
		},
		{name: "error doesn't existing remind", args: args{
			ctx: ctx,
			id:  9999,
			dao: model.TodoUpdateStatusInput{
				Completed:  true,
				FinishedAt: &tn,
			},
		}, wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err = testStorage.UpdateStatus(tt.args.ctx, tt.args.id, tt.args.dao)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			remind, err := testStorage.GetRemindByID(ctx, expectedTodo[0].ID)
			require.NoError(t, err)

			require.Equal(t, tt.args.dao.FinishedAt.Truncate(time.Millisecond), remind.FinishedAt.Truncate(time.Millisecond))
			require.Equal(t, tt.args.dao.Completed, remind.Completed)
		})
	}
}

func TestTodoStorage_SeedUserConfig(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		require.NoError(t, err)
	}()

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "success",
			userID:  "rrdZH9ERxueDxj2m1e1T2vIQKBP2",
			wantErr: false,
		},
		{
			name:    "error",
			userID:  "wrongID",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testStorage.SeedUserConfig()
			require.NoError(t, err)

			if !tt.wantErr {
				require.Equal(t, tt.userID, got)
			}
		})
	}
}
