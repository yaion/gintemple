package rabbitmq

import (
	"shop/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	Conn *amqp.Connection
}

func NewRabbitMQ(cfg *config.Config, logger *zap.Logger) (*RabbitMQ, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		logger.Error("failed to connect to rabbitmq", zap.Error(err))
		return nil, err
	}

	return &RabbitMQ{Conn: conn}, nil
}

func (r *RabbitMQ) Close() error {
	return r.Conn.Close()
}
