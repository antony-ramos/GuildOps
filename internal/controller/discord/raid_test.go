package discordhandler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
		}

		mockRaidUseCase.On("CreateRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(entity.Raid{}, nil)

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
							Value: "random raid",
						},
						{
							Name:  "date",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "03/05/23",
						},
						{
							Name:  "difficulty",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: "Heroic",
						},
					},
				},
			},
		}

		msg, err := discord.CreateRaidHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "Raid successfully created with ID 0")
		mockRaidUseCase.AssertExpectations(t)
	})
}

//nolint:dupl
func TestDiscord_DeleteRaidHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockRaidUseCase := mocks.NewRaidUseCase(t)

		discord := discordHandler.Discord{
			RaidUseCase: mockRaidUseCase,
		}

		mockRaidUseCase.On("DeleteRaid", mock.Anything, mock.Anything).
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

		msg, err := discord.DeleteRaidHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "Raid with ID 1 successfully deleted")
		mockRaidUseCase.AssertExpectations(t)
	})
}

func TestDiscord_ListRaidHandler(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockRaidUseCase := mocks.NewRaidUseCase(t)

		discord := discordHandler.Discord{
			RaidUseCase: mockRaidUseCase,
		}

		mockRaidUseCase.On("ReadRaid", mock.Anything, time.Date(2030, time.September, 5, 0, 0, 0, 0, time.UTC)).
			Return(entity.Raid{
				ID:         1,
				Name:       "random raid",
				Difficulty: "Heroic",
				Date:       time.Date(2030, time.September, 5, 0, 0, 0, 0, time.UTC),
			}, nil).Once()
		mockRaidUseCase.On("ReadRaid", mock.Anything, time.Date(2030, time.September, 30, 0, 0, 0, 0, time.UTC)).
			Return(entity.Raid{
				ID:         1,
				Name:       "random raid",
				Difficulty: "Heroic",
				Date:       time.Date(2030, time.September, 30, 0, 0, 0, 0, time.UTC),
			}, nil).Once()
		mockRaidUseCase.On("ReadRaid", mock.Anything, mock.Anything).Return(entity.Raid{}, errors.New("Not found"))

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
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "from",
							Value: "05/09/30",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "to",
							Value: "05/10/30",
						},
					},
				},
			},
		}

		msg, err := discord.ListRaidHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "Raid List:\n* random raid Thu 05/09/30 Heroic 1\n* random raid Mon 30/09/30 Heroic 1\n")
		mockRaidUseCase.AssertExpectations(t)
	})
}

//nolint:maintidx
func TestDiscord_GenerateRaidsOnRangeHandler(t *testing.T) {
	t.Parallel()

	t.Run("Context is done", func(t *testing.T) {
		t.Parallel()
		mockRaidUseCase := mocks.NewRaidUseCase(t)

		discord := discordHandler.Discord{
			RaidUseCase: mockRaidUseCase,
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

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
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "from",
							Value: "05/09/30",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "to",
							Value: "05/10/30",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "difficulty",
							Value: "Heroic",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "weekday",
							Value: "Monday",
						},
					},
				},
			},
		}

		msg, err := discord.GenerateRaidsOnRangeHandler(ctx, interaction)
		assert.EqualError(t, err, "create multiple raids wait goroutines: context canceled")
		assert.Equal(t, msg, "error while creating multiple raids: context canceled")
		mockRaidUseCase.AssertExpectations(t)
	})

	t.Run("Create a range of raids", func(t *testing.T) {
		t.Parallel()
		mockRaidUseCase := mocks.NewRaidUseCase(t)

		discord := discordHandler.Discord{
			RaidUseCase: mockRaidUseCase,
		}

		mockRaidUseCase.On("CreateRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(entity.Raid{}, nil).Times(8)

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
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "from",
							Value: "05/09/30",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "to",
							Value: "05/10/30",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "difficulty",
							Value: "Heroic",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "weekdays",
							Value: "Monday, Tuesday",
						},
					},
				},
			},
		}

		msg, err := discord.GenerateRaidsOnRangeHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "Raid List:\n*  Mon 01/01/01  0\n*  Mon 01/01/01  0\n*  Mon 01/01/01"+
			"  0\n*  Mon 01/01/01  0\n*  Mon 01/01/01  0\n*  "+
			"Mon 01/01/01  0\n*  Mon 01/01/01  0\n*  Mon 01/01/01  0\n")
		mockRaidUseCase.AssertExpectations(t)
	})

	t.Run("Create a single raid", func(t *testing.T) {
		t.Parallel()
		mockRaidUseCase := mocks.NewRaidUseCase(t)

		discord := discordHandler.Discord{
			RaidUseCase: mockRaidUseCase,
		}

		mockRaidUseCase.On("CreateRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(entity.Raid{}, nil).Times(1)

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
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "from",
							Value: "05/09/30",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "difficulty",
							Value: "Heroic",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "weekdays",
							Value: "Thursday",
						},
					},
				},
			},
		}

		msg, err := discord.GenerateRaidsOnRangeHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "Raid List:\n*  Mon 01/01/01  0\n")
		mockRaidUseCase.AssertExpectations(t)
	})

	t.Run("Create a no raid", func(t *testing.T) {
		t.Parallel()
		mockRaidUseCase := mocks.NewRaidUseCase(t)

		discord := discordHandler.Discord{
			RaidUseCase: mockRaidUseCase,
		}

		mockRaidUseCase.On("CreateRaid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(entity.Raid{}, errors.New("already exists"))

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
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "from",
							Value: "05/09/30",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "difficulty",
							Value: "Heroic",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "weekdays",
							Value: "Thursday",
						},
					},
				},
			},
		}

		msg, err := discord.GenerateRaidsOnRangeHandler(context.Background(), interaction)
		assert.NoError(t, err)
		assert.Equal(t, msg, "no raid created")
		mockRaidUseCase.AssertExpectations(t)
	})

	t.Run("check params: to is before from", func(t *testing.T) {
		t.Parallel()
		mockRaidUseCase := mocks.NewRaidUseCase(t)

		discord := discordHandler.Discord{
			RaidUseCase: mockRaidUseCase,
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
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "from",
							Value: "05/10/30",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "to",
							Value: "05/09/30",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "difficulty",
							Value: "Heroic",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "weekdays",
							Value: "Monday, Tuesday",
						},
					},
				},
			},
		}

		msg, err := discord.GenerateRaidsOnRangeHandler(context.Background(), interaction)
		assert.Error(t, err)
		assert.Equal(t, msg, "error while creating multiple raids: endDate is before startDate")
		mockRaidUseCase.AssertExpectations(t)
	})

	t.Run("check params: difficulty is incorrect", func(t *testing.T) {
		t.Parallel()
		mockRaidUseCase := mocks.NewRaidUseCase(t)

		discord := discordHandler.Discord{
			RaidUseCase: mockRaidUseCase,
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
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "from",
							Value: "05/09/30",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "difficulty",
							Value: "IncorrectDifficulty",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "weekdays",
							Value: "Thursday",
						},
					},
				},
			},
		}

		msg, err := discord.GenerateRaidsOnRangeHandler(context.Background(), interaction)
		assert.Error(t, err)
		assert.Equal(t, msg, "difficulty must be one of: Normal, Heroic, Mythic")
		mockRaidUseCase.AssertExpectations(t)
	})

	t.Run("check params: weekday is incorrect", func(t *testing.T) {
		t.Parallel()
		mockRaidUseCase := mocks.NewRaidUseCase(t)

		discord := discordHandler.Discord{
			RaidUseCase: mockRaidUseCase,
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
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "from",
							Value: "05/09/30",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "difficulty",
							Value: "Heroic",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "weekdays",
							Value: "IncorrectWeekday",
						},
					},
				},
			},
		}

		msg, err := discord.GenerateRaidsOnRangeHandler(context.Background(), interaction)
		assert.Error(t, err)
		assert.Equal(t, msg, "week days must be one of: Monday, Tuesday, "+
			"Wednesday, Thursday, Friday, Saturday, Sunday")
		mockRaidUseCase.AssertExpectations(t)
	})
}
