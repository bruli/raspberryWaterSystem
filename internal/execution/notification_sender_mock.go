// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package execution

import (
	"sync"
)

var (
	lockNotificationSenderMockSend sync.RWMutex
)

// Ensure, that NotificationSenderMock does implement NotificationSender.
// If this is not the case, regenerate this file with moq.
var _ NotificationSender = &NotificationSenderMock{}

// NotificationSenderMock is a mock implementation of NotificationSender.
//
//     func TestSomethingThatUsesNotificationSender(t *testing.T) {
//
//         // make and configure a mocked NotificationSender
//         mockedNotificationSender := &NotificationSenderMock{
//             SendFunc: func(message string) error {
// 	               panic("mock out the Send method")
//             },
//         }
//
//         // use mockedNotificationSender in code that requires NotificationSender
//         // and then make assertions.
//
//     }
type NotificationSenderMock struct {
	// SendFunc mocks the Send method.
	SendFunc func(message string) error

	// calls tracks calls to the methods.
	calls struct {
		// Send holds details about calls to the Send method.
		Send []struct {
			// Message is the message argument value.
			Message string
		}
	}
}

// Send calls SendFunc.
func (mock *NotificationSenderMock) Send(message string) error {
	if mock.SendFunc == nil {
		panic("NotificationSenderMock.SendFunc: method is nil but NotificationSender.Send was just called")
	}
	callInfo := struct {
		Message string
	}{
		Message: message,
	}
	lockNotificationSenderMockSend.Lock()
	mock.calls.Send = append(mock.calls.Send, callInfo)
	lockNotificationSenderMockSend.Unlock()
	return mock.SendFunc(message)
}

// SendCalls gets all the calls that were made to Send.
// Check the length with:
//     len(mockedNotificationSender.SendCalls())
func (mock *NotificationSenderMock) SendCalls() []struct {
	Message string
} {
	var calls []struct {
		Message string
	}
	lockNotificationSenderMockSend.RLock()
	calls = mock.calls.Send
	lockNotificationSenderMockSend.RUnlock()
	return calls
}
