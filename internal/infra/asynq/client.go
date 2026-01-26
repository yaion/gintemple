package asynq

import (
	"shop/internal/config"

	"github.com/hibiken/asynq"
)

func NewAsynqClient(cfg *config.Config) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Asynq.Addr,
		Password: cfg.Asynq.Password,
		DB:       cfg.Asynq.DB,
	})
}
