package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/anicoll/unicom/internal/email"
	"github.com/anicoll/unicom/internal/workflows"
	"github.com/aws/aws-sdk-go-v2/config"
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

const SendAsyncTaskQueue string = "unicom_send_async_task_queue"

func SendAsyncWorker(temporalClient client.Client, es *email.Service) error {
	w := worker.New(temporalClient, SendAsyncTaskQueue, worker.Options{})

	registerOptions := workflow.RegisterOptions{}

	activities := workflows.NewActivities(es)

	w.RegisterWorkflowWithOptions(workflows.SendAsyncWorkflow, registerOptions)

	w.RegisterActivityWithOptions(activities.SendEmail, activity.RegisterOptions{})
	w.RegisterActivityWithOptions(activities.NotifySqs, activity.RegisterOptions{})
	w.RegisterActivityWithOptions(activities.NotifyWebhook, activity.RegisterOptions{})

	return w.Run(worker.InterruptCh())
}

func sendAsyncWorkerAction(args workerArgs) error {
	ctx := context.Background()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer func() { _ = zapLogger.Sync() }()

	logger := logur.LoggerToKV(zapadapter.New(zapLogger))

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(args.region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	emailService := email.NewService(awsConfig)

	temporalClient, err := client.Dial(client.Options{
		HostPort:  args.temporalAddress,
		Namespace: args.temporalNamespace,
		Logger:    logger,
		MetricsHandler: sdktally.NewMetricsHandler(newPrometheusScope(prometheus.Configuration{
			ListenAddress: fmt.Sprintf("0.0.0.0:%d", args.opsPort),
			TimerType:     "histogram",
		})),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporalClient.Close()

	return SendAsyncWorker(temporalClient, emailService)
}
