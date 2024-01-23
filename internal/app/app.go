package app

import (
	"log/slog"

	grpcapp "github.com/IskanderSh/taqwa-auth/internal/app/grpc"
	"github.com/IskanderSh/taqwa-auth/internal/config"
	"github.com/IskanderSh/taqwa-auth/internal/services/auth"
	"github.com/IskanderSh/taqwa-auth/internal/storage/mongoDB"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	dbConfig *config.DB,
	token *config.Token,
) *App {
	storage, err := mongoDB.New(dbConfig)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, token)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
