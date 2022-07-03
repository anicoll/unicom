package main

import (
	"log"
	"os"

	"github.com/anicoll/unicom/cmd/server"
	"github.com/anicoll/unicom/cmd/worker"
	"github.com/urfave/cli/v2"
)

var version = "development"

func main() {
	app := &cli.App{
		Name:    "unicom-public-api-service",
		Usage:   "exposes a 'public' api to other domains",
		Version: version,
		Commands: []*cli.Command{
			server.ServerCommand(),
			worker.SendSyncWorkerCommand(),
			worker.SendAsyncWorkerCommand(),
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
				Name:     "aws-region",
				EnvVars:  []string{"AWS_REGION"},
				Required: false,
				Value:    "eu-west-2",
			},
			&cli.StringFlag{
				Name:     "temporal-server",
				EnvVars:  []string{"TEMPORAL_SERVER"},
				Required: true,
				Value:    "localhost:7233",
			},
			&cli.StringFlag{
				Name:     "temporal-namespace",
				EnvVars:  []string{"TEMPORAL_NAMESPACE"},
				Required: false,
				Value:    "default",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
