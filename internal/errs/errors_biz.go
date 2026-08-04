package errs

import (
	"net/http"
)

type BizError struct {
	http int `json:"-"`
	//BizEventCode       int    `json:"code"`
	Message string `json:"message"`
}

func (e *BizError) Error() string   { return e.Message }
func (e *BizError) StatusCode() int { return e.http }

func BuildBadRequest(msg string) *BizError {
	return &BizError{http: http.StatusBadRequest, Message: msg}
}

func BuildUnauthorized(msg string) *BizError {
	return &BizError{http: http.StatusUnauthorized, Message: msg}
}
func BuildNotFound(msg string) *BizError {
	return &BizError{http: http.StatusNotFound, Message: msg}
}
func BuildInternal(msg string) *BizError {
	return &BizError{http: http.StatusInternalServerError, Message: msg}
}
