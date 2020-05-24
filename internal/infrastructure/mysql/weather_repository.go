package mysql

import "time"

type weatherRepository struct {
	repository *Repository
}

func (w weatherRepository) Write(temp, hum float32) error {
	db, err := w.repository.conn()
	if err != nil {
		return err
	}
	defer db.Close()

	q := "insert into weather (weather_value, created_at, type) values (?, ?, ?)"
	ins, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer ins.Close()
	_, err = ins.Exec(hum, time.Now(), "humidity")
	if err != nil {
		return err
	}
	_, err = ins.Exec(temp, time.Now(), "temperature")
	return err

}

func NewWeatherRepository(repository *Repository) *weatherRepository {
	return &weatherRepository{repository: repository}
}
