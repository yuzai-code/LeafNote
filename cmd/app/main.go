package main

import (
	"fmt"
	"log"

	"leafnote/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	logger *zap.Logger
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化日志系统
	logger, err = config.InitLogger(&cfg.Log)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化 Gin
	r := gin.Default()

	// 初始化数据库
	db, err = config.InitDB(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// 获取底层的 sqlDB
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get underlying *sql.DB", zap.Error(err))
	}
	defer sqlDB.Close()

	// TODO: 初始化路由
	// TODO: 初始化中间件

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Server starting", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatal("Server startup failed", zap.Error(err))
	}
}
