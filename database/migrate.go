package database

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func CreateMigrationsTable(db *sqlx.DB) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logrus.Error("Error creating migrations table")
		panic("err1")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres",
		driver,
	)
	// checking errors if it's not nil then only we will up it
	if err != nil {
		logrus.Fatalf("Error creating migrations table%v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.Errorf("Error creating table: %v", err)
	} else {
		logrus.Info("Migrations applied successfully")
	}
}
