package queue

import (
	"context"
	"time"
)

// TaskQueue defines the common interface for message queues
type TaskQueue interface {
	Enqueue(ctx context.Context, taskName string, payload []byte, opts ...QueueOption) error
	RegisterHandler(taskName string, handler func(ctx context.Context, payload []byte) error)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type QueueOption func(o *QueueOptions)

type QueueOptions struct {
	Delay time.Duration
}

func WithDelay(d time.Duration) QueueOption {
	return func(o *QueueOptions) {
		o.Delay = d
	}
}
