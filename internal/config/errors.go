package config

type errType int

const (
	Unknown = iota
)

type ErrConfigValidation struct {
	Field   string
	Message string
}

func (e *ErrConfigValidation) Error() string {
	return "Config validation error: " + e.Field + " - " + e.Message
}
