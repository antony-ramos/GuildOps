package app

import (
	"context"
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/coven-discord-bot/config"
	"github.com/coven-discord-bot/internal/controller/api"
	discordHandler "github.com/coven-discord-bot/internal/controller/discord"
	"github.com/coven-discord-bot/internal/usecase"
	"github.com/coven-discord-bot/internal/usecase/backend_pg"
	"github.com/coven-discord-bot/pkg/discord"
	"github.com/coven-discord-bot/pkg/postgres"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(ctx context.Context, cfg *config.Config) {

	zap.L().Info("loading backend")

	pg, err := postgres.New(cfg.URL, postgres.MaxPoolSize(cfg.PoolMax), postgres.ConnAttempts(cfg.ConnAttempts), postgres.ConnTimeout(cfg.ConnTimeOut))
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	db := backend_pg.PG{Postgres: pg}
	db.Init(cfg.URL)

	mapHandler := map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error{}

	auc := usecase.NewAbsenceUseCase(&db)
	puc := usecase.NewPlayerUseCase(&db)
	luc := usecase.NewLootUseCase(&db)
	ruc := usecase.NewRaidUseCase(&db)
	suc := usecase.NewStrikeUseCase(&db)

	d := discordHandler.Discord{
		AbsenceUseCase: auc,
		PlayerUseCase:  puc,
		LootUseCase:    luc,
		RaidUseCase:    ruc,
		StrikeUseCase:  suc,
	}

	var inits []func() map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error
	inits = append(inits, d.InitAbsence, d.InitLoot, d.InitPlayer, d.InitRaid, d.InitStrike)
	for _, v := range inits {
		for k, v := range v() {
			mapHandler[k] = v
		}
	}

	var handlers []*discordgo.ApplicationCommand
	handlers = append(handlers, &discordHandler.AbsenceDescriptor[0], &discordHandler.AbsenceDescriptor[1], &discordHandler.AbsenceDescriptor[2])
	handlers = append(handlers, &discordHandler.LootDescriptors[0], &discordHandler.LootDescriptors[1], &discordHandler.LootDescriptors[2], &discordHandler.LootDescriptors[3])
	handlers = append(handlers, &discordHandler.PlayerDescriptors[0], &discordHandler.PlayerDescriptors[1], &discordHandler.PlayerDescriptors[2])
	handlers = append(handlers, &discordHandler.RaidDescriptors[0], &discordHandler.RaidDescriptors[1])
	handlers = append(handlers, &discordHandler.StrikeDescriptors[0], &discordHandler.StrikeDescriptors[1], &discordHandler.StrikeDescriptors[2])

	serve := discord.New(discord.CommandHandlers(mapHandler), discord.Token(cfg.Discord.Token), discord.Command(handlers), discord.GuildID(cfg.Discord.GuildID))

	server := api.API{
		Engine:         gin.Default(),
		AbsenceUseCase: auc,
		PlayerUseCase:  puc,
		LootUseCase:    luc,
		RaidUseCase:    ruc,
		StrikeUseCase:  suc,
	}
	server.Init()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: server.Engine,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	zap.L().Info("starting to serve to discord webhooks")
	err = serve.Run(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown", zap.Error(err))
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		zap.L().Info("timeout of 5 seconds.")
	}
	zap.L().Info("Server exiting")
}
