package worker

import (
	"log"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	"github.com/urfave/cli/v2"
)

type workerArgs struct {
	temporalAddress   string
	temporalNamespace string
	owner             string
	opsPort           int
	dbDsn             string
	migrationAction   string
	name              string
	region            string
	description       string
	version           string
	onesignalAppId    string
	onesignalAuthKey  string
}

func CommunicationWorkerCommand() *cli.Command {
	return &cli.Command{
		Name: "communication-worker",
		Flags: []cli.Flag{
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
			&cli.StringFlag{
				Name:     "onesignal-app-id",
				EnvVars:  []string{"ONESIGNAL_APP_ID"},
				Required: true,
				Value:    "",
			},
			&cli.StringFlag{
				Name:     "onesignal-auth-key",
				EnvVars:  []string{"ONESIGNAL_AUTH_KEY"},
				Required: true,
				Value:    "",
			},
		},
		Action: func(c *cli.Context) error {
			args := workerArgs{
				temporalNamespace: c.String("temporal-namespace"),
				temporalAddress:   c.String("temporal-server"),
				opsPort:           c.Int("ops-port"),
				region:            c.String("aws-region"),
				dbDsn:             c.String("db-dsn"),
				migrationAction:   c.String("migrate-action"),
				onesignalAppId:    c.String("onesignal-app-id"),
				onesignalAuthKey:  c.String("onesignal-auth-key"),
				name:              c.Command.Name,
				description:       c.Command.Description,
				version:           c.App.Version,
				owner:             c.App.Authors[0].Name,
			}
			return communicationWorkerAction(args)
		},
	}
}

func newPrometheusScope(c prometheus.Configuration, prefix string) tally.Scope {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				log.Println("error in prometheus reporter", err)
			},
		},
	)
	if err != nil {
		log.Fatalln("error creating prometheus reporter", err)
	}
	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sanitizeOptions,
		Prefix:          prefix,
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)

	log.Println("prometheus metrics scope created")
	return scope
}

// tally sanitizer options that satisfy Prometheus restrictions.
// This will rename metrics at the tally emission level, so metrics name we
// use maybe different from what gets emitted. In the current implementation
// it will replace - and . with _
var (
	safeCharacters = []rune{'_'}

	sanitizeOptions = tally.SanitizeOptions{
		NameCharacters: tally.ValidCharacters{
			Ranges:     tally.AlphanumericRange,
			Characters: safeCharacters,
		},
		KeyCharacters: tally.ValidCharacters{
			Ranges:     tally.AlphanumericRange,
			Characters: safeCharacters,
		},
		ValueCharacters: tally.ValidCharacters{
			Ranges:     tally.AlphanumericRange,
			Characters: safeCharacters,
		},
		ReplacementCharacter: tally.DefaultReplacementCharacter,
	}
)
