package sqs

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Service struct {
	sqsClient *aws_sqs.Client
}

type Request struct {
	Queue        string
	WorkflowId   string
	Status       string
	ErrorMessage *string
}

func NewService(cfg aws.Config) *Service {
	return &Service{
		sqsClient: aws_sqs.NewFromConfig(cfg),
	}
}

type message struct {
	WorkflowId   string  `json:"name:workflowId"`
	Status       string  `json:"name:status"`
	ErrorMessage *string `json:"name:errorMessage"`
}

func (s *Service) Send(ctx context.Context, req Request) (*string, error) {
	data, err := json.Marshal(message{
		WorkflowId:   req.WorkflowId,
		Status:       req.Status,
		ErrorMessage: req.ErrorMessage,
	})
	if err != nil {
		return nil, err
	}

	response, err := s.sqsClient.SendMessage(ctx, &aws_sqs.SendMessageInput{
		QueueUrl:    aws.String(req.Queue),
		MessageBody: aws.String(string(data)),
	})
	if err != nil {
		return nil, err
	}
	return response.MessageId, nil
}
