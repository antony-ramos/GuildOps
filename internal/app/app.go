package app

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/coven-discord-bot/config"
	discordHandler "github.com/coven-discord-bot/internal/controller/discord"
	"github.com/coven-discord-bot/internal/usecase"
	"github.com/coven-discord-bot/internal/usecase/backend"
	"github.com/coven-discord-bot/pkg/discord"
)

func Run(cfg *config.Config) {
	ctx := context.TODO()

	bk := backend.NewRPC(cfg.Backend.URL)

	mapHandler := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

	absenceusecase := usecase.NewAbsenceUseCase(bk)
	d := discordHandler.Discord{
		absenceusecase,
	}

	mapHandler["absence"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		d.AbsenceHandler(s, i)
	}

	handlers := []*discordgo.ApplicationCommand{}
	handlers = append(handlers, &discordHandler.AbsenceDescriptor)

	serve := discord.New(discord.CommandHandlers(mapHandler), discord.Token(cfg.Discord.Token), discord.Command(handlers), discord.GuildID(cfg.Discord.GuildID))
	err := serve.Run(ctx)
	if err != nil {
		fmt.Print(err)
	}

}
