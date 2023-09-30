package discordhandler_test

import (
	"context"
	"testing"

	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"
	"github.com/antony-ramos/guildops/internal/controller/discord/mocks"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

func TestDiscord_FailOnPlayerHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
			Fake:        true,
		}

		mockFailUseCase.On("CreateFail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
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
							Name:  "name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "Milowenn",
						},
						{
							Name:  "reason",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "why not",
						},
						{
							Name:  "date",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "03/05/23",
						},
					},
				},
			},
		}

		err := discord.FailOnPlayerHandler(context.Background(), session, interaction)
		if err != nil {
			return
		}
		mockFailUseCase.AssertExpectations(t)
	})
}

//nolint:dupl
func TestDiscord_ListFailsOnPlayerHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
			Fake:        true,
		}

		mockFailUseCase.On("ListFailOnPLayer", mock.Anything, mock.Anything).
			Return(nil, nil)

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
							Value: "Milowenn",
						},
					},
				},
			},
		}

		err := discord.ListFailsOnPlayerHandler(context.Background(), session, interaction)
		if err != nil {
			return
		}
		mockFailUseCase.AssertExpectations(t)
	})
}

//nolint:dupl
func TestDiscord_ListFailsOnRaidHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
			Fake:        true,
		}

		mockFailUseCase.On("ListFailOnRaid", mock.Anything, mock.Anything).
			Return(nil, nil)

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
							Name:  "date",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "29/09/23",
						},
					},
				},
			},
		}

		err := discord.ListFailsOnRaidHandler(context.Background(), session, interaction)
		if err != nil {
			return
		}
		mockFailUseCase.AssertExpectations(t)
	})
}

func TestDiscord_DeleteFailOnPlayerHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
			Fake:        true,
		}

		mockFailUseCase.On("DeleteFail", mock.Anything, mock.Anything).
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

		err := discord.DeleteFailHandler(context.Background(), session, interaction)
		if err != nil {
			return
		}
		mockFailUseCase.AssertExpectations(t)
	})
}
