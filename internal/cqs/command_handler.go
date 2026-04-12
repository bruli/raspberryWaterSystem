package cqs

import (
	"context"
	"fmt"
)

type InvalidCommandError struct {
	expected string
	had      string
}

func NewInvalidCommandError(expected, had string) InvalidCommandError {
	return InvalidCommandError{expected: expected, had: had}
}

const errMsgInvalidCommand = "invalid command, expected '%s' but found '%s'"

func (e InvalidCommandError) Error() string {
	return fmt.Sprintf(errMsgInvalidCommand, e.expected, e.had)
}

type Command interface {
	Name() string
}

type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) ([]Event, error)
}

type CommandHandlerFunc func(ctx context.Context, cmd Command) ([]Event, error)

func (f CommandHandlerFunc) Handle(ctx context.Context, cmd Command) ([]Event, error) {
	return f(ctx, cmd)
}
