package execution

//go:generate moq -out relay_manager_mock.go . RelayManager
type RelayManager interface {
	ActivatePins(pins []string) error
	DeactivatePins(pins []string) error
}
