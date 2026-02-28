package dbmigrations

import (
	"database/sql"
	"embed"
	"errors"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var dbMigrationFS embed.FS

func RunMigrations(db *sql.DB, path string, migrationFS ...fs.FS) error {
	var selectedFS fs.FS
	if len(migrationFS) > 1 {
		return errors.New("Only one filesystem allowed")
	} else if len(migrationFS) == 1 {
		selectedFS = migrationFS[0]
	} else {
		selectedFS = dbMigrationFS
	}
	sourceDriver, err := iofs.New(selectedFS, path)
	if err != nil {
		return err
	}
	dbDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", dbDriver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
