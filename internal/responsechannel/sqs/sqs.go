package sqs

import (
	"context"
	"encoding/json"

	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"
	"github.com/anicoll/unicom/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	aws_sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
)

type sqsClient interface {
	SendMessage(ctx context.Context, params *aws_sqs.SendMessageInput, optFns ...func(*aws_sqs.Options)) (*aws_sqs.SendMessageOutput, error)
}

type Service struct {
	sqsClient sqsClient
}

func NewService(client sqsClient) *Service {
	return &Service{
		sqsClient: client,
	}
}

func (s *Service) Send(ctx context.Context, req model.ResponseChannelRequest) (*string, error) {
	data, err := json.Marshal(pb.ResponseEvent{
		WorkflowId:   req.WorkflowId,
		Status:       req.Status,
		ErrorMessage: req.ErrorMessage,
	})
	if err != nil {
		return nil, err
	}

	response, err := s.sqsClient.SendMessage(ctx, &aws_sqs.SendMessageInput{
		QueueUrl:               aws.String(req.Url),
		MessageBody:            aws.String(string(data)),
		MessageDeduplicationId: aws.String(req.WorkflowId),
	})
	if err != nil {
		return nil, err
	}
	return response.MessageId, nil
}
