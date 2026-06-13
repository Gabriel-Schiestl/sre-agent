package data

import (
	"database/sql"
	"fmt"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/config"
)

type DB struct {
	db *sql.DB
}

func Open(config *config.DBConfig) (*DB, error) {
	db, err := sql.Open(config.Driver, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name))
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}