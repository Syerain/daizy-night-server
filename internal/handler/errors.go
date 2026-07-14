package handler

type RegisterValidationError struct {
	Field   string
	Message string
}

func (e *RegisterValidationError) Error() string {
	return "Validation error:" + e.Field + " - " + e.Message
}
