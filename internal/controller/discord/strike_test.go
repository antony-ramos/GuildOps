package discordhandler_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"
	"github.com/antony-ramos/guildops/internal/controller/discord/mocks"
	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDiscord_InitStrike(t *testing.T) {
	t.Parallel()

	t.Run("Is not nil", func(t *testing.T) {
		t.Parallel()
		discord := discordHandler.Discord{}
		assert.NotNil(t, discord.InitStrike())
	})
}

func TestDiscord_CreateStrikeOnPlayerHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockStrikeUseCase := mocks.NewStrikeUseCase(t)

		discord := discordHandler.Discord{
			StrikeUseCase: mockStrikeUseCase,
		}

		mockStrikeUseCase.On("CreateStrike", mock.Anything, mock.Anything, mock.Anything).
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
							Name:  "name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "milowenn",
						},
						{
							Name:  "reason",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "test strike",
						},
					},
				},
			},
		}

		msg, err := discord.StrikeOnPlayerHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "Strike created successfully")
		mockStrikeUseCase.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()
		mockStrikeUseCase := mocks.NewStrikeUseCase(t)

		discord := discordHandler.Discord{
			StrikeUseCase: mockStrikeUseCase,
		}

		mockStrikeUseCase.On("CreateStrike", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error"))

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
							Name:  "name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "milowenn",
						},
						{
							Name:  "reason",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "test strike",
						},
					},
				},
			},
		}

		msg, err := discord.StrikeOnPlayerHandler(context.Background(), interaction)
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile("Error while creating strike: .*"), msg)
		mockStrikeUseCase.AssertExpectations(t)
	})
}

//nolint:dupl
func TestDiscord_DeleteStrikeHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockStrikeUseCase := mocks.NewStrikeUseCase(t)

		discord := discordHandler.Discord{
			StrikeUseCase: mockStrikeUseCase,
		}

		mockStrikeUseCase.On("DeleteStrike", mock.Anything, mock.Anything).
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
							Value: "123456789",
						},
					},
				},
			},
		}

		msg, err := discord.DeleteStrikeHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "Strike deleted successfully")
		mockStrikeUseCase.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()
		mockStrikeUseCase := mocks.NewStrikeUseCase(t)

		discord := discordHandler.Discord{
			StrikeUseCase: mockStrikeUseCase,
		}

		mockStrikeUseCase.On("DeleteStrike", mock.Anything, mock.Anything).
			Return(errors.New("error"))

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
							Value: "123456789",
						},
					},
				},
			},
		}

		msg, err := discord.DeleteStrikeHandler(context.Background(), interaction)
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile("Error while deleting strike: .*"), msg)
		mockStrikeUseCase.AssertExpectations(t)
	})

	t.Run("Incorrect ID format", func(t *testing.T) {
		t.Parallel()
		mockStrikeUseCase := mocks.NewStrikeUseCase(t)

		discord := discordHandler.Discord{
			StrikeUseCase: mockStrikeUseCase,
		}

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
							Value: "i am not a number",
						},
					},
				},
			},
		}

		msg, err := discord.DeleteStrikeHandler(context.Background(), interaction)
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile("Error while deleting strike: .*"), msg)
		mockStrikeUseCase.AssertExpectations(t)
	})
}

func TestDiscord_ListStrikesOnPlayerHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockStrikeUseCase := mocks.NewStrikeUseCase(t)

		discord := discordHandler.Discord{
			StrikeUseCase: mockStrikeUseCase,
		}

		mockStrikeUseCase.On("ReadStrikes", mock.Anything, mock.Anything).
			Return([]entity.Strike{
				{
					ID:     1,
					Reason: "test strike",
					Player: &entity.Player{
						Name: "Milowenn",
					},
				},
				{
					ID:     2,
					Reason: "test strike 2",
					Player: &entity.Player{
						Name: "Milowenn",
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
							Name:  "name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "Milowenn",
						},
					},
				},
			},
		}

		msg, err := discord.ListStrikesOnPlayerHandler(context.Background(), interaction)
		assert.NoError(t, err)

		wantedMsg := `Strikes of milowenn (2) :
* 01/01/01 | test strike | 1
* 01/01/01 | test strike 2 | 2
`

		assert.Equal(t, msg, wantedMsg)
		mockStrikeUseCase.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()
		mockStrikeUseCase := mocks.NewStrikeUseCase(t)

		discord := discordHandler.Discord{
			StrikeUseCase: mockStrikeUseCase,
		}

		mockStrikeUseCase.On("ReadStrikes", mock.Anything, mock.Anything).
			Return(nil, errors.New("error"))

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
							Name:  "name",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "Milowenn",
						},
					},
				},
			},
		}

		msg, err := discord.ListStrikesOnPlayerHandler(context.Background(), interaction)
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile("Error while getting strikes on player: .*"), msg)
		mockStrikeUseCase.AssertExpectations(t)
	})
}
