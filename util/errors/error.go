// gateway-service/util/errors/error.go
package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

type AppError struct {
	Code    codes.Code
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, error: %v", e.Code, e.Message, e.Err)
}

func Wrap(err error, code codes.Code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func ErrNotFound(model string) error {
	return Wrap(nil, codes.NotFound, fmt.Sprintf("the %s resource not found", model))
}

func ErrBadRequest(msg string) error {
	return Wrap(nil, codes.InvalidArgument, msg)
}

func DuplicateTransaction() error {
	return Wrap(nil, codes.AlreadyExists, "Duplicate Transaction")
}
