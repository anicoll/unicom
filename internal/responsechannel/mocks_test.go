package responsechannel_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/mock"
)

type mockSqsClient struct {
	mock.Mock
}

func (db *mockSqsClient) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	args := db.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sqs.SendMessageOutput), args.Error(1)
}
