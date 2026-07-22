package crypto

type ErrEd25519key struct {
	Field   string
	Message string
}

func (e *ErrEd25519key) Error() string {
	return "ed25519 key error: " + e.Field + " - " + e.Message
}
