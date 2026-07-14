package utils

type SaltGenError struct {
	Content string
	Message string
}

func (e *SaltGenError) Error() string {
	return "Failed to generate salt for " + e.Content + ": " + e.Message
}

type SaltVarifyError struct {
	Content string
	Message string
}

func (e *SaltVarifyError) Error() string {
	return "Failed to verify salt for " + e.Content + ": " + e.Message
}
