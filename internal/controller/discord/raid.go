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
		"guildops-raid-create": d.CreateRaidHandler,
		"guildops-raid-del":    d.DeleteRaidHandler,
	}
}

var RaidDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "guildops-raid-create",
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
		Name:        "guildops-raid-del",
		Description: "Supprimer un raid",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "id",
				Description: "ex: 4546646",
				Required:    true,
			},
		},
	},
}

func (d Discord) CreateRaidHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

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
		msg = "Erreur lors de la création du raid: " + HumanReadableError(err)
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
		msg = "Erreur lors de la création du raid: " + HumanReadableError(err)
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
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	var returnErr error
	var msg string

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	raidID := optionMap["raidID"].IntValue()

	err := d.DeleteRaid(ctx, int(raidID))
	if err != nil {
		msg = "Erreur lors de la suppression du joueur: " + HumanReadableError(err)
		returnErr = err
	} else {
		msg = "Joueur " + strconv.Itoa(int(raidID)) + " supprimé avec succès"
	}

	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return returnErr
}
