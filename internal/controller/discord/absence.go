package discordHandler

import (
	"context"
	"fmt"
	"sync"
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
	ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error{
		"coven-absence-create": d.AbsenceHandler,
		"coven-absence-delete": d.AbsenceHandler,
		"coven-absence-list":   d.ListAbsenceHandler,
	}
}

func (d Discord) ListAbsenceHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var msg string
	dates, err := parseDate(optionMap["date"].StringValue())
	if err != nil {
		msg = "Erreur lors de la récupération des absences: " + err.Error()
	} else {
		msg = "Absences pour le " + dates[0].Format("02-01-2006") + ":\n"
		absences, err := d.ListAbsence(ctx, dates[0])
		if err != nil {
			msg = "Erreur lors de la récupération des absences: " + err.Error()
		} else {
			for _, absence := range absences {
				msg += absence.Player.Name + "\n"
			}
		}
	}
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return err
}

func (d Discord) AbsenceHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var user string
	if i.Member != nil {
		user = i.Member.User.Username
	} else {
		user = i.User.Username
	}

	dates, err := parseDate(optionMap["date"].StringValue())
	var wg sync.WaitGroup
	for _, date := range dates {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if i.ApplicationCommandData().Name == "coven-absence-delete" {
				err = d.DeleteAbsence(ctx, user, date)
				if err != nil {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Erreur lors de la suppression de l'absence pour le " + date.Format("02-01-2006"),
						},
					})
				}
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Absence supprimée pour le " + date.Format("02-01-2006"),
					},
				})
				if err != nil {
					fmt.Print(err)
				}
			} else {
				err = d.CreateAbsence(ctx, user, date)
				if err != nil {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Erreur lors de la suppression de l'absence pour le " + date.Format("02-01-2006"),
						},
					})
				}
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Absence créée pour le " + date.Format("02-01-2006"),
					},
				})
				if err != nil {
					fmt.Print(err)
				}
			}
		}()
	}
	wg.Wait()
	return nil
}
