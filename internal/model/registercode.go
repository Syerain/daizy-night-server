package model

import (
	"time"
)

type RegistercodePayload struct {
	Magicword string
	Before    time.Time
}
