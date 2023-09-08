package responsechannel

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"
	"github.com/anicoll/unicom/internal/model"
)

type sqsClient interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type SQSService struct {
	sqsClient sqsClient
}

func NewSQSService(client sqsClient) *SQSService {
	return &SQSService{
		sqsClient: client,
	}
}

func (s *SQSService) Send(ctx context.Context, req model.ResponseChannelRequest) (*string, error) {
	data, err := json.Marshal(pb.ResponseEvent{
		WorkflowId:   req.WorkflowId,
		Status:       req.Status,
		ErrorMessage: req.ErrorMessage,
	})
	if err != nil {
		return nil, err
	}

	response, err := s.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:               aws.String(req.Url),
		MessageBody:            aws.String(string(data)),
		MessageDeduplicationId: aws.String(req.WorkflowId),
	})
	if err != nil {
		return nil, err
	}
	return response.MessageId, nil
}
