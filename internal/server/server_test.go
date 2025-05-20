package server_test

import (
	"context"
	"testing"

	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"
	"github.com/anicoll/unicom/internal/server"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type ServerUnitTestSuite struct {
	suite.Suite
	svc *server.Server
	tc  *mocktemporalClient
	db  *mockpostgres
}

func TestServerUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ServerUnitTestSuite))
}

func (s *ServerUnitTestSuite) TestSendCommunication_Email_Success() {
	req := &pb.SendCommunicationRequest{
		Email:   &pb.EmailRequest{ToAddress: "test@example.com", Subject: "Test", Html: "Hello"},
		IsAsync: false,
		Domain:  "test-domain",
	}
	s.db.EXPECT().CreateCommunication(mock.Anything, mock.Anything).Once().Return(nil)
	s.tc.EXPECT().StartCommunicationWorkflow(mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	s.tc.EXPECT().GetWorkflowResult(mock.Anything, mock.Anything).Once().Return(nil)

	resp, err := s.svc.SendCommunication(context.Background(), req)
	s.NoError(err)
	s.NotEmpty(resp.Id)
}

func (s *ServerUnitTestSuite) TestSendCommunication_Push_Success() {
	req := &pb.SendCommunicationRequest{
		Push: &pb.PushRequest{IdempotencyKey: "Push", Content: &pb.LanguageContent{
			English: "Push body",
		}},
		IsAsync: false,
		Domain:  "test-domain",
	}
	s.db.EXPECT().CreateCommunication(mock.Anything, mock.Anything).Once().Return(nil)
	s.tc.EXPECT().StartCommunicationWorkflow(mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	s.tc.EXPECT().GetWorkflowResult(mock.Anything, mock.Anything).Once().Return(nil)

	resp, err := s.svc.SendCommunication(context.Background(), req)
	s.NoError(err)
	s.NotEmpty(resp.Id)
}

func (s *ServerUnitTestSuite) TestSendCommunication_InvalidRequest_MultipleMediums() {
	req := &pb.SendCommunicationRequest{
		Email: &pb.EmailRequest{ToAddress: "test@example.com"},
		Push:  &pb.PushRequest{IdempotencyKey: "Push"},
	}
	resp, err := s.svc.SendCommunication(context.Background(), req)
	s.Nil(resp)
	s.Error(err)
	s.Contains(err.Error(), "invalid request for multiple notification types")
}

func (s *ServerUnitTestSuite) TestSendCommunication_InvalidRequest_NoMedium() {
	req := &pb.SendCommunicationRequest{}
	resp, err := s.svc.SendCommunication(context.Background(), req)
	s.Nil(resp)
	s.Error(err)
	s.Contains(err.Error(), "invalid request must include any request medium")
}

func (s *ServerUnitTestSuite) TestSendCommunication_DBError() {
	req := &pb.SendCommunicationRequest{
		Email:   &pb.EmailRequest{ToAddress: "test@example.com"},
		IsAsync: false,
		Domain:  "test-domain",
	}
	s.db.EXPECT().CreateCommunication(mock.Anything, mock.Anything).Once().Return(assert.AnError)

	resp, err := s.svc.SendCommunication(context.Background(), req)
	s.Nil(resp)
	s.Error(err)
	s.Contains(err.Error(), "unable to save communication")
}

func (s *ServerUnitTestSuite) TestSendCommunication_WorkflowError() {
	req := &pb.SendCommunicationRequest{
		Email:   &pb.EmailRequest{ToAddress: "test@example.com"},
		IsAsync: false,
		Domain:  "test-domain",
	}
	s.db.EXPECT().CreateCommunication(mock.Anything, mock.Anything).Once().Return(nil)
	s.tc.EXPECT().StartCommunicationWorkflow(mock.Anything, mock.Anything, mock.Anything).Once().Return(assert.AnError)

	resp, err := s.svc.SendCommunication(context.Background(), req)
	s.Nil(resp)
	s.Error(err)
	s.Contains(err.Error(), "unable to send request")
}

func (s *ServerUnitTestSuite) TestSendCommunication_GetWorkflowResultError() {
	req := &pb.SendCommunicationRequest{
		Email:   &pb.EmailRequest{ToAddress: "test@example.com"},
		IsAsync: false,
		Domain:  "test-domain",
	}
	s.db.EXPECT().CreateCommunication(mock.Anything, mock.Anything).Once().Return(nil)
	s.tc.EXPECT().StartCommunicationWorkflow(mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	s.tc.EXPECT().GetWorkflowResult(mock.Anything, mock.Anything).Once().Return(assert.AnError)

	resp, err := s.svc.SendCommunication(context.Background(), req)
	s.Nil(resp)
	s.Error(err)
	s.Contains(err.Error(), "unable to get request result")
}
func (s *ServerUnitTestSuite) SetupSuite() {}

func (s *ServerUnitTestSuite) SetupTest() {
	s.tc = newMocktemporalClient(s.T())
	s.db = newMockpostgres(s.T())
	s.svc = server.New(zap.NewNop(), s.tc, s.db)
}
