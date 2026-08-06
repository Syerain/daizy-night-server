package dbware

import (
	"log/slog"

	"github.com/atomreforge/daizy-night-server/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ProviderDB struct {
	db *gorm.DB
}

func NewDBProvider(ctx struct {
	IsDebugMode bool
	DSN         string
}) (*ProviderDB, error) {
	logLevel := logger.Warn
	if ctx.IsDebugMode {
		logLevel = logger.Info
	}

	gormLogger := logger.NewSlogLogger(slog.Default(), logger.Config{
		LogLevel: logLevel,
	})

	// create or connect db
	db, err := gorm.Open(sqlite.Open(ctx.DSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.RegistercodeRecord{},
	); err != nil {
		return nil, err
	}

	return &ProviderDB{db: db}, nil
}

// check db health
func (p *ProviderDB) Check() error {
	sql, err := p.db.DB()
	if err != nil {
		return err
	}
	return sql.Ping()
}

func (p *ProviderDB) Close() error {
	sql, err := p.db.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}

func (p *ProviderDB) DB() *gorm.DB { return p.db }
