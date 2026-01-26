package bootstrap

import (
	"context"
	"net/http"
	"shop/internal/config"
	"shop/internal/cron"
	"shop/internal/database"
	"shop/internal/handler"
	"shop/internal/infra/asynq"
	"shop/internal/infra/storage/local"
	"shop/internal/middleware"
	"shop/internal/repository"
	"shop/internal/router"
	"shop/internal/server"
	"shop/internal/service"
	"shop/internal/websocket"
	"shop/pkg/idgen"
	"shop/pkg/logger"
	"shop/pkg/queue"
	"shop/pkg/search"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			config.NewConfig,
			logger.NewLogger,
			database.NewDatabase,
			asynq.NewAsynqServer,

			// Interfaces
			ProvideSearchEngine,
			ProvideQueue,

			// Storage provider based on config (currently simplified to always provide local)
			// In a real app, use a factory function to choose provider based on config
			local.NewLocalStorage,
			idgen.NewIDGenerator,
			server.NewServer,
			middleware.NewMiddleware,
			repository.NewUserRepository,
			service.NewUserService,
			service.NewFileService,
			handler.NewUserHandler,
			handler.NewFileHandler,
			cron.NewCronManager,
			websocket.NewHub,
		),
		fx.Invoke(
			router.RegisterRoutes,
			cron.StartCron,
			asynq.StartAsynqServer,
			StartWebSocket,
			StartServer,
		),
	)
}

func ProvideSearchEngine(cfg *config.Config, logger *zap.Logger) (search.Engine, error) {
	// Simple switch based on config presence or a specific flag
	// For example, if Meilisearch host is set, use it, else if ES addresses set, use ES
	if cfg.Meilisearch.Host != "" {
		logger.Info("Using Meilisearch as search engine")
		return search.NewMeiliAdapter(cfg, logger)
	}
	if len(cfg.Elasticsearch.Addresses) > 0 {
		logger.Info("Using Elasticsearch as search engine")
		return search.NewESAdapter(cfg, logger)
	}
	// Default or Error
	return nil, nil // Or return a NoOp engine
}

func ProvideQueue(cfg *config.Config, logger *zap.Logger) (queue.Queue, error) {
	// Similar logic for Queue
	if cfg.RabbitMQ.URL != "" {
		logger.Info("Using RabbitMQ as message queue")
		return queue.NewRabbitMQAdapter(cfg, logger)
	}
	if cfg.Asynq.Addr != "" {
		logger.Info("Using Asynq as message queue")
		return queue.NewAsynqAdapter(cfg, logger)
	}
	return nil, nil
}

func StartWebSocket(lc fx.Lifecycle, hub *websocket.Hub) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go hub.Run()
			return nil
		},
	})
}

func StartServer(lc fx.Lifecycle, r *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting server", zap.String("port", cfg.Server.Port))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("listen: %s\n", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server")
			return srv.Shutdown(ctx)
		},
	})
}
