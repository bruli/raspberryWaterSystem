package cqs

import (
	"context"
	"fmt"
)

// InvalidCommandError should be returned by the implementations of the interface when the handler does not receive the needed command.
type InvalidCommandError struct {
	expected string
	had      string
}

// NewInvalidCommandError is a constructor
func NewInvalidCommandError(expected string, had string) InvalidCommandError {
	return InvalidCommandError{expected: expected, had: had}
}

const errMsgInvalidCommand = "invalid command, expected '%s' but found '%s'"

func (e InvalidCommandError) Error() string {
	return fmt.Sprintf(errMsgInvalidCommand, e.expected, e.had)
}

// Command is the interface for identifying commands by name.
type Command interface {
	Name() string
}

// CommandHandler is self-described
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) ([]Event, error)
}

// CommandHandlerFunc is a function that implements CommandHandler interface.
type CommandHandlerFunc func(ctx context.Context, cmd Command) ([]Event, error)

// Handle is the CommandHandler interface implementation.
func (f CommandHandlerFunc) Handle(ctx context.Context, cmd Command) ([]Event, error) {
	return f(ctx, cmd)
}
