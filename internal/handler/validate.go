package handler

import (
	"unicode/utf8"

	"github.com/atomreforge/daizy-night-server/internal/model"
)

/* you mustn't trust frontend completely XD */

// it only check params and it's not its task to
func ValidateRegisterParams(b *model.RegisterBody) error {
	// username
	if b.Username == "" {
		return &ErrRegisterValidation{Field: "username", Type: NotNull, Message: NotNull.Say()}
	}

	for i := 0; i < len(b.Username); i++ {
		c := b.Username[i]
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9')) {
			return &ErrRegisterValidation{Field: "username", Type: InvalidChar, Message: InvalidChar.Say()}
		}
	}

	if utf8.RuneCountInString(b.Username) > 15 {
		return &ErrRegisterValidation{Field: "username", Type: OverLength, Message: OverLength.Say()}
	}

	// nickname
	if b.Nickname == "" {
		return &ErrRegisterValidation{Field: "nickname", Type: NotNull, Message: NotNull.Say()}
	}

	if utf8.RuneCountInString(b.Nickname) > 15 {
		return &ErrRegisterValidation{Field: "nickname", Type: OverLength, Message: OverLength.Say()}
	}

	// password
	if b.Password == "" {
		return &ErrRegisterValidation{Field: "password", Type: NotNull, Message: NotNull.Say()}
	}

	// registercode
	if b.Registercode == "" {
		return &ErrRegisterValidation{Field: "registercode", Type: NotNull, Message: NotNull.Say()}
	}

	return nil
}
