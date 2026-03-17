package database

import (
	"fmt"
	"sync"
	"time"

	"erp-server/pkg/config"
	"erp-server/pkg/log"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	db         *gorm.DB
	cfg        *config.DatabaseConfig
	mu         sync.RWMutex
)

// Init 初始化数据库连接
func Init(databaseCfg *config.DatabaseConfig) error {
	mu.Lock()
	defer mu.Unlock()

	cfg = databaseCfg // 保存配置，用于后续重连
	return connect()
}

// connect 建立数据库连接
func connect() error {
	var err error

	// 配置 GORM 日志
	gormConfig := &gorm.Config{}
	if config.Get().Server.Mode == "debug" {
		gormConfig.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	} else {
		gormConfig.Logger = gormlogger.Default.LogMode(gormlogger.Silent)
	}

	// 连接数据库
	db, err = gorm.Open(mysql.Open(cfg.GetDSN()), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层的 sql.DB 以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)                      // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                     // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)            // 连接最大存活时间（1小时后重建连接）
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)     // 空闲连接最大存活时间（10分钟不使用就关闭）

	log.Info("数据库连接成功")
	return nil
}

// Get 获取数据库连接（带自动重连）
func Get() *gorm.DB {
	mu.RLock()
	if db != nil {
		// 检查连接是否有效
		sqlDB, err := db.DB()
		if err == nil {
			// Ping 检测连接是否有效
			if err = sqlDB.Ping(); err == nil {
				mu.RUnlock()
				return db
			}
			log.Warn("数据库连接失效，尝试重连", zap.Error(err))
		}
	}
	mu.RUnlock()

	// 连接失效，尝试重连
	mu.Lock()
	defer mu.Unlock()

	// 双重检查，避免并发重复重连
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil && sqlDB.Ping() == nil {
			return db
		}
	}

	// 重连
	log.Info("开始重新连接数据库...")
	if err := connect(); err != nil {
		log.Error("数据库重连失败", zap.Error(err))
		return nil
	}

	log.Info("数据库重连成功")
	return db
}

// Close 关闭数据库连接
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// Ping 检查数据库连接是否正常
func Ping() error {
	mu.RLock()
	defer mu.RUnlock()

	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}
