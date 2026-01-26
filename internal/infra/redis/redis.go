package redis

import (
	"context"
	"shop/internal/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedis(cfg *config.Config, logger *zap.Logger) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Error("failed to connect to redis", zap.Error(err))
		return nil, err
	}

	return rdb, nil
}
