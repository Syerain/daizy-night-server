package model

import "time"

type User struct {
	// basic things:
	ID       int64  `gorm:"unique;not null"`
	Nickname string `gorm:"unique;not null"`
	Email    string `gorm:"not null"`

	// time related:
	// RegisterTime differs from gorm.Model.CreatedAt
	RegisterTime time.Time

	// security:
	Password string `gorm:"not null"`
}
