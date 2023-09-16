package discordHandler

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sync"
	"time"
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
				Description: "exemple : 09/09/23",
				Required:    true,
			},
		},
	},
}

func (d Discord) InitAbsence() map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error{
		"coven-absence-create": d.CreateAbsenceHandler,
		"coven-absence-delete": d.DeleteAbsenceHandler,
		"coven-absence-list":   d.ListAbsenceHandler,
	}
}

func (d Discord) ListAbsenceHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

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

func (d Discord) CreateAbsenceHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	resultCh := make(chan bool)
	errorCh := make(chan error)
	go func() {

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
				select {
				case <-ctx.Done():
					return
				default:
					err = d.CreateAbsence(ctx, user, date)
					if err != nil {
						errorCh <- err
						return
					}

				}
			}()
			wg.Wait()
			resultCh <- true
		}
	}()

	var message string
	select {
	case <-resultCh:
		if optionMap["annuler"] != nil {
			message = "Votre absence a été annulée"
		} else {
			message = "Votre absence a été prise en compte"
		}
	case err := <-errorCh:
		message = err.Error()
	case <-ctx.Done():
		message = "Sorry, backend takes too much time to respond"
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		fmt.Print(err)
	}
	return nil
}

func (d Discord) DeleteAbsenceHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	resultCh := make(chan bool)
	errorCh := make(chan error)
	go func() {

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
				select {
				case <-ctx.Done():
					return
				default:
					err = d.DeleteAbsence(ctx, user, date)
					if err != nil {
						errorCh <- err
						return
					}
				}
			}()
			wg.Wait()
			resultCh <- true
		}
	}()

	var message string
	select {
	case <-resultCh:
		if optionMap["annuler"] != nil {
			message = "Votre absence a été annulée"
		} else {
			message = "Votre absence a été prise en compte"
		}
	case err := <-errorCh:
		message = err.Error()
	case <-ctx.Done():
		message = "Sorry, backend takes too much time to respond"
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		fmt.Print(err)
	}
	return nil

}
