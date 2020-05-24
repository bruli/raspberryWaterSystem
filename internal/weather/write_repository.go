package weather

//go:generate moq -out write_repository_mock.go . WriteRepository
type WriteRepository interface {
	Write(temp, hum float32) error
}
