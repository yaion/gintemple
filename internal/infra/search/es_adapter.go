package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"shop/internal/config"
	"shop/internal/infra/search"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

type ESAdapter struct {
	client *elasticsearch.Client
	logger *zap.Logger
}

func NewESAdapter(cfg *config.Config, logger *zap.Logger) (search.Engine, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Elasticsearch.Addresses,
		Username:  cfg.Elasticsearch.Username,
		Password:  cfg.Elasticsearch.Password,
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, err
	}

	return &ESAdapter{client: client, logger: logger}, nil
}

func (e *ESAdapter) Index(ctx context.Context, indexName string, docID string, doc interface{}) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	res, err := e.client.Index(
		indexName,
		bytes.NewReader(data),
		e.client.Index.WithDocumentID(docID),
		e.client.Index.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return res.String()
	}
	return nil
}

func (e *ESAdapter) Search(ctx context.Context, indexName string, query string, opts ...search.SearchOption) (*search.Result, error) {
	options := &search.SearchOptions{Limit: 20}
	for _, o := range opts {
		o(options)
	}

	// Simplified match query
	queryBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": query,
			},
		},
		"from": options.Offset,
		"size": options.Limit,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(queryBody); err != nil {
		return nil, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
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

	hits := r["hits"].(map[string]interface{})
	total := hits["total"].(map[string]interface{})["value"].(float64)
	hitsList := hits["hits"].([]interface{})

	// Extract source from hits
	var results []interface{}
	for _, hit := range hitsList {
		results = append(results, hit.(map[string]interface{})["_source"])
	}

	return &search.Result{
		Hits:  results,
		Total: int64(total),
	}, nil
}

func (e *ESAdapter) Delete(ctx context.Context, indexName string, docID string) error {
	res, err := e.client.Delete(indexName, docID, e.client.Delete.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
