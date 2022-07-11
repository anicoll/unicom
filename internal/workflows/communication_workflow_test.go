package workflows_test

import (
	"errors"
	"testing"

	"github.com/anicoll/unicom/internal/email"
	"github.com/anicoll/unicom/internal/model"
	"github.com/anicoll/unicom/internal/workflows"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bxcodec/faker"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_ComminucationWorkflow_NoResponseChannels_Success() {
	var activities *workflows.UnicomActivities

	sesMessageId := aws.String(uuid.NewString())
	emailRequest := &email.Request{}

	err := faker.FakeData(&emailRequest)
	s.NoError(err)

	s.env.OnActivity(activities.SendEmail, mock.Anything, *emailRequest).Times(1).Return(sesMessageId, nil)
	s.env.OnActivity(activities.UpdateCommunicationStatus, mock.Anything, mock.Anything, model.Success, sesMessageId).Times(1).Return(nil)

	s.env.ExecuteWorkflow(workflows.CommunicationWorkflow, workflows.Request{
		EmailRequest:  emailRequest,
		SleepDuration: 0,
	})
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) Test_ComminucationWorkflow_WithWebhookResponseChannels_Success() {
	var activities *workflows.UnicomActivities

	sesMessageId := aws.String(uuid.NewString())

	emailRequest := &email.Request{}
	err := faker.FakeData(&emailRequest)
	s.NoError(err)

	webhookResponse := &workflows.ResponseRequest{}
	err = faker.FakeData(&webhookResponse)
	s.NoError(err)
	webhookResponse.Type = model.Webhook

	s.env.OnActivity(activities.SendEmail, mock.Anything, *emailRequest).Times(1).Return(sesMessageId, nil)
	s.env.OnActivity(activities.UpdateCommunicationStatus, mock.Anything, mock.Anything, model.Success, sesMessageId).Times(1).Return(nil)
	s.env.OnActivity(activities.NotifyWebhook, mock.Anything, model.ResponseChannelRequest{
		Url:          webhookResponse.Url,
		WorkflowId:   "default-test-workflow-id",
		Status:       string(workflows.WorkflowActivityComplete),
		ErrorMessage: nil,
	},
	).Times(1).Return(nil, nil)
	s.env.OnActivity(activities.SaveResponseChannelOutcome, mock.Anything, webhookResponse.ID, *sesMessageId, model.Success).Times(1).Return(nil)

	s.env.ExecuteWorkflow(workflows.CommunicationWorkflow, workflows.Request{
		EmailRequest:     emailRequest,
		SleepDuration:    0,
		ResponseRequests: []*workflows.ResponseRequest{webhookResponse},
		Domain:           "test-domain",
	})
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) Test_ComminucationWorkflow_WithTwoSQSResponseChannels_Success() {
	var activities *workflows.UnicomActivities

	sesMessageId := aws.String(uuid.NewString())

	emailRequest := &email.Request{}
	err := faker.FakeData(&emailRequest)
	s.NoError(err)

	sqsResponse1 := &workflows.ResponseRequest{}
	err = faker.FakeData(&sqsResponse1)
	s.NoError(err)
	sqsResponse1.Type = model.Sqs
	sqsMessageId1 := aws.String(uuid.NewString())

	sqsResponse2 := &workflows.ResponseRequest{}
	err = faker.FakeData(&sqsResponse2)
	s.NoError(err)
	sqsResponse2.Type = model.Sqs
	sqsMessageId2 := aws.String(uuid.NewString())

	s.env.OnActivity(activities.SendEmail, mock.Anything, *emailRequest).Times(1).Return(sesMessageId, nil)
	s.env.OnActivity(activities.UpdateCommunicationStatus, mock.Anything, mock.Anything, model.Success, sesMessageId).Times(1).Return(nil)
	s.env.OnActivity(activities.NotifySqs, mock.Anything, model.ResponseChannelRequest{
		Url:          sqsResponse1.Url,
		WorkflowId:   "default-test-workflow-id",
		Status:       string(workflows.WorkflowActivityComplete),
		ErrorMessage: nil,
	},
	).Times(1).Return(sqsMessageId1, nil)
	s.env.OnActivity(activities.SaveResponseChannelOutcome, mock.Anything, sqsResponse1.ID, *sqsMessageId1, model.Success).Times(1).Return(nil)

	s.env.OnActivity(activities.NotifySqs, mock.Anything, model.ResponseChannelRequest{
		Url:          sqsResponse2.Url,
		WorkflowId:   "default-test-workflow-id",
		Status:       string(workflows.WorkflowActivityComplete),
		ErrorMessage: nil,
	},
	).Times(1).Return(sqsMessageId2, nil)
	s.env.OnActivity(activities.SaveResponseChannelOutcome, mock.Anything, sqsResponse2.ID, *sqsMessageId2, model.Success).Times(1).Return(nil)

	s.env.ExecuteWorkflow(workflows.CommunicationWorkflow, workflows.Request{
		EmailRequest:     emailRequest,
		SleepDuration:    0,
		ResponseRequests: []*workflows.ResponseRequest{sqsResponse1, sqsResponse2},
		Domain:           "test-domain",
	})
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) Test_ComminucationWorkflow_NoResponseChannels_ErrorSendingComms() {
	var activities *workflows.UnicomActivities

	sesMessageId := (*string)(nil)
	emailRequest := &email.Request{}

	err := faker.FakeData(&emailRequest)
	s.NoError(err)

	s.env.OnActivity(activities.SendEmail, mock.Anything, *emailRequest).Times(1).Return(nil, errors.New("some failed reason"))
	s.env.OnActivity(activities.UpdateCommunicationStatus, mock.Anything, mock.Anything, model.Failed, sesMessageId).Times(1).Return(nil)

	s.env.ExecuteWorkflow(workflows.CommunicationWorkflow, workflows.Request{
		EmailRequest:  emailRequest,
		SleepDuration: 0,
	})
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}
