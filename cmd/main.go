package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/anicoll/unicom/cmd/dbinit"
	"github.com/anicoll/unicom/cmd/server"
	"github.com/anicoll/unicom/cmd/worker"
)

var (
	version = "development"
	author  = ""
)

func main() {
	app := &cli.Command{
		Name:    "unicom-public-api-service",
		Usage:   "exposes a 'public' api to other domains",
		Version: version,
		Authors: []any{author},
		Commands: []*cli.Command{
			server.ServerCommand(),
			worker.CommunicationWorkerCommand(),
			dbinit.DatabaseCreationCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "log-level",
				Sources:  cli.NewValueSourceChain(cli.EnvVar("LOG_LEVEL")),
				Required: false,
				Value:    "DEBUG",
			},
			&cli.IntFlag{
				Name:     "ops-port",
				Sources:  cli.NewValueSourceChain(cli.EnvVar("OPS_PORT")),
				Required: false,
				Value:    8081,
			},
			&cli.StringFlag{
				Name:     "db-dsn",
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_DSN")),
				Required: false,
				Value:    "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
			},
			&cli.StringFlag{
				Name:     "migrate-action",
				Sources:  cli.NewValueSourceChain(cli.EnvVar("MIGRATE_ACTION")),
				Required: false,
				Value:    "up",
				Usage:    "to indicate either up/down for migrations",
			},
		},
	}
	ctx := context.Background()
	err := app.Run(ctx, os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
