package discordHandler

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var StrikeDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "guildops-strike-create", // Tested
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
		Name:        "guildops-strike-list",
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
		Name:        "guildops-strike-del",
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

func (d Discord) InitStrike() map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error{
		"guildops-strike-create": d.StrikeOnPlayerHandler,
		"guildops-strike-del":    d.DeleteStrikeHandler,
		"guildops-strike-list":   d.ListStrikesOnPlayerHandler,
	}
}

func (d Discord) StrikeOnPlayerHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
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
		msg = "Erreurs lors de la création du strike: " + err.Error()
		returnErr = err
	} else {
		msg = "Strike créé avec succès"
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return returnErr
}

func (d Discord) ListStrikesOnPlayerHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var msg string
	playerName := optionMap["name"].StringValue()

	strikes, err := d.ReadStrikes(ctx, playerName)
	if err != nil {
		msg = "Erreurs lors de la récupération des strikes: " + err.Error()
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return err
	}

	msg = "Strikes de " + playerName + ":\n"
	for _, strike := range strikes {
		msg += strike.Date.String() + " | " + strike.Reason + "\n"
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}

func (d Discord) DeleteStrikeHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var msg string
	idString := optionMap["id"].StringValue()
	id, err := strconv.ParseInt(idString, 10, 64)
	returnErr := error(nil)
	if err != nil {
		msg = "Erreurs lors de la suppression du strike: " + err.Error()
		returnErr = err
	} else {
		err = d.DeleteStrike(ctx, int(id))
		if err != nil {
			msg = "Erreurs lors de la suppression du strike: " + err.Error()
			returnErr = err
		} else {
			msg = "Strike supprimé avec succès"
		}
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return returnErr
}
