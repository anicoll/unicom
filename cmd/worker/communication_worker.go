package worker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anicoll/unicom/internal/database"
	"github.com/anicoll/unicom/internal/email"
	"github.com/anicoll/unicom/internal/responsechannel/sqs"
	"github.com/anicoll/unicom/internal/responsechannel/webhook"
	"github.com/anicoll/unicom/internal/workflows"
	"github.com/aws/aws-sdk-go-v2/config"
	ses "github.com/aws/aws-sdk-go-v2/service/sesv2"
	aws_sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/uber-go/tally/v4/prometheus"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	sdktally "go.temporal.io/sdk/contrib/tally"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
	zapadapter "logur.dev/adapter/zap"
	"logur.dev/logur"
)

const CommunicationTaskQueue string = "unicom_task_queue"

func CommunicationWorker(temporalClient client.Client, emailClient *email.Service, sqsClient *sqs.Service, webhookClient *webhook.Service, db *database.Postgres) error {
	w := worker.New(temporalClient, CommunicationTaskQueue, worker.Options{})

	registerOptions := workflow.RegisterOptions{}

	activities := workflows.NewActivities(emailClient, sqsClient, webhookClient, db)

	w.RegisterWorkflowWithOptions(workflows.CommunicationWorkflow, registerOptions)

	w.RegisterActivityWithOptions(activities.SendEmail, activity.RegisterOptions{})
	w.RegisterActivityWithOptions(activities.NotifySqs, activity.RegisterOptions{})
	w.RegisterActivityWithOptions(activities.NotifyWebhook, activity.RegisterOptions{})

	w.RegisterActivityWithOptions(activities.UpdateCommunicationStatus, activity.RegisterOptions{})
	w.RegisterActivityWithOptions(activities.SaveResponseChannelOutcome, activity.RegisterOptions{})

	return w.Run(worker.InterruptCh())
}

func communicationWorkerAction(args workerArgs) error {
	ctx := context.Background()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer func() { _ = zapLogger.Sync() }()

	conn, err := pgxpool.Connect(ctx, args.DbDsn)
	if err != nil {
		return err
	}
	migrations := database.NewMigrations(args.DbDsn, database.MigrateUp)
	err = migrations.Execute()
	if err != nil {
		return err
	}
	db := database.New(conn)

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(args.region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// TODO: add status Checkers
	sqsClient := aws_sqs.NewFromConfig(awsConfig)
	sqsService := sqs.NewService(sqsClient)

	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}
	// TODO: add status Checkers
	webhookClient := webhook.NewService(httpClient)

	// TODO: add status Checkers
	sesClient := ses.NewFromConfig(awsConfig)
	emailService := email.NewService(sesClient)

	temporalClient, err := client.Dial(client.Options{
		HostPort:  args.temporalAddress,
		Namespace: args.temporalNamespace,
		Logger:    logur.LoggerToKV(zapadapter.New(zapLogger)),
		MetricsHandler: sdktally.NewMetricsHandler(newPrometheusScope(prometheus.Configuration{
			ListenAddress: fmt.Sprintf("0.0.0.0:%d", args.opsPort),
			TimerType:     "histogram",
		}, args.owner)),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporalClient.Close()

	return CommunicationWorker(temporalClient, emailService, sqsService, webhookClient, db)
}
