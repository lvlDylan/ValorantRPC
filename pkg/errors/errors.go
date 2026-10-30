package errors

import (
	"errors"
	"fmt"
)

var (
	ErrLocaleNotFound  = errors.New("requested language file not found")
	ErrLocaleMalformed = errors.New("requested language file malformed")

	ErrEnvironmentVars = errors.New("environment variables not set")

	ErrLoggerCannotOpenFile = errors.New("cannot open log file")

	ErrLockfileNotFound  = errors.New("requested lockfile not found")
	ErrLockfileMalformed = errors.New("requested lockfile malformed")
)

type AppError struct {
	Op  string
	Err error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s]: %v", e.Op, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(op string, err error) error {
	if err == nil {
		return nil
	}

	return &AppError{
		Op:  op,
		Err: err,
	}
}
