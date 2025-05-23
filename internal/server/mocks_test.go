// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package server_test

import (
	"context"

	"github.com/anicoll/unicom/internal/model"
	"github.com/anicoll/unicom/internal/workflows"
	mock "github.com/stretchr/testify/mock"
)

// newMocktemporalClient creates a new instance of mocktemporalClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMocktemporalClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *mocktemporalClient {
	mock := &mocktemporalClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// mocktemporalClient is an autogenerated mock type for the temporalClient type
type mocktemporalClient struct {
	mock.Mock
}

type mocktemporalClient_Expecter struct {
	mock *mock.Mock
}

func (_m *mocktemporalClient) EXPECT() *mocktemporalClient_Expecter {
	return &mocktemporalClient_Expecter{mock: &_m.Mock}
}

// GetWorkflowResult provides a mock function for the type mocktemporalClient
func (_mock *mocktemporalClient) GetWorkflowResult(ctx context.Context, workflowId string) error {
	ret := _mock.Called(ctx, workflowId)

	if len(ret) == 0 {
		panic("no return value specified for GetWorkflowResult")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = returnFunc(ctx, workflowId)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// mocktemporalClient_GetWorkflowResult_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWorkflowResult'
type mocktemporalClient_GetWorkflowResult_Call struct {
	*mock.Call
}

// GetWorkflowResult is a helper method to define mock.On call
//   - ctx
//   - workflowId
func (_e *mocktemporalClient_Expecter) GetWorkflowResult(ctx interface{}, workflowId interface{}) *mocktemporalClient_GetWorkflowResult_Call {
	return &mocktemporalClient_GetWorkflowResult_Call{Call: _e.mock.On("GetWorkflowResult", ctx, workflowId)}
}

func (_c *mocktemporalClient_GetWorkflowResult_Call) Run(run func(ctx context.Context, workflowId string)) *mocktemporalClient_GetWorkflowResult_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mocktemporalClient_GetWorkflowResult_Call) Return(err error) *mocktemporalClient_GetWorkflowResult_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *mocktemporalClient_GetWorkflowResult_Call) RunAndReturn(run func(ctx context.Context, workflowId string) error) *mocktemporalClient_GetWorkflowResult_Call {
	_c.Call.Return(run)
	return _c
}

// GetWorkflowStatus provides a mock function for the type mocktemporalClient
func (_mock *mocktemporalClient) GetWorkflowStatus(ctx context.Context, req workflows.StatusRequest) (string, error) {
	ret := _mock.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetWorkflowStatus")
	}

	var r0 string
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, workflows.StatusRequest) (string, error)); ok {
		return returnFunc(ctx, req)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, workflows.StatusRequest) string); ok {
		r0 = returnFunc(ctx, req)
	} else {
		r0 = ret.Get(0).(string)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, workflows.StatusRequest) error); ok {
		r1 = returnFunc(ctx, req)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// mocktemporalClient_GetWorkflowStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWorkflowStatus'
type mocktemporalClient_GetWorkflowStatus_Call struct {
	*mock.Call
}

// GetWorkflowStatus is a helper method to define mock.On call
//   - ctx
//   - req
func (_e *mocktemporalClient_Expecter) GetWorkflowStatus(ctx interface{}, req interface{}) *mocktemporalClient_GetWorkflowStatus_Call {
	return &mocktemporalClient_GetWorkflowStatus_Call{Call: _e.mock.On("GetWorkflowStatus", ctx, req)}
}

func (_c *mocktemporalClient_GetWorkflowStatus_Call) Run(run func(ctx context.Context, req workflows.StatusRequest)) *mocktemporalClient_GetWorkflowStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(workflows.StatusRequest))
	})
	return _c
}

func (_c *mocktemporalClient_GetWorkflowStatus_Call) Return(s string, err error) *mocktemporalClient_GetWorkflowStatus_Call {
	_c.Call.Return(s, err)
	return _c
}

func (_c *mocktemporalClient_GetWorkflowStatus_Call) RunAndReturn(run func(ctx context.Context, req workflows.StatusRequest) (string, error)) *mocktemporalClient_GetWorkflowStatus_Call {
	_c.Call.Return(run)
	return _c
}

// StartCommunicationWorkflow provides a mock function for the type mocktemporalClient
func (_mock *mocktemporalClient) StartCommunicationWorkflow(ctx context.Context, req workflows.Request, workflowId string) error {
	ret := _mock.Called(ctx, req, workflowId)

	if len(ret) == 0 {
		panic("no return value specified for StartCommunicationWorkflow")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, workflows.Request, string) error); ok {
		r0 = returnFunc(ctx, req, workflowId)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// mocktemporalClient_StartCommunicationWorkflow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StartCommunicationWorkflow'
type mocktemporalClient_StartCommunicationWorkflow_Call struct {
	*mock.Call
}

// StartCommunicationWorkflow is a helper method to define mock.On call
//   - ctx
//   - req
//   - workflowId
func (_e *mocktemporalClient_Expecter) StartCommunicationWorkflow(ctx interface{}, req interface{}, workflowId interface{}) *mocktemporalClient_StartCommunicationWorkflow_Call {
	return &mocktemporalClient_StartCommunicationWorkflow_Call{Call: _e.mock.On("StartCommunicationWorkflow", ctx, req, workflowId)}
}

func (_c *mocktemporalClient_StartCommunicationWorkflow_Call) Run(run func(ctx context.Context, req workflows.Request, workflowId string)) *mocktemporalClient_StartCommunicationWorkflow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(workflows.Request), args[2].(string))
	})
	return _c
}

func (_c *mocktemporalClient_StartCommunicationWorkflow_Call) Return(err error) *mocktemporalClient_StartCommunicationWorkflow_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *mocktemporalClient_StartCommunicationWorkflow_Call) RunAndReturn(run func(ctx context.Context, req workflows.Request, workflowId string) error) *mocktemporalClient_StartCommunicationWorkflow_Call {
	_c.Call.Return(run)
	return _c
}

// newMockpostgres creates a new instance of mockpostgres. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockpostgres(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockpostgres {
	mock := &mockpostgres{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// mockpostgres is an autogenerated mock type for the postgres type
type mockpostgres struct {
	mock.Mock
}

type mockpostgres_Expecter struct {
	mock *mock.Mock
}

func (_m *mockpostgres) EXPECT() *mockpostgres_Expecter {
	return &mockpostgres_Expecter{mock: &_m.Mock}
}

// CreateCommunication provides a mock function for the type mockpostgres
func (_mock *mockpostgres) CreateCommunication(ctx context.Context, comm *model.Communication) error {
	ret := _mock.Called(ctx, comm)

	if len(ret) == 0 {
		panic("no return value specified for CreateCommunication")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *model.Communication) error); ok {
		r0 = returnFunc(ctx, comm)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// mockpostgres_CreateCommunication_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateCommunication'
type mockpostgres_CreateCommunication_Call struct {
	*mock.Call
}

// CreateCommunication is a helper method to define mock.On call
//   - ctx
//   - comm
func (_e *mockpostgres_Expecter) CreateCommunication(ctx interface{}, comm interface{}) *mockpostgres_CreateCommunication_Call {
	return &mockpostgres_CreateCommunication_Call{Call: _e.mock.On("CreateCommunication", ctx, comm)}
}

func (_c *mockpostgres_CreateCommunication_Call) Run(run func(ctx context.Context, comm *model.Communication)) *mockpostgres_CreateCommunication_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Communication))
	})
	return _c
}

func (_c *mockpostgres_CreateCommunication_Call) Return(err error) *mockpostgres_CreateCommunication_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *mockpostgres_CreateCommunication_Call) RunAndReturn(run func(ctx context.Context, comm *model.Communication) error) *mockpostgres_CreateCommunication_Call {
	_c.Call.Return(run)
	return _c
}
