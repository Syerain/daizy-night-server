package crypto

type errType int

const (
	Unknown errType = iota
	RegistercodeFormat
	RegistercodeFormatInvalidLength
	RegistercodeFormatInvalidChar

	RegistercodeAuthen
	RegistercodeAuthenFailed

	RegistercodeUnusable
	RegistercodeUnusableOutdated
	RegistercodeUnusableUsed
	RegistercodeUnusableRepeatedUsername
)

func (t errType) Say() string {
	switch t {
	case Unknown:
		return "unknown error"
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
	}
	return "undefined error type"
}

type ErrRegister struct {
	StatusCode int
	Message    string
	Field      string
	Type       errType
}

func (e *ErrRegister) Error() string {
	return "Register error:" + e.Field + " - " + e.Message
}
