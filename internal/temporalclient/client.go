package temporalclient

import (
	"context"

	"github.com/anicoll/unicom/cmd/worker"
	"github.com/anicoll/unicom/internal/workflows"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

type Client struct {
	temporalClient client.Client
}

func New(tc client.Client) *Client {
	return &Client{
		temporalClient: tc,
	}
}

func (c *Client) newWorkflowId() string {
	return uuid.NewString()
}

func (c *Client) StartSendSyncWorkflow(ctx context.Context, req workflows.Request) (string, error) {
	workflowId := c.newWorkflowId()

	wf, err := c.temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		TaskQueue: worker.SendSyncTaskQueue,
		ID:        workflowId,
	}, workflows.SendSyncWorkflow, req)
	if err != nil {
		return "", err
	}
	err = wf.Get(ctx, nil)
	if err != nil {
		return "", err
	}
	return workflowId, nil
}

func (c *Client) StartSendAsyncWorkflow(ctx context.Context, req workflows.Request) (string, error) {
	workflowId := c.newWorkflowId()
	_, err := c.temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		TaskQueue: worker.SendAsyncTaskQueue,
		ID:        workflowId,
	}, workflows.SendAsyncWorkflow, req)
	if err != nil {
		return "", err
	}
	return workflowId, nil
}

func (c *Client) GetWorkflowStatus(ctx context.Context, req workflows.StatusRequest) (string, error) {
	queryResponse, err := c.temporalClient.QueryWorkflowWithOptions(ctx, &client.QueryWorkflowWithOptionsRequest{
		WorkflowID: req.WorkflowId,
		QueryType:  "current_state",
	})
	if err != nil {
		return "", err
	}
	respo := workflows.WorkflowState{}
	err = queryResponse.QueryResult.Get(&respo)
	if err != nil {
		return "", err
	}
	return string(respo.Status), nil
}
