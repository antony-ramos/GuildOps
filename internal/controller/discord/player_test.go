package discordhandler_test

import (
	"context"
	"errors"
	"testing"

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

		mockPlayerUseCase.On("ReadPlayer", mock.Anything, "playerone").
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

		mockPlayerUseCase.On("ReadPlayer", mock.Anything, "playerone").Return(entity.Player{
			Name: "playerone",
		}, nil)

		mockPlayerUseCase.On("LinkPlayer", mock.Anything, "playerone", mock.Anything).Return(nil)

		msg, err := discord.GenerateLinkPlayerMsg(context.Background(), "playerone", "playerone")

		assert.NoError(t, err)
		assert.Equal(t, "Vous êtes maintenant lié à ce joueur : \nName : **playerone**\nID : **0**\n", msg)
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

		mockPlayerUseCase.On("ReadPlayer", mock.Anything, "playerone").Return(entity.Player{
			Name: "playerone",
		}, nil)

		mockPlayerUseCase.On("LinkPlayer", mock.Anything, "playerone", mock.Anything).Return(errors.New("Link Failed"))

		msg, err := discord.GenerateLinkPlayerMsg(context.Background(), "playerone", "playerone")

		assert.Error(t, err)
		assert.Equal(t, "Error while linking player: Link Failed", msg)
		mockPlayerUseCase.AssertExpectations(t)
	})
}
