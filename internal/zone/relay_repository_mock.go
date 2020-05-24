// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package zone

import (
	"sync"
)

var (
	lockRelayRepositoryMockGet sync.RWMutex
)

// Ensure, that RelayRepositoryMock does implement RelayRepository.
// If this is not the case, regenerate this file with moq.
var _ RelayRepository = &RelayRepositoryMock{}

// RelayRepositoryMock is a mock implementation of RelayRepository.
//
//     func TestSomethingThatUsesRelayRepository(t *testing.T) {
//
//         // make and configure a mocked RelayRepository
//         mockedRelayRepository := &RelayRepositoryMock{
//             GetFunc: func() []string {
// 	               panic("mock out the Get method")
//             },
//         }
//
//         // use mockedRelayRepository in code that requires RelayRepository
//         // and then make assertions.
//
//     }
type RelayRepositoryMock struct {
	// GetFunc mocks the Get method.
	GetFunc func() []string

	// calls tracks calls to the methods.
	calls struct {
		// Get holds details about calls to the Get method.
		Get []struct {
		}
	}
}

// Get calls GetFunc.
func (mock *RelayRepositoryMock) Get() []string {
	if mock.GetFunc == nil {
		panic("RelayRepositoryMock.GetFunc: method is nil but RelayRepository.Get was just called")
	}
	callInfo := struct {
	}{}
	lockRelayRepositoryMockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	lockRelayRepositoryMockGet.Unlock()
	return mock.GetFunc()
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//     len(mockedRelayRepository.GetCalls())
func (mock *RelayRepositoryMock) GetCalls() []struct {
} {
	var calls []struct {
	}
	lockRelayRepositoryMockGet.RLock()
	calls = mock.calls.Get
	lockRelayRepositoryMockGet.RUnlock()
	return calls
}
