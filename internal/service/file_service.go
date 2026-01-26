package service

import (
	"context"
	"io"
	"mime/multipart"
	"path/filepath"
	"shop/internal/infra/storage"
	"shop/pkg/utils"
	"time"
)

type FileService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
	InitiateMultipart(ctx context.Context, filename string, folder string) (string, string, error)
	UploadPart(ctx context.Context, key string, uploadID string, partNumber int, file io.Reader, size int64) (string, error)
	CompleteMultipart(ctx context.Context, key string, uploadID string, parts []storage.Part) (string, error)
}

type fileService struct {
	provider storage.Provider
}

func NewFileService(provider storage.Provider) FileService {
	return &fileService{provider: provider}
}

func (s *fileService) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := utils.GenerateUUID() + ext
	key := filepath.Join(folder, time.Now().Format("20060102"), filename)

	return s.provider.PutObject(ctx, key, src, file.Size)
}

func (s *fileService) InitiateMultipart(ctx context.Context, filename string, folder string) (string, string, error) {
	ext := filepath.Ext(filename)
	newFilename := utils.GenerateUUID() + ext
	key := filepath.Join(folder, time.Now().Format("20060102"), newFilename)

	uploadID, err := s.provider.InitiateMultipartUpload(ctx, key)
	if err != nil {
		return "", "", err
	}
	return key, uploadID, nil
}

func (s *fileService) UploadPart(ctx context.Context, key string, uploadID string, partNumber int, file io.Reader, size int64) (string, error) {
	return s.provider.UploadPart(ctx, key, uploadID, partNumber, file, size)
}

func (s *fileService) CompleteMultipart(ctx context.Context, key string, uploadID string, parts []storage.Part) (string, error) {
	return s.provider.CompleteMultipartUpload(ctx, key, uploadID, parts)
}
