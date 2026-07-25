package crypto

type errType int

const (
	Unknown errType = iota
	UnmatchedUserPassword
	Ed25519key
)

type ErrEd25519key struct {
	Field   string
	Message string
	ErrType errType
}

func (e *ErrEd25519key) Error() string {
	return "ed25519 key error: " + e.Field + " - " + e.Message
}

type ErrUser struct {
	Field   string
	Message string
	ErrType errType
}

func (e *ErrUser) Error() string {
	return "user login error: " + e.Field + " - " + e.Message
}
