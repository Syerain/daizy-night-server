package model

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	AtomID    int    `gorm:"not null"`
	TokenHash string `gorm:"unique; not null"`
	RevokedAt *time.Time
}
