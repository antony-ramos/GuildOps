package main

import (
	"context"
	"fmt"
	"github.com/coven-discord-bot/config"
	"github.com/coven-discord-bot/internal/app"
	"github.com/coven-discord-bot/pkg/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	spanName := "main function"

	// Starting Log
	logger := zap.NewExample()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(logger)

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	zap.L().Info("replaced zap's global loggers")

	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error while loading config : %s", err.Error()))
		return
	}

	// Tracing
	zap.L().Info("Starting telemetry")

	// Tracing
	shutdown, err := tracing.InstallExportPipeline(ctx, cfg.Name, cfg.Version)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	ctx, span := otel.Tracer(spanName).Start(ctx, "main", trace.WithTimestamp(time.Now()))

	// Metrics
	zap.L().Info(fmt.Sprintf("Starting metrics server on port %s", cfg.Metrics.Port))
	http.Handle("/metrics", promhttp.Handler())
	go func() {

		err = http.ListenAndServe(":"+cfg.Metrics.Port, nil)
		if err != nil {
			span.RecordError(err)
			zap.L().Fatal(err.Error())
			span.End(trace.WithTimestamp(time.Now()))
		}
	}()

	// Run
	zap.L().Info("Starting app")
	app.Run(ctx, cfg)
	span.End(trace.WithTimestamp(time.Now()))

}
