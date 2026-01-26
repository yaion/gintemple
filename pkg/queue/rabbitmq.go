package queue

import (
	"context"
	"shop/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQAdapter struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger *zap.Logger
}

func NewRabbitMQAdapter(cfg *config.Config, logger *zap.Logger) (Queue, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQAdapter{conn: conn, ch: ch, logger: logger}, nil
}

func (r *RabbitMQAdapter) Publish(ctx context.Context, topic string, payload []byte, options *PublishOptions) error {
	// Declare a queue (idempotent)
	q, err := r.ch.QueueDeclare(
		topic, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = r.ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		})
	return err
}

func (r *RabbitMQAdapter) Subscribe(topic string, handler func(ctx context.Context, payload []byte) error) error {
	msgs, err := r.ch.Consume(
		topic, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			handler(context.Background(), d.Body)
		}
	}()
	return nil
}
