package dbware

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestLegacyTokenHashColumnIsDropped simulates a database created before the
// TokenHash removal (NOT NULL token_hash column) and asserts that opening it
// through NewProviderDB drops the legacy column without losing rows. Keeping
// the NOT NULL column would break every insert of the new schema.
func TestLegacyTokenHashColumnIsDropped(t *testing.T) {
	dsn := "file:dntest_legacy_tokenhash?mode=memory&cache=shared"

	// simulate the legacy database: refresh_tokens still has token_hash
	legacy, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open legacy db: %v", err)
	}
	if err := legacy.Exec("CREATE TABLE `refresh_tokens` (\n"+
		"`id` integer PRIMARY KEY AUTOINCREMENT,\n"+
		"`created_at` datetime,\n"+
		"`updated_at` datetime,\n"+
		"`deleted_at` datetime,\n"+
		"`uid` numeric NOT NULL,\n"+
		"`token_hash` text NOT NULL UNIQUE,\n"+
		"`lookup_hash` text NOT NULL,\n"+
		"`revoked_at` datetime\n"+
		")").Error; err != nil {
		t.Fatalf("create legacy table: %v", err)
	}
	if err := legacy.Exec(`INSERT INTO refresh_tokens
		(uid, token_hash, lookup_hash, created_at, updated_at)
		VALUES (1, 'legacy-hash', 'legacy-lookup', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`).Error; err != nil {
		t.Fatalf("insert legacy row: %v", err)
	}

	// open through the provider: AutoMigrate + legacy column drop
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

	// the legacy row must survive the column drop
	var n int64
	if err := p.DB().Raw(`SELECT count(*) FROM refresh_tokens WHERE lookup_hash = 'legacy-lookup'`).Scan(&n).Error; err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if n != 1 {
		t.Fatalf("legacy row lost after migration, count = %d", n)
	}
}
