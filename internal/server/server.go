package server

import (
	"context"
	"io"
	"time"

	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"
	"github.com/anicoll/unicom/internal/model"
	"github.com/anicoll/unicom/internal/workflows"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type temporalClient interface {
	StartCommunicationWorkflow(ctx context.Context, req workflows.Request, workflowId string) error
	GetWorkflowStatus(ctx context.Context, req workflows.StatusRequest) (string, error)
	GetWorkflowResult(ctx context.Context, workflowId string) error
}

type postgres interface {
	CreateCommunication(ctx context.Context, comm *model.Communication) error
}

type Server struct {
	tc     temporalClient
	db     postgres
	logger *zap.Logger
}

func New(logger *zap.Logger, tc temporalClient, db postgres) *Server {
	return &Server{
		tc:     tc,
		logger: logger,
		db:     db,
	}
}

func (s *Server) SendCommunication(ctx context.Context, req *pb.SendCommunicationRequest) (*pb.SendResponse, error) {
	return s.sendCommunication(ctx, req)
}

func (s *Server) sendCommunication(ctx context.Context, req *pb.SendCommunicationRequest) (*pb.SendResponse, error) {
	err := s.validateRequest(req)
	if err != nil {
		s.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	emailRequest, err := mapEmailRequestIn(req.GetEmail())
	if err != nil {
		s.logger.Error(err.Error(), zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "unable to map email request")
	}
	pushRequest := mapPushNotificationIn(req.GetPush())

	workflowRequest := workflows.Request{
		EmailRequest:     emailRequest,
		PushRequest:      pushRequest,
		SleepDuration:    time.Duration(0),
		ResponseRequests: make([]*workflows.ResponseRequest, 0, len(req.GetResponseChannels())),
		Domain:           req.GetDomain(),
	}
	if req.IsAsync {
		for _, responseChannal := range req.GetResponseChannels() {
			switch responseChannal.Schema {
			case pb.ResponseSchema_HTTP:
				workflowRequest.ResponseRequests = append(workflowRequest.ResponseRequests, &workflows.ResponseRequest{
					Type: model.Webhook,
					Url:  responseChannal.GetUrl(),
					ID:   uuid.NewString(),
				})
			case pb.ResponseSchema_SQS:
				workflowRequest.ResponseRequests = append(workflowRequest.ResponseRequests, &workflows.ResponseRequest{
					Type: model.Sqs,
					Url:  responseChannal.GetUrl(),
					ID:   uuid.NewString(),
				})
			case pb.ResponseSchema_EVENT_BRIDGE:
				workflowRequest.ResponseRequests = append(workflowRequest.ResponseRequests, &workflows.ResponseRequest{
					Type: model.EventBridge,
					Url:  responseChannal.GetUrl(),
					ID:   uuid.NewString(),
				})
			}
		}
		now := time.Now()
		requestTime := req.SendAt.AsTime()
		if requestTime.After(now) {
			workflowRequest.SleepDuration = requestTime.Sub(now)
		}
	}

	workflowId := uuid.NewString()
	err = s.db.CreateCommunication(ctx, mapWorkflowRequestToModel(workflowId, workflowRequest))
	if err != nil {
		s.logger.Error(err.Error(), zap.Error(err))
		return nil, status.Error(codes.Internal, "unable to save communication")
	}
	err = s.tc.StartCommunicationWorkflow(ctx, workflowRequest, workflowId)
	if err != nil {
		s.logger.Error(err.Error(), zap.Error(err))
		return nil, status.Error(codes.Internal, "unable to send request")
	}
	if !req.IsAsync {
		err = s.tc.GetWorkflowResult(ctx, workflowId)
		if err != nil {
			s.logger.Error(err.Error(), zap.Error(err))
			return nil, status.Error(codes.Internal, "unable to get request result")
		}
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

func (s *Server) StreamCommunication(stream pb.Unicom_StreamCommunicationServer) error {
	ctx := stream.Context()
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		response, err := s.sendCommunication(ctx, &pb.SendCommunicationRequest{
			IsAsync: false,
			Domain:  in.GetDomain(),
			Email:   in.GetEmail(),
			Push:    in.GetPush(),
		})
		if err != nil {
			return err
		}
		err = stream.Send(response)
		if err != nil {
			return err
		}
	}
}

func (s *Server) validateRequest(req *pb.SendCommunicationRequest) error {
	if req.GetEmail() != nil && req.GetPush() != nil {
		return status.Error(codes.InvalidArgument, "invalid request for multiple notification types, please only send comms for a single medium")
	}
	if req.GetEmail() == nil && req.GetPush() == nil {
		return status.Error(codes.InvalidArgument, "invalid request must include any request medium")
	}
	return nil
}
