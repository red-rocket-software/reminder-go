package storage

import (
	"context"
	"log"
	"testing"
	"time"

	model "github.com/red-rocket-software/reminder-go/internal/reminder/domain"
	"github.com/stretchr/testify/require"
)

func TestStorageTodo_UpdateUserConfig(t *testing.T) {
	defer func() {
		err := Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	expectedUserID, err := SeedUserConfig()
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
		err = testConfigStorage.UpdateUserConfig(context.Background(), expectedUserID, updateConfigInput)
		require.NoError(t, err)
	})
	t.Run("empty input", func(t *testing.T) {
		err = testConfigStorage.UpdateUserConfig(context.Background(), expectedUserID, model.UserConfigs{})
		require.NoError(t, err)
	})

	t.Run("user configs not found", func(t *testing.T) {
		err = testConfigStorage.UpdateUserConfig(context.Background(), "0", updateConfigInput)
		require.Error(t, err)
	})

}

func TestStorage_CreateUserConfigs(t *testing.T) {
	defer func() {
		err := Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	expectedUserID, err := SeedUserConfig()
	if err != nil {
		log.Fatal("error truncate config")
	}
	require.NoError(t, err)

	expectedUserConfig := model.UserConfigs{
		ID:           "1",
		Notification: false,
		Period:       2,
	}

	t.Run("success", func(t *testing.T) {
		got, err := testConfigStorage.CreateUserConfigs(context.Background(), expectedUserConfig.ID)
		require.NoError(t, err)
		require.Equal(t, got.ID, expectedUserConfig.ID)
		require.Equal(t, got.Notification, expectedUserConfig.Notification)
		require.Equal(t, got.Period, expectedUserConfig.Period)
	})

	t.Run("fail user already exist", func(t *testing.T) {
		got, err := testConfigStorage.CreateUserConfigs(context.Background(), expectedUserID)
		require.Error(t, err)
		require.Empty(t, got)
	})
}

func TestStorage_GetUserConfigs(t *testing.T) {
	defer func() {
		err := Truncate()
		if err != nil {
			log.Fatal("error truncate table")
		}
	}()

	expectedUserID, err := SeedUserConfig()
	if err != nil {
		log.Fatal("error truncate config")
	}
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		config, err := testConfigStorage.GetUserConfigs(context.Background(), expectedUserID)
		require.NoError(t, err)
		require.Equal(t, config.ID, expectedUserID)
		require.Equal(t, config.Notification, true)
		require.Equal(t, config.Period, 2)
	})
	t.Run("no rows in result set", func(t *testing.T) {
		config, err := testConfigStorage.GetUserConfigs(context.Background(), "0")
		require.Empty(t, err)
		require.Empty(t, config)
	})
}

func TestTodoStorage_SeedUserConfig(t *testing.T) {
	defer func() {
		err := Truncate()
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SeedUserConfig()

			require.NoError(t, err)
			require.Equal(t, tt.userID, got)
		})
	}
}
