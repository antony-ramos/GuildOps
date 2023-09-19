package discord

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
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
	var err error
	d.s, err = discordgo.New("Bot " + d.token)
	if err != nil {
		return fmt.Errorf("discord - Run - discordgo.New: %w", err)
	}

	d.s.AddHandler(func(session *discordgo.Session, r *discordgo.Ready) {
	})
	err = d.s.Open()
	if err != nil {
		return fmt.Errorf("discord - Run - d.s.Open: %w", err)
	}

	d.s.AddHandler(func(session *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := d.commandHandlers[i.ApplicationCommandData().Name]; ok {
			err := h(ctx, session, i)
			if err != nil {
				zap.L().Error(fmt.Sprintf("Error while handling command %s : %s", i.ApplicationCommandData().Name, err.Error()))
			}
		}
	})

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
				return // S'arrête immédiatement si un autre goroutine a signalé une erreur
			default:
				cmd, err := d.s.ApplicationCommandCreate(d.s.State.User.ID, strconv.Itoa(d.guildID), command)
				if err != nil {
					errCh <- err
					close(stopCh) // Ferme le canal pour signaler aux autres goroutines de s'arrêter
				}
				registeredCommands[commandName] = cmd
			}
		}()
	}
	waitGroup.Wait()
	close(errCh)

	defer func(session *discordgo.Session) {
		err := session.Close()
		if err != nil {
			zap.L().Error(err.Error())
		}
	}(d.s)

	<-ctx.Done()
	if d.DeleteCommands {
		for _, v := range registeredCommands {
			err := d.s.ApplicationCommandDelete(d.s.State.User.ID, strconv.Itoa(d.guildID), v.ID)
			if err != nil {
				return fmt.Errorf("discord - Run - d.s.ApplicationCommandDelete: %w", err)
			}
		}
	}
	return nil
}
