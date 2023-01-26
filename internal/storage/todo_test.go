package storage

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/stretchr/testify/require"
)

func TestStorageTodo_CreateRemind(t *testing.T) {
	defer func() {
		err := testStorage.Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

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

	insertTodo := model.Todo{
		Description: "test",
		DeadlineAt:  date,
		CreatedAt:   time.Now(),
	}

	id, _ := testStorage.CreateRemind(context.Background(), insertTodo)

	got, err := testStorage.GetRemindByID(context.Background(), id)

	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.Equal(t, insertTodo.Description, got.Description)
	require.Equal(t, insertTodo.DeadlineAt, got.DeadlineAt)

}

// func TestStorageTodo_GetNewReminds(t *testing.T) {
// 	defer func() {
// 		err := testStorage.Truncate()
// 		if err != nil {
// 			log.Fatal("error truncate table")
// 		}
// 	}()

// 	expectedToto, err := testStorage.SeedTodos()

// 	var nextCursor int
// 	if len(expectedToto) > 0 {
// 		nextCursor = expectedToto[len(expectedToto)-2].ID
// 	}

// 	if err != nil {
// 		log.Fatal("error seed todos")
// 	}

// 	type args struct {
// 		ctx    context.Context
// 		params FetchParam
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    int
// 		want1   int
// 		wantErr bool
// 	}{
// 		{name: "success", args: args{context.Background(), FetchParam{
// 			Limit: 3,
// 		}},
// 			want:    3,
// 			want1:   nextCursor,
// 			wantErr: false},
// 		{name: "error no limit", args: args{context.Background(), FetchParam{}},
// 			want:    0,
// 			want1:   0,
// 			wantErr: false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, got1, err := testStorage.GetNewReminds(tt.args.ctx, tt.args.params)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetNewReminds() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(len(got), tt.want) {
// 				t.Errorf("GetNewReminds() got = %v, want %v", len(got), tt.want)
// 			}
// 			if got1 != tt.want1 {
// 				t.Errorf("GetNewReminds() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}
// }

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
		fetchParams FetchParam
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Todo
		want1   int
		wantErr bool
	}{
		{name: "success", args: args{context.Background(), FetchParam{
			Limit: 2,
		}},
			want:    []model.Todo{expectedTodo[1], expectedTodo[0]},
			want1:   nextCursor,
			wantErr: false},
		{name: "error no limit", args: args{context.Background(), FetchParam{}},
			want:    nil,
			want1:   0,
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := testStorage.GetAllReminds(tt.args.ctx, tt.args.fetchParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllReminds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllReminds() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetAllReminds() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestStorageTodo_GetComplitedReminds(t *testing.T) {
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
		params FetchParam
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Todo
		want1   int
		wantErr bool
	}{
		{name: "success", args: args{context.Background(), FetchParam{
			Limit: 5,
		}},
			want:    []model.Todo{expectedTodo[1]},
			want1:   nextCursor,
			wantErr: false},
		{name: "error no limit", args: args{context.Background(), FetchParam{}},
			want:    nil,
			want1:   0,
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := testStorage.GetComplitedReminds(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetComplitedReminds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetComplitedReminds() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetComplitedReminds() got1 = %v, want %v", got1, tt.want1)
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

	updateInput := model.TodoUpdate{
		Description: "New text",
		FinishedAt:  nil,
		Completed:   true,
	}

	err = testStorage.UpdateRemind(context.Background(), expectedTodo[1].ID, updateInput)
	require.NoError(t, err)

	newTodo, _ := testStorage.GetRemindByID(context.Background(), expectedTodo[1].ID)
	require.Equal(t, updateInput.Description, newTodo.Description)
	require.Equal(t, updateInput.Completed, newTodo.Completed)
}

func TestStorageTodo_SeedTodos(t *testing.T) {
	todos, err := testStorage.SeedTodos()

	require.NoError(t, err)
	require.Equal(t, len(todos), 5)
}
