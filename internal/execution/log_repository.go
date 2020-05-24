package execution

//go:generate moq -out log_repository_mock.go . LogRepository
type LogRepository interface {
	Get() (*Logs, error)
	Save(l Log) error
}
