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

func TestCreateAbsence(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		raids := []entity.Raid{
			{
				ID:         1,
				Name:       "raid name",
				Difficulty: "normal",
			},
		}
		mockBackend.On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(raids, nil)
		mockBackend.On("CreateAbsence", mock.Anything, mock.Anything).Return(entity.Absence{}, nil)

		err := absenceUseCase.CreateAbsence(context.Background(), "PlayerName", time.Now())

		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Error while searching player", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{}, errors.New("error while searching player"))

		err := absenceUseCase.CreateAbsence(context.Background(), "PlayerName", time.Now())

		assert.Error(t, err, "no player found")
		mockBackend.AssertExpectations(t)
		t.Cleanup(func() {
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("PlayerName Invalid", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		raids := []entity.Raid{
			{
				ID:   1,
				Name: "RaidName",
			},
		}
		mockBackend.On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(raids, nil)

		// Call the function you want to test
		err := absenceUseCase.CreateAbsence(context.Background(), "", time.Now())

		// Assert that the function behaves as expected
		assert.Error(t, err, "name cannot be empty")
		mockBackend.AssertExpectations(t)
	})

	t.Run("Player Not Found", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{}, nil)

		err := absenceUseCase.CreateAbsence(context.Background(), "PlayerName", time.Now())

		// Assert that the function behaves as expected
		assert.Error(t, err, "no player found")
		mockBackend.AssertExpectations(t)
		t.Cleanup(func() {
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("Absence already exists", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "PlayerName",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		raids := []entity.Raid{
			{
				ID:   1,
				Name: "RaidName",
			},
		}
		mockBackend.On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(raids, nil)

		err := absenceUseCase.CreateAbsence(context.Background(), "playername", time.Now())

		// Assert that the function behaves as expected
		assert.Error(t, err, "absence already exists")
		mockBackend.AssertExpectations(t)
		t.Cleanup(func() {
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("Error while searching raid", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "PlayerName",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		mockBackend.On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("no raid found"))

		err := absenceUseCase.CreateAbsence(context.Background(), "PlayerName", time.Now())

		// Assert that the function behaves as expected
		assert.Error(t, err, "no raid found on this date")
		mockBackend.AssertExpectations(t)
		t.Cleanup(func() {
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("No raid on this date", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "PlayerName",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		mockBackend.On("SearchRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		date := time.Now()
		err := absenceUseCase.CreateAbsence(context.Background(), "PlayerName", date)

		assert.Error(t, err, "CreateAbsence:  no raid found on %s", date)
		mockBackend.AssertExpectations(t)
		t.Cleanup(func() {
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("Context is Done", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		date := time.Now()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := absenceUseCase.CreateAbsence(ctx, "PlayerName", date)

		assert.Error(t, err, "AbsenceUseCase - CreateAbsence:  ctx.Done: request took too much time to be proceed")
		mockBackend.AssertExpectations(t)
	})
}

func TestDeleteAbsence(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		absences := []entity.Absence{
			{
				ID:     1,
				Player: &player,
			},
		}
		mockBackend.On("SearchAbsence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(absences, nil)
		mockBackend.On("DeleteAbsence", mock.Anything, mock.Anything).Return(nil)

		err := absenceUseCase.DeleteAbsence(context.Background(), "PlayerName", time.Now())

		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Error while searching player", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{}, errors.New("error while searching player"))

		err := absenceUseCase.DeleteAbsence(context.Background(), "PlayerName", time.Now())

		// Assert that the function behaves as expected
		assert.Error(t, err, "no player found")
		mockBackend.AssertExpectations(t)
		t.Cleanup(func() {
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("Error while searching absence", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		mockBackend.On("SearchAbsence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("error while searching absence"))

		err := absenceUseCase.DeleteAbsence(context.Background(), "PlayerName", time.Now())

		// Assert that the function behaves as expected
		assert.Error(t, err, "no player found")
		mockBackend.AssertExpectations(t)
		t.Cleanup(func() {
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("Player Not Found", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{}, nil)

		err := absenceUseCase.DeleteAbsence(context.Background(), "PlayerName", time.Now())

		assert.Error(t, err, "no player found")
		mockBackend.AssertExpectations(t)
		t.Cleanup(func() {
			mockBackend.AssertExpectations(t)
		})
	})

	t.Run("No absence on date", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		mockBackend.On("SearchAbsence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		err := absenceUseCase.DeleteAbsence(context.Background(), "PlayerName", time.Now())

		assert.Error(t, err, "no absence found")
		mockBackend.AssertExpectations(t)
	})

	t.Run("Error deleting absence", func(t *testing.T) {
		t.Parallel()

		// Delete a new instance of the MockBackend
		mockBackend := mocks.NewBackend(t)

		// Delete an instance of your AbsenceUseCase with the mock backend
		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)
		absences := []entity.Absence{
			{
				ID:     1,
				Player: &player,
			},
		}
		mockBackend.On("SearchAbsence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(absences, nil)
		mockBackend.On("DeleteAbsence", mock.Anything, mock.Anything).Return(errors.New("error deleting absence"))

		err := absenceUseCase.DeleteAbsence(context.Background(), "PlayerName", time.Now())

		assert.Error(t, err, "error deleting absence")
		mockBackend.AssertExpectations(t)
	})

	t.Run("Context is Done", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		absenceUseCase := usecase.NewAbsenceUseCase(mockBackend)

		date := time.Now()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := absenceUseCase.DeleteAbsence(ctx, "PlayerName", date)

		assert.Error(t, err, "AbsenceUseCase - DeleteAbsence - ctx.Done: request took too much time to be proceed")
		mockBackend.AssertExpectations(t)
	})
}
