package storage

import (
	"context"
	"io"
)

// Part represents a part of a multipart upload
type Part struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
}

// Provider defines the interface for file storage
type Provider interface {
	// Simple upload
	PutObject(ctx context.Context, key string, data io.Reader, size int64) (string, error)

	// Multipart upload
	InitiateMultipartUpload(ctx context.Context, key string) (string, error)
	UploadPart(ctx context.Context, key string, uploadID string, partNumber int, data io.Reader, size int64) (string, error)
	CompleteMultipartUpload(ctx context.Context, key string, uploadID string, parts []Part) (string, error)

	// Utility
	GetURL(key string) string
}
