package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/internal/usecase"
	"github.com/antony-ramos/guildops/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPlayerUseCase_CreatePlayer(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		playerUseCase := usecase.NewPlayerUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("CreatePlayer", mock.Anything, mock.Anything, mock.Anything).Return(player, nil)

		id, err := playerUseCase.CreatePlayer(context.Background(), "playername")
		assert.Equalf(t, 1, id, "CreatePlayer(%v)", "playername")

		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Player name is empty", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		playerUseCase := usecase.NewPlayerUseCase(mockBackend)

		_, err := playerUseCase.CreatePlayer(context.Background(), "")
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		playerUseCase := usecase.NewPlayerUseCase(mockBackend)

		mockBackend.On("CreatePlayer", mock.Anything, mock.Anything, mock.Anything).
			Return(entity.Player{}, errors.New("Backend Error"))

		id, err := playerUseCase.CreatePlayer(context.Background(), "playername")
		assert.Equalf(t, -1, id, "CreatePlayer(%v)", "playername")

		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Context is done", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		playerUseCase := usecase.NewPlayerUseCase(mockBackend)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		id, err := playerUseCase.CreatePlayer(ctx, "playername")
		assert.Equalf(t, -1, id, "CreatePlayer(%v)", "playername")

		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})
}