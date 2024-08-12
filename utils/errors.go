package utils

import "fmt"

type CustomError struct {
	Code    string
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewError(code, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}
