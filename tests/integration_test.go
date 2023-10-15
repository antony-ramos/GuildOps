//nolint:paralleltest,maintidx
package e2e_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"

	"github.com/antony-ramos/guildops/internal/usecase"
	"github.com/antony-ramos/guildops/internal/usecase/postgresbackend"
	"github.com/antony-ramos/guildops/pkg/logger"
	"github.com/antony-ramos/guildops/pkg/postgres"
	"github.com/bwmarrin/discordgo"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // used by migrator
	_ "github.com/golang-migrate/migrate/v4/source/file"       // used by migrator
	_ "github.com/jackc/pgx/v4/stdlib"                         // used by migrator
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DBName = "test_db"
	DBUser = "test_user"
	DBPass = "test_password"
)

var discord discordHandler.Discord

func guildOpsInfo(discordName string) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
				},
			},
			Data: discordgo.ApplicationCommandInteractionData{
				ID:       "mock",
				Name:     "guildops-player-info",
				TargetID: "mock",
				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
				Options:  []*discordgo.ApplicationCommandInteractionDataOption{},
			},
		},
	}
}

func init() {
	ctx := context.Background()
	url, err := createContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	pgHandler, err := postgres.New(
		ctx,
		url,
		postgres.MaxPoolSize(1),
		postgres.ConnAttempts(5),
		postgres.ConnTimeout(2*time.Second))
	if err != nil {
		logger.FromContext(ctx).Fatal(err.Error())
	}

	backend := postgresbackend.PG{Postgres: pgHandler}
	err = backend.Init(ctx, url, nil)
	if err != nil {
		return
	}

	auc := usecase.NewAbsenceUseCase(&backend)
	puc := usecase.NewPlayerUseCase(&backend)
	luc := usecase.NewLootUseCase(&backend)
	ruc := usecase.NewRaidUseCase(&backend)
	suc := usecase.NewStrikeUseCase(&backend)
	fuc := usecase.NewFailUseCase(&backend)

	discord = discordHandler.Discord{
		AbsenceUseCase: auc,
		PlayerUseCase:  puc,
		LootUseCase:    luc,
		RaidUseCase:    ruc,
		StrikeUseCase:  suc,
		FailUseCase:    fuc,
	}
}

func TestPlayer(t *testing.T) {
	name := "testPlayer"
	discordName := "testPlayerDiscord"

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
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
						Value: name,
					},
				},
			},
		},
	}

	t.Run(fmt.Sprintf("create user %s", name), func(t *testing.T) {
		msg, _ := discord.PlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("Player %s created successfully: ID 1", strings.ToLower(name)), msg)
	})
	t.Run("try to wrongly recreate it", func(t *testing.T) {
		msg, _ := discord.PlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("Player %s already exists", strings.ToLower(name)), msg)
	})

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
				},
			},
			Data: discordgo.ApplicationCommandInteractionData{
				ID:       "mock",
				Name:     "guildops-player-link",
				TargetID: "mock",
				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "name",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: name,
					},
				},
			},
		},
	}

	t.Run("link discord to previously created player", func(t *testing.T) {
		msg, _ := discord.LinkPlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("You are now linked to this player : \nName : **%s**\nDiscord Name : **%s**\n",
			strings.ToLower(name), strings.ToLower(discordName)), msg)
	})

	t.Run("try to wrongly link discord again", func(t *testing.T) {
		msg, _ := discord.LinkPlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("Error while linking player: discord account already linked to player name %s. "+
			"Contact Staff for modification", strings.ToLower(name)), msg)
	})

	t.Run("get info on player", func(t *testing.T) {
		msg, _ := discord.GetPlayerHandler(context.Background(), guildOpsInfo(discordName))
		assert.Equal(t, fmt.Sprintf("Name : **%s**\nID : **1**\nDiscord ID : **%s**\n",
			strings.ToLower(name), strings.ToLower(discordName)), msg)
	})

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
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
						Value: name,
					},
				},
			},
		},
	}

	t.Run("delete player", func(t *testing.T) {
		msg, _ := discord.PlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("Player %s deleted successfully", strings.ToLower(name)), msg)
	})

	t.Run("get info on deleted linked player", func(t *testing.T) {
		msg, _ := discord.GetPlayerHandler(context.Background(), guildOpsInfo(discordName))
		assert.Equal(t, fmt.Sprintf("Error while getting player infos: "+
			"didn't find a player linked to this discord user named %s", strings.ToLower(discordName)), msg)
	})

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: "test",
				},
			},
			Data: discordgo.ApplicationCommandInteractionData{
				ID:       "mock",
				Name:     "guildops-player-get",
				TargetID: "mock",
				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "name",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: "testplayer",
					},
				},
			},
		},
	}
	t.Run("get info on deleted player", func(t *testing.T) {
		msg, _ := discord.GetPlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("Error while getting player infos: player %s not found", strings.ToLower(name)), msg)
	})
}

func TestAbsence(t *testing.T) {
	name := "testAbsence"
	discordName := "testAbsenceDiscord"

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
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
						Value: name,
					},
				},
			},
		},
	}

	t.Run(fmt.Sprintf("create user %s", name), func(t *testing.T) {
		msg, _ := discord.PlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("Player %s created successfully: ID 3", strings.ToLower(name)), msg)
	})

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
				},
			},
			Data: discordgo.ApplicationCommandInteractionData{
				ID:       "mock",
				Name:     "guildops-player-link",
				TargetID: "mock",
				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "name",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: name,
					},
				},
			},
		},
	}
	t.Run("link discord to previously created player", func(t *testing.T) {
		msg, _ := discord.LinkPlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("You are now linked to this player :"+
			" \nName : **%s**\nDiscord Name : **%s**\n", strings.ToLower(name), strings.ToLower(discordName)), msg)
	})

	// Create a range of raid for 5 days in september 2030
	for index := 1; index <= 5; index++ {
		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionApplicationCommand,
				Member: &discordgo.Member{
					User: &discordgo.User{
						Username: discordName,
					},
				},
				Data: discordgo.ApplicationCommandInteractionData{
					ID:       "mock",
					Name:     "guildops-raid-create",
					TargetID: "mock",
					Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "name",
							Value: "raidName",
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "difficulty",
							Value: "normal",
						},
						{
							Name:  "date",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: fmt.Sprintf("%02d/09/30", index),
						},
					},
				},
			},
		}
		t.Run(fmt.Sprintf("Create Raid on %02d/09/30", index), func(t *testing.T) {
			msg, _ := discord.CreateRaidHandler(context.Background(), interaction)
			assert.Equal(t, fmt.Sprintf("Raid successfully created with ID %d", index), msg)
		})
	}

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
				},
			},
			Data: discordgo.ApplicationCommandInteractionData{
				ID:       "mock",
				Name:     "guildops-absence-create",
				TargetID: "mock",
				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "from",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: "01/09/30",
					},
					{
						Name:  "to",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: "03/09/30",
					},
				},
			},
		},
	}

	t.Run("Create Absence from 01/09/30 to 03/09/30", func(t *testing.T) {
		msg, _ := discord.AbsenceHandler(context.Background(), interaction)
		assert.Equal(t, "Absence(s) created for :\n* Sun 01/09/30\n* Mon 02/09/30\n* Tue 03/09/30\n", msg)
	})

	t.Run("Check if absences appears in player info", func(t *testing.T) {
		msg, _ := discord.GetPlayerHandler(context.Background(), guildOpsInfo(discordName))
		assert.Equal(t, fmt.Sprintf("Name : **%s**\nID : **3**\n"+
			"Discord ID : **%s**\n**Absences (3) :**\n*  01/09/30 | normal | raidname\n"+
			"*  02/09/30 | normal | raidname\n"+
			"*  03/09/30 | normal | raidname\n", strings.ToLower(name), strings.ToLower(discordName)), msg)
	})

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
				},
			},
			Data: discordgo.ApplicationCommandInteractionData{
				ID:       "mock",
				Name:     "guildops-raid-delete",
				TargetID: "mock",
				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "id",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: "2",
					},
				},
			},
		},
	}

	t.Run("Delete Raid on 02/09/30", func(t *testing.T) {
		msg, _ := discord.DeleteRaidHandler(context.Background(), interaction)
		assert.Equal(t, "Raid with ID 2 successfully deleted", msg)
	})

	t.Run("Check if absences appears in player info for deleted raid", func(t *testing.T) {
		msg, _ := discord.GetPlayerHandler(context.Background(), guildOpsInfo(discordName))
		assert.Equal(t, fmt.Sprintf("Name : **%s**\nID : **3**\n"+
			"Discord ID : **%s**\n**Absences (2) :**\n"+
			"*  01/09/30 | normal | raidname\n*  03/09/30 | normal | raidname\n",
			strings.ToLower(name), strings.ToLower(discordName)), msg)
	})

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
				},
			},
			Data: discordgo.ApplicationCommandInteractionData{
				ID:       "mock",
				Name:     "guildops-absence-delete",
				TargetID: "mock",
				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "from",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: "03/09/30",
					},
					{
						Name:  "to",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: "04/09/30",
					},
				},
			},
		},
	}

	t.Run("Delete Absence from 03/09/30 to 04/09/30", func(t *testing.T) {
		msg, _ := discord.AbsenceHandler(context.Background(), interaction)
		assert.Equal(t, "Absence(s) deleted for :\n* Tue 03/09/30\n", msg)
	})

	t.Run("Check if absences appears in player info for deleted absences", func(t *testing.T) {
		msg, _ := discord.GetPlayerHandler(context.Background(), guildOpsInfo(discordName))
		assert.Equal(t, fmt.Sprintf("Name : **%s**\nID : **3**\nDiscord ID : **%s**\n**Absences (1) :**\n"+
			"*  01/09/30 | normal | raidname\n", strings.ToLower(name), strings.ToLower(discordName)), msg)
	})

	for index := 1; index <= 3; index++ {
		if index == 2 {
			continue
		}
		interaction = &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionApplicationCommand,
				Member: &discordgo.Member{
					User: &discordgo.User{
						Username: discordName,
					},
				},
				Data: discordgo.ApplicationCommandInteractionData{
					ID:       "mock",
					Name:     "guildops-raid-delete",
					TargetID: "mock",
					Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Name:  "id",
							Type:  discordgo.ApplicationCommandOptionString,
							Value: fmt.Sprintf("%d", index),
						},
					},
				},
			},
		}
		t.Run(fmt.Sprintf("Delete Raid on %02d/09/30", index), func(t *testing.T) {
			msg, _ := discord.DeleteRaidHandler(context.Background(), interaction)
			assert.Equal(t, fmt.Sprintf("Raid with ID %d successfully deleted", index), msg)
		})
	}

	t.Run("Check if absences appears after remove all raids", func(t *testing.T) {
		msg, _ := discord.GetPlayerHandler(context.Background(), guildOpsInfo(discordName))
		assert.Equal(t, fmt.Sprintf("Name : **%s**\nID : **3**\nDiscord ID : **%s**\n",
			strings.ToLower(name), strings.ToLower(discordName)), msg)
	})

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
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
						Value: name,
					},
				},
			},
		},
	}

	t.Run("delete player and finish this test", func(t *testing.T) {
		msg, _ := discord.PlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("Player %s deleted successfully", strings.ToLower(name)), msg)
	})
}

func TestStrike(t *testing.T) {
	name := "testStrike"
	discordName := "testStrikeDiscord"

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
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
						Value: name,
					},
				},
			},
		},
	}

	t.Run(fmt.Sprintf("create user %s", name), func(t *testing.T) {
		msg, _ := discord.PlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("Player %s created successfully: ID 4", strings.ToLower(name)), msg)
	})

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
				},
			},
			Data: discordgo.ApplicationCommandInteractionData{
				ID:       "mock",
				Name:     "guildops-player-link",
				TargetID: "mock",
				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "name",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: name,
					},
				},
			},
		},
	}

	t.Run("link discord to previously created player", func(t *testing.T) {
		msg, _ := discord.LinkPlayerHandler(context.Background(), interaction)
		assert.Equal(t,
			fmt.Sprintf("You are now linked to this player : \nName : **%s**\nDiscord Name : **%s**\n",
				strings.ToLower(name), strings.ToLower(discordName)), msg)
	})

	// Create two strike

	for i := 1; i <= 2; i++ {
		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionApplicationCommand,
				Member: &discordgo.Member{
					User: &discordgo.User{
						Username: discordName,
					},
				},
				Data: discordgo.ApplicationCommandInteractionData{
					ID:       "mock",
					Name:     "guildops-strike-create",
					TargetID: "mock",
					Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "name",
							Value: name,
						},
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "reason",
							Value: "testReason",
						},
					},
				},
			},
		}
		t.Run("create a strike", func(t *testing.T) {
			msg, _ := discord.StrikeOnPlayerHandler(context.Background(), interaction)
			assert.Equal(t, "Strike created successfully", msg)
		})
	}

	// Get PLayer info
	t.Run("check if strikes are showed in player info", func(t *testing.T) {
		msg, _ := discord.GetPlayerHandler(context.Background(), guildOpsInfo(discordName))
		assert.Equal(t, fmt.Sprintf("Name : **%s**\nID : **4**\nDiscord ID : **%s**\n**Strikes (2) :**\n"+
			"*  "+time.Now().Format("02/01/06")+" | testReason | DF/S2 | 1\n*  "+
			time.Now().Format("02/01/06")+" | testReason | DF/S2 | 2\n",
			strings.ToLower(name), strings.ToLower(discordName)), msg)
	})

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
				},
			},
			Data: discordgo.ApplicationCommandInteractionData{
				ID:       "mock",
				Name:     "guildops-strike-list",
				TargetID: "mock",
				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Type:  discordgo.ApplicationCommandOptionString,
						Name:  "name",
						Value: name,
					},
				},
			},
		},
	}
	t.Run("use command to show all strikes on player", func(t *testing.T) {
		msg, _ := discord.ListStrikesOnPlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("Strikes of %s (2) :\n"+
			"* "+time.Now().Format("02/01/06")+" | testReason | 1\n* "+time.Now().Format("02/01/06")+" | testReason | 2\n",
			strings.ToLower(name)), msg)
	})

	// Delete a strike
	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
				},
			},
			Data: discordgo.ApplicationCommandInteractionData{
				ID:       "mock",
				Name:     "guildops-strike-delete",
				TargetID: "mock",
				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Type:  discordgo.ApplicationCommandOptionString,
						Name:  "id",
						Value: "1",
					},
				},
			},
		},
	}
	t.Run("delete a strike", func(t *testing.T) {
		msg, _ := discord.DeleteStrikeHandler(context.Background(), interaction)
		assert.Equal(t, "Strike deleted successfully", msg)
	})

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
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
						Value: name,
					},
				},
			},
		},
	}

	// Get PLayer info
	t.Run("show if deleted strike is visible in player info", func(t *testing.T) {
		msg, _ := discord.GetPlayerHandler(context.Background(), guildOpsInfo(discordName))
		assert.Equal(t, fmt.Sprintf("Name : **%s**\nID : **4**\nDiscord ID : **%s**\n**Strikes (1) :**\n"+
			"*  "+time.Now().Format("02/01/06")+" | testReason | DF/S2 | 2\n",
			strings.ToLower(name), strings.ToLower(discordName)), msg)
	})

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
				},
			},
			Data: discordgo.ApplicationCommandInteractionData{
				ID:       "mock",
				Name:     "guildops-strike-list",
				TargetID: "mock",
				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{},
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Type:  discordgo.ApplicationCommandOptionString,
						Name:  "name",
						Value: name,
					},
				},
			},
		},
	}
	t.Run("check if deleted strike is visible with list strike on player", func(t *testing.T) {
		msg, _ := discord.ListStrikesOnPlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("Strikes of %s (1) :\n* "+
			time.Now().Format("02/01/06")+" | testReason | 2\n", strings.ToLower(name)), msg)
	})

	interaction = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: discordName,
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
						Value: name,
					},
				},
			},
		},
	}

	t.Run("delete player and finish test", func(t *testing.T) {
		msg, _ := discord.PlayerHandler(context.Background(), interaction)
		assert.Equal(t, fmt.Sprintf("Player %s deleted successfully", strings.ToLower(name)), msg)
	})
}

func createContainer(ctx context.Context) (string, error) {
	env := map[string]string{
		"POSTGRES_PASSWORD": DBPass,
		"POSTGRES_USER":     DBUser,
		"POSTGRES_DB":       DBName,
	}
	port := "5432"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	openPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return "", fmt.Errorf("failed to get container external port: %w", err)
	}

	log.Println("postgres container ready and running at port: ", openPort.Port())
	return fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		DBUser, DBName, DBPass, "localhost", openPort.Port()), nil
}
