package crypto

import (
	"crypto/ed25519"
	"encoding/hex"
	"time"
)

var Enckey ed25519.PrivateKey
var Deckey ed25519.PublicKey

// transfer hex into bytes; set global keys.
func Init(enckeyHex string, deckeyHex string) error {
	enckeyBytes, err1 := hex.DecodeString(enckeyHex)
	if err1 != nil {
		return err1
	}
	Enckey = ed25519.PrivateKey(enckeyBytes)

	deckeyBytes, err2 := hex.DecodeString(deckeyHex)
	if err2 != nil {
		return err2
	}
	Deckey = ed25519.PublicKey(deckeyBytes)
	return nil
}

func ValidateSignature(signature string) bool {
	dateStr := time.Now().Format("20230630")
	dateBytes := []byte(dateStr)
	sigBytes := []byte(signature)
	isValid := ed25519.Verify(Deckey, dateBytes, sigBytes)
	return isValid
}
