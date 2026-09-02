package dbware

import (
	"log/slog"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _ abstract.InterfaceProviderDB = (*ProviderDB)(nil)

type ProviderDB struct {
	db *gorm.DB
}

func NewProviderDB(ctx struct {
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
		Logger:         gormLogger,
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.RegistercodeRecord{},
		&model.Calendar{}, // AutoMigrate also creates the associated calendar_items table
	); err != nil {
		return nil, err
	}

	// legacy migration: databases created before the TokenHash removal keep a
	// NOT NULL token_hash column; leaving it behind would break every insert.
	// gorm's generic HasColumn probes INFORMATION_SCHEMA, which does not exist
	// on SQLite (it silently reports false), so the check goes through pragma.
	var cols []string
	if err := db.Raw(`SELECT name FROM pragma_table_info('refresh_tokens')`).Scan(&cols).Error; err != nil {
		return nil, err
	}
	for _, c := range cols {
		if c == "token_hash" {
			if err := db.Migrator().DropColumn(&model.RefreshToken{}, "token_hash"); err != nil {
				return nil, err
			}
			break
		}
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
