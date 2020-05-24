package execution

import "github.com/bruli/raspberryWaterSystem/internal/logger"

const MaxLogs = 30

type ReadLogs struct {
	repository LogRepository
	logger     logger.Logger
}

func NewReadLogs(repository LogRepository, logger logger.Logger) *ReadLogs {
	return &ReadLogs{repository: repository, logger: logger}
}

func (r *ReadLogs) Read() (*Logs, error) {
	l, err := r.repository.Get()
	if err != nil {
		r.logger.Fatalf("failed to get executions log: %w", err)
		return nil, err
	}

	if MaxLogs > len(*l) {
		return l, nil
	}

	logs := Logs{}
	max := len(*l) - MaxLogs
	m := *l
	for _, j := range m[max:] {
		logs.Add(j)
	}
	return &logs, nil
}
