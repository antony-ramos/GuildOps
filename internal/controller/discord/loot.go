package discordhandler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var LootDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "coven-loot-attribute",
		Description: "Attribuer un Loot à un joueur",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "loot-name",
				Description: "ex: Tête de Nefarian",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "raid-id",
				Description: "ex: 4488766425",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player-name",
				Description: "ex: Milowenn",
				Required:    true,
			},
		},
	},
	{
		Name:        "coven-loot-list",
		Description: "Donner un Loot à un joueur",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player-name",
				Description: "ex: Milowenn",
				Required:    true,
			},
		},
	},
	{
		Name:        "coven-loot-delete",
		Description: "Supprimer un Loot à un joueur",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "ex: 5465-4444-5557",
				Required:    true,
			},
		},
	},
	{
		Name:        "coven-loot-selector",
		Description: "Donner la liste des joueurs qui peuvent avoir un loot",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player-list",
				Description: "ex: Milowenn,Arthas,Jailer",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "difficulty",
				Description: "ex: HM",
				Required:    true,
			},
		},
	},
}

func (d Discord) InitLoot() map[string]func(
	ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error{
		"coven-loot-attribute": d.AttributeLootHandler,
		"coven-loot-list":      d.ListLootsOnPlayerHandler,
		"coven-loot-delete":    d.DeleteLootHandler,
		"coven-loot-selector":  d.LootCounterCheckerHandler,
	}
}

func (d Discord) AttributeLootHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	lootName := optionMap["loot-name"].StringValue()
	raidID := optionMap["raid-id"].IntValue()
	playerName := optionMap["player-name"].StringValue()

	err := d.LootUseCase.CreateLoot(ctx, lootName, int(raidID), playerName)
	if err != nil {
		msg := "Erreur lors de l'attribution du loot: " + HumanReadableError(err)
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return fmt.Errorf("discord - AttributeLootHandler - d.LootUseCase.CreateLoot: %w", err)
	}
	msg := "Loot attribué avec succès"
	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}

func (d Discord) ListLootsOnPlayerHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	playerName := optionMap["player-name"].StringValue()

	lootList, err := d.LootUseCase.ListLootOnPLayer(ctx, playerName)
	if err != nil {
		msg := "Erreur lors de la récupération des loots: " + HumanReadableError(err)
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return fmt.Errorf("discord - ListLootsOnPlayerHandler - d.LootUseCase.ListLootOnPLayer: %w", err)
	}
	msg := "Tous les loots de " + playerName + ":\n"
	for _, loot := range lootList {
		msg += loot.Name + " " + loot.Raid.Date.String() + " " + loot.Raid.Difficulty + "\n"
	}
	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}

func (d Discord) DeleteLootHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	id, err := strconv.Atoi(optionMap["id"].StringValue())
	if err != nil {
		return fmt.Errorf("discord - DeleteLootHandler - strconv.Atoi: %w", err)
	}

	err = d.LootUseCase.DeleteLoot(ctx, id)
	if err != nil {
		msg := "Erreur lors de la suppression du loot: " + HumanReadableError(err)
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return fmt.Errorf("discord - DeleteLootHandler - d.LootUseCase.DeleteLoot: %w", err)
	}
	msg := "Loot supprimé avec succès"
	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}

func (d Discord) LootCounterCheckerHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	playerNames := strings.Split(optionMap["player-list"].StringValue(), ",")
	difficulty := optionMap["difficulty"].StringValue()

	player, err := d.LootUseCase.SelectPlayerToAssign(ctx, playerNames, difficulty)
	if err != nil {
		msg := "Erreur lors de l'assignation du loot: " + HumanReadableError(err)
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return fmt.Errorf("discord - LootCounterCheckerHandler - d.LootUseCase.SelectPlayerToAssign: %w", err)
	}

	msg := "Le joueur " + player.Name + " a été sélectionné pour recevoir le loot"

	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}
