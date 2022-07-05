package workflows

import (
	"time"

	"github.com/anicoll/unicom/internal/email"
	"github.com/anicoll/unicom/internal/model"
	"github.com/anicoll/unicom/internal/sqs"
	"go.temporal.io/sdk/workflow"
)

type Request struct {
	EmailRequest     *email.Request
	ResponseRequests []*ResponseRequest
	SleepDuration    time.Duration
	Domain           string
	IsAsync          bool
}

type ResponseRequest struct {
	ID   string
	Type model.ResponseChannelType
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

func CommunicationWorkflow(ctx workflow.Context, request Request) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var activities *UnicomActivities
	info := workflow.GetInfo(ctx)

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

	if request.EmailRequest != nil {
		err = workflow.ExecuteActivity(ctx,
			activities.SendEmail,
			request.EmailRequest,
		).Get(ctx, &messageId)
		if err != nil {
			logger.Error("Activity failed.", "activities.SendEmail", "Error", err)
			currentState.Status = WorkflowError
			currentState.Error = err
			err = workflow.ExecuteActivity(ctx,
				activities.MarkCommunicationAsFailed,
				info.WorkflowExecution.ID,
			).Get(ctx, nil)
			if err != nil {
				logger.Error("Activity failed.", "activities.MarkCommunicationAsFailed", "Error", err)
			}
		}
		err = workflow.ExecuteActivity(ctx,
			activities.MarkCommunicationAsSent,
			info.WorkflowExecution.ID,
		).Get(ctx, nil)
		if err != nil {
			return err
		}
	}
	currentState.Status = WorkflowActivityComplete

	for _, responseRequest := range request.ResponseRequests {
		switch responseRequest.Type {
		case model.Sqs:
			var sqsMessageId *string
			err = workflow.ExecuteActivity(ctx,
				activities.NotifySqs,
				sqs.Request{
					Queue:        responseRequest.Url,
					WorkflowId:   info.WorkflowExecution.ID,
					Status:       string(currentState.Status),
					ErrorMessage: messageFromError(currentState.Error),
				},
			).Get(ctx, &sqsMessageId)
			// dont return, save as failed
			if err != nil {
				return err
			}
		case model.Webhook:
			err = workflow.ExecuteActivity(ctx,
				activities.NotifyWebhook,
			).Get(ctx, nil)
			// dont return, save as failed
			if err != nil {
				return err
			}
		}
		err = workflow.ExecuteActivity(ctx,
			activities.SaveResponseChannelOutcome,
			responseRequest.ID,
			stringFromPtr(messageId),
			statusFromError(err),
		).Get(ctx, nil)
	}

	currentState.Status = WorkflowComplete
	logger.Info("SendSyncWorkflow completed.")
	return nil
}

func statusFromError(err error) model.Status {
	if err != nil {
		return model.Failed
	}
	return model.Success
}

func stringFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ptrFromString(s string) *string {
	return &s
}

func messageFromError(err error) *string {
	if err != nil {
		return ptrFromString(err.Error())
	}
	return nil
}
