package storage

import (
	"context"
	"testing"
	"time"

	"github.com/red-rocket-software/reminder-go/internal/app/model"
	"github.com/stretchr/testify/require"
)

func TestStorageTodo_CreateRemind(t *testing.T) {

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
