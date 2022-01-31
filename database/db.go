package database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Todo *sqlx.DB

type SSLMode string

const (
	SSLModeEnable  SSLMode = "enable"
	SSLModeDisable SSLMode = "disable"
)

func Connect(host, port, dbname, user, password string, sslMode SSLMode) error {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s  dbname=%s sslmode=%s", host, port, user, password, dbname, SSLModeDisable)
	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	Todo = db
	return migrateStart(db)

}

func migrateStart(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, NewErr := migrate.NewWithDatabaseInstance("file://database/migration", "postgres", driver)
	if NewErr != nil {
		return err
	}
	if MigrateErr := m.Up(); err != nil && MigrateErr != migrate.ErrNoChange { //up(): will migrate all the way up
		return err
	}
	return nil
}
