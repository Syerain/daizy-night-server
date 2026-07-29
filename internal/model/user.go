package model

import (
	"daizynight/internal/constants"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	AtomID       int    `gorm:"unique;not null"`
	Username     string `gorm:"unique;not null"`
	Nickname     string `gorm:"not null"`
	Email        string `gorm:"not null"`
	Telephone    string
	Registercode RegisterCode[]

	RegisterTime time.Time

	PasswordHash string `gorm:"not null" json:"-"`

	Role constants.Role `gorm:"not null;default:'user'"`

	// Github OAuth:
	GitHubID    int64  `gorm:"unique"`
	GitHubLogin string `gorm:"unique"` // github login username for presentation
}

type RegisterCode struct {
	gorm.Model

	AtomID int `gorm:"not null"`
	Code   string
	Before time.Time
}
