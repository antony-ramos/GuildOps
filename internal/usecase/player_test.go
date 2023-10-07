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

func TestPlayerUseCase_LinkPlayer(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		playerUseCase := usecase.NewPlayerUseCase(mockBackend)

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, "titi").
			Return(nil, nil)

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, "toto", mock.Anything).
			Return([]entity.Player{{Name: "toto"}}, nil)

		mockBackend.On("UpdatePlayer", mock.Anything, entity.Player{Name: "toto", DiscordName: "titi"}).
			Return(nil)

		err := playerUseCase.LinkPlayer(context.Background(), "toto", "titi")
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Context is done", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		playerUseCase := usecase.NewPlayerUseCase(mockBackend)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := playerUseCase.LinkPlayer(ctx, "toto", "titi")
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})
}

func TestPlayerUseCase_DeletePlayer(t *testing.T) {
	t.Parallel()

	t.Run("Delete player with no strikes, fails or absences", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		playerUseCase := usecase.NewPlayerUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		mockBackend.On("SearchStrike", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil)
		mockBackend.On("SearchFail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil)
		mockBackend.On("SearchAbsence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil)
		mockBackend.On("DeletePlayer", mock.Anything, mock.Anything).
			Return(nil)

		err := playerUseCase.DeletePlayer(context.Background(), "playername")
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Delete player with strikes", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		playerUseCase := usecase.NewPlayerUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		mockBackend.On("SearchStrike", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Strike{
				{
					ID:     1,
					Player: &player,
					Date:   time.Now(),
				},
				{
					ID:     1,
					Player: &player,
					Date:   time.Now(),
				},
			}, nil)
		mockBackend.On("DeleteStrike", mock.Anything, mock.Anything).
			Return(nil).Twice()
		mockBackend.On("SearchFail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil)
		mockBackend.On("SearchAbsence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil)
		mockBackend.On("DeletePlayer", mock.Anything, mock.Anything).
			Return(nil)

		err := playerUseCase.DeletePlayer(context.Background(), "playername")
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Delete player with fails", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		playerUseCase := usecase.NewPlayerUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		mockBackend.On("SearchStrike", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil)
		mockBackend.On("SearchFail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Fail{
				{
					ID:     1,
					Player: &player,
					Raid: &entity.Raid{
						ID:   1,
						Name: "raidname",
						Date: time.Now(),
					},
				},
				{
					ID:     1,
					Player: &player,
					Raid: &entity.Raid{
						ID:   2,
						Name: "raidname",
						Date: time.Now(),
					},
				},
			}, nil)

		mockBackend.On("DeleteFail", mock.Anything, mock.Anything).
			Return(nil).Twice()
		mockBackend.On("SearchAbsence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil)
		mockBackend.On("DeletePlayer", mock.Anything, mock.Anything).
			Return(nil)

		err := playerUseCase.DeletePlayer(context.Background(), "playername")
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Delete player with absences", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		playerUseCase := usecase.NewPlayerUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		mockBackend.On("SearchStrike", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil)
		mockBackend.On("SearchFail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil)
		mockBackend.On("SearchAbsence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Absence{
				{
					ID:     1,
					Player: &player,
					Raid: &entity.Raid{
						ID:   1,
						Name: "raidname",
						Date: time.Now(),
					},
				},
				{
					ID:     1,
					Player: &player,
					Raid: &entity.Raid{
						ID:   2,
						Name: "raidname",
						Date: time.Now(),
					},
				},
			}, nil)
		mockBackend.On("DeleteAbsence", mock.Anything, mock.Anything).
			Return(nil).Twice()

		mockBackend.On("DeletePlayer", mock.Anything, mock.Anything).
			Return(nil)

		err := playerUseCase.DeletePlayer(context.Background(), "playername")
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})
}
