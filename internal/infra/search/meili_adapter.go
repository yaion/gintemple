package meilisearch

import (
	"context"
	"shop/internal/config"
	"shop/internal/infra/search"

	meilisearchgo "github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

type MeiliAdapter struct {
	client meilisearchgo.ServiceManager
	logger *zap.Logger
}

func NewMeiliAdapter(cfg *config.Config, logger *zap.Logger) (search.Engine, error) {
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

	return &MeiliAdapter{client: client, logger: logger}, nil
}

func (m *MeiliAdapter) Index(ctx context.Context, indexName string, docID string, doc interface{}) error {
	index := m.client.Index(indexName)
	// Meilisearch expects a list of documents or a single document map
	// If doc is struct, we might need to wrap it in array for AddDocuments
	task, err := index.AddDocuments([]interface{}{doc})
	if err != nil {
		return err
	}
	m.logger.Debug("Indexed document", zap.Int64("taskUID", task.TaskUID))
	return nil
}

func (m *MeiliAdapter) Search(ctx context.Context, indexName string, query string, opts ...search.SearchOption) (*search.Result, error) {
	options := &search.SearchOptions{Limit: 20}
	for _, o := range opts {
		o(options)
	}

	req := &meilisearchgo.SearchRequest{
		Limit:  int64(options.Limit),
		Offset: int64(options.Offset),
		Filter: options.Filter,
	}

	resp, err := m.client.Index(indexName).Search(query, req)
	if err != nil {
		return nil, err
	}

	return &search.Result{
		Hits:  resp.Hits,
		Total: resp.EstimatedTotalHits,
	}, nil
}

func (m *MeiliAdapter) Delete(ctx context.Context, indexName string, docID string) error {
	_, err := m.client.Index(indexName).DeleteDocument(docID)
	return err
}
