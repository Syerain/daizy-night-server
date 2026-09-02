package errs

type errType int

// ErrType is the exported alias of errType.
// It lets external packages (handler/service/crypto/dbware...) pass errType
// constants (e.g. ValidationKeyNull) to the BuildErr* constructors, whose
// parameter type must be visible outside the package.
type ErrType = errType

// we use negative expressions to describe errors.
const (
	Unknown errType = iota
	Undefined

	ValidationKeyNull
	ValidationKeyOverLength
	ValidationKeyInvalidLength
	ValidationKeyInvalidChar
	ValidationKeyDuplicatedValue
	ValidationKeyBadFormat

	ValidationCryptoUnexpectedSigningMethod

	// differs from RepeatedUserName
	UserExists
	UserUnknown

	UserLoginParamsIncorrect
	UserLoginParamsIncorrectPassword

	UserLoginTokenExpired
	UserLoginTokenInvalid

	RegistercodeFormat
	RegistercodeFormatInvalidLength
	RegistercodeFormatInvalidChar

	RegistercodeAuthen
	RegistercodeAuthenFailed

	RegistercodeUnusable
	RegistercodeUnusableOutdated
	RegistercodeUnusableUsed
	RegistercodeUnusableRepeatedUsername //http409

	RegistercodeUnmarshalFailed

	FeatureUnsupported

	DbRecordNotFound
	DbRecordUsernameNotFound
	DbRecordEmailNotFound

	JwtAccessTokenExpired
	JwtAccessTokenInvalid
	JwtAccessTokenNotFound

	JwtRefreshTokenInvalid
	JwtRefreshTokenRevoked
	JwtRefreshTokenNotFound
	JwtRefreshTokenUsed
	JwtRefreshTokenExpired
)

/*
func (t *errType) Say() string {
	switch *t {
	case Unknown:
		return "unknown error"
	case Undefined:
		return "undefined error"
	case ValidationKeyNull:
		return "must not be null"
	case ValidationKeyOverLength:
		return "length should be shorter than 15 chars"
	case ValidationKeyInvalidLength:
		return "invalid length"
	case ValidationKeyInvalidChar:
		return "must contains digits and english letters only"
	case ValidationKeyDuplicatedValue:
		return "duplicated value"

	case RegistercodeFormat:
		return "invalid registercode format"
	case RegistercodeFormatInvalidLength:
		return "registercode length should be shorter than 15 chars"
	case RegistercodeFormatInvalidChar:
		return "registercode must contains digits and english letters only"

	case RegistercodeAuthen:
		return "registercode authentication failed"
	case RegistercodeAuthenFailed:
		return "registercode signature verification failed"

	case RegistercodeUnusable:
		return "unusable registercode"
	case RegistercodeUnusableOutdated:
		return "outdated registercode"
	case RegistercodeUnusableUsed:
		return "used registercode"
	case RegistercodeUnusableRepeatedUsername:
		return "repeated username in registercode"
	case RegistercodeUnmarshalFailed:
		return "failed to unmarshal regcode to golang struct"
	case FeatureUnsupported:
		return "unsupported feature"
	case UserLoginParamsIncorrect:
		return "incorrect user login params"
	case UserLoginParamsIncorrectPassword:
		return "incorrect password"
	case UserLoginTokenExpired:
		return "token expired"
	case UserLoginTokenInvalid:
		return "invalid token"
	case DbRecordNotFound:
		return "db record not found"
	case DbRecordUsernameNotFound:
		return "unknown user"

	case ValidationCryptoUnexpectedSigningMethod:
		return "expected signing method of public key"
	case ValidationKeyBadFormat:
		return "key bad format"
	}
	return "error description undefined"
}*/

var errorMessages = map[errType]string{
	Unknown:   "unknown error",
	Undefined: "undefined error",

	ValidationCryptoUnexpectedSigningMethod: "expected signing method of public key",

	ValidationKeyNull:            "must not be null",
	ValidationKeyOverLength:      "length should be shorter than 15 chars",
	ValidationKeyInvalidLength:   "invalid length", // length should be 32 chars
	ValidationKeyInvalidChar:     "invalid chars",
	ValidationKeyDuplicatedValue: "duplicated value",
	ValidationKeyBadFormat:       "key bad format",

	RegistercodeFormat:              "invalid registercode format",
	RegistercodeFormatInvalidLength: "registercode length should be shorter than 15 chars",
	RegistercodeFormatInvalidChar:   "registercode must contains digits and english letters only",

	RegistercodeAuthen:       "registercode authentication failed",
	RegistercodeAuthenFailed: "registercode signature verification failed",

	RegistercodeUnusable:                 "unusable registercode",
	RegistercodeUnusableOutdated:         "outdated registercode",
	RegistercodeUnusableUsed:             "used registercode",
	RegistercodeUnusableRepeatedUsername: "repeated username in registercode",
	RegistercodeUnmarshalFailed:          "failed to unmarshal regcode to golang struct",

	FeatureUnsupported: "unsupported feature",

	UserLoginParamsIncorrect:         "incorrect user login params",
	UserLoginParamsIncorrectPassword: "incorrect password",
	UserLoginTokenExpired:            "token expired",
	UserLoginTokenInvalid:            "invalid token",

	DbRecordNotFound:         "db record not found",
	DbRecordUsernameNotFound: "unknown user",
	DbRecordEmailNotFound:    "unknown email",

	JwtAccessTokenExpired:  "jwt access token expired",
	JwtAccessTokenInvalid:  "jwt access token invalid",
	JwtAccessTokenNotFound: "jwt access token not found",

	JwtRefreshTokenInvalid:  "jwt refresh token invalid",
	JwtRefreshTokenRevoked:  "jwt refresh token revoked",
	JwtRefreshTokenNotFound: "jwt refresh token not found",
	JwtRefreshTokenUsed:     "jwt refresh token used; user may has faced or is facing an attack;",
	JwtRefreshTokenExpired:  "jwt refresh token expired",
}

/*

var errorMessages = map[errType]string{
    Unknown:                         "unknown error",
    Undefined:                       "undefined error",
    ValidationKeyNull:               "must not be null",
    ValidationKeyOverLength:         "length should be shorter than 15 chars",
}*/

func (t *errType) Say() string {
	if msg, ok := errorMessages[*t]; ok {
		return msg
	}
	return "error description undefined"
}
