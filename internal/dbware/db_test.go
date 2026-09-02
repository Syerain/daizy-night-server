package dbware

import (
	"testing"
	"time"

	"github.com/atomreforge/daizy-night-server/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestLegacyTokenHashColumnIsDropped simulates a database created before the
// TokenHash removal. The legacy table declares token_hash through a
// table-level UNIQUE constraint plus a separate unique index for lookup_hash
// — exactly what earlier gorm AutoMigrate runs produced — and asserts that
// opening it through NewProviderDB rebuilds the table without the column,
// keeping the rows and staying writable afterwards.
func TestLegacyTokenHashColumnIsDropped(t *testing.T) {
	dsn := "file:dntest_legacy_tokenhash?mode=memory&cache=shared"

	// simulate the legacy database
	legacy, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open legacy db: %v", err)
	}
	if err := legacy.Exec("CREATE TABLE `refresh_tokens` ("+
		"`id` integer PRIMARY KEY AUTOINCREMENT,"+
		"`created_at` datetime,"+
		"`updated_at` datetime,"+
		"`deleted_at` datetime,"+
		"`uid` integer NOT NULL,"+
		"`token_hash` text NOT NULL,"+
		"`lookup_hash` text NOT NULL,"+
		"`revoked_at` datetime,"+
		"CONSTRAINT `uni_refresh_tokens_token_hash` UNIQUE (`token_hash`))").Error; err != nil {
		t.Fatalf("create legacy table: %v", err)
	}
	if err := legacy.Exec("CREATE UNIQUE INDEX `idx_refresh_tokens_lookup_hash` ON `refresh_tokens`(`lookup_hash`)").Error; err != nil {
		t.Fatalf("create legacy index: %v", err)
	}
	if err := legacy.Exec("INSERT INTO `refresh_tokens`" +
		"(`uid`, `token_hash`, `lookup_hash`, `created_at`, `updated_at`)" +
		"VALUES (1, 'legacy-hash', 'legacy-lookup', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)").Error; err != nil {
		t.Fatalf("insert legacy row: %v", err)
	}

	// open through the provider: AutoMigrate + legacy rebuild
	p, err := NewProviderDB(struct {
		IsDebugMode bool
		DSN         string
	}{false, dsn})
	if err != nil {
		t.Fatalf("provider on legacy db: %v", err)
	}
	defer p.Close()

	var cols []string
	if err := p.DB().Raw(`SELECT name FROM pragma_table_info('refresh_tokens')`).Scan(&cols).Error; err != nil {
		t.Fatalf("read columns: %v", err)
	}
	for _, c := range cols {
		if c == "token_hash" {
			t.Fatalf("token_hash column still present, columns = %v", cols)
		}
	}

	// the legacy row must survive the rebuild
	var n int64
	if err := p.DB().Raw(`SELECT count(*) FROM refresh_tokens WHERE lookup_hash = 'legacy-lookup'`).Scan(&n).Error; err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if n != 1 {
		t.Fatalf("legacy row lost after migration, count = %d", n)
	}

	// and the rebuilt table must be writable through the new schema
	cfg := &config.Config{}
	cfg.Security.JwtRefreshTokenExpireTime = 168 * time.Hour
	cfg.Security.JwtRevokedTokensRetainTime = 72 * time.Hour
	rt := &RepoToken{pDB: p, cfg: cfg}
	if err := rt.SaveRefreshToken(1, "fresh-token"); err != nil {
		t.Fatalf("insert into rebuilt table: %v", err)
	}
}
