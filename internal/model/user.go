package model

import (
	"time"

	"github.com/atomreforge/daizy-night-server/internal/consts"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// notice that:
	// 26.8 'User.Atomid' was replaced by 'gorm.ID' due to async problems
	// 26.8.17 switched to 'User.UserID' from 'gorm.ID' now. the values are now generated randomly.
	UserID    uint   `gorm:"unique;not null"`
	Username  string `gorm:"unique;not null"`
	Nickname  string `gorm:"not null"`
	Email     string
	Telephone string

	RegisterTime time.Time

	Registercode RegistercodeRawHex `gorm:"not null;unique"`
	PasswordHash string             `gorm:"not null" json:"-"`

	Role consts.Role `gorm:"not null;default:'user'"`

	// Github OAuth:
	GithubID    *int64  `gorm:"unique"`
	GithubLogin *string `gorm:"unique"` // github login username for presentation
}
