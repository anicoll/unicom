package database

import (
	"embed"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
)

//go:embed migrations
var migrations embed.FS

type MigrationAction string

const (
	MigrateUp   MigrationAction = "up"
	MigrateDown MigrationAction = "down"
)

type Migration struct {
	dbUrl  string
	action MigrationAction
}

func NewMigrations(url string, action MigrationAction) *Migration {
	return &Migration{
		dbUrl:  url,
		action: action,
	}
}

func (m *Migration) Execute() error {
	if m.action == MigrateUp {
		return m.migrateUp(m.dbUrl)
	}
	if m.action == MigrateDown {
		return m.migrateDown(m.dbUrl)
	}
	return errors.New("please select a migration direction")
}

func (m *Migration) migrateUp(url string) error {
	source, _ := httpfs.New(http.FS(migrations), "migrations")

	migrations, err := migrate.NewWithSourceInstance("httpfs", source, url)
	if err != nil {
		return err
	}
	err = migrations.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	fmt.Println("migrated up")
	return nil
}

func (m *Migration) migrateDown(url string) error {
	source, _ := httpfs.New(http.FS(migrations), "migrations")

	migrations, err := migrate.NewWithSourceInstance("httpfs", source, url)
	if err != nil {
		return err
	}
	err = migrations.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	fmt.Println("migrated down")
	return nil
}
