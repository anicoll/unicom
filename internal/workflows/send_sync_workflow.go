package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func SendSyncWorkflow(ctx workflow.Context, request Request) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	var activities *UnicomActivities

	err := workflow.ExecuteActivity(ctx,
		activities.SendEmail,
	).Get(ctx, nil)
	if err != nil {
		logger.Error("Activity failed.", "activities.SendEmail", "Error", err)
		return err
	}

	logger.Info("SendSyncWorkflow completed.")
	return nil
}
