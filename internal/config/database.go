package config

import (
	"leafnote/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB(cfg *DatabaseConfig) (*gorm.DB, error) {
	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接数据库
	db, err := gorm.Open(sqlite.Open(cfg.Name), gormConfig)
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库结构
	err = db.AutoMigrate(
		&model.Note{},
		&model.Category{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
