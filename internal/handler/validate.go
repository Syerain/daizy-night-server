package handler

import (
	"net/http"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/atomreforge/daizy-night-server/internal/utils"
)

/* you mustn't trust frontend completely XD */

// it only check params and it's not its task to
func ValidateRegisterParams(b model.RegisterBody) error {
	// username
	if _, err := validateNonNull(b.Username); err != nil {
		return errs.BuildErrValidation(errs.ValidationKeyNull, http.StatusBadRequest, string(consts.ExprUsername), string(consts.ExprBlank))
	}

	if _, err := validateCharValid(b.Username); err != nil {
		return errs.BuildErrValidation(errs.ValidationKeyInvalidChar, http.StatusBadRequest, string(consts.ExprUsername), b.Username)
	}
	if _, err := validateLengthValid(b.Username, 1, 15); err != nil {
		return errs.BuildErrValidation(errs.ValidationKeyInvalidLength, http.StatusBadRequest, string(consts.ExprUsername), b.Username)
	}

	// nickname
	if _, err := validateNonNull(b.Nickname); err != nil {
		return errs.BuildErrValidation(errs.ValidationKeyNull, http.StatusBadRequest, string(consts.ExprNickname), string(consts.ExprBlank))
	}
	if _, err := validateLengthValid(b.Nickname, 1, 15); err != nil {
		return errs.BuildErrValidation(errs.ValidationKeyInvalidLength, http.StatusBadRequest, string(consts.ExprNickname), b.Nickname)
	}

	// password
	if _, err := validateNonNull(b.Password); err != nil {
		return errs.BuildErrValidation(errs.ValidationKeyNull, http.StatusBadRequest, string(consts.ExprPassword), string(consts.ExprBlank))
	}
	if _, err := validateLengthValid(b.Password, 6, 128); err != nil {
		return errs.BuildErrValidation(errs.ValidationKeyInvalidLength, http.StatusBadRequest, string(consts.ExprPassword), b.Password)
	}

	// registercode
	if ok, err := ValidateFormatRegistercode(b.Registercode); !ok {
		return err
	}

	/*// registercode
	if _, err := validateNonNull(b.Registercode); err != nil {
		return errs.BuildErrValidation(errs.ValidationKeyNull, http.StatusBadRequest, string(consts.ExprRegistercode), string(consts.ExprBlank))
	}*/

	return nil
}

func ValidateLoginParams(b *model.LoginBody) error {
	switch b.Loginway {
	case consts.LoginLegacy:
		//
		if _, err := validateNonNull(b.Username); err != nil {
			return errs.BuildErrValidation(errs.ValidationKeyNull, http.StatusBadRequest, string(consts.ExprUsername), string(consts.ExprBlank))
		}
		if _, err := validateNonNull(b.Password); err != nil {
			return errs.BuildErrValidation(errs.ValidationKeyNull, http.StatusBadRequest, string(consts.ExprPassword), string(consts.ExprBlank))
		}

	case consts.LoginGithub:
		return errs.BuildErrSupport(errs.FeatureUnsupported, http.StatusBadRequest)
		/* under constructing*/
	default:
		return errs.BuildErrValidation(errs.ValidationKeyNull, http.StatusBadRequest, string(consts.ExprLogin), string(b.Loginway))
	}

	return nil
}

func validateLengthValid(s string, min int, max int) (bool, error) {
	l := utf8.RuneCountInString(s)
	if l < min || l > max {
		return false, errs.BuildErrValidation(errs.ValidationKeyInvalidLength, http.StatusBadRequest, "", s)
	}
	return true, nil

}

func validateNonNull(v any) (bool, error) {
	if v == nil {
		return false, errs.BuildErrValidation(errs.ValidationKeyNull, http.StatusBadRequest, "", string(consts.ExprNull))
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		if rv.String() == "" {
			return false, errs.BuildErrValidation(errs.ValidationKeyNull, http.StatusBadRequest, "", string(consts.ExprBlank))
		}
	case reflect.Slice, reflect.Map:
		if rv.Len() == 0 {
			return false, errs.BuildErrValidation(errs.ValidationKeyNull, http.StatusBadRequest, "", string(consts.ExprBlank))
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
			return false, errs.BuildErrValidation(errs.ValidationKeyInvalidChar, http.StatusBadRequest, "", s)
		}
	}
	return true, nil
}

func ValidateFormatRegistercode(c model.RegistercodeRawHex) (bool, error) {
	if c == "" {
		return false, errs.BuildErrValidation(errs.ValidationKeyNull, http.StatusBadRequest, string(consts.ExprRegistercode), string(consts.ExprNull))
	}
	// 用 Split after last "." 确保恰好两部分
	i := strings.LastIndex(string(c), ".")
	if i <= 0 || i == len(c)-1 {
		return false, errs.BuildErrValidation(errs.ValidationKeyBadFormat, http.StatusBadRequest, string(consts.ExprRegistercode), c.String())
	}
	payloadHex, sigHex := string(c)[:i], string(c)[i+1:]
	if !(utils.IsHex(payloadHex) && utils.IsHex(sigHex)) {
		return false, errs.BuildErrValidation(errs.ValidationKeyBadFormat, http.StatusBadRequest, string(consts.ExprRegistercode), c.String())
	}
	return true, nil
}
