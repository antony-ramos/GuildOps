package discordhandler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"
	"github.com/antony-ramos/guildops/internal/controller/discord/mocks"
	"github.com/antony-ramos/guildops/internal/entity"

	"github.com/bwmarrin/discordgo"
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
		assert.Equal(t, "Absence(s) created for :\n* "+time.Now().AddDate(0, 0, 1).Format("Mon 02/01/06")+"\n", msg)
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
			"* "+time.Now().AddDate(0, 0, 1).Format("Mon 02/01/06")+"\n"+
			"* "+time.Now().AddDate(0, 0, 2).Format("Mon 02/01/06")+" Absence already exists\n", msg)
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
		assert.Equal(t, "Absence(s) deleted for :\n* "+time.Now().AddDate(0, 0, 1).Format("Mon 02/01/06")+"\n", msg)
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

func TestDiscord_ListAbsenceHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		AbsenceUseCase := mocks.NewAbsenceUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: AbsenceUseCase,
		}

		AbsenceUseCase.On("ListAbsence", mock.Anything, mock.Anything).
			Return([]entity.Absence{
				{
					Player: &entity.Player{
						Name: "Paragon",
					},
					Raid: &entity.Raid{
						Date: time.Date(2030, time.September, 29, 0, 0, 0, 0, time.UTC),
					},
				},
				{
					Player: &entity.Player{
						Name: "Paragon",
					},
					Raid: &entity.Raid{
						Date: time.Date(2030, time.September, 30, 0, 0, 0, 0, time.UTC),
					},
				},
			}, nil)

		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionApplicationCommand,
				Member: &discordgo.Member{
					User: &discordgo.User{
						Username: "test",
					},
				},
				Data: discordgo.ApplicationCommandInteractionData{
					ID:       "mock",
					Name:     "mock",
					TargetID: "mock",
					Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Name:  "date",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "29/09/23",
						},
					},
				},
			},
		}

		msg, err := discord.ListAbsenceHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, "29/09/23 absences :\n* Paragon\n* Paragon\n", msg)
		AbsenceUseCase.AssertExpectations(t)
	})
}
