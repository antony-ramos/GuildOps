package discordhandler

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/bwmarrin/discordgo"
)

var PlayerDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "guildops-player-create",
		Description: "Créer un joueur",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Milowenn",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-player-delete",
		Description: "Supprimer un joueur",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Milowenn",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-player-get",
		Description: "Infos sur le joueur",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Milowenn",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-player-link",
		Description: "link your discord account to your player name",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: milowenn",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-player-info",
		Description: "Show info about yourself",
	},
}

func (d Discord) InitPlayer() map[string]func(
	ctx context.Context, interaction *discordgo.InteractionCreate) (string, error) {
	return map[string]func(ctx context.Context, interaction *discordgo.InteractionCreate) (string, error){
		"guildops-player-create": d.PlayerHandler,
		"guildops-player-delete": d.PlayerHandler,
		"guildops-player-get":    d.GetPlayerHandler,
		"guildops-player-link":   d.LinkPlayerHandler,
		"guildops-player-info":   d.GetPlayerHandler,
	}
}

// PlayerHandler call an usecase to create or delete a player
// and return a message to the user.
// It requires a player name field to be passed in the interaction.
func (d Discord) PlayerHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Player/PlayerHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	name := optionMap["name"].StringValue()
	span.SetAttributes(
		attribute.String("player", name),
	)

	if interaction.ApplicationCommandData().Name == "guildops-player-create" {
		playerID, err := d.CreatePlayer(ctx, name)
		if err != nil {
			alreadyExists := regexp.MustCompile(".*player already exists.*")
			if alreadyExists.MatchString(err.Error()) {
				return "Player " + strings.ToLower(name) + " already exists", err
			}
			msg := "Error while creating player: " + HumanReadableError(err)
			return msg, fmt.Errorf("call create player usecase : %w", err)
		}
		return "Player " + strings.ToLower(name) + " created successfully: ID " + strconv.Itoa(playerID), nil
	}

	if interaction.ApplicationCommandData().Name == "guildops-player-delete" {
		err := d.DeletePlayer(ctx, name)
		if err != nil {
			msg := "Error while deleting player: " + HumanReadableError(err)
			return msg, fmt.Errorf("call delete player usecase: %w", err)
		} else {
			return "Player " + strings.ToLower(name) + " deleted successfully", nil
		}
	}

	return "error while handling player command", nil
}

// GetPlayerHandler call an usecase to get player infos
// and return a message to the user.
// It requires a player name field to be passed in the interaction for admin
// Or it will catch infos about the user who called the command.
func (d Discord) GetPlayerHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Player/GetPlayerHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	player, err := func() (entity.Player, error) {
		switch interaction.ApplicationCommandData().Name {
		case "guildops-player-info":
			return d.ReadPlayer(ctx, "", interaction.Member.User.Username)
		default:
			name := optionMap["name"].StringValue()
			span.SetAttributes(
				attribute.String("player", name))
			return d.ReadPlayer(ctx, name, "")
		}
	}()
	if err != nil {
		msg := "Error while getting player infos: " + HumanReadableError(err)
		return msg, fmt.Errorf("call read player usecase : %w", err)
	}

	msg := "Name : **" + player.Name + "**\n"
	msg += "ID : **" + strconv.Itoa(player.ID) + "**\n"
	if player.DiscordName != "" {
		msg += "Discord ID : **" + player.DiscordName + "**\n"
	}

	lootCounter := make(map[string]int)
	for _, loot := range player.Loots {
		lootCounter[loot.Raid.Difficulty]++
	}
	if len(lootCounter) > 0 {
		msg += "**Loots Count:**\n"
		for difficulty, count := range lootCounter {
			msg += "*  " + difficulty + " | " + strconv.Itoa(count) + " loots \n"
		}
	}

	if len(player.Strikes) > 0 {
		msg += "**Strikes (" + strconv.Itoa(len(player.Strikes)) + ") :**\n"
		for _, strike := range player.Strikes {
			msg += "*  " + strike.Date.Format("02/01/06") +
				" | " + strike.Reason + " | " + strike.Season + " | " + strconv.Itoa(strike.ID) + "\n"
		}
	}

	if len(player.MissedRaids) > 0 {
		msg += "**Absences (" + strconv.Itoa(len(player.MissedRaids)) + ") :**\n"
		for _, raid := range player.MissedRaids {
			msg += "*  " + raid.Date.Format("02/01/06") +
				" | " + raid.Difficulty +
				" | " + raid.Name + "\n"
		}
	}

	if len(player.Loots) > 0 {
		msg += "**Loots (" + strconv.Itoa(len(player.Loots)) + ") :**\n"
		for _, loot := range player.Loots {
			msg += "*  " + loot.Raid.Date.Format("02/01/06") +
				" | " + loot.Raid.Difficulty +
				" | " + loot.Name + "\n"
		}
	}

	if len(player.Fails) > 0 {
		msg += "**Fails (" + strconv.Itoa(len(player.Fails)) + ") :**\n"
		for _, fail := range player.Fails {
			msg += "*  " + fail.Raid.Date.Format("02/01/06") +
				" | " + fail.Reason + "\n"
		}
	}

	return msg, nil
}

// LinkPlayerHandler call an usecase to link a discord account to a player name
// and return a message to the user.
// It requires a player name field to be passed in the interaction.
func (d Discord) LinkPlayerHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Player/LinkPlayerHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	playerName := optionMap["name"].StringValue()
	discordName := interaction.Member.User.Username
	span.SetAttributes(
		attribute.String("player", playerName),
		attribute.String("discord_name", discordName),
	)

	err := d.LinkPlayer(ctx, playerName, discordName)
	if err != nil {
		msg := "Error while linking player: " + HumanReadableError(err)
		return msg, fmt.Errorf("call link player usecase : %w", err)
	}

	msg := "You are now linked to this player : \n"
	msg += "Name : **" + strings.ToLower(playerName) + "**\n"
	msg += "Discord Name : **" + strings.ToLower(discordName) + "**\n"

	return msg, nil
}
