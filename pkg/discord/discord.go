package discord

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/antony-ramos/guildops/pkg/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
	"go.opentelemetry.io/otel"
)

type Discord struct {
	token           string
	guildID         int
	DeleteCommands  bool
	commands        []*discordgo.ApplicationCommand
	commandHandlers map[string]func(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error
	s               *discordgo.Session
}

func New(opts ...Option) *Discord {
	d := &Discord{}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func (d *Discord) Run(ctx context.Context) error {
	logger.FromContext(ctx).Info("create discord session")
	session, err := discordgo.New("Bot " + d.token)
	if err != nil {
		return errors.Wrap(err, "new discord session")
	}
	d.s = session

	logger.FromContext(ctx).Debug("adding handler to discord ready event")
	d.s.AddHandler(func(session *discordgo.Session, r *discordgo.Ready) {
	})
	err = d.s.Open()
	if err != nil {
		return errors.Wrap(err, "add handler to discord ready event")
	}

	logger.FromContext(ctx).Debug("create handlers to discord interaction create event")
	loggerHandler := logger.FromContext(ctx)
	d.s.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if handler, ok := d.commandHandlers[interaction.ApplicationCommandData().Name]; ok {
			ctx := context.Background()
			ctx = logger.AddLoggerToContext(ctx, loggerHandler)

			logger.FromContext(ctx).Debug("handling command " + interaction.ApplicationCommandData().Name)
			ctx, span := otel.Tracer("discordHandler").Start(ctx, interaction.ApplicationCommandData().Name)
			ctx = logger.AddLoggerToContext(ctx, logger.FromContext(ctx).
				With(zap.String("discordHandler", interaction.ApplicationCommandData().Name)))
			defer span.End()
			err := handler(ctx, session, interaction)
			if err != nil {
				logger.FromContext(ctx).Error(
					fmt.Sprintf("handle command %s : %s", interaction.ApplicationCommandData().Name, err.Error()))
			}
		}
	})

	logger.FromContext(ctx).Debug("register commands to discord")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(d.commands))
	var waitGroup sync.WaitGroup
	errCh := make(chan error, len(d.commands))
	stopCh := make(chan struct{})
	for i, v := range d.commands {
		commandName := i
		waitGroup.Add(1)
		command := v
		go func() {
			defer waitGroup.Done()
			select {
			case <-stopCh:
				return
			default:
				logger.FromContext(ctx).Info("register command " + command.Name)
				cmd, err := d.s.ApplicationCommandCreate(d.s.State.User.ID, strconv.Itoa(d.guildID), command)
				if err != nil {
					errCh <- err
					close(stopCh)
				}
				registeredCommands[commandName] = cmd
				logger.FromContext(ctx).Info("command " + command.Name + " registered")
			}
		}()
	}
	waitGroup.Wait()
	close(errCh)

	defer func(session *discordgo.Session) {
		logger.FromContext(ctx).Info("close discord session")
		err := session.Close()
		if err != nil {
			logger.FromContext(ctx).Error(err.Error())
		}
	}(d.s)

	<-ctx.Done()

	logger.FromContext(ctx).Info("delete commands")
	if d.DeleteCommands {
		for _, value := range registeredCommands {
			logger.FromContext(ctx).Info("delete command " + value.Name)
			err := d.s.ApplicationCommandDelete(d.s.State.User.ID, strconv.Itoa(d.guildID), value.ID)
			if err != nil {
				return fmt.Errorf("discord - Run - d.s.ApplicationCommandDelete: %w", err)
			}
			logger.FromContext(ctx).Info("command " + value.Name + " deleted")
		}
	}
	return nil
}
