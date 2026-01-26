package meilisearch

import (
	"shop/internal/config"

	meilisearchgo "github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

func NewMeilisearchClient(cfg *config.Config, logger *zap.Logger) meilisearchgo.ServiceManager {
	client := meilisearchgo.New(
		cfg.Meilisearch.Host,
		meilisearchgo.WithAPIKey(cfg.Meilisearch.APIKey),
	)

	// Optional: Health check
	if _, err := client.Health(); err != nil {
		logger.Warn("Failed to connect to Meilisearch", zap.Error(err))
	} else {
		logger.Info("Connected to Meilisearch", zap.String("host", cfg.Meilisearch.Host))
	}

	return client
}
