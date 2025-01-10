package testutil

import (
	"leafnote/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// StringPtr 返回字符串的指针
func StringPtr(s string) *string {
	return &s
}

// SetupTestDB 设置测试数据库
func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(
		&model.Note{},
		&model.Tag{},
		&model.Category{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
