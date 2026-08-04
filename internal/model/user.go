package model

import (
	"time"

	"github.com/atomreforge/daizy-night-server/internal/consts"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	AtomID    uint   `gorm:"unique;not null"`
	Username  string `gorm:"unique;not null"`
	Nickname  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Telephone string

	RegisterTime time.Time

	Registercode string
	PasswordHash string `gorm:"not null" json:"-"`

	Role consts.Role `gorm:"not null;default:'user'"`

	// Github OAuth:
	GitHubID    *int64  `gorm:"unique"`
	GitHubLogin *string `gorm:"unique"` // github login username for presentation
}
