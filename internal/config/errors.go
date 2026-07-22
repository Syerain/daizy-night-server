package config

type ErrConfigValidation struct {
	Field   string
	Message string
}

func (e *ErrConfigValidation) Error() string {
	return "Config validation error: " + e.Field + " - " + e.Message
}
