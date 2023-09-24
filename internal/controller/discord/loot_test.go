package discordhandler_test

import (
	"context"
	"testing"

	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"
	"github.com/antony-ramos/guildops/internal/controller/discord/mocks"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

func TestDiscord_AttributeLootHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockLootUseCase := mocks.NewLootUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: nil,
			PlayerUseCase:  nil,
			StrikeUseCase:  nil,
			LootUseCase:    mockLootUseCase,
			RaidUseCase:    nil,
			Fake:           true,
		}

		mockLootUseCase.On("CreateLoot", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
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
							Name:  "loot-name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "TestLoot",
						},
						{
							Name:  "raid-id",
							Type:  discordgo.ApplicationCommandOptionInteger,
							Value: float64(123),
						},
						{
							Name:  "player-name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "TestPlayer",
						},
					},
				},
			},
		}

		err := discord.AttributeLootHandler(context.Background(), session, interaction)
		if err != nil {
			return
		}
		mockLootUseCase.AssertExpectations(t)
	})
}
