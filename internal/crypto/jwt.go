package crypto

import (
	"daizynight/internal/config"
	"daizynight/internal/constants"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessTokenExpireTime time.Duration
var refreshTokenExpireTime time.Duration

func InitJwt(cfg *config.Config) {
	accessTokenExpireTime = cfg.Http.JwtAccessTokenExpireTime
	refreshTokenExpireTime = cfg.Http.JwtRefreshTokenExpireTime
}

type JwtAccessTokenPayload struct {
	jwt.RegisteredClaims

	// business
	AtomID   int            `json:"atomid"`
	Username string         `json:"username"`
	Role     constants.Role `json:"role"`
}

type JwtRefreshTokenPayload struct {
	jwt.RegisteredClaims

	AtomID   int    `json:"atomid"`
	Username string `json:"username"`
}

func SignAccessToken(payload JwtAccessTokenPayload) (string, error) {
	payload.ExpiresAt = jwt.NewNumericDate(time.Now().Add(accessTokenExpireTime))
	payload.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, payload)
	return token.SignedString(jwtAccessTokenEnckey)
}

func SignRefreshToken(payload JwtRefreshTokenPayload) (string, error) {
	payload.ExpiresAt = jwt.NewNumericDate(time.Now().Add(refreshTokenExpireTime))
	payload.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, payload)
	return token.SignedString(jwtAccessTokenEnckey)
}

func VerifyAccessToken(tokenStr string) (*JwtAccessTokenPayload, error) {
	var payload JwtAccessTokenPayload
	token, err := jwt.ParseWithClaims(tokenStr, &payload, func(t *jwt.Token) (interface{}, error) {
		return jwtAccessTokenDeckey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return &payload, nil
}

func VerifyRefreshToken(tokenStr string) (*JwtRefreshTokenPayload, error) {
	var payload JwtRefreshTokenPayload
	token, err := jwt.ParseWithClaims(tokenStr, &payload, func(t *jwt.Token) (interface{}, error) {
		return jwtRefreshTokenDeckey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return &payload, nil
}
