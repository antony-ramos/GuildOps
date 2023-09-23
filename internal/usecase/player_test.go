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

		strikeUseCase := usecase.NewPlayerUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("CreatePlayer", mock.Anything, mock.Anything, mock.Anything).Return(player, nil)

		id, err := strikeUseCase.CreatePlayer(context.Background(), "playername")
		assert.Equalf(t, 1, id, "CreatePlayer(%v)", "playername")

		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Player name is empty", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewPlayerUseCase(mockBackend)

		_, err := strikeUseCase.CreatePlayer(context.Background(), "")
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewPlayerUseCase(mockBackend)

		mockBackend.On("CreatePlayer", mock.Anything, mock.Anything, mock.Anything).
			Return(entity.Player{}, errors.New("Backend Error"))

		id, err := strikeUseCase.CreatePlayer(context.Background(), "playername")
		assert.Equalf(t, -1, id, "CreatePlayer(%v)", "playername")

		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})
}
