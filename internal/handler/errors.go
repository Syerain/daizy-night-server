package handler

type ErrRegisterValidation struct {
	Field   string
	Message string
}

func (e *ErrRegisterValidation) Error() string {
	return "Validation error:" + e.Field + " - " + e.Message
}
