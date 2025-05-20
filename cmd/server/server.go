package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urfave/cli/v3"
	"github.com/utilitywarehouse/go-operational/op"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	zapadapter "logur.dev/adapter/zap"
	"logur.dev/logur"

	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"
	"github.com/anicoll/unicom/internal/database"
	"github.com/anicoll/unicom/internal/server"
	"github.com/anicoll/unicom/internal/temporalclient"
)

func ServerCommand() *cli.Command {
	return &cli.Command{
		Name:        "server",
		Description: "unicom public api service to expose a public interface for communications",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "http-port",
				Sources:  cli.NewValueSourceChain(cli.EnvVar("HTTP_PORT")),
				Required: false,
				Value:    8080,
			},
			&cli.IntFlag{
				Name:     "grpc-port",
				Sources:  cli.NewValueSourceChain(cli.EnvVar("GRPC_PORT")),
				Required: false,
				Value:    8090,
			},
			&cli.StringFlag{
				Name:     "aws-region",
				Sources:  cli.NewValueSourceChain(cli.EnvVar("AWS_REGION")),
				Required: false,
				Value:    "eu-west-2",
			},
			&cli.StringFlag{
				Name:     "temporal-server",
				Sources:  cli.NewValueSourceChain(cli.EnvVar("TEMPORAL_SERVER")),
				Required: true,
				Value:    "localhost:7233",
			},
			&cli.StringFlag{
				Name:     "temporal-namespace",
				Sources:  cli.NewValueSourceChain(cli.EnvVar("TEMPORAL_NAMESPACE")),
				Required: false,
				Value:    "default",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			args := serverArgs{
				grpcPort:          c.Int("grpc-port"),
				httpPort:          c.Int("http-port"),
				opsPort:           c.Int("ops-port"),
				dbDsn:             c.String("db-dsn"),
				migrationAction:   c.String("migrate-action"),
				temporalNamespace: c.String("temporal-namespace"),
				temporalAddress:   c.String("temporal-server"),
				name:              c.Name,
				description:       c.Description,
				version:           c.Version,
				owner:             c.Authors[0].(string),
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
	dbDsn             string
	migrationAction   string
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

	parsedCfg, err := pgxpool.ParseConfig(args.dbDsn)
	if err != nil {
		return err
	}
	conn, err := pgxpool.NewWithConfig(ctx, parsedCfg)
	if err != nil {
		return err
	}

	migrationAction := database.MigrateUp
	if args.migrationAction == "down" {
		migrationAction = database.MigrateDown
	}
	migrations := database.NewMigrations(args.dbDsn, migrationAction)
	err = migrations.Execute()
	if err != nil {
		return err
	}

	db := database.New(conn, logger)

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

		opts := []logging.Option{
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			// Add any other option (check functions starting with logging.With).
		}

		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
				prometheus.NewServerMetrics().UnaryServerInterceptor(),
			))
		pb.RegisterUnicomServiceServer(s, server)
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
