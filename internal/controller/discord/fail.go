package discordhandler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antony-ramos/guildops/internal/entity"

	"github.com/bwmarrin/discordgo"
)

var FailDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "guildops-fail-create", // Tested
		Description: "Générer un Fail sur un joueur",
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
				Description: "ex: Erreur P3 Sarkareth",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "ex: 03/05/2023",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-fail-list-player",
		Description: "Lister les fails sur un joueur",
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
		Name:        "guildops-fail-list-raid",
		Description: "Lister les fails sur un raid",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "ex: 03/05/2021",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-fail-delete",
		Description: "Supprimer un fail via son ID (ListFails pour l'avoir)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "ex: qzdq-qzdqz-qddq",
				Required:    true,
			},
		},
	},
}

func (d Discord) InitFail() map[string]func(
	ctx context.Context, interaction *discordgo.InteractionCreate) (string, error) {
	return map[string]func(ctx context.Context, interaction *discordgo.InteractionCreate) (string, error){
		"guildops-fail-create":      d.CreateFailHandler,
		"guildops-fail-delete":      d.DeleteFailHandler,
		"guildops-fail-list-player": d.ListFailsOnPlayerHandler,
		"guildops-fail-list-raid":   d.ListFailsOnRaidHandler,
	}
}

func (d Discord) CreateFailHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Fail/CreateFailHandler")
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
	reason := optionMap["reason"].StringValue()
	span.SetAttributes(
		attribute.String("player", name),
		attribute.String("reason", reason),
		attribute.String("date", optionMap["date"].StringValue()),
	)

	raidDate, err := ParseDate(optionMap["date"].StringValue(), "")
	if err != nil {
		msg := "Error while creating fail: " + HumanReadableError(err)
		return msg, fmt.Errorf("create fail parse date: %w", err)
	}

	err = d.CreateFail(ctx, reason, raidDate[0], name)
	if err != nil {
		msg := "Error while creating fail: " + HumanReadableError(err)
		return msg, fmt.Errorf("create fail usecase: %w", err)
	}
	return "Fail created successfully", nil
}

func (d Discord) ListFailsOnPlayerHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Fail/ListFailsOnPlayerHandler")
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
		attribute.String("player_name", playerName),
	)

	fails, err := d.ListFailOnPLayer(ctx, playerName)
	if err != nil {
		return "Fail to list fails on player", fmt.Errorf("error while listing fails on player: %w", err)
	}

	if len(fails) == 0 {
		return "No fails found for " + playerName, nil
	}

	msg := "Fails of " + playerName + " (" + strconv.Itoa(len(fails)) + ") :\n"
	for _, fail := range fails {
		msg += "* " + fail.Raid.Date.Format("02/01/06") + " - " + fail.Reason + " - " + strconv.Itoa(fail.ID) + "\n"
	}

	return msg, nil
}

func (d Discord) ListFailsOnRaidHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Fail/ListFailsOnRaidHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	span.SetAttributes(
		attribute.String("date", optionMap["date"].StringValue()),
	)
	raidDate, err := ParseDate(optionMap["date"].StringValue(), "")
	if err != nil {
		msg := "Error while getting fails: " + HumanReadableError(err)
		return msg, fmt.Errorf("list fail parse date: %w", err)
	}

	fails, err := d.ListFailOnRaid(ctx, raidDate[0])
	if err != nil {
		msg := "Error while getting fails: " + HumanReadableError(err)
		return msg, fmt.Errorf("list fail call usecase: %w", err)
	}

	if len(fails) == 0 {
		return "No fails found for " + raidDate[0].Format("02/01/06"), nil
	}

	msg := "Fails for " + raidDate[0].Format("02/01/06") + " (" + strconv.Itoa(len(fails)) + ") :\n"
	var players []entity.Player
	for _, fail := range fails {
		players = append(players, *fail.Player)
		players[len(players)-1].Fails = append(players[len(players)-1].Fails, fail)
	}
	for _, player := range players {
		for _, fail := range player.Fails {
			msg += "* " + player.Name + " - " + fail.Reason + "\n"
		}
	}
	return msg, nil
}

func (d Discord) DeleteFailHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Fail/DeleteFailHandle")
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
	failID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		msg := "Error while deleting fail: " + HumanReadableError(err)
		return msg, fmt.Errorf("delete fail parse id: %w", err)
	}

	err = d.DeleteFail(ctx, int(failID))
	if err != nil {
		msg := "Error while deleting fail: " + HumanReadableError(err)
		return msg, fmt.Errorf("delete fail usecase: %w", err)
	}

	return "Fail successfully deleted", nil
}
