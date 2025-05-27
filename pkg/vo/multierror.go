package vo

import (
	"errors"
	"strings"
)

var _ error = &MultiError{}

const (
	multiErrorPrefix    = "multi error: "
	multiErrorSeparator = "; "
)

// NewMultiError is a constructor
func NewMultiError() *MultiError {
	return &MultiError{
		errors: make([]error, 0),
	}
}

// MultiError is self described
type MultiError struct {
	errors []error
}

// Error implements the Error interface
func (e MultiError) Error() string {
	if len(e.errors) == 0 {
		return ""
	}

	errorMsgs := make([]string, len(e.errors))
	for i, e := range e.errors {
		errorMsgs[i] = e.Error()
	}
	return multiErrorPrefix + strings.Join(errorMsgs, multiErrorSeparator)
}

// ErrResult returns nil if not errors have been added
func (e MultiError) ErrResult() error {
	if len(e.errors) == 0 {
		return nil
	}
	return &e
}

// Unwrap returns the underlying error
func (e MultiError) Unwrap() error {
	// Avoid infinite loop if a MultiError instance is used to call `errors.Is` function.
	// See https://github.com/golang/go/blob/master/src/errors/wrap.go#L55
	return nil
}

// Add adds a new error
func (e *MultiError) Add(err error) {
	e.errors = append(e.errors, err)
}

// Is provides the `errors.Is` feature
func (e MultiError) Is(err error) bool {
	for _, e := range e.errors {
		if errors.Is(e, err) {
			return true
		}
	}
	return false
}

// Errors return all underlying errors
func (e MultiError) Errors() []error {
	return e.errors
}

// HasErrors is self-described
func (e MultiError) HasErrors() bool {
	return len(e.errors) > 0
}
