package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"shop/internal/config"
	"shop/internal/infra/storage"
	"shop/pkg/utils"

	"go.uber.org/zap"
)

type LocalStorage struct {
	baseDir string
	baseURL string
	logger  *zap.Logger
}

func NewLocalStorage(cfg *config.Config, logger *zap.Logger) storage.Provider {
	// Create upload directory if not exists
	baseDir := cfg.Storage.Local.Path
	if baseDir == "" {
		baseDir = "./uploads"
	}
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		logger.Fatal("failed to create storage directory", zap.Error(err))
	}

	// Create temp directory for multipart uploads
	tempDir := filepath.Join(baseDir, "temp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		logger.Fatal("failed to create storage temp directory", zap.Error(err))
	}

	return &LocalStorage{
		baseDir: baseDir,
		baseURL: cfg.Storage.Local.URL, // e.g., "http://localhost:8080/uploads"
		logger:  logger,
	}
}

func (s *LocalStorage) PutObject(ctx context.Context, key string, data io.Reader, size int64) (string, error) {
	fullPath := filepath.Join(s.baseDir, key)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, data); err != nil {
		return "", err
	}

	return s.GetURL(key), nil
}

func (s *LocalStorage) InitiateMultipartUpload(ctx context.Context, key string) (string, error) {
	// For local storage, uploadID can just be a unique ID
	// We will create a directory for this uploadID in temp
	uploadID := utils.GenerateUUID()
	tempPath := filepath.Join(s.baseDir, "temp", uploadID)
	if err := os.MkdirAll(tempPath, 0755); err != nil {
		return "", err
	}
	return uploadID, nil
}

func (s *LocalStorage) UploadPart(ctx context.Context, key string, uploadID string, partNumber int, data io.Reader, size int64) (string, error) {
	tempPath := filepath.Join(s.baseDir, "temp", uploadID, fmt.Sprintf("%d", partNumber))

	out, err := os.Create(tempPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, data); err != nil {
		return "", err
	}

	// Return path or a simple success indicator as ETag
	return fmt.Sprintf("part-%d", partNumber), nil
}

func (s *LocalStorage) CompleteMultipartUpload(ctx context.Context, key string, uploadID string, parts []storage.Part) (string, error) {
	finalPath := filepath.Join(s.baseDir, key)

	// Ensure final directory exists
	dir := filepath.Dir(finalPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	out, err := os.Create(finalPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Combine parts
	for _, part := range parts {
		partPath := filepath.Join(s.baseDir, "temp", uploadID, fmt.Sprintf("%d", part.PartNumber))

		in, err := os.Open(partPath)
		if err != nil {
			return "", err
		}

		if _, err := io.Copy(out, in); err != nil {
			in.Close()
			return "", err
		}
		in.Close()
	}

	// Cleanup temp dir
	os.RemoveAll(filepath.Join(s.baseDir, "temp", uploadID))

	return s.GetURL(key), nil
}

func (s *LocalStorage) GetURL(key string) string {
	// If key starts with slash, remove it to avoid double slash
	if len(key) > 0 && key[0] == '/' {
		key = key[1:]
	}
	return fmt.Sprintf("%s/%s", s.baseURL, key)
}
