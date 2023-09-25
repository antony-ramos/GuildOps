package discordhandler

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var AbsenceDescriptor = []discordgo.ApplicationCommand{
	{
		Name:        "coven-absence-create",
		Description: "Créer une absence pour un ou plusieurs raids",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "(ex: 11-05-23 | ou 11-05-23 au 13-05-23)",
				Required:    true,
			},
		},
	},
	{
		Name:        "coven-absence-delete",
		Description: "Supprimer une absence pour un ou plusieurs raids",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "(ex: 11-05-23 | ou 11-05-23 au 13-05-23)",
				Required:    true,
			},
		},
	},
	{
		Name:        "coven-absence-list",
		Description: "Lister les absences pour un raid",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "example : 09/09/23",
				Required:    true,
			},
		},
	},
}

func (d Discord) InitAbsence() map[string]func(
	ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error{
		"coven-absence-create": d.AbsenceHandler,
		"coven-absence-delete": d.AbsenceHandler,
		"coven-absence-list":   d.ListAbsenceHandler,
	}
}

func (d Discord) GenerateListAbsenceHandlerMsg(ctx context.Context, date string) (string, error) {
	errorMsg := "Error while listing absences" +
		": "

	select {
	case <-ctx.Done():
		return ctxError,
			fmt.Errorf("discord - GenerateListAbsenceHandlerMsg - ctx.Done: %w", ctx.Err())
	default:
		var msg string
		dates, err := parseDate(date)
		if err != nil {
			msg = errorMsg + HumanReadableError(err)
		} else {
			absences, err := d.ListAbsence(ctx, dates[0])
			if len(absences) == 0 {
				msg = "Aucune absence pour le " + date + "\n"
				return msg, err
			}
			msg = "Absence(s) pour le " + dates[0].Format("02-01-2006") + " :\n"
			if err != nil {
				msg = errorMsg + HumanReadableError(err)
				return msg, err
			} else {
				for _, absence := range absences {
					msg += "* " + absence.Player.Name + "\n"
				}
			}
		}
		return msg, err
	}
}

func (d Discord) ListAbsenceHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	msg, err := d.GenerateListAbsenceHandlerMsg(ctx, optionMap["date"].StringValue())
	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return err
}

func (d Discord) GenerateAbsenceHandlerMsg(
	ctx context.Context, user string, dates string, created bool,
) (string, error) {
	errorMsg := "Error while creating absence: "
	msg := "Absence(s) créée(s) pour le(s) :\n"
	if !created {
		errorMsg = "Error while deleting absence: "
		msg = "Absence(s) supprimée(s) pour le(s) :\n"
	}
	select {
	case <-ctx.Done():
		return ctxError,
			fmt.Errorf("discord - GenerateAbsenceHandlerMsg - ctx.Done: %w", ctx.Err())
	default:
		dates, err := parseDate(dates)
		if err != nil {
			return errorMsg + HumanReadableError(err), err
		}
		for _, date := range dates {
			date := date
			if !created {
				err = d.DeleteAbsence(ctx, user, date)
				if err != nil {
					return errorMsg + HumanReadableError(err), err
				}
				msg += "* " + date.Format("02-01-2006") + "\n"
			} else {
				err = d.CreateAbsence(ctx, user, date)
				if err != nil {
					return errorMsg + HumanReadableError(err), err
				}
				msg += "* " + date.Format("02-01-2006") + "\n"
			}
		}
	}
	return msg, nil
}

func (d Discord) AbsenceHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var user string
	if interaction.Member != nil {
		user = interaction.Member.User.Username
	} else {
		user = interaction.User.Username
	}

	msg, err := d.GenerateAbsenceHandlerMsg(
		ctx, user, optionMap["date"].StringValue(), interaction.ApplicationCommandData().Name == "coven-absence-create")
	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return err
}
