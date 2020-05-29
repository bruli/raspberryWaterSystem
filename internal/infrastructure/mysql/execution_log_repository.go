package mysql

import (
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"time"
)

const layout = "2006-01-02 15:04:05"
type executionLogRepository struct {
	repository *Repository
}

func (e *executionLogRepository) Get() (*execution.Logs, error) {
	db, err := e.repository.conn()
	if err != nil {
		return nil, err
	}
	query := "select seconds, zone, created_at from executions order by created_at desc limit 30"
	results, err := db.Query(query)
	defer results.Close()
	if err != nil {
		return nil, err
	}
	logs := execution.Logs{}
	for results.Next() {
		lo := execution.Log{}
		var createdValue string
		err = results.Scan(&lo.Seconds, &lo.Zone, &createdValue)
		if err != nil {
			return nil, err
		}
		createdAt, err := time.Parse(layout, createdValue)
		if err != nil {
			return nil, err
		}
		lo.CreatedAt = createdAt
		logs.Add(&lo)
	}
	return &logs, nil
}

func NewExecutionLogRepository(repository *Repository) *executionLogRepository {
	return &executionLogRepository{repository: repository}
}

func (e *executionLogRepository) Save(l execution.Log) error {
	db, err := e.repository.conn()
	if err != nil {
		return err
	}
	q := "insert into executions (created_at, seconds, zone) values(?, ?, ?)"
	ins, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer ins.Close()
	_, err = ins.Exec(l.CreatedAt.Format(layout), l.Seconds, l.Zone)
	return err
}
