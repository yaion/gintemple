package search

import (
	"context"
	"shop/internal/config"

	meilisearchgo "github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

type MeiliAdapter struct {
	client meilisearchgo.ServiceManager
	logger *zap.Logger
}

func NewMeiliAdapter(cfg *config.Config, logger *zap.Logger) (Engine, error) {
	client := meilisearchgo.New(
		cfg.Meilisearch.Host,
		meilisearchgo.WithAPIKey(cfg.Meilisearch.APIKey),
	)

	return &MeiliAdapter{client: client, logger: logger}, nil
}

func (m *MeiliAdapter) Index(ctx context.Context, indexName string, docID string, doc interface{}) error {
	_, err := m.client.Index(indexName).AddDocuments([]interface{}{doc}, nil)
	return err
}

func (m *MeiliAdapter) Delete(ctx context.Context, indexName string, docID string) error {
	_, err := m.client.Index(indexName).DeleteDocument(docID, nil)
	return err
}

func (m *MeiliAdapter) Search(ctx context.Context, indexName string, query string, options *SearchOptions) (*SearchResult, error) {
	req := &meilisearchgo.SearchRequest{
		Limit:  int64(options.Limit),
		Offset: int64(options.Offset),
	}
	// Note: Filter handling in Meilisearch might require string parsing/building

	resp, err := m.client.Index(indexName).Search(query, req)
	if err != nil {
		return nil, err
	}

	// Convert Hits (which is []interface{} in latest SDK or custom type) to []interface{}
	// The SDK defines Hits as []interface{} usually, but let's be safe
	var hits []interface{}
	for _, hit := range resp.Hits {
		hits = append(hits, hit)
	}

	return &SearchResult{
		Hits:  hits,
		Total: resp.EstimatedTotalHits,
	}, nil
}
