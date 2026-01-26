package search

import (
	"bytes"
	"context"
	"encoding/json"
	"shop/internal/config"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

type ESAdapter struct {
	client *elasticsearch.Client
	logger *zap.Logger
}

func NewESAdapter(cfg *config.Config, logger *zap.Logger) (Engine, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Elasticsearch.Addresses,
		Username:  cfg.Elasticsearch.Username,
		Password:  cfg.Elasticsearch.Password,
	})
	if err != nil {
		return nil, err
	}
	return &ESAdapter{client: es, logger: logger}, nil
}

func (e *ESAdapter) Index(ctx context.Context, indexName string, docID string, doc interface{}) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	_, err = e.client.Index(indexName, bytes.NewReader(data), e.client.Index.WithDocumentID(docID))
	return err
}

func (e *ESAdapter) Delete(ctx context.Context, indexName string, docID string) error {
	_, err := e.client.Delete(indexName, docID)
	return err
}

func (e *ESAdapter) Search(ctx context.Context, indexName string, query string, options *SearchOptions) (*SearchResult, error) {
	// Simplified ES search query
	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": query,
			},
		},
		"from": options.Offset,
		"size": options.Limit,
	}
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := e.client.Search(
		e.client.Search.WithIndex(indexName),
		e.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	// Parse hits... this is a simplified example
	hits := r["hits"].(map[string]interface{})
	total := hits["total"].(map[string]interface{})["value"].(float64)
	hitsList := hits["hits"].([]interface{})

	var results []interface{}
	for _, hit := range hitsList {
		source := hit.(map[string]interface{})["_source"]
		results = append(results, source)
	}

	return &SearchResult{
		Hits:  results,
		Total: int64(total),
	}, nil
}
