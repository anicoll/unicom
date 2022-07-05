package migrations

import (
	"embed"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/urfave/cli/v2"
)

//go:embed migrations
var migrations embed.FS

type migrationArgs struct {
	dbUrl   string
	migrate string
}

func MigrationCommand() *cli.Command {
	return &cli.Command{
		Name: "migrate",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "migrate",
				EnvVars:  []string{"MIGRATE"},
				Required: false,
				Value:    "up",
				Usage:    "up/down to migrate up or down",
			},
		},
		Action: func(c *cli.Context) error {
			args := migrationArgs{
				dbUrl:   c.String("db-dsn"),
				migrate: c.String("migrate"),
			}
			return migrateAction(args)
		},
	}
}

func migrateAction(args migrationArgs) error {
	if args.migrate == "up" {
		return migrateUp(args.dbUrl)
	}
	if args.migrate == "down" {
		return migrateDown(args.dbUrl)
	}
	return errors.New("please select a migration direction")
}

func migrateUp(url string) error {
	source, _ := httpfs.New(http.FS(migrations), "migrations")

	m, err := migrate.NewWithSourceInstance("httpfs", source, url)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	fmt.Println("migrated up")
	return nil
}

func migrateDown(url string) error {
	source, _ := httpfs.New(http.FS(migrations), "migrations")

	m, err := migrate.NewWithSourceInstance("httpfs", source, url)
	if err != nil {
		return err
	}
	err = m.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	fmt.Println("migrated down")
	return nil
}
