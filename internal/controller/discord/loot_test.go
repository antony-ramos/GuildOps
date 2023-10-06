package discordhandler_test

import (
	"context"
	"testing"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"

	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"
	"github.com/antony-ramos/guildops/internal/controller/discord/mocks"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
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
		}

		mockLootUseCase.On("CreateLoot", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

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
							Name:  "loot-name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "TestLoot",
						},
						{
							Name:  "raid-date",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "01/10/23",
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

		msg, err := discord.AttributeLootHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "Loot successfully attributed")
		mockLootUseCase.AssertExpectations(t)
	})
}

func TestDiscord_ListLootsOnPlayerHandler(t *testing.T) {
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
		}

		mockLootUseCase.On("ListLootOnPLayer", mock.Anything, mock.Anything).
			Return([]entity.Loot{
				{
					ID:   1,
					Name: "TestLoot",
					Raid: &entity.Raid{
						Date:       time.Now(),
						Difficulty: "Heroic",
					},
				},
				{
					ID:   2,
					Name: "TestLoot2",
					Raid: &entity.Raid{
						Date:       time.Now(),
						Difficulty: "Heroic",
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
							Name:  "player-name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "TestPlayer",
						},
					},
				},
			},
		}

		msg, err := discord.ListLootsOnPlayerHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "All loots of TestPlayer:\n"+
			"TestLoot "+time.Now().Format("02/01/06")+" Heroic\n"+
			"TestLoot2 "+time.Now().Format("02/01/06")+" Heroic\n")
		mockLootUseCase.AssertExpectations(t)
	})
}

// DeleteLootHandler.
func TestDiscord_DeleteLootHandler(t *testing.T) {
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
		}

		mockLootUseCase.On("DeleteLoot", mock.Anything, mock.Anything).
			Return(nil)

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
							Name:  "id",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "1",
						},
					},
				},
			},
		}

		msg, err := discord.DeleteLootHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "Loot successfully deleted")
		mockLootUseCase.AssertExpectations(t)
	})
}
