package execution

//go:generate moq -out notification_sender_mock.go . NotificationSender
type NotificationSender interface {
	Send(message string) error
}
