package responsechannel_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/anicoll/unicom/internal/model"
	"github.com/anicoll/unicom/internal/responsechannel"
)

type ServiceTestSuite struct {
	suite.Suite
	svc       *responsechannel.SQSService
	sqsClient *mocksqsClient
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) SetupTest() {
	s.sqsClient = newMocksqsClient(s.T())
	s.svc = responsechannel.NewSQSService(s.sqsClient)
}

func (s *ServiceTestSuite) TestService_SendMessage_Success() {
	ctx := context.Background()

	req := model.ResponseChannelRequest{}
	err := faker.FakeData(&req)
	s.NoError(err)

	expectedResponse := sqs.SendMessageOutput{}
	err = faker.FakeData(&expectedResponse)
	s.NoError(err)

	s.sqsClient.EXPECT().SendMessage(ctx, mock.Anything, mock.Anything).Return(&expectedResponse, nil)

	assert := assert.New(s.T())

	resp, err := s.svc.Send(ctx, req)

	assert.Equal(expectedResponse.MessageId, resp)
	assert.NoError(err)
}
