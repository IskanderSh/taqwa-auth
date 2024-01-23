package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/IskanderSh/taqwa-auth/internal/app"
	"github.com/IskanderSh/taqwa-auth/internal/config"
	"github.com/IskanderSh/taqwa-auth/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"

	debugLvl = "DEBUG"
	infoLvl  = "INFO"
	warnLvl  = "WARN"
	errorLvl = "ERROR"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg)

	log.Info("logger initialled successfully")

	application := app.New(log, cfg.GRPC.Port, &cfg.DB, &cfg.Token)
	log.Info("application created")

	// TODO: start gRPC server
	go application.GRPCServer.MustRun()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))
	application.GRPCServer.Stop()

	log.Info("application stopped")
}

func setupLogger(cfg *config.Config) *slog.Logger {
	var log *slog.Logger

	switch cfg.Env {
	case envLocal:
		log = setupPrettySlog(getLogLevel(cfg.LogLevel))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: getLogLevel(cfg.LogLevel)}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: getLogLevel(cfg.LogLevel)}))
	}

	return log
}

func setupPrettySlog(level slog.Level) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

func getLogLevel(lvl string) slog.Level {
	var res slog.Level

	switch strings.ToUpper(lvl) {
	case debugLvl:
		res = slog.LevelDebug
	case infoLvl:
		res = slog.LevelInfo
	case warnLvl:
		res = slog.LevelWarn
	case errorLvl:
		res = slog.LevelError
	}

	return res
}
