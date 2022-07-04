package workflows

import (
	"context"

	"github.com/anicoll/unicom/internal/email"
	"github.com/anicoll/unicom/internal/sqs"
)

type emailService interface {
	Send(ctx context.Context, args email.Request) (*string, error)
}

type sqsService interface {
	Send(ctx context.Context, args sqs.Request) (*string, error)
}

type UnicomActivities struct {
	emailService emailService
	sqsService   sqsService
}

func NewActivities(es emailService, sqs sqsService) *UnicomActivities {
	return &UnicomActivities{
		emailService: es,
		sqsService:   sqs,
	}
}

func (a *UnicomActivities) SendEmail(ctx context.Context, req email.Request) (*string, error) {
	return a.emailService.Send(ctx, req)
}

func (a *UnicomActivities) NotifySqs(ctx context.Context, req sqs.Request) (*string, error) {
	return a.sqsService.Send(ctx, req)
}

func (a *UnicomActivities) NotifyWebhook(ctx context.Context) error {

	return nil
}
