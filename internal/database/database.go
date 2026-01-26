package database

import (
	"shop/internal/config"
	"shop/internal/model"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migrate
	if err := db.AutoMigrate(&model.User{}); err != nil {
		logger.Error("failed to auto migrate", zap.Error(err))
		return nil, err
	}

	return db, nil
}
