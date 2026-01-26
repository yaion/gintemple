package handler

import (
	"net/http"
	"shop/internal/infra/storage"
	"shop/internal/service"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	service service.FileService
}

func NewFileHandler(service service.FileService) *FileHandler {
	return &FileHandler{service: service}
}

// UploadSimple handles small file uploads
func (h *FileHandler) UploadSimple(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	folder := c.DefaultPostForm("folder", "default")
	url, err := h.service.UploadFile(c.Request.Context(), file, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// InitiateMultipart starts a multipart upload
func (h *FileHandler) InitiateMultipart(c *gin.Context) {
	var req struct {
		Filename string `json:"filename" binding:"required"`
		Folder   string `json:"folder"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Folder == "" {
		req.Folder = "default"
	}

	key, uploadID, err := h.service.InitiateMultipart(c.Request.Context(), req.Filename, req.Folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":       key,
		"upload_id": uploadID,
	})
}

// UploadPart uploads a chunk
func (h *FileHandler) UploadPart(c *gin.Context) {
	// We expect form-data with 'file', 'key', 'upload_id', 'part_number'
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file chunk is required"})
		return
	}

	var req struct {
		Key        string `form:"key" binding:"required"`
		UploadID   string `form:"upload_id" binding:"required"`
		PartNumber int    `form:"part_number" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	etag, err := h.service.UploadPart(c.Request.Context(), req.Key, req.UploadID, req.PartNumber, src, file.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"etag": etag})
}

// CompleteMultipart finishes the upload
func (h *FileHandler) CompleteMultipart(c *gin.Context) {
	var req struct {
		Key      string         `json:"key" binding:"required"`
		UploadID string         `json:"upload_id" binding:"required"`
		Parts    []storage.Part `json:"parts" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.service.CompleteMultipart(c.Request.Context(), req.Key, req.UploadID, req.Parts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
