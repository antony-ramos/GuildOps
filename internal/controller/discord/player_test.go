package discordhandler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"

	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"
	"github.com/antony-ramos/guildops/internal/controller/discord/mocks"
	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDiscord_PlayerHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success Create Player", func(t *testing.T) {
		t.Parallel()
		mockPlayerUseCase := mocks.NewPlayerUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: nil,
			PlayerUseCase:  mockPlayerUseCase,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		mockPlayerUseCase.On("CreatePlayer", mock.Anything, mock.Anything).
			Return(1, nil)

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
					Name:     "guildops-player-create",
					TargetID: "mock",
					Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Name:  "name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "TestPlayer",
						},
					},
				},
			},
		}

		msg, err := discord.PlayerHandler(context.Background(), interaction)
		mockPlayerUseCase.AssertExpectations(t)
		assert.Equal(t, msg, "Player TestPlayer created successfully: ID 1")
		assert.NoError(t, err)
	})

	t.Run("Success Delete Player", func(t *testing.T) {
		t.Parallel()
		mockPlayerUseCase := mocks.NewPlayerUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: nil,
			PlayerUseCase:  mockPlayerUseCase,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		mockPlayerUseCase.On("DeletePlayer", mock.Anything, mock.Anything).
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
					Name:     "guildops-player-delete",
					TargetID: "mock",
					Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Name:  "name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "TestPlayer",
						},
					},
				},
			},
		}

		msg, err := discord.PlayerHandler(context.Background(), interaction)
		assert.Equal(t, msg, "Player TestPlayer deleted successfully")
		mockPlayerUseCase.AssertExpectations(t)
		assert.NoError(t, err)
	})
}

func TestDiscord_GetPlayerHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success Get Player", func(t *testing.T) {
		t.Parallel()
		mockPlayerUseCase := mocks.NewPlayerUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: nil,
			PlayerUseCase:  mockPlayerUseCase,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		player := entity.Player{
			Name:        "TestPlayer",
			DiscordName: "TestDiscordName",
			ID:          1,
		}

		strikes := []entity.Strike{
			{
				ID:     1,
				Reason: "TestReason",
				Season: "DF/S2",
				Date:   time.Now(),
			},
			{
				ID:     1,
				Reason: "TestReason",
				Season: "DF/S2",
				Date:   time.Now(),
			},
		}
		player.Strikes = strikes

		loots := []entity.Loot{
			{
				ID:   1,
				Name: "TestLoot",
				Raid: &entity.Raid{
					ID:         1,
					Name:       "TestRaid",
					Difficulty: "TestDifficulty",
					Date:       time.Now(),
				},
			},
		}

		player.Loots = loots

		missedRaids := []entity.Raid{
			{
				ID:         1,
				Name:       "TestRaid",
				Difficulty: "TestDifficulty",
				Date:       time.Now(),
			},
		}

		player.MissedRaids = missedRaids

		fails := []entity.Fail{
			{
				ID:     1,
				Reason: "TestReason",
				Raid: &entity.Raid{
					ID:         1,
					Name:       "TestRaid",
					Difficulty: "TestDifficulty",
					Date:       time.Now(),
				},
			},
		}

		player.Fails = fails

		mockPlayerUseCase.On("ReadPlayer", mock.Anything, mock.Anything, "").
			Return(player, nil)

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
					Name:     "guildops-player-create",
					TargetID: "mock",
					Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Name:  "name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "TestPlayer",
						},
					},
				},
			},
		}

		exceptMsg := `Name : **TestPlayer**
ID : **1**
Discord ID : **TestDiscordName**
**Loots Count:**
*  TestDifficulty | 1 loots 
**Strikes (2) :**
*  ` + time.Now().Format("02/01/06") + ` | TestReason | DF/S2 | 1
*  ` + time.Now().Format("02/01/06") + ` | TestReason | DF/S2 | 1
**Absences (1) :**
*  ` + time.Now().Format("02/01/06") + ` | TestDifficulty | TestRaid
**Loots (1) :**
*  ` + time.Now().Format("02/01/06") + ` | TestDifficulty | TestLoot
**Fails (1) :**
*  ` + time.Now().Format("02/01/06") + ` | TestReason
`
		msg, err := discord.GetPlayerHandler(context.Background(), interaction)
		mockPlayerUseCase.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, msg, exceptMsg)
	})

	t.Run("Player doesnt exist", func(t *testing.T) {
		t.Parallel()
		mockPlayerUseCase := mocks.NewPlayerUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: nil,
			PlayerUseCase:  mockPlayerUseCase,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		mockPlayerUseCase.On("ReadPlayer", mock.Anything, mock.Anything, "").
			Return(entity.Player{}, errors.New("player not found"))

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
					Name:     "guildops-player-create",
					TargetID: "mock",
					Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Name:  "name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "TestPlayer",
						},
					},
				},
			},
		}

		msg, err := discord.GetPlayerHandler(context.Background(), interaction)

		mockPlayerUseCase.AssertExpectations(t)
		assert.Error(t, err)
		assert.Equal(t, msg, "Error while getting player infos: player not found")
	})
}
