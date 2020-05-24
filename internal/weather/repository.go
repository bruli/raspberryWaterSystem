package weather

//go:generate moq -out repository_mock.go . Repository
type Repository interface {
	Read() (temp float32, hum float32, err error)
}
