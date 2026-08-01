package confx

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	var msgs []string
	for field, msg := range e.Errors {
		msgs = append(msgs, fmt.Sprintf("%s: %s", field, msg))
	}
	return fmt.Sprintf("validation failed:\n - %s", strings.Join(msgs, "\n - "))
}
