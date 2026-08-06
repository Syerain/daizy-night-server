package model

import (
	"time"

	"gorm.io/gorm"
)

// Registercode == message + ed25519-Signature
// Signature(64B) -> R(commitment)(32B) + S(scalar)(32B)
// Publickey(64B)
type RegistercodePayload struct {
	Magicword string
	Before    time.Time
}

// the whole vanilla string of a Registercode, in hex.
type RegistercodeRawHex string

func (c RegistercodeRawHex) String() string { return string(c) }

// the payload part of a registercode, in hex
type RegistercodePayloadRawHex string

// the signature part of a registercode, in hex
type RegistercodeSigRawHex string

type RegistercodeRecord struct {
	gorm.Model
	RawHex RegistercodeRawHex `gorm:"uniqueIndex;not null;type:text"`
}
