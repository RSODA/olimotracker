package migrator

import (
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db: db,
	}
}

func (m *Migrator) Up() error {
	err := goose.Up(m.db, "/app/migrations")
	if err != nil {
		log.Print("err migration up: ", err)
		return err
	}

	return nil
}
