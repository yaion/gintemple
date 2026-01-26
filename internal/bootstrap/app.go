package bootstrap

import (
	"context"
	"net/http"
	"shop/internal/config"
	"shop/internal/cron"
	"shop/internal/database"
	"shop/internal/handler"
	"shop/internal/infra/elasticsearch"
	"shop/internal/infra/rabbitmq"
	"shop/internal/infra/redis"
	"shop/internal/middleware"
	"shop/internal/repository"
	"shop/internal/router"
	"shop/internal/server"
	"shop/internal/service"
	"shop/internal/websocket"
	"shop/pkg/idgen"
	"shop/pkg/logger"

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
			redis.NewRedis,
			rabbitmq.NewRabbitMQ,
			elasticsearch.NewElasticsearch,
			idgen.NewIDGenerator,
			server.NewServer,
			middleware.NewMiddleware,
			repository.NewUserRepository,
			service.NewUserService,
			handler.NewUserHandler,
			cron.NewCronManager,
			websocket.NewHub,
		),
		fx.Invoke(
			router.RegisterRoutes,
			cron.StartCron,
			StartWebSocket,
			StartServer,
		),
	)
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
