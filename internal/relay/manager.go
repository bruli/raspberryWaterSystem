package relay

//go:generate moq -out manager_mock.go . Manager
type Manager interface {
	DeactivatePins(pins []string) error
}
