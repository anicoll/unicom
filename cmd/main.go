package main

import (
	"log"
	"os"

	"github.com/anicoll/unicom/cmd/dbinit"
	"github.com/anicoll/unicom/cmd/server"
	"github.com/anicoll/unicom/cmd/worker"
	"github.com/urfave/cli/v2"
)

var version = "development"
var author = ""

func main() {
	app := &cli.App{
		Name:    "unicom-public-api-service",
		Usage:   "exposes a 'public' api to other domains",
		Version: version,
		Authors: []*cli.Author{
			{
				Name: author,
			},
		},
		Commands: []*cli.Command{
			server.ServerCommand(),
			worker.CommunicationWorkerCommand(),
			dbinit.DatabaseCreationCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "log-level",
				EnvVars:  []string{"LOG_LEVEL"},
				Required: false,
				Value:    "DEBUG",
			},
			&cli.IntFlag{
				Name:     "ops-port",
				EnvVars:  []string{"OPS_PORT"},
				Required: false,
				Value:    8081,
			},
			&cli.StringFlag{
				Name:     "db-dsn",
				EnvVars:  []string{"DB_DSN"},
				Required: false,
				Value:    "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
			},
			&cli.StringFlag{
				Name:     "migrate-action",
				EnvVars:  []string{"MIGRATE_ACTION"},
				Required: false,
				Value:    "up",
				Usage:    "to indicate either up/down for migrations",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
