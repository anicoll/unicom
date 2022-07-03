package workflows

import (
	"context"
)

type emailService interface {
}

type UnicomActivities struct {
	emailService emailService
}

func NewActivities(es emailService) *UnicomActivities {
	return &UnicomActivities{
		emailService: es,
	}
}

func (a *UnicomActivities) SendEmail(ctx context.Context) error {

	return nil
}

func (a *UnicomActivities) NotifySqs(ctx context.Context) error {

	return nil
}

func (a *UnicomActivities) NotifyWebhook(ctx context.Context) error {

	return nil
}
