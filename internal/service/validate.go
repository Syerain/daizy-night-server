package service

import (
	"daizynight/internal/crypto"
	"daizynight/internal/model"
	"unicode/utf8"
)

/* you mustn't trust frontend completely XD */

func ValidateRegisterParams(b *model.RegisterBody) error {
	if b.Username == "" {
		return &RegisterError{Message: "username shouldnt be null"}
	}

	for i := 0; i < len(b.Username); i++ {
		c := b.Username[i]
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9')) {
			return &RegisterError{
				Message: "username must contains english-letters and digits only",
			}
		}
	}

	if utf8.RuneCountInString(b.Username) > 15 {
		return &RegisterError{Message: "username length should be shorter than 15 chars"}
	}

	if b.Nickname == "" {
		return &RegisterError{Message: "nickname shouldnt be null"}
	}

	if utf8.RuneCountInString(b.Nickname) > 15 {
		return &RegisterError{Message: "nickname length should be shorter than 15 chars"}
	}

	if !crypto.ValidateSignature(b.Registercode) {
		return &RegisterError{Message: "invalid registercode"}
	}

	return nil
}
