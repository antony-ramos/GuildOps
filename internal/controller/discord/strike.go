package discordhandler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/bwmarrin/discordgo"
)

var StrikeDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "guildops-strike-create",
		Description: "Generate a strike for a player",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Milowenn",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "ex: Retard de 5min",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-strike-list",
		Description: "list off strikes on a player",
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
		Name:        "guildops-strike-delete",
		Description: "Delete a strike",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "ex: 123456789",
				Required:    true,
			},
		},
	},
}

// InitStrike return lists of func which can be processed by discord bot.
// They are all strike related.
func (d Discord) InitStrike() map[string]func(
	ctx context.Context, interaction *discordgo.InteractionCreate) (string, error) {
	return map[string]func(ctx context.Context, interaction *discordgo.InteractionCreate) (string, error){
		"guildops-strike-create": d.StrikeOnPlayerHandler,
		"guildops-strike-delete": d.DeleteStrikeHandler,
		"guildops-strike-list":   d.ListStrikesOnPlayerHandler,
	}
}

// StrikeOnPlayerHandler call an usecase to create a strike
// and return a message to the user.
// It requires a player name and a reason field to be passed in the interaction.
func (d Discord) StrikeOnPlayerHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)

	ctx, span := otel.Tracer("Discord").Start(ctx, "Strike/StrikeOnPlayerHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	defer cancel()

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	name := optionMap["name"].StringValue()
	reason := optionMap["reason"].StringValue()
	span.SetAttributes(
		attribute.String("player", name),
		attribute.String("reason", reason),
	)

	err := d.CreateStrike(ctx, reason, name)
	if err != nil {
		msg := "Error while creating strike: " + HumanReadableError(err)
		return msg, fmt.Errorf("create strike call usecase: %w", err)
	}
	return "Strike created successfully", nil
}

// ListStrikesOnPlayerHandler call an usecase to get strikes on a player
// and return a message to the user.
// It requires a player name field to be passed in the interaction.
func (d Discord) ListStrikesOnPlayerHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Strike/ListStrikesOnPlayerHandler")
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
	span.SetAttributes(
		attribute.String("player", playerName),
	)

	strikes, err := d.ReadStrikes(ctx, playerName)
	if err != nil {
		msg := "Error while getting strikes on player: " + HumanReadableError(err)
		return msg, fmt.Errorf("database - ListStrikesOnPlayerHandler - r.ReadStrikes: %w", err)
	}

	msg := "Strikes of " + playerName + " (" + strconv.Itoa(len(strikes)) + ") :\n"
	for _, strike := range strikes {
		msg += "* " + strike.Date.Format("02/01/06") + " | " + strike.Reason + " | " + strconv.Itoa(strike.ID) + "\n"
	}
	return msg, nil
}

// DeleteStrikeHandler call an usecase to delete a strike
// and return a message to the user.
// It requires an id field to be passed in the interaction.
//
//nolint:dupl
func (d Discord) DeleteStrikeHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Strike/DeleteStrikeHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	idString := optionMap["id"].StringValue()
	span.SetAttributes(
		attribute.String("id", idString),
	)
	strikeID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		msg := "Error while deleting strike: " + HumanReadableError(err)
		return msg, fmt.Errorf("delete strike parse id: %w", err)
	}

	err = d.DeleteStrike(ctx, int(strikeID))
	if err != nil {
		msg := "Error while deleting strike: " + HumanReadableError(err)
		return msg, fmt.Errorf("delete strike usecase: %w", err)
	}

	return "Strike deleted successfully", nil
}
