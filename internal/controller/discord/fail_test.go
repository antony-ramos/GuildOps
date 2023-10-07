package discordhandler_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/stretchr/testify/assert"

	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"
	"github.com/antony-ramos/guildops/internal/controller/discord/mocks"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

func TestDiscord_CreateFailHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
		}

		mockFailUseCase.On("CreateFail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
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

		msg, err := discord.CreateFailHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "Fail created successfully")
		mockFailUseCase.AssertExpectations(t)
	})

	t.Run("Wrong date format", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
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
							Value: "03-05-23",
						},
					},
				},
			},
		}

		msg, err := discord.CreateFailHandler(context.Background(), interaction)
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile("Error while creating fail: .*"), msg)
		mockFailUseCase.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
		}

		mockFailUseCase.On("CreateFail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("backend Error"))

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

		msg, err := discord.CreateFailHandler(context.Background(), interaction)
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile("Error while creating fail: .*"), msg)
		mockFailUseCase.AssertExpectations(t)
	})
}

func TestDiscord_ListFailsOnPlayerHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
		}

		mockFailUseCase.On("ListFailOnPLayer", mock.Anything, mock.Anything).
			Return([]entity.Fail{
				{
					ID: 1,
					Raid: &entity.Raid{
						Date: time.Now(),
					},
					Reason: "why not",
				},
				{
					ID: 1,
					Raid: &entity.Raid{
						Date: time.Now(),
					},
					Reason: "why not 2",
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

		msg, err := discord.ListFailsOnPlayerHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, "Fails of Milowenn (2) :\n* "+
			time.Now().Format("02/01/06")+" - why not - 1\n* "+
			time.Now().Format("02/01/06")+" - why not 2 - 1\n",
			msg)
		mockFailUseCase.AssertExpectations(t)
	})

	t.Run("Success with no fails", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
		}

		mockFailUseCase.On("ListFailOnPLayer", mock.Anything, mock.Anything).
			Return([]entity.Fail{}, nil)

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

		msg, err := discord.ListFailsOnPlayerHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, "No fails found for Milowenn", msg)
		mockFailUseCase.AssertExpectations(t)
	})
}

func TestDiscord_ListFailsOnRaidHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
		}

		mockFailUseCase.On("ListFailOnRaid", mock.Anything, mock.Anything).
			Return([]entity.Fail{
				{
					ID: 1,
					Raid: &entity.Raid{
						Date: time.Now(),
					},
					Player: &entity.Player{
						Name: "Paragon",
					},
					Reason: "why not",
				},
				{
					ID: 1,
					Raid: &entity.Raid{
						Date: time.Now(),
					},
					Player: &entity.Player{
						Name: "Milowenn",
					},
					Reason: "why not 2",
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

		msg, err := discord.ListFailsOnRaidHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "Fails for 29/09/23 (2) :\n* Paragon - why not\n* Milowenn - why not 2\n")
		mockFailUseCase.AssertExpectations(t)
	})

	t.Run("Success with no fails", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
		}

		mockFailUseCase.On("ListFailOnRaid", mock.Anything, mock.Anything).
			Return([]entity.Fail{}, nil)

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

		msg, err := discord.ListFailsOnRaidHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "No fails found for 29/09/23")
		mockFailUseCase.AssertExpectations(t)
	})
}

//nolint:dupl
func TestDiscord_DeleteFailOnPlayerHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
		}

		mockFailUseCase.On("DeleteFail", mock.Anything, mock.Anything).
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

		msg, err := discord.DeleteFailHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "Fail successfully deleted")
		mockFailUseCase.AssertExpectations(t)
	})

	t.Run("Wrong ID format", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
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
							Value: "ABC",
						},
					},
				},
			},
		}

		msg, err := discord.DeleteFailHandler(context.Background(), interaction)
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile("Error while deleting fail: .*"), msg)
		mockFailUseCase.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()
		mockFailUseCase := mocks.NewFailUseCase(t)

		discord := discordHandler.Discord{
			FailUseCase: mockFailUseCase,
		}

		mockFailUseCase.On("DeleteFail", mock.Anything, mock.Anything).
			Return(errors.New("Backend Error"))

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

		msg, err := discord.DeleteFailHandler(context.Background(), interaction)
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile("Error while deleting fail: .*"), msg)
		mockFailUseCase.AssertExpectations(t)
	})
}
