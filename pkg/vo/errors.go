package vo

import (
	"fmt"
)

// NotFoundError is self-described
type NotFoundError struct {
	id string
}

// NewNotFoundError is a constructor
func NewNotFoundError(id string) NotFoundError {
	return NotFoundError{id: id}
}

// Error implements Error interface
func (n NotFoundError) Error() string {
	return fmt.Sprintf("not found %q", n.id)
}
