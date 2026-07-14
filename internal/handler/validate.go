package handler

import (
	"daizynight/internal/crypto"
	"daizynight/internal/model"
	"unicode/utf8"
)

/* you mustn't trust frontend completely XD */

func ValidateRegisterParams(b *model.RegisterBody) error {
	if b.Username == "" {
		return &RegisterValidationError{Field: "username", Message: "shouldnt be null"}
	}

	for i := 0; i < len(b.Username); i++ {
		c := b.Username[i]
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9')) {
			return &RegisterValidationError{Field: "username", Message: "must contains english-letters and digits only"}
		}
	}

	if utf8.RuneCountInString(b.Username) > 15 {
		return &RegisterValidationError{Field: "username", Message: "length should be shorter than 15 chars"}
	}

	if b.Nickname == "" {
		return &RegisterValidationError{Field: "nickname", Message: "shouldnt be null"}
	}

	if utf8.RuneCountInString(b.Nickname) > 15 {
		return &RegisterValidationError{Field: "nickname", Message: "length should be shorter than 15 chars"}
	}

	if !crypto.ValidateSignature(b.Registercode) {
		return &RegisterValidationError{Field: "registercode", Message: "invalid format or outdated code"}
	}

	return nil
}
