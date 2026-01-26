package cron

import (
	"context"

	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// CronManager handles background tasks
type CronManager struct {
	scheduler *cron.Cron
	logger    *zap.Logger
}

func NewCronManager(logger *zap.Logger) *CronManager {
	// Create a new cron scheduler with second-level precision
	c := cron.New(cron.WithSeconds())
	return &CronManager{
		scheduler: c,
		logger:    logger,
	}
}

// RegisterJobs registers all cron jobs
func (m *CronManager) RegisterJobs() {
	// Example: Run every 10 seconds
	_, err := m.scheduler.AddFunc("*/10 * * * * *", func() {
		m.logger.Info("Executing example cron job")
	})

	if err != nil {
		m.logger.Error("Failed to register example job", zap.Error(err))
	}
}

// StartCron starts the cron scheduler using Fx Lifecycle
func StartCron(lc fx.Lifecycle, m *CronManager) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			m.logger.Info("Starting cron scheduler")
			m.RegisterJobs()
			m.scheduler.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			m.logger.Info("Stopping cron scheduler")
			ctx = m.scheduler.Stop() // Returns a context that is done when all jobs have completed
			<-ctx.Done()
			return nil
		},
	})
}
