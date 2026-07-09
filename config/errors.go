package config

type ConfigValidationError struct {
	Field   string
	Message string
}

func (e *ConfigValidationError) Error() string {
	return "Config validation error: " + e.Field + " - " + e.Message
}
