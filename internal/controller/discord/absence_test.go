package discordhandler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"
	"github.com/antony-ramos/guildops/internal/entity"

	"github.com/antony-ramos/guildops/internal/controller/discord/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDiscord_GenerateAbsenceHandlerMsg(t *testing.T) {
	t.Parallel()

	t.Run("context is done", func(t *testing.T) {
		t.Parallel()

		mockAbsenceUseCase := mocks.NewAbsenceUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: mockAbsenceUseCase,
			PlayerUseCase:  nil,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := discord.GenerateAbsenceHandlerMsg(ctx, "playerone", "01/01/21", "", true)

		assert.Error(t, err)
		mockAbsenceUseCase.AssertExpectations(t)
	})

	t.Run("Invalid date", func(t *testing.T) {
		t.Parallel()

		mockAbsenceUseCase := mocks.NewAbsenceUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: mockAbsenceUseCase,
			PlayerUseCase:  nil,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		_, err := discord.GenerateAbsenceHandlerMsg(
			context.Background(), "playerone", time.Now().AddDate(0, 0, -2).Format("02/01/06"), "", true)

		assert.Error(t, err)
		mockAbsenceUseCase.AssertExpectations(t)
	})

	t.Run("Create Absence", func(t *testing.T) {
		t.Parallel()

		mockAbsenceUseCase := mocks.NewAbsenceUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: mockAbsenceUseCase,
			PlayerUseCase:  nil,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		mockAbsenceUseCase.On("CreateAbsence", mock.Anything, "playerone", mock.Anything).Return(nil)

		msg, err := discord.GenerateAbsenceHandlerMsg(
			context.Background(), "playerone", time.Now().AddDate(0, 0, 1).Format("02/01/06"), "", true)

		assert.NoError(t, err)
		assert.Equal(t, "Absence(s) created for :\n* "+time.Now().AddDate(0, 0, 1).Format("02/01/06")+"\n", msg)
		mockAbsenceUseCase.AssertExpectations(t)
	})

	t.Run("Create Absence Over Range with already exist", func(t *testing.T) {
		t.Parallel()

		mockAbsenceUseCase := mocks.NewAbsenceUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: mockAbsenceUseCase,
			PlayerUseCase:  nil,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		mockAbsenceUseCase.On("CreateAbsence", mock.Anything, "playerone", mock.Anything).Return(nil).Once()
		mockAbsenceUseCase.
			On("CreateAbsence", mock.Anything, "playerone", mock.Anything).
			Return(errors.New(" absence already exist ")).Once()
		mockAbsenceUseCase.
			On("CreateAbsence", mock.Anything, "playerone", mock.Anything).
			Return(errors.New("no raid found")).Once()

		msg, err := discord.GenerateAbsenceHandlerMsg(
			context.Background(), "playerone",
			time.Now().AddDate(0, 0, 1).Format("02/01/06"),
			time.Now().AddDate(0, 0, 3).Format("02/01/06"), true)

		assert.NoError(t, err)
		assert.Equal(t, "Absence(s) created for :\n"+
			"* "+time.Now().AddDate(0, 0, 1).Format("02/01/06")+"\n"+
			"* "+time.Now().AddDate(0, 0, 2).Format("02/01/06")+"\n", msg)
		mockAbsenceUseCase.AssertExpectations(t)
	})

	t.Run("Delete Absence", func(t *testing.T) {
		t.Parallel()

		mockAbsenceUseCase := mocks.NewAbsenceUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: mockAbsenceUseCase,
			PlayerUseCase:  nil,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		mockAbsenceUseCase.On("DeleteAbsence", mock.Anything, "playerone", mock.Anything).Return(nil)

		msg, err := discord.GenerateAbsenceHandlerMsg(
			context.Background(), "playerone", time.Now().AddDate(0, 0, 1).Format("02/01/06"), "", false)

		assert.NoError(t, err)
		assert.Equal(t, "Absence(s) deleted for :\n* "+time.Now().AddDate(0, 0, 1).Format("02/01/06")+"\n", msg)
		mockAbsenceUseCase.AssertExpectations(t)
	})

	t.Run("Backend Error Create", func(t *testing.T) {
		t.Parallel()

		mockAbsenceUseCase := mocks.NewAbsenceUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: mockAbsenceUseCase,
			PlayerUseCase:  nil,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		mockAbsenceUseCase.On("CreateAbsence", mock.Anything, "playerone", mock.Anything).Return(errors.New("Backend Error"))

		msg, err := discord.GenerateAbsenceHandlerMsg(
			context.Background(), "playerone", time.Now().AddDate(0, 0, 1).Format("02/01/06"), "", true)

		assert.Error(t, err)
		assert.Equal(t, "Error while creating absence: Backend Error", msg)
		mockAbsenceUseCase.AssertExpectations(t)
	})

	t.Run("Backend Error Delete", func(t *testing.T) {
		t.Parallel()

		mockAbsenceUseCase := mocks.NewAbsenceUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: mockAbsenceUseCase,
			PlayerUseCase:  nil,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		mockAbsenceUseCase.On("DeleteAbsence", mock.Anything, "playerone", mock.Anything).Return(errors.New("Backend Error"))

		msg, err := discord.GenerateAbsenceHandlerMsg(
			context.Background(), "playerone", time.Now().AddDate(0, 0, 1).Format("02/01/06"), "", false)

		assert.Error(t, err)
		assert.Equal(t, "Error while deleting absence: Backend Error", msg)
		mockAbsenceUseCase.AssertExpectations(t)
	})
}

func TestDiscord_GenerateListAbsenceHandlerMsg(t *testing.T) {
	t.Parallel()

	t.Run("context is done", func(t *testing.T) {
		t.Parallel()

		mockAbsenceUseCase := mocks.NewAbsenceUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: mockAbsenceUseCase,
			PlayerUseCase:  nil,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := discord.GenerateListAbsenceHandlerMsg(ctx, "01/01/21")

		assert.Error(t, err)
		mockAbsenceUseCase.AssertExpectations(t)
	})

	t.Run("Invalid date", func(t *testing.T) {
		t.Parallel()

		mockAbsenceUseCase := mocks.NewAbsenceUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: mockAbsenceUseCase,
			PlayerUseCase:  nil,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		_, err := discord.GenerateListAbsenceHandlerMsg(context.Background(), "01-01-2021")

		assert.Error(t, err)
		mockAbsenceUseCase.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockAbsenceUseCase := mocks.NewAbsenceUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: mockAbsenceUseCase,
			PlayerUseCase:  nil,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		mockAbsenceUseCase.On("ListAbsence", mock.Anything, mock.Anything).Return([]entity.Absence{
			{
				Player: &entity.Player{
					Name: "playerone",
				},
				Raid: &entity.Raid{
					Name:       "raidname",
					Difficulty: "normal",
					Date:       time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		}, nil)

		msg, err := discord.GenerateListAbsenceHandlerMsg(context.Background(), "01/01/21")

		assert.NoError(t, err)
		assert.Equal(t, "Absence(s) pour le 01-01-2021 :\n* playerone\n", msg)
		mockAbsenceUseCase.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()

		mockAbsenceUseCase := mocks.NewAbsenceUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: mockAbsenceUseCase,
			PlayerUseCase:  nil,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		mockAbsenceUseCase.On("ListAbsence", mock.Anything, mock.Anything).Return(nil, errors.New("Backend Error"))

		msg, err := discord.GenerateListAbsenceHandlerMsg(context.Background(), "01/01/21")

		assert.Error(t, err)
		assert.Equal(t, "Aucune absence pour le 01/01/21\n", msg)
		mockAbsenceUseCase.AssertExpectations(t)
	})
}
