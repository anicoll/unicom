package workflows

import (
	"time"

	"github.com/anicoll/unicom/internal/email"
	"go.temporal.io/sdk/workflow"
)

type Request struct {
	EmailRequest    email.SendEmailRequest
	ResponseRequest ResponseRequest
	SleepDuration   time.Duration
}

type ResponseType string

const (
	SqsResponseType         ResponseType = "SQS"
	WebhookResponseType     ResponseType = "HTTP"
	EventBridgeResponseType ResponseType = "HTTP"
)

type ResponseRequest struct {
	Type ResponseType
	Url  string
}

type Status string

const (
	WorkflowStarted          Status = "STARTED"
	WorkflowWaiting          Status = "WAITING"
	WorkflowError            Status = "ERROR"
	WorkflowCancelled        Status = "CANCELLED"
	WorkflowActivityComplete Status = "ACTIVITY_COMPLETE"
	WorkflowResponding       Status = "RESPONDING"
	WorkflowComplete         Status = "COMPLETE"
)

type WorkflowState struct {
	Status Status
	Error  error
}

type StatusRequest struct {
	WorkflowId string
}

func SendAsyncWorkflow(ctx workflow.Context, request Request) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var activities *UnicomActivities

	currentState := &WorkflowState{
		Status: WorkflowStarted,
	}
	queryType := "current_state"
	err := workflow.SetQueryHandler(ctx, queryType, func() (*WorkflowState, error) {
		return currentState, nil
	})

	if err != nil {
		currentState.Status = WorkflowError
		currentState.Error = err
		return err
	}

	currentState.Status = WorkflowWaiting
	// send in the future
	err = workflow.Sleep(ctx, request.SleepDuration)
	if err != nil {
		currentState.Status = WorkflowError
		currentState.Error = err
		return err
	}

	err = workflow.ExecuteActivity(ctx,
		activities.SendEmail,
	).Get(ctx, nil)
	if err != nil {
		logger.Error("Activity failed.", "activities.SendEmail", "Error", err)
		currentState.Status = WorkflowError
		currentState.Error = err
		return err
	}
	currentState.Status = WorkflowActivityComplete

	switch request.ResponseRequest.Type {
	case SqsResponseType:
		workflow.ExecuteActivity(ctx,
			activities.NotifySqs,
		).Get(ctx, nil)
	case WebhookResponseType:
		workflow.ExecuteActivity(ctx,
			activities.NotifyWebhook,
		).Get(ctx, nil)
	}

	currentState.Status = WorkflowComplete
	logger.Info("SendSyncWorkflow completed.")
	return nil
}
