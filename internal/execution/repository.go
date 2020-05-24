package execution

//go:generate moq -out repository_mock.go . Repository
type Repository interface {
	Save(e Execution) error
	GetExecutions() (*Execution, error)
}
