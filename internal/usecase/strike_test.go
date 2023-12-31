//nolint:dupl
package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/internal/usecase"
	"github.com/antony-ramos/guildops/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStrikeUseCase_CreateStrike(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)

		mockBackend.On("CreateStrike", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := strikeUseCase.CreateStrike(context.Background(), "valid reason", "playername")

		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("invalid strike", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		err := strikeUseCase.CreateStrike(context.Background(),
			"ZBNSZVmKQwgZCBU9KjsbEOEewrPl5U1XkH10K4uXYVTuZiZiWzcydA1ISnH7iapcneGp"+
				"m4CjbdMd1FdDyxuQ4eluwy3jP7kfrLhT"+
				"Wcm6Pbj2DbMnd4J71OzqqPmntmWd5wyiUFoVtcVNthJXFO23rQIg6MrT25DI4V1LLHmZ9dcMJUbcdaGlJ60nLTgmKnBUhYzYC0roBXeC"+
				"jBCStg16teOgFS23m6j1Yrejjba9Eyro1YOi2ETX6sCesMvKfG2N0", "playername")

		assert.Error(t, err, "invalid strike")
		mockBackend.AssertExpectations(t)
	})

	t.Run("bug SearchPlayer", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("bug SearchPlayer"))

		err := strikeUseCase.CreateStrike(context.Background(), "valid reason", "playername")

		assert.Error(t, err, "bug SearchPlayer")
		mockBackend.AssertExpectations(t)
	})

	t.Run("player not found", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		err := strikeUseCase.CreateStrike(context.Background(), "valid reason", "playername")

		assert.Error(t, err, "player not found")
		mockBackend.AssertExpectations(t)
	})

	t.Run("bug Create Strike", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)

		mockBackend.On("CreateStrike", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("bug Create Strike"))

		err := strikeUseCase.CreateStrike(context.Background(), "valid reason", "playername")

		assert.Error(t, err, "bug Create Strike")
		mockBackend.AssertExpectations(t)
	})

	t.Run("Context is done", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := strikeUseCase.CreateStrike(ctx, "valid reason", "playername")

		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})
}

func TestStrikeUseCase_ReadStrike(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)

		strikes := []entity.Strike{
			{
				ID:     1,
				Reason: "valid reason",
				Date:   time.Now(),
				Season: "DF/S2",
			},
			{
				ID:     2,
				Reason: "valid reason 2",
				Date:   time.Now(),
				Season: "DF/S2",
			},
		}

		mockBackend.On("SearchStrike",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(strikes, nil)

		strikes, err := strikeUseCase.ReadStrikes(context.Background(), "playername")

		assert.NoError(t, err)
		assert.Equal(t, strikes, strikes)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Context is done", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := strikeUseCase.ReadStrikes(ctx, "playername")

		assert.Error(t, err)
	})

	t.Run("error on search player", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("error on search player"))

		_, err := strikeUseCase.ReadStrikes(context.Background(), "playername")

		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})
	t.Run("error on search strike", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)

		mockBackend.On("SearchStrike",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("error on search strike"))

		_, err := strikeUseCase.ReadStrikes(context.Background(), "playername")

		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("player not found", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil)

		_, err := strikeUseCase.ReadStrikes(context.Background(), "playername")

		assert.Equal(t, "player not found", err.Error())
		mockBackend.AssertExpectations(t)
	})

	t.Run("no strikes found", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)

		mockBackend.On("SearchStrike",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil)

		_, err := strikeUseCase.ReadStrikes(context.Background(), "playername")

		assert.Equal(t, "no strikes found", err.Error())
		mockBackend.AssertExpectations(t)
	})
}

func TestStrikeUseCase_DeleteStrike(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		mockBackend.On("DeleteStrike", mock.Anything, mock.Anything).Return(nil)

		err := strikeUseCase.DeleteStrike(context.Background(), 1)

		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		mockBackend.On("DeleteStrike", mock.Anything, mock.Anything).Return(errors.New("Backend Error"))

		err := strikeUseCase.DeleteStrike(context.Background(), 1)

		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Context is done", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := strikeUseCase.DeleteStrike(ctx, 1)

		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})
}
