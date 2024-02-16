package errs

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not_found")
)

type UserError struct {
	str string
}

func (e *UserError) Error() string {
	return e.str
}

func NewUserError(str string) *UserError {
	return &UserError{str: str}
}
