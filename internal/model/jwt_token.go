package model

import (
	"time"

	"github.com/atomreforge/daizy-night-server/internal/consts"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	AtomID    uint   `gorm:"not null"`
	TokenHash string `gorm:"unique; not null"`
	RevokedAt *time.Time
}

type JwtAccessTokenPayload struct {
	jwt.RegisteredClaims

	// business
	AtomID   uint        `json:"atomid"`
	Username string      `json:"username"`
	Role     consts.Role `json:"role"`
}

type JwtRefreshTokenPayload struct {
	jwt.RegisteredClaims

	AtomID   uint   `json:"atomid"`
	Username string `json:"username"`
}
