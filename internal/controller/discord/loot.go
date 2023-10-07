package discordhandler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/bwmarrin/discordgo"
)

var LootDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "guildops-loot-attribute",
		Description: "Attribuer un Loot à un joueur",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "loot-name",
				Description: "ex: Tête de Nefarian",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "raid-date",
				Description: "(ex: 02/10/23)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player-name",
				Description: "(ex: milowenn)",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-loot-list-on-player",
		Description: "List loot on player",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player-name",
				Description: "(ex: milowenn)",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-loot-list-on-raid",
		Description: "List loot on raid",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "(ex: 03/10/23)",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-loot-delete",
		Description: "Supprimer un Loot à un joueur",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "(ex: 123456789)",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-loot-selector",
		Description: "Donner la liste des joueurs qui peuvent avoir un loot",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player-list",
				Description: "(ex: arthas,jailer,garrosh)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "difficulty",
				Description: "(ex: mythic, heroic, normal)",
				Required:    true,
			},
		},
	},
}

func (d Discord) InitLoot() map[string]func(
	ctx context.Context, interaction *discordgo.InteractionCreate) (string, error) {
	return map[string]func(ctx context.Context, interaction *discordgo.InteractionCreate) (string, error){
		"guildops-loot-attribute":      d.AttributeLootHandler,
		"guildops-loot-list-on-player": d.ListLootsOnPlayerHandler,
		"guildops-loot-list-on-raid":   d.ListLootsOnRaidHandler,
		"guildops-loot-delete":         d.DeleteLootHandler,
		"guildops-loot-selector":       d.LootCounterCheckerHandler,
	}
}

func (d Discord) AttributeLootHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Loot/ListLootsOnPlayerHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	lootName := optionMap["loot-name"].StringValue()
	raidDate := optionMap["raid-date"].StringValue()
	playerName := optionMap["player-name"].StringValue()
	span.SetAttributes(
		attribute.String("loot_name", lootName),
		attribute.String("raid_date", raidDate),
		attribute.String("player_name", playerName),
	)

	raidDates, err := ParseDate(raidDate, "")
	if err != nil {
		return "invalid date", fmt.Errorf("discord - AttributeLootHandler - ParseDate: %w", err)
	}

	err = d.LootUseCase.CreateLoot(ctx, lootName, raidDates[0], playerName)
	if err != nil {
		msg := "Error while proceeding loot attribution: " + HumanReadableError(err)
		return msg, fmt.Errorf("discord - AttributeLootHandler - d.LootUseCase.CreateLoot: %w", err)
	}
	msg := "Loot successfully attributed"
	return msg, nil
}

func (d Discord) ListLootsOnPlayerHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Loot/ListLootsOnPlayerHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	playerName := optionMap["player-name"].StringValue()
	span.SetAttributes(
		attribute.String("player_name", playerName),
	)

	lootList, err := d.LootUseCase.ListLootOnPLayer(ctx, playerName)
	if err != nil {
		msg := "Error while getting loot for player: " + HumanReadableError(err)
		return msg, fmt.Errorf("discord - ListLootsOnPlayerHandler - d.LootUseCase.ListLootOnPLayer: %w", err)
	}
	if len(lootList) == 0 {
		return "no loot for " + playerName, nil
	}
	msg := "All loots of " + playerName + ":\n"
	for _, loot := range lootList {
		msg += "* " + loot.Name + " " + loot.Raid.Date.Format("02/01/06") + " " +
			loot.Raid.Difficulty + " " + strconv.Itoa(loot.ID) + "\n"
	}
	return msg, nil
}

func (d Discord) ListLootsOnRaidHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Loot/ListLootsOnRaidHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	dateString := optionMap["date"].StringValue()
	span.SetAttributes(
		attribute.String("date", dateString),
	)
	date, err := ParseDate(dateString, "")
	if err != nil {
		return "invalid date", fmt.Errorf("discord - ListLootsOnRaidHandler - ParseDate: %w", err)
	}

	lootList, err := d.LootUseCase.ListLootOnRaid(ctx, date[0])
	if err != nil {
		msg := "Error while listing loot for raid: " + HumanReadableError(err)
		return msg, fmt.Errorf("discord - ListLootsOnPlayerHandler - d.LootUseCase.ListLootOnPLayer: %w", err)
	}
	if len(lootList) == 0 {
		return "no loot for " + date[0].Format("02/01/06"), nil
	}
	msg := "All loots of  " + date[0].Format("02/01/06") + ":\n"
	for _, loot := range lootList {
		msg += "* " + loot.Name + " " + loot.Player.Name + " " + loot.Raid.Difficulty + " " + strconv.Itoa(loot.ID) + "\n"
	}
	return msg, nil
}

func (d Discord) DeleteLootHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Loot/DeleteLootHandler")
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
		attribute.String("id", optionMap["id"].StringValue()))
	id, err := strconv.Atoi(optionMap["id"].StringValue())
	if err != nil {
		return "id format is invalid", fmt.Errorf("discord - DeleteLootHandler - strconv.Atoi: %w", err)
	}

	err = d.LootUseCase.DeleteLoot(ctx, id)
	if err != nil {
		msg := "Error while deleting loot: " + HumanReadableError(err)
		return msg, fmt.Errorf("discord - DeleteLootHandler - d.LootUseCase.DeleteLoot: %w", err)
	}
	msg := "Loot successfully deleted"
	return msg, nil
}

func (d Discord) LootCounterCheckerHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Loot/LootCounterCheckerHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	playerList := optionMap["player-list"].StringValue()
	playerList = strings.ReplaceAll(playerList, " ", "")
	playerNames := strings.Split(playerList, ",")
	difficulty := optionMap["difficulty"].StringValue()
	span.SetAttributes(
		attribute.String("player_list", optionMap["player-list"].StringValue()),
		attribute.String("difficulty", difficulty),
	)

	player, err := d.LootUseCase.SelectPlayerToAssign(ctx, playerNames, difficulty)
	if err != nil {
		msg := "Error while searching a player to attribute loot: " + HumanReadableError(err)
		return msg, fmt.Errorf("discord - LootCounterCheckerHandler - d.LootUseCase.SelectPlayerToAssign: %w", err)
	}

	msg := player.Name + " have been selected to receive the loot"
	return msg, nil
}
