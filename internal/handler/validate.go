package handler

import (
	"reflect"
	"unicode/utf8"

	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/model"
)

/* you mustn't trust frontend completely XD */

// it only check params and it's not its task to
func ValidateRegisterParams(b *model.RegisterBody) error {
	// username
	if _, err := validateNonNull(b.Username); err != nil {
		return &errs.ErrValidation{
			Type:  errs.ValidationKeyNull,
			Field: string(consts.ExprUsername),
			Value: string(consts.ExprBlank),
		}
	}

	if _, err := validateCharValid(b.Username); err != nil {
		return &errs.ErrValidation{
			Type:  errs.ValidationKeyInvalidChar,
			Field: string(consts.ExprUsername),
			Value: b.Username,
		}
	}
	if _, err := validateLengthValid(b.Username, 1, 15); err != nil {
		return &errs.ErrValidation{
			Type:  errs.ValidationKeyInvalidLength,
			Field: string(consts.ExprUsername),
			Value: b.Username,
		}
	}

	// nickname
	if _, err := validateNonNull(b.Nickname); err != nil {
		return &errs.ErrValidation{
			Type:  errs.ValidationKeyNull,
			Field: string(consts.ExprNickname),
			Value: string(consts.ExprBlank),
		}
	}
	if _, err := validateLengthValid(b.Nickname, 1, 15); err != nil {
		return &errs.ErrValidation{
			Type:  errs.ValidationKeyInvalidLength,
			Field: string(consts.ExprNickname),
			Value: b.Nickname,
		}
	}

	// password
	if _, err := validateNonNull(b.Password); err != nil {
		return &errs.ErrValidation{
			Type:  errs.ValidationKeyNull,
			Field: string(consts.ExprPassword),
			Value: string(consts.ExprBlank),
		}
	}
	if _, err := validateLengthValid(b.Password, 6, 128); err != nil {
		return &errs.ErrValidation{
			Type:  errs.ValidationKeyInvalidLength,
			Field: string(consts.ExprPassword),
			Value: b.Password,
		}
	}

	// registercode
	if _, err := validateNonNull(b.Registercode); err != nil {
		return &errs.ErrValidation{
			Type:  errs.ValidationKeyNull,
			Field: string(consts.ExprRegistercode),
			Value: string(consts.ExprBlank),
		}
	}

	return nil
}

func ValidateLoginParams(b *model.LoginBody) error {
	switch b.Loginway {
	case consts.LoginLegacy:
		//
		if _, err := validateNonNull(b.Username); err != nil {
			return &errs.ErrValidation{
				Type:  errs.ValidationKeyNull,
				Field: string(consts.ExprUsername),
				Value: string(consts.ExprBlank),
			}
		}
		if _, err := validateNonNull(b.Password); err != nil {
			return &errs.ErrValidation{
				Type:  errs.ValidationKeyNull,
				Field: string(consts.ExprPassword),
				Value: string(consts.ExprBlank),
			}
		}

	case consts.LoginGithub:
		return &errs.ErrSupport{
			Type: errs.FeatureUnsupported,
		}
		/* under constructing
		// OAuth 登录: atomid 非零
		if b.Atomid == 0 {
			return &errs.ErrValidation{
				Type:  errs.ValidationKeyNull,
				Field: string(constants.ExprAtomid),
				Value: string(constants.ExprNull),
			}
		} */

	default:
		return &errs.ErrValidation{
			Type:  errs.ValidationKeyNull,
			Field: string(consts.ExprLogin),
			Value: string(b.Loginway),
		}
	}

	return nil
}

func validateLengthValid(s string, min int, max int) (bool, error) {
	l := utf8.RuneCountInString(s)
	if l < min || l > max {
		return false, &errs.ErrValidation{
			Type:  errs.ValidationKeyInvalidLength,
			Value: s,
		}
	}
	return true, nil

}

func validateNonNull(v any) (bool, error) {
	if v == nil {
		return false, &errs.ErrValidation{
			Type:  errs.ValidationKeyNull,
			Value: string(consts.ExprNull),
		}
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		if rv.String() == "" {
			return false, &errs.ErrValidation{
				Type:  errs.ValidationKeyNull,
				Value: string(consts.ExprBlank),
			}
		}
	case reflect.Slice, reflect.Map:
		if rv.Len() == 0 {
			return false, &errs.ErrValidation{
				Type:  errs.ValidationKeyNull,
				Value: string(consts.ExprBlank),
			}
		}
	}
	return true, nil
}

func validateCharValid(s string) (bool, error) {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9')) {
			return false, &errs.ErrValidation{
				Type:  errs.ValidationKeyInvalidChar,
				Value: s,
			}
		}
	}
	return true, nil
}
