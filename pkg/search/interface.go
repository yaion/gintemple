package search

import "context"

// Engine defines the unified interface for search operations
type Engine interface {
	// Index adds or updates a document in the specified index
	Index(ctx context.Context, indexName string, docID string, doc interface{}) error

	// Delete removes a document from the specified index
	Delete(ctx context.Context, indexName string, docID string) error

	// Search performs a search query
	// This is a simplified search interface. Real-world cases might need more complex query builders.
	Search(ctx context.Context, indexName string, query string, options *SearchOptions) (*SearchResult, error)
}

type SearchOptions struct {
	Limit  int
	Offset int
	Filter interface{}
}

type SearchResult struct {
	Hits  []interface{} `json:"hits"`
	Total int64         `json:"total"`
}
