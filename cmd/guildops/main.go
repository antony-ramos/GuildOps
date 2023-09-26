package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	logger "github.com/antony-ramos/guildops/pkg/logger"

	"github.com/antony-ramos/guildops/config"
	"github.com/antony-ramos/guildops/internal/app"
	"github.com/antony-ramos/guildops/pkg/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	// Configuration
	configPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("error while loading config : %s", err.Error())
		return
	}

	// Logs
	atom := zap.NewAtomicLevel()
	atom.SetLevel(LogLevels[cfg.Log.Level])

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	zapLog := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(zapLog)

	zapLog = zapLog.With(zap.String("service", cfg.Name), zap.String("version", cfg.Version), zap.String("env", cfg.Env))

	ctx = logger.AddLoggerToContext(ctx, zapLog)

	// Tracing
	logger.FromContext(ctx).Info("Starting telemetry")
	logger.FromContext(ctx).Info("Starting telemetry")

	shutdown, err := tracing.InstallExportPipeline(ctx, cfg.Name, cfg.Version)
	if err != nil {
		logger.FromContext(ctx).Error(err.Error())
		return
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Metrics
	logger.FromContext(ctx).Info(fmt.Sprintf("Starting metrics server on port %s", cfg.Metrics.Port))
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		server := &http.Server{
			Addr:              ":" + cfg.Metrics.Port,
			ReadHeaderTimeout: 3 * time.Second,
		}
		err := server.ListenAndServe()
		if err != nil {
			logger.FromContext(ctx).Fatal(err.Error())
		}
	}()

	// Run
	logger.FromContext(ctx).Info("Starting app")
	app.Run(ctx, cfg)
}
