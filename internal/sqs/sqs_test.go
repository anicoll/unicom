package sqs_test

import (
	"context"
	"testing"

	"github.com/anicoll/unicom/internal/sqs"
	aws_sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	svc       *sqs.Service
	sqsClient *mockSqsClient
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.sqsClient = &mockSqsClient{}
	suite.svc = sqs.NewService(suite.sqsClient)
}

func (suite *ServiceTestSuite) AfterTest() {
	suite.sqsClient.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestService_SendMessage_Success() {
	ctx := context.Background()

	req := sqs.Request{}
	faker.FakeData(&req)

	expectedResponse := aws_sqs.SendMessageOutput{}
	faker.FakeData(&expectedResponse)

	suite.sqsClient.On("SendMessage", ctx, mock.Anything, mock.Anything).Return(&expectedResponse, nil)

	assert := assert.New(suite.T())

	resp, err := suite.svc.Send(ctx, req)

	assert.Equal(expectedResponse.MessageId, resp)
	assert.NoError(err)

}
