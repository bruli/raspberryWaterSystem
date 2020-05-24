package zone

//go:generate moq -out relay_repository_mock.go . RelayRepository
type RelayRepository interface {
	Get() []string
}
