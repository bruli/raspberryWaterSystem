package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	host, port, user, password, database string
}

func NewConfig(host string, port string, user string, password string, database string) *Config {
	return &Config{host: host, port: port, user: user, password: password, database: database}
}

type Repository struct {
	config *Config
}

func NewRepository(config *Config) *Repository {
	return &Repository{config: config}
}

func (r *Repository) conn() (*sql.DB, error) {
	return sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", r.config.user, r.config.password, r.config.host, r.config.port, r.config.database))
}
