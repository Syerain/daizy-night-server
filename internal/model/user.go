package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// basic things:
	Username     string `gorm:"unique;not null"`
	Nickname     string `gorm:"unique;not null"`
	Email        string
	Tele         string
	Registercode string `gorm:"unique;not null"`

	// time related:
	RegisterTime time.Time // differs from gorm.Model.CreatedAt

	// Security:
	SaltedPassword string `gorm:"not null"`

	// Permission:
	/* notice that despite the existence of PermissionWeight,
	still currently using legacy permission roles. */
	IsAdmin          bool  `gorm:"not null;default:false"`
	PermissionWeight int16 `gorm:"default:20"`

	// Github OAuth:
	GitHubID    int64  `gorm:"unique"`
	GitHubLogin string `gorm:"unique"` // github login username for presentation
}
