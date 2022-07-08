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

	webhookResponse := &workflows.ResponseRequest{}
	err := faker.FakeData(&webhookResponse)
	s.NoError(err)
	webhookResponse.Type = model.Webhook

	s.env.OnActivity(activities.SendEmail, mock.Anything, *emailRequest).Times(1).Return(sesMessageId, nil)
	s.env.OnActivity(activities.UpdateCommunicationStatus, mock.Anything, mock.Anything, model.Success, sesMessageId).Times(1).Return(nil)
	s.env.OnActivity(activities.NotifyWebhook, mock.Anything).Times(1).Return(nil)
	s.env.OnActivity(activities.SaveResponseChannelOutcome, mock.Anything, webhookResponse.ID, *sesMessageId, model.Success).Times(1).Return(nil)

	s.env.ExecuteWorkflow(workflows.CommunicationWorkflow, workflows.Request{
		EmailRequest:     emailRequest,
		SleepDuration:    0,
		ResponseRequests: []*workflows.ResponseRequest{webhookResponse},
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

// func (s *UnitTestSuite) Test_AgreementInPrincipalSubmissionWorkflow_Failed_Doc_Generation() {
// 	var activities *AgreementInPrincipalSubmissionActivities
// 	aipId, customerId := uuid.NewString(), uuid.NewString()
// 	status := Approved

// 	submissionDate := time.Now().Format(time.RFC3339)

// 	generateDocError := errors.New("it FAILED to generate")

// 	s.env.OnActivity(activities.GenerateAipDocument, mock.Anything, GenerateAipDocumentInput{
// 		UserId: customerId,
// 		MipId:  aipId,
// 	}).Times(10).Return(nil, generateDocError)

// 	s.env.ExecuteWorkflow(AgreementInPrincipalSubmissionWorkflow, SubmitAipWorkflowRequest{
// 		CustomerId:     customerId,
// 		SubmissionDate: submissionDate,
// 		AipId:          aipId,
// 		Status:         status,
// 	})

// 	s.True(s.env.IsWorkflowCompleted())
// 	err := s.env.GetWorkflowError()
// 	var applicationErr *temporal.ApplicationError
// 	s.ErrorAs(err, &applicationErr)
// 	s.Equal(generateDocError.Error(), applicationErr.Error())
// }
