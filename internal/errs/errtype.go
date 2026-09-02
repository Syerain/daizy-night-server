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
)

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
}
