package search

import "context"

// Engine defines the common interface for search engines
type Engine interface {
	Index(ctx context.Context, indexName string, docID string, doc interface{}) error
	Search(ctx context.Context, indexName string, query string, opts ...SearchOption) (*Result, error)
	Delete(ctx context.Context, indexName string, docID string) error
}

type SearchOption func(o *SearchOptions)

type SearchOptions struct {
	Limit  int
	Offset int
	Filter string
}

type Result struct {
	Hits  []interface{} `json:"hits"`
	Total int64         `json:"total"`
}

func WithLimit(limit int) SearchOption {
	return func(o *SearchOptions) {
		o.Limit = limit
	}
}

func WithOffset(offset int) SearchOption {
	return func(o *SearchOptions) {
		o.Offset = offset
	}
}
