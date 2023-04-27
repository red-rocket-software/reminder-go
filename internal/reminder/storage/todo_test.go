package storage

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"

	model "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/red-rocket-software/reminder-go/pkg/pagination"
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

func TestStorageTodo_GetNewReminds(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	expectedToto, err := testStorage.SeedTodos()

	var nextCursor int
	if len(expectedToto) > 0 {
		nextCursor = expectedToto[len(expectedToto)-2].ID
	}

	if err != nil {
		log.Fatal("error seed todos")
	}

	type args struct {
		ctx    context.Context
		params pagination.Page
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   int
		want2   int
		wantErr bool
	}{
		{name: "success", args: args{context.Background(),
			pagination.Page{Limit: 3, Filter: "DeadlineAt", FilterOption: "ASC"},
			expectedToto[0].UserID,
		},
			want:    3,
			want1:   4,
			want2:   nextCursor,
			wantErr: false},
		{name: "error no limit", args: args{context.Background(), pagination.Page{
			Filter:       "DeadlineAt",
			FilterOption: "ASC",
		}, expectedToto[0].UserID},
			want:    0,
			want1:   0,
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := testStorage.GetNewReminds(tt.args.ctx, tt.args.params, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNewReminds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetNewReminds() got = %v, want %v", len(got), tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetNewReminds() got1 = %v, want1 %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("GetNewReminds() got2 = %v, want2 %v", got1, tt.want2)
			}
		})
	}
}

func TestStorageTodo_GetAllReminds(t *testing.T) {
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

	var nextCursor int
	if len(expectedTodo) > 0 {
		nextCursor = expectedTodo[len(expectedTodo)-5].ID
	}

	type args struct {
		ctx         context.Context
		fetchParams pagination.Page
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
		{name: "success", args: args{context.Background(), pagination.Page{
			Limit:        2,
			Filter:       "DeadlineAt",
			FilterOption: "ASC",
		}, expectedTodo[0].UserID},
			want:    []model.Todo{expectedTodo[1], expectedTodo[0]},
			want1:   5,
			want2:   nextCursor,
			wantErr: false},
		{name: "error no limit", args: args{context.Background(), pagination.Page{Filter: "DeadlineAt", FilterOption: "ASC"}, expectedTodo[0].UserID},
			want:    []model.Todo{},
			want1:   0,
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := testStorage.GetAllReminds(tt.args.ctx, tt.args.fetchParams, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllReminds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllReminds() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetAllReminds() got1 = %v, want1 %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("GetAllReminds() got2 = %v, want2 %v", got2, tt.want2)
			}
		})
	}
}

func TestStorageTodo_GetCompletedReminds(t *testing.T) {
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

	var nextCursor int
	if len(expectedTodo) > 0 {
		nextCursor = expectedTodo[len(expectedTodo)-4].ID
	}

	type args struct {
		ctx    context.Context
		params Params
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Todo
		want1   int
		want2   int
		wantErr bool
	}{
		{name: "success", args: args{context.Background(), Params{
			Page: pagination.Page{
				Limit:        5,
				Filter:       "DeadlineAt",
				FilterOption: "ASC",
			},
			TimeRangeFilter: TimeRangeFilter{},
		},
			expectedTodo[0].UserID,
		},
			want:    []model.Todo{expectedTodo[1]},
			want1:   1,
			want2:   nextCursor,
			wantErr: false},
		{name: "error no limit", args: args{context.Background(), Params{
			Page: pagination.Page{Filter: "DeadlineAt", FilterOption: "ASC"},
		}, expectedTodo[0].UserID},

			want:    []model.Todo{},
			want1:   0,
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := testStorage.GetCompletedReminds(tt.args.ctx, tt.args.params, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetComplitedReminds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetComplitedReminds() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetComplitedReminds() got1 = %v, want1 %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("GetComplitedReminds() got2 = %v, want2 %v", got2, tt.want2)
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
