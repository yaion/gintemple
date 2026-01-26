package asynq

import (
	"context"
	"shop/internal/config"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AsynqServer struct {
	Server *asynq.Server
	Mux    *asynq.ServeMux
	Logger *zap.Logger
}

func NewAsynqServer(cfg *config.Config, logger *zap.Logger) *AsynqServer {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Asynq.Addr,
			Password: cfg.Asynq.Password,
			DB:       cfg.Asynq.DB,
		},
		asynq.Config{
			Concurrency: cfg.Asynq.Concurrency,
			// You can add logger adapter here if needed
			// Logger: logger,
		},
	)

	mux := asynq.NewServeMux()

	return &AsynqServer{
		Server: srv,
		Mux:    mux,
		Logger: logger,
	}
}

// StartAsynqServer starts the asynq server
func StartAsynqServer(lc fx.Lifecycle, srv *AsynqServer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				srv.Logger.Info("Starting asynq server")
				if err := srv.Server.Run(srv.Mux); err != nil {
					srv.Logger.Fatal("could not run asynq server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.Logger.Info("Stopping asynq server")
			srv.Server.Stop()
			srv.Server.Shutdown()
			return nil
		},
	})
}
