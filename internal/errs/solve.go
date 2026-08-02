package errs

import (
	"errors"
)

// "Error As"
func Eas[T error](err error) bool {
	var target T
	return errors.As(err, &target)
}

func Easx[T error](err error) (T, bool) {
	var target T
	if errors.As(err, &target) {
		return target, true
	}
	return target, false
}
