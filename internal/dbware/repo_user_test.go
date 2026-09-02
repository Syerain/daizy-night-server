package dbware

import (
	"fmt"
	"testing"
	"time"

	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newUserTestEnv spins up an isolated in-memory database with the full schema
// migrated and a RepoUser bound to it.
func newUserTestEnv(t *testing.T) *RepoUser {
	t.Helper()

	dsn := fmt.Sprintf("file:dntest_user_%d?mode=memory&cache=shared", time.Now().UnixNano())
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := gdb.AutoMigrate(&model.User{}, &model.RefreshToken{}, &model.RegistercodeRecord{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return &RepoUser{pDB: &ProviderDB{db: gdb}}
}

func mustSeedRegcode(t *testing.T, r *RepoUser, rawHex string, used bool) {
	t.Helper()
	if err := r.pDB.DB().Create(&model.RegistercodeRecord{RawHex: model.RegistercodeRawHex(rawHex), Used: used}).Error; err != nil {
		t.Fatalf("seed regcode %s: %v", rawHex, err)
	}
}

func TestCreateUserSuccessClaimsRegcode(t *testing.T) {
	r := newUserTestEnv(t)
	mustSeedRegcode(t, r, "aa.bb", false)

	b := &model.User{Username: "alice", Nickname: "Alice", PasswordHash: "hash", Registercode: "aa.bb"}
	if err := r.CreateUser(b); err != nil {
		t.Fatalf("create: %v", err)
	}

	if b.UserID < 1_000_000 || b.UserID > 9_999_999 {
		t.Fatalf("uid out of the 7-digit range: %d", b.UserID)
	}
	if b.RegisterTime.IsZero() {
		t.Fatal("register time not set")
	}

	var rec model.RegistercodeRecord
	if err := r.pDB.DB().Where("raw_hex = ?", "aa.bb").First(&rec).Error; err != nil {
		t.Fatalf("load regcode: %v", err)
	}
	if !rec.Used || rec.UserID == nil || *rec.UserID != b.UserID {
		t.Fatalf("regcode not claimed properly: used=%v user_id=%v", rec.Used, rec.UserID)
	}
}

// GetUserByUid must answer 404 for a missing user record, per the API docs
// (the /me endpoint surfaces this status directly).
func TestGetUserByUidNotFoundIs404(t *testing.T) {
	r := newUserTestEnv(t)

	_, err := r.GetUserByUid(4242424)
	if err == nil {
		t.Fatal("expected error for unknown uid")
	}
	errapp, ok := errs.Easx[*errs.ErrDbRecord](err)
	if !ok {
		t.Fatalf("expected *errs.ErrDbRecord, got %T: %v", err, err)
	}
	if errapp.StatusCode() != 404 {
		t.Fatalf("status = %d, want 404", errapp.StatusCode())
	}

	// the login flow's getter keeps 400 (mapped to a generic login error)
	_, err = r.GetUserByUsername("ghost")
	errapp2, ok := errs.Easx[*errs.ErrDbRecord](err)
	if !ok {
		t.Fatalf("expected *errs.ErrDbRecord, got %T: %v", err, err)
	}
	if errapp2.StatusCode() != 400 {
		t.Fatalf("status = %d, want 400", errapp2.StatusCode())
	}
}

func TestCreateUserUsernameDuplicate(t *testing.T) {
	r := newUserTestEnv(t)
	mustSeedRegcode(t, r, "aa.bb", false)

	first := &model.User{Username: "alice", Nickname: "A", PasswordHash: "h", Registercode: "aa.bb"}
	if err := r.CreateUser(first); err != nil {
		t.Fatalf("first create: %v", err)
	}

	second := &model.User{Username: "alice", Nickname: "B", PasswordHash: "h", Registercode: "cc.dd"}
	// the service layer records the regcode before calling CreateUser
	mustSeedRegcode(t, r, "cc.dd", false)
	err := r.CreateUser(second)
	if err == nil {
		t.Fatal("duplicate username unexpectedly accepted")
	}
	errapp, ok := errs.Easx[*errs.ErrValidation](err)
	if !ok {
		t.Fatalf("expected *errs.ErrValidation, got %T: %v", err, err)
	}
	if errapp.Type != errs.ValidationKeyDuplicatedValue || errapp.StatusCode() != 400 {
		t.Fatalf("type=%v status=%d, want duplicated-value/400", errapp.Type, errapp.StatusCode())
	}

	// the failed registration must not consume the second regcode
	var rec model.RegistercodeRecord
	if err := r.pDB.DB().Where("raw_hex = ?", "cc.dd").First(&rec).Error; err != nil {
		t.Fatalf("load regcode: %v", err)
	}
	if rec.Used {
		t.Fatal("regcode of the failed registration must stay unused")
	}
}

func TestCreateUserUsedRegistercodeRollsBack(t *testing.T) {
	r := newUserTestEnv(t)
	mustSeedRegcode(t, r, "aa.bb", true) // already claimed

	b := &model.User{Username: "bob", Nickname: "B", PasswordHash: "h", Registercode: "aa.bb"}
	err := r.CreateUser(b)
	if err == nil {
		t.Fatal("used registercode unexpectedly accepted")
	}
	errapp, ok := errs.Easx[*errs.ErrRegistercode](err)
	if !ok {
		t.Fatalf("expected *errs.ErrRegistercode, got %T: %v", err, err)
	}
	if errapp.Type != errs.RegistercodeUnusableUsed || errapp.StatusCode() != 400 {
		t.Fatalf("type=%v status=%d, want unusable-used/400", errapp.Type, errapp.StatusCode())
	}

	// the whole transaction (user insert included) must roll back
	var count int64
	if err := r.pDB.DB().Model(&model.User{}).Count(&count).Error; err != nil {
		t.Fatalf("count users: %v", err)
	}
	if count != 0 {
		t.Fatalf("user insert was not rolled back, count = %d", count)
	}
}
