package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/antony-ramos/guildops/config"
	"github.com/antony-ramos/guildops/internal/app"
	"github.com/antony-ramos/guildops/pkg/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LogLevels = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
	"fatal": zap.FatalLevel,
	"panic": zap.PanicLevel,
}

func main() {
	ctx := context.Background()
	spanName := "main function"

	// Configuration
	configPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("error while loading config : %s", err.Error())
		return
	}

	// Starting Log
	atom := zap.NewAtomicLevel()
	atom.SetLevel(LogLevels[cfg.Log.Level])

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(logger)

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	zap.L().Info("replaced zap's global loggers")

	// Setup Zap Log Level

	// Tracing
	zap.L().Info("Starting telemetry")

	// Tracing
	shutdown, err := tracing.InstallExportPipeline(ctx, cfg.Name, cfg.Version)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	ctx, span := otel.Tracer(spanName).Start(ctx, "main")

	// Metrics
	zap.L().Info(fmt.Sprintf("Starting metrics server on port %s", cfg.Metrics.Port))
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		server := &http.Server{
			Addr:              ":" + cfg.Metrics.Port,
			ReadHeaderTimeout: 3 * time.Second,
		}
		err := server.ListenAndServe()
		if err != nil {
			span.RecordError(err)
			zap.L().Fatal(err.Error())
			span.End()
		}
	}()

	// Run
	zap.L().Info("Starting app")
	app.Run(ctx, cfg)
	span.End()
}
