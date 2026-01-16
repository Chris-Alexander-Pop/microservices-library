package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // Internal error, not serialized
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}


func NotFound(msg string, err error) *AppError {
	return New(http.StatusNotFound, msg, err)
}

func InvalidArgument(msg string, err error) *AppError {
	return New(http.StatusBadRequest, msg, err)
}

func Unauthorized(msg string, err error) *AppError {
	return New(http.StatusUnauthorized, msg, err)
}

func Internal(err error) *AppError {
	return New(http.StatusInternalServerError, "Internal Server Error", err)
}
