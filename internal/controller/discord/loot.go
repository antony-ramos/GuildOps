package discordHandler

import (
	"context"
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
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "raid-id",
				Description: "ex: Milowenn",
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
	ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error{
		"coven-loot-attribute": d.AttributeLootHandler,
		"coven-loot-list":      d.ListLootsOnPlayerHandler,
		"coven-loot-delete":    d.DeleteLootHandler,
		"coven-loot-selector":  d.LootCounterCheckerHandler,
	}
}

func (d Discord) AttributeLootHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	lootName := optionMap["loot-name"].StringValue()
	raidID, err := strconv.Atoi(optionMap["raid-id"].StringValue())
	if err != nil {
		msg := "Erreur lors de l'attribution du loot: " + err.Error()
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return err
	}
	playerName := optionMap["player-name"].StringValue()

	err = d.LootUseCase.CreateLoot(ctx, lootName, raidID, playerName)
	if err != nil {
		msg := "Erreur lors de l'attribution du loot: " + err.Error()
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return err
	}
	msg := "Loot attribué avec succès"
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}

func (d Discord) ListLootsOnPlayerHandler(
	ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate,
) error {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	playerName := optionMap["player-name"].StringValue()

	lootList, err := d.LootUseCase.ListLootOnPLayer(ctx, playerName)
	if err != nil {
		msg := "Erreur lors de la récupération des loots: " + err.Error()
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return err
	}
	msg := "Tous les loots de " + playerName + ":\n"
	for _, loot := range lootList {
		msg += loot.Name + " " + loot.Raid.Date.String() + " " + loot.Raid.Difficulty + "\n"
	}
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}

func (d Discord) DeleteLootHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	id, err := strconv.Atoi(optionMap["id"].StringValue())
	if err != nil {
		return err
	}

	err = d.LootUseCase.DeleteLoot(ctx, id)
	if err != nil {
		msg := "Erreur lors de la suppression du loot: " + err.Error()
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return err
	}
	msg := "Loot supprimé avec succès"
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}

func (d Discord) LootCounterCheckerHandler(
	ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate,
) error {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	playerNames := strings.Split(optionMap["player-list"].StringValue(), ",")
	difficulty := optionMap["difficulty"].StringValue()

	player, err := d.LootUseCase.SelectPlayerToAssign(ctx, playerNames, difficulty)
	if err != nil {
		msg := "Erreur lors de l'assignation du loot: " + err.Error()
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return err
	}

	msg := "Le joueur " + player.Name + " a été sélectionné pour recevoir le loot"

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}
