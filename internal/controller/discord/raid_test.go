package discordhandler_test

import (
	"context"
	"testing"

	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"
	"github.com/antony-ramos/guildops/internal/controller/discord/mocks"
	"github.com/antony-ramos/guildops/internal/entity"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

func TestDiscord_CreateRaidHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockRaidUseCase := mocks.NewRaidUseCase(t)

		discord := discordHandler.Discord{
			RaidUseCase: mockRaidUseCase,
			Fake:        true,
		}

		mockRaidUseCase.On("CreateRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(entity.Raid{}, nil)

		session := &discordgo.Session{StateEnabled: true, State: discordgo.NewState()}
		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionApplicationCommand,
				Data: discordgo.ApplicationCommandInteractionData{
					ID:       "mock",
					Name:     "mock",
					TargetID: "mock",
					Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Name:  "name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "random raid",
						},
						{
							Name:  "date",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "03/05/23",
						},
						{
							Name:  "difficulté",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "Heroic",
						},
					},
				},
			},
		}

		err := discord.CreateRaidHandler(context.Background(), session, interaction)
		if err != nil {
			return
		}
		mockRaidUseCase.AssertExpectations(t)
	})
}

func TestDiscord_DeleteRaidHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockRaidUseCase := mocks.NewRaidUseCase(t)

		discord := discordHandler.Discord{
			RaidUseCase: mockRaidUseCase,
			Fake:        true,
		}

		mockRaidUseCase.On("DeleteRaid", mock.Anything, mock.Anything).
			Return(nil)

		session := &discordgo.Session{StateEnabled: true, State: discordgo.NewState()}
		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionApplicationCommand,
				Data: discordgo.ApplicationCommandInteractionData{
					ID:       "mock",
					Name:     "mock",
					TargetID: "mock",
					Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Name:  "id",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "1",
						},
					},
				},
			},
		}

		err := discord.DeleteRaidHandler(context.Background(), session, interaction)
		if err != nil {
			return
		}
		mockRaidUseCase.AssertExpectations(t)
	})
}
