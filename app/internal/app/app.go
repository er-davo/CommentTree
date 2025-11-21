package app

import (
	"context"

	"comment-tree/internal/config"
	"comment-tree/internal/database"
	"comment-tree/internal/handler"
	"comment-tree/internal/repository"
	"comment-tree/internal/service"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

type CommentsTreeApp struct {
	cfg *config.Config

	engine *ginext.Engine

	log *zlog.Zerolog
}

func NewCommentsTreeApp(cfg *config.Config, log *zlog.Zerolog) (*CommentsTreeApp, error) {
	r := ginext.New("release")

	strategy := retry.Strategy{
		Attempts: cfg.Retry.Attempts,
		Delay:    cfg.Retry.Delay,
		Backoff:  cfg.Retry.Backoff,
	}

	opts := &dbpg.Options{
		MaxOpenConns:    cfg.DB.MaxOpenConns,
		MaxIdleConns:    cfg.DB.MaxIdleConns,
		ConnMaxLifetime: cfg.DB.ConnMaxLifetime,
	}

	db, err := database.Connect(cfg.DB.URL, []string{}, opts)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to connect to database")
		return nil, err
	}

	comRepo := repository.NewCommentsRepository(db, strategy)

	comService := service.NewCommentsService(comRepo, log)

	comHandler := handler.NewCommentsHandler(comService, log)

	r.GET("/", func(c *ginext.Context) {
		c.File("public/index.html")
	})

	r.Engine.Use(ginext.Logger())
	r.Engine.Use(ginext.Recovery())

	comHandler.RegisterRoutes(r)

	return &CommentsTreeApp{
		cfg:    cfg,
		engine: r,
		log:    log,
	}, nil
}

func (a *CommentsTreeApp) Run(ctx context.Context) {
	if err := a.engine.Run(":" + a.cfg.App.Port); err != nil {
		a.log.Error().
			Err(err).
			Msg("failed to run server")
		return
	}
}
