package discordHandler

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/coven-discord-bot/internal/entity"
	"github.com/coven-discord-bot/internal/usecase"
	"github.com/coven-discord-bot/pkg/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
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
			Description: "(ex: 11/0523 ou 11/05/23 au 13/05/23)",
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

// TODO : Function is too long. Should be cut into pieces
func (d Discord) AbsenceHandler(ctx context.Context, l logger.Logger, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	_, span := otel.Tracer("").Start(ctx, "Controller is handling request", trace.WithTimestamp(time.Now()), trace.WithAttributes(attribute.KeyValue{
		Key:   "HandlerType",
		Value: attribute.StringValue("absence"),
	}))
	defer span.End(trace.WithTimestamp(time.Now()))

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	resultCh := make(chan bool)
	errorCh := make(chan error)
	go func() {
		_, span := otel.Tracer("").Start(ctx, "InteractWithBackend", trace.WithTimestamp(time.Now()))
		defer span.End(trace.WithTimestamp(time.Now()))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		var user string
		if i.Member != nil {
			user = i.Member.User.Username
		} else {
			user = i.User.Username
		}

		l = l.With(zap.String("Date", optionMap["date"].StringValue()))
		l.Info("Checking date validity")
		_, s := otel.Tracer("").Start(ctx, "Controller is checking date validity", trace.WithTimestamp(time.Now()), trace.WithAttributes(attribute.KeyValue{
			Key:   "HandlerType",
			Value: attribute.StringValue("absence"),
		}))
		dates, err := parseDate(optionMap["date"].StringValue())
		if err != nil {
			s.RecordError(err)
			errorCh <- err
			return
		}
		s.End(trace.WithTimestamp(time.Now()))

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
						l = l.With(zap.String("Cancel", "true"))
						l.Info("Absence request is a cancellation")
						err = d.RemoveAbsence(ctx, l, abs)
						if err != nil {
							errorCh <- err
							return
						}
					} else {
						err = d.AddAbsence(ctx, l, abs)
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
		l.Error(err.Error())
		span.RecordError(err)
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
		l.Error(err.Error())
		span.RecordError(err)
		return err
	}
	return nil
}
