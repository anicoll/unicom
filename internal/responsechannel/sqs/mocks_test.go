package sqs_test

import (
	"context"

	aws_sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/mock"
)

type mockSqsClient struct {
	mock.Mock
}

func (db *mockSqsClient) SendMessage(ctx context.Context, params *aws_sqs.SendMessageInput, optFns ...func(*aws_sqs.Options)) (*aws_sqs.SendMessageOutput, error) {
	args := db.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*aws_sqs.SendMessageOutput), args.Error(1)
}
