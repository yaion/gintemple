package queue

import (
	"context"
	"shop/internal/config"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type AsynqAdapter struct {
	client *asynq.Client
	logger *zap.Logger
}

func NewAsynqAdapter(cfg *config.Config, logger *zap.Logger) (Queue, error) {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Asynq.Addr,
		Password: cfg.Asynq.Password,
		DB:       cfg.Asynq.DB,
	})
	return &AsynqAdapter{client: client, logger: logger}, nil
}

func (a *AsynqAdapter) Publish(ctx context.Context, topic string, payload []byte, options *PublishOptions) error {
	task := asynq.NewTask(topic, payload)
	var opts []asynq.Option
	if options != nil && options.Delay > 0 {
		opts = append(opts, asynq.ProcessIn(time.Duration(options.Delay)*time.Second))
	}
	_, err := a.client.Enqueue(task, opts...)
	return err
}

func (a *AsynqAdapter) Subscribe(topic string, handler func(ctx context.Context, payload []byte) error) error {
	// Note: Asynq server handles subscription in a different way (Mux).
	// This generic interface might need adjustment or the adapter needs to hold a reference to the Server Mux
	// For simplicity, we assume this is called during setup to register handlers to a global/passed Mux
	// Or we simply log that dynamic subscription isn't fully supported in this simplified adapter
	a.logger.Warn("Subscribe called on AsynqAdapter - ensure handlers are registered with AsynqServer Mux")
	return nil
}
