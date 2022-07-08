package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"
	"github.com/anicoll/unicom/internal/database"
	"github.com/anicoll/unicom/internal/op"
	"github.com/anicoll/unicom/internal/server"
	"github.com/anicoll/unicom/internal/temporalclient"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jackc/pgx/v4/pgxpool"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	zapadapter "logur.dev/adapter/zap"

	"github.com/urfave/cli/v2"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"logur.dev/logur"
)

func ServerCommand() *cli.Command {
	return &cli.Command{
		Name:        "server",
		Description: "unicom public api service to expose a public interface for communications",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "http-port",
				EnvVars:  []string{"HTTP_PORT"},
				Required: false,
				Value:    8080,
			},
			&cli.IntFlag{
				Name:     "grpc-port",
				EnvVars:  []string{"GRPC_PORT"},
				Required: false,
				Value:    8090,
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
		Action: func(c *cli.Context) error {
			args := serverArgs{
				grpcPort:          c.Int("grpc-port"),
				httpPort:          c.Int("http-port"),
				opsPort:           c.Int("ops-port"),
				DbDsn:             c.String("db-dsn"),
				temporalNamespace: c.String("temporal-namespace"),
				temporalAddress:   c.String("temporal-server"),
				name:              c.Command.Name,
				description:       c.Command.Description,
				version:           c.App.Version,
				owner:             c.App.Authors[0].Name,
			}
			return run(args)
		},
	}
}

type serverArgs struct {
	grpcPort          int
	httpPort          int
	opsPort           int
	owner             string
	temporalAddress   string
	temporalNamespace string
	name              string
	DbDsn             string
	description       string
	version           string
}

func run(args serverArgs) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	status := op.NewStatus(args.name, args.description).
		SetRevision(args.version)

	eg, ctx := errgroup.WithContext(ctx)

	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, err := zapConfig.Build()
	if err != nil {
		return err
	}

	// nolint: errcheck
	defer logger.Sync()

	conn, err := pgxpool.Connect(ctx, args.DbDsn)
	if err != nil {
		return err
	}
	db := database.New(conn)

	status.AddChecker("database", func(cr *op.CheckResponse) {
		if err := db.Ping(ctx); err != nil {
			cr.Unhealthy("database unavailable", "check database connection/network", "service wont be able to record expressions-of-interest")
		} else {
			cr.Healthy("healthy")
		}
	})

	tClient, err := client.Dial(client.Options{
		HostPort:  args.temporalAddress,
		Namespace: args.temporalNamespace,
		Logger:    logur.LoggerToKV(zapadapter.New(logger)),
		// MetricsHandler: sdktally.NewMetricsHandler(newPrometheusScope(prometheus.Configuration{
		// 	ListenAddress: fmt.Sprintf("0.0.0.0:%d", args.opsPort),
		// 	TimerType:     "histogram",
		// }, args.owner)),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer tClient.Close()

	status.AddChecker("temporal-client", func(cr *op.CheckResponse) {
		if _, err := tClient.CheckHealth(ctx, &client.CheckHealthRequest{}); err != nil {
			cr.Unhealthy("temporal-client unresponsive", "check server connection/network", "service wont be able to initiate any new workflows")
		} else {
			cr.Healthy("healthy")
		}
	})

	tc := temporalclient.New(tClient)

	server := server.New(logger, tc, db)

	eg.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", args.grpcPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer(
			grpc.UnaryInterceptor(
				grpc_middleware.ChainUnaryServer(
					grpc_ctxtags.UnaryServerInterceptor(),
					grpc_prometheus.UnaryServerInterceptor,
					grpc_zap.UnaryServerInterceptor(logger),
				)))
		pb.RegisterUnicomServer(s, server)
		logger.Info("serving GRPC", zap.Int("port", args.grpcPort))
		return s.Serve(lis)
	})

	eg.Go(func() error {
		logger.Info("serving HTTP", zap.Int("port", args.httpPort))
		return runHTTPGateway(context.Background(), args.httpPort, args.grpcPort)
	})

	eg.Go(func() error {
		logger.Info("serving ops status", zap.Int("port", args.opsPort))
		http.Handle("/__/", op.NewHandler(status.ReadyUseHealthCheck()))
		return http.ListenAndServe(fmt.Sprintf(":%d", args.opsPort), nil)
	})

	return eg.Wait()
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
