package handler

type errType int

// we use negative expressions to describe errors.
const (
	Unknown errType = iota
	NotNull
	OverLength
	InvalidChar
	InvalidRegistercode
)

func (t errType) Say() string {
	switch t {
	case Unknown:
		return "unknown error"
	case NotNull:
		return "must not be null"
	case OverLength:
		return "length should be shorter than 15 chars"
	case InvalidChar:
		return "must contains digits and english letters only"
	case InvalidRegistercode:
		return "invalid format or outdated registercode"
	}
	return "undefined error type"
}

type ErrRegisterValidation struct {
	StatusCode int
	Field      string
	Message    string
	Type       errType
}

func (e *ErrRegisterValidation) Error() string {
	return "Validation error:" + e.Field + " - " + e.Message
}
