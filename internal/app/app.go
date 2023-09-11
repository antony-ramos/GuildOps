package app

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/coven-discord-bot/config"
	discordHandler "github.com/coven-discord-bot/internal/controller/discord"
	"github.com/coven-discord-bot/internal/usecase"
	"github.com/coven-discord-bot/internal/usecase/backend"
	"github.com/coven-discord-bot/pkg/discord"
	"github.com/coven-discord-bot/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

var (
	eventsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "coven_processed_total",
		Help: "The total number of discord webhooks processed",
	})
	eventsSucceededProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "coven_event_succeeded_total",
		Help: "The total number of discord succeeded events",
	})
	eventsFailedProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "coven_event_failed_total",
		Help: "The total number of discord failed events",
	})
)

func Run(ctx context.Context, cfg *config.Config, log *logger.Factory) {
	l := log.For(ctx).With(zap.String("backend", cfg.Backend.URL))

	l.Info("loading backend")
	bk := backend.NewRPC(cfg.Backend.URL)

	mapHandler := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

	auc := usecase.NewAbsenceUseCase(bk)
	d := discordHandler.Discord{
		AbsenceUseCase: auc,
	}

	mapHandler["absence"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		eventsProcessed.Inc()
		_, span := otel.Tracer("Discord").Start(ctx, "Receiving a webhook from Discord", trace.WithTimestamp(time.Now()), trace.WithAttributes(attribute.KeyValue{
			Key:   "Member",
			Value: attribute.StringValue(i.Member.User.Username),
		}))
		defer span.End(trace.WithTimestamp(time.Now()))
		l = l.With(zap.String("Member", i.Member.User.Username))
		l.Info("Receive request")
		err := d.AbsenceHandler(ctx, l, s, i)
		if err != nil {
			l.Error("Failed to proceed request", zap.Error(err))
			span.RecordError(err)
			eventsFailedProcessed.Inc()
			return
		}
		eventsSucceededProcessed.Inc()
	}

	var handlers []*discordgo.ApplicationCommand
	handlers = append(handlers, &discordHandler.AbsenceDescriptor)

	serve := discord.New(discord.CommandHandlers(mapHandler), discord.Token(cfg.Discord.Token), discord.Command(handlers), discord.GuildID(cfg.Discord.GuildID))

	l.Info("starting to serve to discord webhooks")
	err := serve.Run(ctx)
	if err != nil {
		l.Error(err.Error())
		return
	}
}
