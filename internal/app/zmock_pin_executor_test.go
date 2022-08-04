// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package app_test

import (
	"context"
	"github.com/bruli/raspberryWaterSystem/internal/app"
	"sync"
)

// Ensure, that PinExecutorMock does implement app.PinExecutor.
// If this is not the case, regenerate this file with moq.
var _ app.PinExecutor = &PinExecutorMock{}

// PinExecutorMock is a mock implementation of app.PinExecutor.
//
// 	func TestSomethingThatUsesPinExecutor(t *testing.T) {
//
// 		// make and configure a mocked app.PinExecutor
// 		mockedPinExecutor := &PinExecutorMock{
// 			ExecuteFunc: func(ctx context.Context, seconds uint, pins []string) error {
// 				panic("mock out the Execute method")
// 			},
// 		}
//
// 		// use mockedPinExecutor in code that requires app.PinExecutor
// 		// and then make assertions.
//
// 	}
type PinExecutorMock struct {
	// ExecuteFunc mocks the Execute method.
	ExecuteFunc func(ctx context.Context, seconds uint, pins []string) error

	// calls tracks calls to the methods.
	calls struct {
		// Execute holds details about calls to the Execute method.
		Execute []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Seconds is the seconds argument value.
			Seconds uint
			// Pins is the pins argument value.
			Pins []string
		}
	}
	lockExecute sync.RWMutex
}

// Execute calls ExecuteFunc.
func (mock *PinExecutorMock) Execute(ctx context.Context, seconds uint, pins []string) error {
	if mock.ExecuteFunc == nil {
		panic("PinExecutorMock.ExecuteFunc: method is nil but PinExecutor.Execute was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Seconds uint
		Pins    []string
	}{
		Ctx:     ctx,
		Seconds: seconds,
		Pins:    pins,
	}
	mock.lockExecute.Lock()
	mock.calls.Execute = append(mock.calls.Execute, callInfo)
	mock.lockExecute.Unlock()
	return mock.ExecuteFunc(ctx, seconds, pins)
}

// ExecuteCalls gets all the calls that were made to Execute.
// Check the length with:
//     len(mockedPinExecutor.ExecuteCalls())
func (mock *PinExecutorMock) ExecuteCalls() []struct {
	Ctx     context.Context
	Seconds uint
	Pins    []string
} {
	var calls []struct {
		Ctx     context.Context
		Seconds uint
		Pins    []string
	}
	mock.lockExecute.RLock()
	calls = mock.calls.Execute
	mock.lockExecute.RUnlock()
	return calls
}
