package discordhandler_test

import (
	"context"
	"errors"
	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"
	"testing"

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

		_, err := discord.GenerateAbsenceHandlerMsg(ctx, "playerone", "01/01/21", true)

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

		msg, err := discord.GenerateAbsenceHandlerMsg(context.Background(), "playerone", "01/01/21", true)

		assert.NoError(t, err)
		assert.Equal(t, "Absence(s) créée(s) pour le(s) :\n* 01-01-2021\n", msg)
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

		msg, err := discord.GenerateAbsenceHandlerMsg(context.Background(), "playerone", "01/01/21", false)

		assert.NoError(t, err)
		assert.Equal(t, "Absence(s) supprimée(s) pour le(s) :\n* 01-01-2021\n", msg)
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

		msg, err := discord.GenerateAbsenceHandlerMsg(context.Background(), "playerone", "01/01/21", true)

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

		msg, err := discord.GenerateAbsenceHandlerMsg(context.Background(), "playerone", "01/01/21", false)

		assert.Error(t, err)
		assert.Equal(t, "Error while deleting absence: Backend Error", msg)
		mockAbsenceUseCase.AssertExpectations(t)
	})
}
