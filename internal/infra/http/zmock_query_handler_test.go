package http_test

import (
	"context"
	"sync"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

// Ensure, that QueryMock does implement cqs.Query.
// If this is not the case, regenerate this file with moq.
var _ cqs.Query = &QueryMock{}

// QueryMock is a mock implementation of cqs.Query.
//
//	func TestSomethingThatUsesQuery(t *testing.T) {
//
//		// make and configure a mocked cqs.Query
//		mockedQuery := &QueryMock{
//			NameFunc: func() string {
//				panic("mock out the Name method")
//			},
//		}
//
//		// use mockedQuery in code that requires cqs.Query
//		// and then make assertions.
//
//	}
type QueryMock struct {
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
func (mock *QueryMock) Name() string {
	if mock.NameFunc == nil {
		panic("QueryMock.NameFunc: method is nil but Query.Name was just called")
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
//	len(mockedQuery.NameCalls())
func (mock *QueryMock) NameCalls() []struct{} {
	var calls []struct{}
	mock.lockName.RLock()
	calls = mock.calls.Name
	mock.lockName.RUnlock()
	return calls
}

// Ensure, that QueryHandlerMock does implement cqs.QueryHandler.
// If this is not the case, regenerate this file with moq.
var _ cqs.QueryHandler = &QueryHandlerMock{}

// QueryHandlerMock is a mock implementation of cqs.QueryHandler.
//
//	func TestSomethingThatUsesQueryHandler(t *testing.T) {
//
//		// make and configure a mocked cqs.QueryHandler
//		mockedQueryHandler := &QueryHandlerMock{
//			HandleFunc: func(ctx context.Context, query cqs.Query) (any, error) {
//				panic("mock out the Handle method")
//			},
//		}
//
//		// use mockedQueryHandler in code that requires cqs.QueryHandler
//		// and then make assertions.
//
//	}
type QueryHandlerMock struct {
	// HandleFunc mocks the Handle method.
	HandleFunc func(ctx context.Context, query cqs.Query) (any, error)

	// calls tracks calls to the methods.
	calls struct {
		// Handle holds details about calls to the Handle method.
		Handle []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Query is the query argument value.
			Query cqs.Query
		}
	}
	lockHandle sync.RWMutex
}

// Handle calls HandleFunc.
func (mock *QueryHandlerMock) Handle(ctx context.Context, query cqs.Query) (any, error) {
	if mock.HandleFunc == nil {
		panic("QueryHandlerMock.HandleFunc: method is nil but QueryHandler.Handle was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Query cqs.Query
	}{
		Ctx:   ctx,
		Query: query,
	}
	mock.lockHandle.Lock()
	mock.calls.Handle = append(mock.calls.Handle, callInfo)
	mock.lockHandle.Unlock()
	return mock.HandleFunc(ctx, query)
}

// HandleCalls gets all the calls that were made to Handle.
// Check the length with:
//
//	len(mockedQueryHandler.HandleCalls())
func (mock *QueryHandlerMock) HandleCalls() []struct {
	Ctx   context.Context
	Query cqs.Query
} {
	var calls []struct {
		Ctx   context.Context
		Query cqs.Query
	}
	mock.lockHandle.RLock()
	calls = mock.calls.Handle
	mock.lockHandle.RUnlock()
	return calls
}
