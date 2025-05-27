package http_test

import (
	"context"
	"sync"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

// Ensure, that CommandHandlerMock does implement cqs.CommandHandler.
// If this is not the case, regenerate this file with moq.
var _ cqs.CommandHandler = &CommandHandlerMock{}

type CommandHandlerMock struct {
	// HandleFunc mocks the Handle method.
	HandleFunc func(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error)

	// calls tracks calls to the methods.
	calls struct {
		// Handle holds details about calls to the Handle method.
		Handle []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Cmd is the cmd argument value.
			Cmd cqs.Command
		}
	}
	lockHandle sync.RWMutex
}

// Handle calls HandleFunc.
func (mock *CommandHandlerMock) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	if mock.HandleFunc == nil {
		panic("CommandHandlerMock.HandleFunc: method is nil but CommandHandler.Handle was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Cmd cqs.Command
	}{
		Ctx: ctx,
		Cmd: cmd,
	}
	mock.lockHandle.Lock()
	mock.calls.Handle = append(mock.calls.Handle, callInfo)
	mock.lockHandle.Unlock()
	return mock.HandleFunc(ctx, cmd)
}

// HandleCalls gets all the calls that were made to Handle.
// Check the length with:
//
//	len(mockedCommandHandler.HandleCalls())
func (mock *CommandHandlerMock) HandleCalls() []struct {
	Ctx context.Context
	Cmd cqs.Command
} {
	var calls []struct {
		Ctx context.Context
		Cmd cqs.Command
	}
	mock.lockHandle.RLock()
	calls = mock.calls.Handle
	mock.lockHandle.RUnlock()
	return calls
}

// Ensure, that CommandMock does implement cqs.Command.
// If this is not the case, regenerate this file with moq.
var _ cqs.Command = &CommandMock{}

type CommandMock struct {
	// NameFunc mocks the Name method.
	NameFunc func() string

	// calls tracks calls to the methods.
	calls struct {
		// Name holds details about calls to the Name method.
		Name []struct{}
	}
	lockName sync.RWMutex
}

// Name calls NameFunc.
func (mock *CommandMock) Name() string {
	if mock.NameFunc == nil {
		panic("CommandMock.NameFunc: method is nil but Command.Name was just called")
	}
	callInfo := struct{}{}
	mock.lockName.Lock()
	mock.calls.Name = append(mock.calls.Name, callInfo)
	mock.lockName.Unlock()
	return mock.NameFunc()
}

// NameCalls gets all the calls that were made to Name.
// Check the length with:
//
//	len(mockedCommand.NameCalls())
func (mock *CommandMock) NameCalls() []struct{} {
	var calls []struct{}
	mock.lockName.RLock()
	calls = mock.calls.Name
	mock.lockName.RUnlock()
	return calls
}
