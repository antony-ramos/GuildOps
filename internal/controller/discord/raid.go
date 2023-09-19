package discordhandler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (d Discord) InitRaid() map[string]func(
	ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error{
		"coven-raid-create": d.CreateRaidHandler,
		"coven-raid-del":    d.DeleteRaidHandler,
	}
}

var RaidDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "coven-raid-create",
		Description: "Créer un raid",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Raid Milo",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "ex: 11/05/23",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "difficulté",
				Description: "ex: HM",
				Required:    true,
			},
		},
	},
	{
		Name:        "coven-raid-list",
		Description: "Lister les raids",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Milowenn",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "ex: Milowenn",
				Required:    false,
			},
		},
	},
	{
		Name:        "coven-raid-del",
		Description: "Supprimer un raid",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "ex: 4444-4444-4444",
				Required:    true,
			},
		},
	},
}

func (d Discord) CreateRaidHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	returnErr := error(nil)
	var msg string

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	name := optionMap["name"].StringValue()
	date, err := parseDate(optionMap["date"].StringValue())
	if err != nil {
		msg = "Erreur lors de la création du raid: " + err.Error()
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return fmt.Errorf("discord - CreateRaidHandler - parseDate: %w", err)
	}
	difficulty := optionMap["difficulté"].StringValue()

	raid, err := d.CreateRaid(ctx, name, difficulty, date[0])
	if err != nil {
		msg = "Erreur lors de la création du raid: " + err.Error()
		returnErr = err
	} else {
		msg = "Raid " + strconv.Itoa(raid.ID) + " créé avec succès"
	}

	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return returnErr
}

func (d Discord) DeleteRaidHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var returnErr error
	var msg string

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	raidID, err := strconv.Atoi(optionMap["raidID"].StringValue())
	if err != nil {
		return fmt.Errorf("discord - DeleteRaidHandler - strconv.Atoi: %w", err)
	}

	err = d.DeleteRaid(ctx, raidID)
	if err != nil {
		msg = "Erreur lors de la suppression du joueur: " + err.Error()
		returnErr = err
	} else {
		msg = "Joueur " + strconv.Itoa(raidID) + " supprimé avec succès"
	}

	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return returnErr
}
