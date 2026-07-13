package crypto

type Ed25519keyError struct {
	Field   string
	Message string
}

func (e *Ed25519keyError) Error() string {
	return "ed25519 key error: " + e.Field + " - " + e.Message
}
