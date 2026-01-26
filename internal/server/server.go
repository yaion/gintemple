package server

import (
	"shop/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewServer(cfg *config.Config, logger *zap.Logger) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger()) // Basic gin logger

	return r
}
