package discordhandler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var StrikeDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "coven-strike-create", // Tested
		Description: "Générer un Strike sur un joueur",
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
		Name:        "coven-strike-list",
		Description: "Lister les strikes sur un joueur",
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
		Name:        "coven-strike-del",
		Description: "Supprimer un strike via son ID (ListStrikes pour l'avoir)",
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

func (d Discord) InitStrike() map[string]func(
	ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error{
		"coven-strike-create": d.StrikeOnPlayerHandler,
		"coven-strike-del":    d.DeleteStrikeHandler,
		"coven-strike-list":   d.ListStrikesOnPlayerHandler,
	}
}

func (d Discord) StrikeOnPlayerHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var msg string
	name := optionMap["name"].StringValue()
	reason := optionMap["reason"].StringValue()
	err := d.CreateStrike(ctx, reason, name)
	returnErr := error(nil)
	if err != nil {
		msg = "Erreurs lors de la création du strike: " + HumanReadableError(err)
		returnErr = err
	} else {
		msg = "Strike créé avec succès"
	}

	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return returnErr
}

func (d Discord) ListStrikesOnPlayerHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var msg string
	playerName := optionMap["name"].StringValue()

	strikes, err := d.ReadStrikes(ctx, playerName)
	if err != nil {
		msg = "Erreurs lors de la récupération des strikes: " + HumanReadableError(err)
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return fmt.Errorf("database - ListStrikesOnPlayerHandler - r.ReadStrikes: %w", err)
	}

	msg = "Strikes de " + playerName + ":\n"
	for _, strike := range strikes {
		msg += strike.Date.String() + " | " + strike.Reason + "\n"
	}

	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}

func (d Discord) DeleteStrikeHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var msg string
	idString := optionMap["id"].StringValue()
	strikeID, err := strconv.ParseInt(idString, 10, 64)
	returnErr := error(nil)
	if err != nil {
		msg = "Erreurs lors de la suppression du strike: " + HumanReadableError(err)
		returnErr = err
	} else {
		err = d.DeleteStrike(ctx, int(strikeID))
		if err != nil {
			msg = "Erreurs lors de la suppression du strike: " + HumanReadableError(err)
			returnErr = err
		} else {
			msg = "Strike supprimé avec succès"
		}
	}

	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return returnErr
}
