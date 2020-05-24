package zone

import "fmt"

type InvalidRelay struct {
	error string
}

func NewInvalidRelay(relay string) *InvalidRelay {
	return &InvalidRelay{error: fmt.Sprintf("'%s' is not a valid relay id", relay)}
}

func (i *InvalidRelay) Error() string {
	return i.error
}

type CreateError struct {
	message string
}

func NewCreateError(message string) CreateError {
	return CreateError{message: message}
}

func (z CreateError) Error() string {
	return z.message
}
