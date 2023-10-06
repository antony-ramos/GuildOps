package discord

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alitto/pond"
	"github.com/antony-ramos/guildops/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type Discord struct {
	token           string
	guildID         int
	DeleteCommands  bool
	commands        []*discordgo.ApplicationCommand
	commandHandlers map[string]func(ctx context.Context, interaction *discordgo.InteractionCreate) (string, error)
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

	logger.FromContext(ctx).Debug("add handler to discord ready event")
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
			msg, err := handler(ctx, interaction)
			if err != nil {
				logger.FromContext(ctx).Error(
					fmt.Sprintf("handle command %s : %s", interaction.ApplicationCommandData().Name, err.Error()))
			}
			data := discordgo.InteractionResponseData{
				Content: msg,
			}
			if interaction.ApplicationCommandData().Name == "guildops-player-info" {
				data = discordgo.InteractionResponseData{
					Content: msg,
					Flags:   discordgo.MessageFlagsEphemeral,
				}
			}
			_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &data,
			})
		}
	})

	logger.FromContext(ctx).Debug("register commands to discord")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(d.commands))

	pool := pond.New(100, 1000)
	group, _ := pool.GroupContext(ctx)

	for i, v := range d.commands {
		commandName := i
		command := v
		group.Submit(func() error {
			logger.FromContext(ctx).Info("register command " + command.Name)
			cmd, err := d.s.ApplicationCommandCreate(d.s.State.User.ID, strconv.Itoa(d.guildID), command)
			if err != nil {
				return errors.Wrap(err, "try to create command "+command.Name)
			}
			registeredCommands[commandName] = cmd
			logger.FromContext(ctx).Info("command " + command.Name + " registered")
			return nil
		})
	}

	defer func(session *discordgo.Session) {
		logger.FromContext(ctx).Info("close discord session")
		err := session.Close()
		if err != nil {
			logger.FromContext(ctx).Error(err.Error())
		}
	}(d.s)

	logger.FromContext(ctx).Info("ready to handle commands")
	err = group.Wait()
	if err != nil {
		return fmt.Errorf("wait command creation: %w", err)
	}
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
