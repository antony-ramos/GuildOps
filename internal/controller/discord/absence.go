package discordhandler

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
)

var AbsenceDescriptor = []discordgo.ApplicationCommand{
	{
		Name:        "guildops-absence-create",
		Description: "Create an absence for a raid or multiple raids",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "from",
				Description: "11/05/23",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "to",
				Description: "15/05/23",
				Required:    false,
			},
		},
	},
	{
		Name:        "guildops-absence-delete",
		Description: "Delete an absence for a raid or multiple raids",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "from",
				Description: "(ex: 11/05/23)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "to",
				Description: "(ex: 15/05/23)",
				Required:    false,
			},
		},
	},
	{
		Name:        "guildops-absence-list",
		Description: "List all absences for a raid",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "(ex: 09/09/23)",
				Required:    true,
			},
		},
	},
}

func (d Discord) InitAbsence() map[string]func(
	ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error{
		"guildops-absence-create": d.AbsenceHandler,
		"guildops-absence-delete": d.AbsenceHandler,
		"guildops-absence-list":   d.ListAbsenceHandler,
	}
}

func (d Discord) GenerateListAbsenceHandlerMsg(ctx context.Context, date string) (string, error) {
	errorMsg := "Error while listing absences" +
		": "

	select {
	case <-ctx.Done():
		return ctxError,
			fmt.Errorf("discord - GenerateListAbsenceHandlerMsg - ctx.Done: request took too much time to be proceed")
	default:
		var msg string
		dates, err := ParseDate(date, "")
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
	ctx context.Context, user string, fromDate string, toDate string, created bool,
) (string, error) {
	errorMsg := "Error while creating absence: "
	msg := "Absence(s) created for :\n"
	if !created {
		errorMsg = "Error while deleting absence: "
		msg = "Absence(s) deleted for :\n"
	}
	select {
	case <-ctx.Done():
		return ctxError,
			fmt.Errorf("discord - GenerateAbsenceHandlerMsg - ctx.Done: request took too much time to be proceed")
	default:
		dates, err := ParseDate(fromDate, toDate)

		if dates[0].Before(
			time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())) {
			return "You can't create a absence in the past",
				errors.New("discord - GenerateAbsenceHandlerMsg: can't create a absence in the past")
		}

		if err != nil {
			return errorMsg + HumanReadableError(err), err
		}

		RaidNotFound := regexp.MustCompile(".*no raid found.*")
		AbsenceAlreadyExist := regexp.MustCompile(".*absence already exist.*")
		AbsenceNotFound := regexp.MustCompile(".*absence not found.*")
		for _, date := range dates {
			date := date
			if !created {
				err = d.DeleteAbsence(ctx, user, date)
				if err != nil {
					errorRegex := fmt.Sprintf("(%s|%s)", AbsenceNotFound, RaidNotFound)
					matched, _ := regexp.MatchString(errorRegex, err.Error())
					if len(dates) == 1 || !matched {
						return errorMsg + HumanReadableError(err), err
					} else {
						matched = RaidNotFound.MatchString(err.Error())
						if !matched {
							msg += "* " + date.Format("Mon 02/01/06") + "\n"
						}
					}
				} else {
					msg += "* " + date.Format("Mon 02/01/06") + "\n"
				}
			} else {
				err = d.CreateAbsence(ctx, user, date)
				if err != nil {
					errorRegex := fmt.Sprintf("(%s|%s)", RaidNotFound, AbsenceAlreadyExist)
					matched, _ := regexp.MatchString(errorRegex, err.Error())
					if len(dates) == 1 || !matched {
						return errorMsg + HumanReadableError(err), err
					} else {
						matched = RaidNotFound.MatchString(err.Error())
						if !matched {
							msg += "* " + date.Format("Mon 02/01/06") + "\n"
						}
					}
				} else {
					msg += "* " + date.Format("Mon 02/01/06") + "\n"
				}
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

	from := optionMap["from"].StringValue()
	toDate := ""
	if len(optionMap) > 1 {
		toDate = optionMap["to"].StringValue()
	}

	msg, err := d.GenerateAbsenceHandlerMsg(
		ctx, user, from, toDate, interaction.ApplicationCommandData().Name == "guildops-absence-create")
	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return err
}
