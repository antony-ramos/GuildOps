package main

import (
	"context"
	"fmt"
	"github.com/coven-discord-bot/config"
	"github.com/coven-discord-bot/internal/app"
	logger "github.com/coven-discord-bot/pkg/log"
	"github.com/coven-discord-bot/pkg/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()
	spanName := "main function"

	// Starting Log
	atom, factory, err := logger.Start()
	if err != nil {
		_ = fmt.Sprint(err)
		return
	}
	l := factory.For(ctx)

	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		l.Error(err.Error())
		return
	}

	// Modify Log Level
	level, ok := logger.LogLevels[strings.ToLower(cfg.Level)]
	if !ok {
		l.Error(fmt.Sprintf("Log level invalid : %s", cfg.Level))
		return
	}
	atom.SetLevel(level)

	// Tracing
	l.Info("Starting telemetry")

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
	l.Info(fmt.Sprintf("Starting metrics server on port %s", cfg.Metrics.Port))
	http.Handle("/metrics", promhttp.Handler())
	go func() {

		err = http.ListenAndServe(":"+cfg.Metrics.Port, nil)
		if err != nil {
			span.RecordError(err)
			l.Fatal(err.Error())
			span.End(trace.WithTimestamp(time.Now()))
		}
	}()

	// Run
	l.Info("Starting app")
	app.Run(ctx, cfg, &factory)
	span.End(trace.WithTimestamp(time.Now()))

}
