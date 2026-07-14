package db

import (
	"daizynight/internal/config"
	"daizynight/internal/model"
	"log/slog"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.Config) error {
	// adapte gorm logger
	logLevel := logger.Warn
	if cfg.Database.IsDebugMode {
		logLevel = logger.Info
	}
	gormLogger := logger.NewSlogLogger(slog.Default(), logger.Config{
		LogLevel: logLevel,
	})

	// create or connect db
	db, err := gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return err
	}
	DB = db
	return nil
}

// check db health
func CheckDB() error {
	sql, err := DB.DB()
	if err != nil {
		return err
	}
	return sql.Ping()
}

func CloseDB() error {
	sql, err := DB.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}
