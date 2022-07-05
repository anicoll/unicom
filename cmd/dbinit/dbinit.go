package dbinit

import (
	"context"
	"fmt"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type initArgs struct {
	databaseToCreate    string
	passwordForDatabase string
	dbDsn               string
}

func DatabaseCreationCommand() *cli.Command {
	return &cli.Command{
		Name: "database-init",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "database-to-create",
				EnvVars:  []string{"DATABASE_TO_CREATE"},
				Required: true,
				Value:    "",
				Usage:    "database to create in DB",
			},
			&cli.StringFlag{
				Name:     "password-for-database",
				EnvVars:  []string{"PASSWORD_FOR_DATABASE"},
				Required: true,
				Value:    "",
				Usage:    "Password to use to create database tables",
			},
		},

		Action: func(c *cli.Context) error {
			args := initArgs{
				databaseToCreate:    c.String("database-to-create"),
				passwordForDatabase: c.String("password-for-database"),
				dbDsn:               c.String("db-dsn"),
			}
			return initDatabaseAction(args)
		},
	}
}

func initDatabaseAction(args initArgs) error {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, args.dbDsn)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", args.databaseToCreate))
	if err != nil && !strings.Contains(err.Error(), "(SQLSTATE 42P04)") {
		return errors.WithMessage(err, "database creation")
	}
	_, err = conn.Exec(ctx, fmt.Sprintf("CREATE USER %s WITH ENCRYPTED PASSWORD '%s'", args.databaseToCreate, args.passwordForDatabase))
	if err != nil && !strings.Contains(err.Error(), "(SQLSTATE 42710)") {
		return errors.WithMessage(err, "user creation")
	}
	_, err = conn.Exec(ctx, fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", args.databaseToCreate, args.databaseToCreate))
	if err != nil {
		return errors.WithMessage(err, "privileges grant")

	}
	return nil
}
