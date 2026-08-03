package errs

type errType int

// we use negative expressions to describe errors.
const (
	Unknown errType = iota

	ValidationKeyNull
	ValidationKeyOverLength
	ValidationKeyInvalidLength
	ValidationKeyInvalidChar
	ValidationKeyDuplicatedValue

	// differs from RepeatedUserName
	UserExists

	UserLoginParamsPasswordIncorrect

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

	DbRecordUsernameNotFound
)

func (t *errType) Say() string {
	switch *t {
	case Unknown:
		return "unknown error"

	case ValidationKeyNull:
		return "must not be null"
	case ValidationKeyOverLength:
		return "length should be shorter than 15 chars"
	case ValidationKeyInvalidLength:
		return "invalid length"
	case ValidationKeyInvalidChar:
		return "must contains digits and english letters only"

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
	case DbRecordUsernameNotFound:
		return "unknown user"
	case UserLoginParamsPasswordIncorrect:
		return "incorrect password"
	}
	return "undefined error type"
}
