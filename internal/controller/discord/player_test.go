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

func TestDiscord_GenerateLinkPlayerMsg(t *testing.T) {
	t.Parallel()

	t.Run("context is done", func(t *testing.T) {
		t.Parallel()

		mockPlayerUseCase := mocks.NewPlayerUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: nil,
			PlayerUseCase:  mockPlayerUseCase,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		msg, err := discord.GenerateLinkPlayerMsg(ctx, "playerone", "playerone")

		assert.Error(t, err)
		assert.Equal(t, "Error because request took too much time to complete", msg)
		mockPlayerUseCase.AssertExpectations(t)
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

		mockPlayerUseCase.On("ReadPlayer", mock.Anything, "playerone", "").
			Return(entity.Player{}, errors.New("Player doesnt exist"))

		msg, err := discord.GenerateLinkPlayerMsg(context.Background(), "playerone", "playerone")

		assert.Error(t, err)
		assert.Equal(t, "Error while reading player: Player doesnt exist", msg)
		mockPlayerUseCase.AssertExpectations(t)
	})

	t.Run("Link Successfully", func(t *testing.T) {
		t.Parallel()

		mockPlayerUseCase := mocks.NewPlayerUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: nil,
			PlayerUseCase:  mockPlayerUseCase,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		mockPlayerUseCase.On("ReadPlayer", mock.Anything, "playerone", "").Return(entity.Player{
			Name: "playerone",
		}, nil)

		mockPlayerUseCase.On("LinkPlayer", mock.Anything, "playerone", mock.Anything).Return(nil)

		msg, err := discord.GenerateLinkPlayerMsg(context.Background(), "playerone", "playerone")

		assert.NoError(t, err)
		assert.Equal(t, "You are now linked to this player : \nName : **playerone**\nID : **0**\n", msg)
		mockPlayerUseCase.AssertExpectations(t)
	})

	t.Run("Link Failed", func(t *testing.T) {
		t.Parallel()

		mockPlayerUseCase := mocks.NewPlayerUseCase(t)

		discord := discordHandler.Discord{
			AbsenceUseCase: nil,
			PlayerUseCase:  mockPlayerUseCase,
			StrikeUseCase:  nil,
			LootUseCase:    nil,
			RaidUseCase:    nil,
		}

		mockPlayerUseCase.On("ReadPlayer", mock.Anything, "playerone", "").Return(entity.Player{
			Name: "playerone",
		}, nil)

		mockPlayerUseCase.On("LinkPlayer", mock.Anything, "playerone", mock.Anything).Return(errors.New("Link Failed"))

		msg, err := discord.GenerateLinkPlayerMsg(context.Background(), "playerone", "playerone")

		assert.Error(t, err)
		assert.Equal(t, "Error while linking player: Link Failed", msg)
		mockPlayerUseCase.AssertExpectations(t)
	})
}

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
			Fake:           true,
		}

		mockPlayerUseCase.On("CreatePlayer", mock.Anything, mock.Anything).
			Return(1, nil)

		session := &discordgo.Session{StateEnabled: true, State: discordgo.NewState()}
		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionApplicationCommand,
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

		err := discord.PlayerHandler(context.Background(), session, interaction)
		if err != nil {
			return
		}
		mockPlayerUseCase.AssertExpectations(t)
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
			Fake:           true,
		}

		mockPlayerUseCase.On("DeletePlayer", mock.Anything, mock.Anything).
			Return(nil)

		session := &discordgo.Session{StateEnabled: true, State: discordgo.NewState()}
		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionApplicationCommand,
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

		err := discord.PlayerHandler(context.Background(), session, interaction)
		if err != nil {
			return
		}
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
			Fake:           true,
		}

		player := entity.Player{
			Name: "TestPlayer",
			ID:   1,
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

		mockPlayerUseCase.On("ReadPlayer", mock.Anything, mock.Anything, "").
			Return(player, nil)

		session := &discordgo.Session{StateEnabled: true, State: discordgo.NewState()}
		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionApplicationCommand,
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

		err := discord.GetPlayerHandler(context.Background(), session, interaction)
		if err != nil {
			return
		}
		mockPlayerUseCase.AssertExpectations(t)
		assert.NoError(t, err)
	})
}
