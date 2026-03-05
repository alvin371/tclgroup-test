package errors

import (
	"errors"
	"fmt"
)

// Wrap wraps an error with additional context.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// New creates a new error.
func New(msg string) error {
	return errors.New(msg)
}
