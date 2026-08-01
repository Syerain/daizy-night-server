package model

import (
	"time"
)

type RegistercodePayload struct {
	AtomID    uint `gorm:"not null"`
	Magicword string
	Before    time.Time
}
