package temporalclient

import (
	"context"

	"go.temporal.io/sdk/client"

	"github.com/anicoll/unicom/cmd/worker"
	"github.com/anicoll/unicom/internal/workflows"
)

type Client struct {
	temporalClient client.Client
}

func New(tc client.Client) *Client {
	return &Client{
		temporalClient: tc,
	}
}

func (c *Client) GetWorkflowResult(ctx context.Context, workflowId string) error {
	workflowRun := c.temporalClient.GetWorkflow(ctx, workflowId, "")
	return workflowRun.Get(ctx, nil)
}

func (c *Client) StartCommunicationWorkflow(ctx context.Context, req workflows.Request, workflowId string) error {
	_, err := c.temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		TaskQueue: worker.CommunicationTaskQueue,
		ID:        workflowId,
	}, workflows.CommunicationWorkflow, req)
	return err
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
