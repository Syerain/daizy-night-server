package utils

import (
	"crypto/rand"
	"encoding/binary"
)

func GenUid() (uint, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	return uint(binary.BigEndian.Uint64(b[:])), nil
}
