package rain

//go:generate moq -out repository_mock.go . Repository
type Repository interface {
	Get() (Rain, error)
}
