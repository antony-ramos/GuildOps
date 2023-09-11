package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type Discord struct {
	token           string
	guildID         string
	removeCommands  bool
	commands        []*discordgo.ApplicationCommand
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
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
	_, span := otel.Tracer("").Start(ctx, "Run", trace.WithTimestamp(time.Now()))
	defer span.End()

	var err error
	d.s, err = discordgo.New("Bot " + d.token)
	if err != nil {
		return err
	}

	d.s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
	})
	err = d.s.Open()
	if err != nil {
		return err
	}

	d.s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.GuildID != d.guildID {
			return
		}
		if h, ok := d.commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
			return
		}
	})

	registeredCommands := make([]*discordgo.ApplicationCommand, len(d.commands))
	for i, v := range d.commands {
		cmd, err := d.s.ApplicationCommandCreate(d.s.State.User.ID, d.guildID, v)
		if err != nil {
			return err
		}
		registeredCommands[i] = cmd
	}
	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {
			fmt.Print(err)
		}
	}(d.s)

	<-ctx.Done()
	if d.removeCommands {
		for _, v := range registeredCommands {
			err := d.s.ApplicationCommandDelete(d.s.State.User.ID, d.guildID, v.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil

}
