package workflows

import (
	"time"

	"github.com/anicoll/unicom/internal/email"
	"github.com/anicoll/unicom/internal/sqs"
	"go.temporal.io/sdk/workflow"
)

type Request struct {
	EmailRequest     *email.Request
	ResponseRequests []*ResponseRequest
	SleepDuration    time.Duration
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

	var messageId *string

	err = workflow.ExecuteActivity(ctx,
		activities.SendEmail,
		request.EmailRequest,
	).Get(ctx, &messageId)
	if err != nil {
		logger.Error("Activity failed.", "activities.SendEmail", "Error", err)
		currentState.Status = WorkflowError
		currentState.Error = err
		return err
	}
	currentState.Status = WorkflowActivityComplete

	info := workflow.GetInfo(ctx)

	for _, responseRequest := range request.ResponseRequests {
		switch responseRequest.Type {
		case SqsResponseType:
			var sqsMessageId *string
			err = workflow.ExecuteActivity(ctx,
				activities.NotifySqs,
				sqs.Request{
					Queue:      responseRequest.Url,
					WorkflowId: info.WorkflowExecution.ID,
					Status:     string(currentState.Status),
				},
			).Get(ctx, &sqsMessageId)
			if err != nil {
				return err
			}
		case WebhookResponseType:
			err = workflow.ExecuteActivity(ctx,
				activities.NotifyWebhook,
			).Get(ctx, nil)
			if err != nil {
				return err
			}
		}
	}

	currentState.Status = WorkflowComplete
	logger.Info("SendSyncWorkflow completed.")
	return nil
}
