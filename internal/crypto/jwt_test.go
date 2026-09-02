package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

// newTestProvider builds a provider with fresh key pairs. refreshTTL is the
// refresh-token lifetime; a negative value makes every issued refresh token
// already expired at signing time.
func newTestProvider(t *testing.T, refreshTTL time.Duration) *ProviderCrypto {
	t.Helper()

	pubAccess, privAccess, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("gen access keypair: %v", err)
	}
	pubRefresh, privRefresh, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("gen refresh keypair: %v", err)
	}

	return &ProviderCrypto{
		JwtAccessTokenEnckey:   privAccess,
		JwtAccessTokenDeckey:   pubAccess,
		JwtRefreshTokenEnckey:  privRefresh,
		JwtRefreshTokenDeckey:  pubRefresh,
		AccessTokenExpireTime:  15 * time.Minute,
		RefreshTokenExpireTime: refreshTTL,
	}
}

// assertTokenLoginErr asserts that err is the app-level 401 business error.
func assertTokenLoginErr(t *testing.T, err error, wantType errs.ErrType) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	errapp, ok := errs.Easx[*errs.ErrUserLogin](err)
	if !ok {
		t.Fatalf("expected *errs.ErrUserLogin, got %T: %v", err, err)
	}
	if errapp.Type != wantType {
		t.Fatalf("err type = %v, want %v", errapp.Type, wantType)
	}
	if code := errapp.StatusCode(); code != 401 {
		t.Fatalf("status code = %d, want 401", code)
	}
}

func TestVerifyRefreshTokenHappyPath(t *testing.T) {
	p := newTestProvider(t, time.Hour)

	raw, err := p.SignRefreshToken(model.JwtRefreshTokenPayload{Uid: 7, Username: "alice"})
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	payload, err := p.VerifyRefreshToken(raw)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if payload.Uid != 7 || payload.Username != "alice" {
		t.Fatalf("payload mismatch: %+v", payload)
	}
	if payload.ID == "" {
		t.Fatal("refresh token must carry a random jti")
	}
}

func TestVerifyRefreshTokenExpired(t *testing.T) {
	p := newTestProvider(t, -time.Minute) // signed already expired

	raw, err := p.SignRefreshToken(model.JwtRefreshTokenPayload{Uid: 7})
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	_, err = p.VerifyRefreshToken(raw)
	assertTokenLoginErr(t, err, errs.UserLoginTokenExpired)
}

func TestVerifyRefreshTokenTampered(t *testing.T) {
	p := newTestProvider(t, time.Hour)

	raw, err := p.SignRefreshToken(model.JwtRefreshTokenPayload{Uid: 7})
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	// corrupt the signature segment
	tampered := raw[:len(raw)-4] + "AAAA"
	_, err = p.VerifyRefreshToken(tampered)
	assertTokenLoginErr(t, err, errs.UserLoginTokenInvalid)

	// garbage input must stay the invalid-token business error as well
	_, err = p.VerifyRefreshToken("not-a-token")
	assertTokenLoginErr(t, err, errs.UserLoginTokenInvalid)
}

func TestVerifyRefreshTokenWrongKey(t *testing.T) {
	p := newTestProvider(t, time.Hour)

	// a token signed with the access key must not verify as a refresh token
	access, err := p.SignAccessToken(model.JwtAccessTokenPayload{Uid: 7})
	if err != nil {
		t.Fatalf("sign access: %v", err)
	}
	_, err = p.VerifyRefreshToken(access)
	assertTokenLoginErr(t, err, errs.UserLoginTokenInvalid)
}

func TestVerifyAccessTokenExpired(t *testing.T) {
	p := newTestProvider(t, time.Hour)
	p.AccessTokenExpireTime = -time.Minute

	raw, err := p.SignAccessToken(model.JwtAccessTokenPayload{Uid: 7})
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	_, err = p.VerifyAccessToken(raw)
	assertTokenLoginErr(t, err, errs.UserLoginTokenExpired)

	// sanity: a valid one still parses with the jwt lib's own claims intact
	p.AccessTokenExpireTime = time.Minute
	raw, err = p.SignAccessToken(model.JwtAccessTokenPayload{Uid: 7})
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	token, perr := jwt.ParseWithClaims(raw, &model.JwtAccessTokenPayload{}, NewJWTKeyFunc(p.JwtAccessTokenDeckey))
	if perr != nil || !token.Valid {
		t.Fatalf("raw jwt parse failed: %v", perr)
	}
}
