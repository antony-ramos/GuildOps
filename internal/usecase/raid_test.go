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

func TestRaidUseCase_CreateRaid(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		raidUseCase := usecase.NewRaidUseCase(mockBackend)

		raid := entity.Raid{
			Name:       "raid name",
			Difficulty: "normal",
			Date:       time.Now(),
		}

		mockBackend.On("CreateRaid", mock.Anything, mock.Anything).
			Return(raid, nil)

		r, err := raidUseCase.CreateRaid(context.Background(), raid.Name, raid.Difficulty, raid.Date)

		assert.NoError(t, err)
		assert.Equal(t, raid, r)
		mockBackend.AssertExpectations(t)
	})
}

func TestRaidUseCase_DeleteRaid(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		raidUseCase := usecase.NewRaidUseCase(mockBackend)

		mockBackend.On("DeleteRaid", mock.Anything, mock.Anything).
			Return(nil)

		err := raidUseCase.DeleteRaid(context.Background(), 1)

		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})
	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		raidUseCase := usecase.NewRaidUseCase(mockBackend)

		mockBackend.On("DeleteRaid", mock.Anything, mock.Anything).
			Return(errors.New("Backend Error"))

		err := raidUseCase.DeleteRaid(context.Background(), 1)

		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})
	t.Run("ctx have been canceled", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		raidUseCase := usecase.NewRaidUseCase(mockBackend)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := raidUseCase.DeleteRaid(ctx, 1)

		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})
}

func TestRaidUseCase_ReadRaid(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		raidUseCase := usecase.NewRaidUseCase(mockBackend)

		mockBackend.On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Raid{
				{
					ID:         1,
					Name:       "raid name",
					Difficulty: "normal",
					Date:       time.Now(),
				},
			}, nil)

		r, err := raidUseCase.ReadRaid(context.Background(), time.Now())

		assert.NoError(t, err)
		assert.Equal(t, 1, r.ID)
		mockBackend.AssertExpectations(t)
	})
	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		raidUseCase := usecase.NewRaidUseCase(mockBackend)

		mockBackend.On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("Backend Error"))

		r, err := raidUseCase.ReadRaid(context.Background(), time.Now())

		assert.Error(t, err)
		assert.Equal(t, entity.Raid{}, r)
		mockBackend.AssertExpectations(t)
	})
	t.Run("ctx have been canceled", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		raidUseCase := usecase.NewRaidUseCase(mockBackend)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		r, err := raidUseCase.ReadRaid(ctx, time.Now())

		assert.Error(t, err)
		assert.Equal(t, entity.Raid{}, r)
		mockBackend.AssertExpectations(t)
	})
}
