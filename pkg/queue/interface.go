package queue

import "context"

// Queue defines the unified interface for message queue operations
type Queue interface {
	// Publish sends a task/message to the queue
	Publish(ctx context.Context, topic string, payload []byte, options *PublishOptions) error

	// Subscribe registers a handler for a topic
	// Note: In some implementations (like Asynq), this might be done at server startup rather than dynamically
	Subscribe(topic string, handler func(ctx context.Context, payload []byte) error) error
}

type PublishOptions struct {
	Delay int64 // Delay in seconds
}
