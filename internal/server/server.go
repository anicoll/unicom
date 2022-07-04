package server

import (
	"context"
	"time"

	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"
	"github.com/anicoll/unicom/internal/workflows"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type temporalClient interface {
	StartSendSyncWorkflow(ctx context.Context, req workflows.Request) (string, error)
	StartSendAsyncWorkflow(ctx context.Context, req workflows.Request) (string, error)
	GetWorkflowStatus(ctx context.Context, req workflows.StatusRequest) (string, error)
}

type Server struct {
	tc     temporalClient
	logger *zap.Logger
}

func New(tc temporalClient, logger *zap.Logger) *Server {
	return &Server{
		tc:     tc,
		logger: logger,
	}
}

func (s *Server) SendAsync(ctx context.Context, req *pb.SendAsyncRequest) (*pb.SendResponse, error) {
	workflowRequest := workflows.Request{
		EmailRequest:     mapEmailRequestIn(req),
		SleepDuration:    time.Duration(0),
		ResponseRequests: make([]*workflows.ResponseRequest, 0, len(req.GetResponseChannels())),
	}

	for _, responseChannal := range req.GetResponseChannels() {
		switch responseChannal.Schema {
		case pb.ResponseSchema_HTTP:
			workflowRequest.ResponseRequests = append(workflowRequest.ResponseRequests, &workflows.ResponseRequest{
				Type: workflows.WebhookResponseType,
				Url:  responseChannal.GetUrl(),
			})
		case pb.ResponseSchema_SQS:
			workflowRequest.ResponseRequests = append(workflowRequest.ResponseRequests, &workflows.ResponseRequest{
				Type: workflows.SqsResponseType,
				Url:  responseChannal.GetUrl(),
			})
		case pb.ResponseSchema_EVENT_BRIDGE:
			workflowRequest.ResponseRequests = append(workflowRequest.ResponseRequests, &workflows.ResponseRequest{
				Type: workflows.EventBridgeResponseType,
				Url:  responseChannal.GetUrl(),
			})
		}
	}

	now := time.Now()
	requestTime := req.SendAt.AsTime()
	if requestTime.After(now) {
		workflowRequest.SleepDuration = requestTime.Sub(now)
	}
	workflowId, err := s.tc.StartSendAsyncWorkflow(ctx, workflowRequest)
	if err != nil {
		s.logger.Error(err.Error(), zap.Error(err))
		return nil, status.Error(codes.Internal, "unable to send request")
	}
	return &pb.SendResponse{
		Id: workflowId,
	}, nil
}

func (s *Server) SendSync(ctx context.Context, req *pb.SendSyncRequest) (*pb.SendResponse, error) {
	workflowId, err := s.tc.StartSendSyncWorkflow(ctx, workflows.Request{
		EmailRequest:  mapEmailRequestIn(req),
		SleepDuration: time.Duration(0),
	})
	if err != nil {
		s.logger.Error(err.Error(), zap.Error(err))
		return nil, status.Error(codes.Internal, "unable to send request")
	}
	return &pb.SendResponse{
		Id: workflowId,
	}, nil
}

func (s *Server) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	workflowStatus, err := s.tc.GetWorkflowStatus(ctx, workflows.StatusRequest{
		WorkflowId: req.GetId(),
	})
	if err != nil {
		s.logger.Error(err.Error(), zap.Error(err))
		return nil, status.Error(codes.Internal, "unable to query result")
	}
	return &pb.GetStatusResponse{
		Status: string(workflowStatus),
	}, nil
}
