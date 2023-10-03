package app

import (
	"context"

	"github.com/antony-ramos/guildops/pkg/logger"
	"github.com/pkg/errors"

	"github.com/antony-ramos/guildops/config"
	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"
	"github.com/antony-ramos/guildops/internal/usecase"
	"github.com/antony-ramos/guildops/internal/usecase/postgresbackend"
	"github.com/antony-ramos/guildops/pkg/discord"
	"github.com/antony-ramos/guildops/pkg/postgres"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func Run(ctx context.Context, cfg *config.Config) {
	logger.FromContext(ctx).Info("loading backend")

	pgHandler, err := postgres.New(
		ctx,
		cfg.URL,
		postgres.MaxPoolSize(cfg.PoolMax),
		postgres.ConnAttempts(cfg.ConnAttempts),
		postgres.ConnTimeout(cfg.ConnTimeOut))
	if err != nil {
		logger.FromContext(ctx).Fatal(err.Error())
	}

	ctx = logger.AddLoggerToContext(ctx, logger.FromContext(ctx).With(zap.String("backend", "postgres")))

	backend := postgresbackend.PG{Postgres: pgHandler}
	err = backend.Init(ctx, cfg.URL, nil)
	if err != nil {
		logger.FromContext(ctx).Fatal(err.Error())
		return
	}

	mapHandler := map[string]func(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error{}

	auc := usecase.NewAbsenceUseCase(&backend)
	puc := usecase.NewPlayerUseCase(&backend)
	luc := usecase.NewLootUseCase(&backend)
	ruc := usecase.NewRaidUseCase(&backend)
	suc := usecase.NewStrikeUseCase(&backend)
	fuc := usecase.NewFailUseCase(&backend)

	disc := discordHandler.Discord{
		AbsenceUseCase: auc,
		PlayerUseCase:  puc,
		LootUseCase:    luc,
		RaidUseCase:    ruc,
		StrikeUseCase:  suc,
		FailUseCase:    fuc,
		Fake:           false,
	}

	var inits []func() map[string]func(
		ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error
	inits = append(inits,
		disc.InitAbsence, disc.InitAdmin, disc.InitLoot,
		disc.InitPlayer, disc.InitRaid, disc.InitStrike, disc.InitFail)
	for _, v := range inits {
		for k, v := range v() {
			mapHandler[k] = v
		}
	}

	var handlers []*discordgo.ApplicationCommand
	handlers = append(handlers,
		&discordHandler.AbsenceDescriptor[0], &discordHandler.AbsenceDescriptor[1], &discordHandler.AbsenceDescriptor[2])
	handlers = append(handlers,
		&discordHandler.LootDescriptors[0], &discordHandler.LootDescriptors[1],
		&discordHandler.LootDescriptors[2], &discordHandler.LootDescriptors[3])
	handlers = append(handlers,
		&discordHandler.PlayerDescriptors[0], &discordHandler.PlayerDescriptors[1],
		&discordHandler.PlayerDescriptors[2], &discordHandler.PlayerDescriptors[3], &discordHandler.PlayerDescriptors[4])
	handlers = append(handlers,
		&discordHandler.RaidDescriptors[0], &discordHandler.RaidDescriptors[1], &discordHandler.RaidDescriptors[2])
	handlers = append(handlers,
		&discordHandler.StrikeDescriptors[0], &discordHandler.StrikeDescriptors[1], &discordHandler.StrikeDescriptors[2])
	handlers = append(handlers,
		&discordHandler.FailDescriptors[0], &discordHandler.FailDescriptors[1],
		&discordHandler.FailDescriptors[2], &discordHandler.FailDescriptors[3])
	handlers = append(handlers,
		&discordHandler.AdminDescriptor[0], &discordHandler.AdminDescriptor[1])

	serve := discord.New(
		discord.CommandHandlers(mapHandler),
		discord.Token(cfg.Discord.Token),
		discord.Command(handlers),
		discord.GuildID(cfg.Discord.GuildID),
		discord.DeleteCommands(cfg.Discord.DeleteCommands))

	logger.FromContext(ctx).Info("start guildOps")
	err = serve.Run(ctx)
	if err != nil {
		logger.FromContext(ctx).Error(errors.Wrap(err, "run discord").Error())
		return
	}
}
