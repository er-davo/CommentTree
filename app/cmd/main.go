package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"comment-tree/internal/app"
	"comment-tree/internal/config"
	"comment-tree/internal/database"

	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.InitConsole()
	log := zlog.Logger

	configFilePath := os.Getenv("CONFIG_PATH")
	if configFilePath == "" {
		log.Fatal().Msg("CONFIG_PATH environment variable is not set")
	}

	cfg, err := config.Load(configFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	err = database.Migrate(cfg.App.MigrationDir, cfg.DB.URL)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to migrate database")
	}

	urlShrtApp, err := app.NewCommentsTreeApp(cfg, &log)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to create app")
	}

	log.Info().Msg("comments tree is running")
	urlShrtApp.Run(ctx)

	<-ctx.Done()
}
