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

func TestFailUseCase_CreateFail(t *testing.T) {
	t.Parallel()
	t.Run("context is done", func(t *testing.T) {
		t.Parallel()
		mockBackend := mocks.NewBackend(t)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := FailUseCase.CreateFail(ctx, "failreason", time.Now(), "playerone")
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("raid not found", func(t *testing.T) {
		t.Parallel()
		mockBackend := mocks.NewBackend(t)

		mockBackend.
			On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{{ID: 1}}, nil)
		mockBackend.
			On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Raid{}, nil)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		ctx := context.Background()
		err := FailUseCase.CreateFail(ctx, "failreason", time.Now(), "playerone")
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("player not found", func(t *testing.T) {
		t.Parallel()
		mockBackend := mocks.NewBackend(t)
		mockBackend.
			On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{}, nil)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		ctx := context.Background()
		err := FailUseCase.CreateFail(ctx, "failreason", time.Now(), "playerone")
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("fail validate", func(t *testing.T) {
		t.Parallel()
		mockBackend := mocks.NewBackend(t)
		mockBackend.
			On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{{ID: 1}}, nil)
		mockBackend.
			On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Raid{{ID: 1}}, nil)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		ctx := context.Background()
		err := FailUseCase.CreateFail(ctx, "", time.Now(), "playerone")
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("backend error", func(t *testing.T) {
		t.Parallel()
		mockBackend := mocks.NewBackend(t)
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]entity.Player{{
			ID:   1,
			Name: "playername",
		}}, nil)
		mockBackend.On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]entity.Raid{{
			ID:         1,
			Name:       "raidname",
			Difficulty: "normal",
		}}, nil)
		mockBackend.
			On("CreateFail", mock.Anything, mock.Anything).
			Return(entity.Fail{}, errors.New("Backend Error"))

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		ctx := context.Background()
		err := FailUseCase.CreateFail(ctx, "failreason", time.Now(), "playerone")
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockBackend := mocks.NewBackend(t)
		mockBackend.
			On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{{
				ID:   1,
				Name: "playername",
			}}, nil)
		mockBackend.
			On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Raid{{
				ID:         1,
				Name:       "raidname",
				Difficulty: "normal",
			}}, nil)
		mockBackend.On("CreateFail", mock.Anything, mock.Anything).Return(entity.Fail{}, nil)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		ctx := context.Background()
		err := FailUseCase.CreateFail(ctx, "failreason", time.Now(), "playerone")
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})
}

func TestFailUseCase_DeleteFail(t *testing.T) {
	t.Parallel()
	t.Run("context is done", func(t *testing.T) {
		t.Parallel()
		mockBackend := mocks.NewBackend(t)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := FailUseCase.DeleteFail(ctx, 1)
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		mockBackend.On("DeleteFail", mock.Anything, mock.Anything).Return(nil)

		err := FailUseCase.DeleteFail(context.Background(), 1)
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		mockBackend.On("DeleteFail", mock.Anything, mock.Anything).Return(nil)

		err := FailUseCase.DeleteFail(context.Background(), 1)
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})
}

//nolint:dupl
func TestFailUseCase_ListFailOnPLayer(t *testing.T) {
	t.Parallel()
	t.Run("context is done", func(t *testing.T) {
		t.Parallel()
		t.Run("context is done", func(t *testing.T) {
			t.Parallel()
			mockBackend := mocks.NewBackend(t)

			FailUseCase := usecase.NewFailUseCase(mockBackend)

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err := FailUseCase.ListFailOnPLayer(ctx, "playerone")
			assert.Error(t, err)
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		mockBackend.
			On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{{ID: 1}}, nil)
		mockBackend.
			On("SearchFail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Fail{{ID: 1}}, nil)

		_, err := FailUseCase.ListFailOnPLayer(context.Background(), "playerone")
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})
}

//nolint:dupl
func TestFailUseCase_ListFailOnRaid(t *testing.T) {
	t.Parallel()
	t.Run("context is done", func(t *testing.T) {
		t.Parallel()
		t.Run("context is done", func(t *testing.T) {
			t.Parallel()
			mockBackend := mocks.NewBackend(t)

			FailUseCase := usecase.NewFailUseCase(mockBackend)

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err := FailUseCase.ListFailOnRaid(ctx, "raidone")
			assert.Error(t, err)
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		mockBackend.
			On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Raid{{ID: 1}}, nil)
		mockBackend.
			On("SearchFail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Fail{{ID: 1}}, nil)

		_, err := FailUseCase.ListFailOnRaid(context.Background(), "raidone")
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})
}

func TestFailUseCase_ListFailOnRaidAndPlayer(t *testing.T) {
	t.Parallel()
	t.Run("context is done", func(t *testing.T) {
		t.Parallel()
		t.Run("context is done", func(t *testing.T) {
			t.Parallel()
			mockBackend := mocks.NewBackend(t)

			FailUseCase := usecase.NewFailUseCase(mockBackend)

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err := FailUseCase.ListFailOnRaidAndPlayer(ctx, "raidone", "playerone")
			assert.Error(t, err)
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		mockBackend.
			On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Raid{{ID: 1}}, nil)
		mockBackend.
			On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{{ID: 1}}, nil)
		mockBackend.
			On("SearchFail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Fail{{ID: 1}}, nil)

		_, err := FailUseCase.ListFailOnRaidAndPlayer(context.Background(), "raidone", "playerone")
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})
}

func TestFailUseCase_ReadFail(t *testing.T) {
	t.Parallel()
	t.Run("context is done", func(t *testing.T) {
		t.Parallel()
		t.Run("context is done", func(t *testing.T) {
			t.Parallel()
			mockBackend := mocks.NewBackend(t)

			FailUseCase := usecase.NewFailUseCase(mockBackend)

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err := FailUseCase.ReadFail(ctx, 1)
			assert.Error(t, err)
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		mockBackend.On("ReadFail", mock.Anything, mock.Anything).Return(entity.Fail{}, nil)

		_, err := FailUseCase.ReadFail(context.Background(), 1)
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		mockBackend.On("ReadFail", mock.Anything, mock.Anything).Return(entity.Fail{}, errors.New("Backend Error"))

		_, err := FailUseCase.ReadFail(context.Background(), 1)
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})
}

func TestFailUseCase_UpdateFail(t *testing.T) {
	t.Parallel()
	t.Run("context is done", func(t *testing.T) {
		t.Parallel()
		t.Run("context is done", func(t *testing.T) {
			t.Parallel()
			mockBackend := mocks.NewBackend(t)

			FailUseCase := usecase.NewFailUseCase(mockBackend)

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			err := FailUseCase.UpdateFail(ctx, 1, "failreason")
			assert.Error(t, err)
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		mockBackend.On("ReadFail", mock.Anything, mock.Anything).Return(entity.Fail{
			ID:     1,
			Reason: "reason",
			Player: &entity.Player{
				ID:   1,
				Name: "playername",
			},
			Raid: &entity.Raid{
				ID:         1,
				Name:       "raidname",
				Difficulty: "normal",
			},
		}, nil)
		mockBackend.On("UpdateFail", mock.Anything, mock.Anything).Return(nil)

		err := FailUseCase.UpdateFail(context.Background(), 1, "failreason")
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		FailUseCase := usecase.NewFailUseCase(mockBackend)

		mockBackend.On("ReadFail", mock.Anything, mock.Anything).Return(entity.Fail{
			ID:     1,
			Reason: "reason",
			Player: &entity.Player{
				ID:   1,
				Name: "playername",
			},
			Raid: &entity.Raid{
				ID:         1,
				Name:       "raidname",
				Difficulty: "normal",
			},
		}, nil)
		mockBackend.On("UpdateFail", mock.Anything, mock.Anything).Return(errors.New("Backend Error"))

		err := FailUseCase.UpdateFail(context.Background(), 1, "failreason")
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})
}
