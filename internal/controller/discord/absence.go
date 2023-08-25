package discordHandler

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/coven-discord-bot/internal/entity"
	"github.com/coven-discord-bot/internal/usecase"
	"strings"
	"sync"
	"time"
)

type Discord struct {
	*usecase.AbsenceUseCase
}

func parseDate(dateStr string) ([]time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	dateParts := strings.Split(dateStr, " au ")

	var dates []time.Time

	for _, datePart := range dateParts {
		date, err := time.Parse("02/01/06", datePart)
		if err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}

	return dates, nil
}

var AbsenceDescriptor = discordgo.ApplicationCommand{
	Name:        "absence",
	Description: "Créer une absence pour un ou plusieurs raids",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "date",
			Description: "(ex: 11-05-23 | ou 11-05-23 au 13-05-23)",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "annuler",
			Description: "Annuler à la place de créer",
			Required:    false,
		},
	},
}

func (d Discord) AbsenceHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
					abs := entity.Absence{
						Name: user,
						Date: date,
					}
					if optionMap["annuler"] != nil {
						err = d.RemoveAbsence(ctx, abs)
						if err != nil {
							errorCh <- err
							return
						}
					} else {
						err = d.AddAbsence(ctx, abs)
						if err != nil {
							errorCh <- err
							return
						}
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
}
