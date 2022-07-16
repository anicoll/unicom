package server_test

import (
	"context"

	"github.com/anicoll/unicom/internal/model"
	"github.com/anicoll/unicom/internal/workflows"
	"github.com/stretchr/testify/mock"
)

type mockDatabase struct {
	mock.Mock
}

func (db *mockDatabase) CreateCommunication(ctx context.Context, comm *model.Communication) error {
	args := db.Called(ctx, comm)
	return args.Error(0)
}

type mockTemporalClient struct {
	mock.Mock
}

func (tc *mockTemporalClient) StartCommunicationWorkflow(ctx context.Context, req workflows.Request, workflowId string) error {
	args := tc.Called(ctx, req, workflowId)
	return args.Error(0)
}

func (tc *mockTemporalClient) GetWorkflowStatus(ctx context.Context, req workflows.StatusRequest) (string, error) {
	args := tc.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (tc *mockTemporalClient) GetWorkflowResult(ctx context.Context, workflowId string) error {
	args := tc.Called(ctx, workflowId)
	return args.Error(0)
}
