package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

type Discord struct {
	token           string
	guildID         int
	DeleteCommands  bool
	commands        []*discordgo.ApplicationCommand
	commandHandlers map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error
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
		return err
	}

	d.s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
	})
	err = d.s.Open()
	if err != nil {
		return err
	}

	d.s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := d.commandHandlers[i.ApplicationCommandData().Name]; ok {
			err := h(ctx, s, i)
			if err != nil {
				zap.L().Error(fmt.Sprintf("Error while handling command %s : %s", i.ApplicationCommandData().Name, err.Error()))
			}
		}

	})

	registeredCommands := make([]*discordgo.ApplicationCommand, len(d.commands))
	var wg sync.WaitGroup
	errCh := make(chan error, len(d.commands))
	stopCh := make(chan struct{})
	for i, v := range d.commands {
		i := i
		wg.Add(1)
		v := v
		go func() {
			defer wg.Done()
			select {
			case <-stopCh:
				return // S'arrête immédiatement si un autre goroutine a signalé une erreur
			default:
				cmd, err := d.s.ApplicationCommandCreate(d.s.State.User.ID, strconv.Itoa(d.guildID), v)
				if err != nil {
					errCh <- err
					close(stopCh) // Ferme le canal pour signaler aux autres goroutines de s'arrêter
				}
				registeredCommands[i] = cmd
			}
		}()
	}
	wg.Wait()
	close(errCh)

	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {
			zap.L().Error(err.Error())
		}
	}(d.s)

	<-ctx.Done()
	if d.DeleteCommands {
		for _, v := range registeredCommands {
			err := d.s.ApplicationCommandDelete(d.s.State.User.ID, strconv.Itoa(d.guildID), v.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil

}
