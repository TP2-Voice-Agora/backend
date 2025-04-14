package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(sourceName string) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", sourceName)

	if err != nil {
		return nil, err
	}
	return &Postgres{db: db}, nil

}
