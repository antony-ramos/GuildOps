package discordhandler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"

	"github.com/bwmarrin/discordgo"
)

var FailDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "coven-fail-create", // Tested
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
		Name:        "coven-fail-list-player",
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
		Name:        "coven-fail-list-raid",
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
		Name:        "coven-fail-del",
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
	ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error{
		"coven-fail-create":      d.FailOnPlayerHandler,
		"coven-fail-del":         d.DeleteFailHandler,
		"coven-fail-list-player": d.ListFailsOnPlayerHandler,
		"coven-fail-list-raid":   d.ListFailsOnRaidHandler,
	}
}

func (d Discord) FailOnPlayerHandler(
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
	raidDate, err := parseDate(optionMap["date"].StringValue())
	if err != nil {
		msg = "Erreurs lors de la création du fail: " + HumanReadableError(err)
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return fmt.Errorf("database - FailOnPlayerHandler - parseDate: %w", err)
	}
	err = d.CreateFail(ctx, reason, raidDate[0], name)
	returnErr := error(nil)
	if err != nil {
		msg = "Erreurs lors de la création du fail: " + HumanReadableError(err)
		returnErr = err
	} else {
		msg = "Fail créé avec succès"
	}

	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return returnErr
}

func (d Discord) ListFailsOnPlayerHandler(
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

	fails, err := d.ListFailOnPLayer(ctx, playerName)
	if err != nil {
		msg = "Erreurs lors de la récupération des fails: " + HumanReadableError(err)
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return fmt.Errorf("database - ListFailsOnPlayerHandler - r.ReadFails: %w", err)
	}

	msg = "Fails de " + playerName + " (" + strconv.Itoa(len(fails)) + ") :\n"
	for _, fail := range fails {
		msg += "* " + fail.Raid.Date.Format("02-01-2006") + " - " + fail.Reason + "\n"
	}

	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}

func (d Discord) ListFailsOnRaidHandler(
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
	raidDate, err := parseDate(optionMap["date"].StringValue())
	if err != nil {
		msg = "Erreurs lors de la récupération des fails: " + HumanReadableError(err)
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return fmt.Errorf("database - ListFailsOnRaid - parseDate: %w", err)
	}

	fails, err := d.ListFailOnRaid(ctx, raidDate[0])
	if err != nil {
		msg = "Erreurs lors de la récupération des fails: " + HumanReadableError(err)
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return fmt.Errorf("database - ListFailsOnPlayerHandler - r.ReadFails: %w", err)
	}

	msg = "Fails du " + raidDate[0].Format("02/01/2006") + " (" + strconv.Itoa(len(fails)) + ") :\n"
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

	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}

//nolint:dupl
func (d Discord) DeleteFailHandler(
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
	failID, err := strconv.ParseInt(idString, 10, 64)
	returnErr := error(nil)
	if err != nil {
		msg = "Erreurs lors de la suppression du fail: " + HumanReadableError(err)
		returnErr = err
	} else {
		err = d.DeleteFail(ctx, int(failID))
		if err != nil {
			msg = "Erreurs lors de la suppression du fail: " + HumanReadableError(err)
			returnErr = err
		} else {
			msg = "Fail supprimé avec succès"
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
