package workflows

import (
	"context"

	"github.com/anicoll/unicom/internal/email"
	"github.com/anicoll/unicom/internal/model"
)

type emailService interface {
	Send(ctx context.Context, args email.Request) (*string, error)
}

type notificationService interface {
	Send(ctx context.Context, args model.ResponseChannelRequest) (*string, error)
}

type postgres interface {
	SetCommunicationStatus(ctx context.Context, workflowId string, status model.Status, externalId *string) error
	CreateResponseChannel(ctx context.Context, channel model.ResponseChannel) error
	SetResponseChannelStatus(ctx context.Context, id, externalId string, status model.Status) error
}

type UnicomActivities struct {
	emailService emailService
	sqsService   notificationService
	database     postgres
}

func NewActivities(es emailService, sqs, webhook notificationService, db postgres) *UnicomActivities {
	return &UnicomActivities{
		emailService: es,
		sqsService:   sqs,
		database:     db,
	}
}

func (a *UnicomActivities) SendEmail(ctx context.Context, req email.Request) (*string, error) {
	return a.emailService.Send(ctx, req)
}

func (a *UnicomActivities) NotifySqs(ctx context.Context, req model.ResponseChannelRequest) (*string, error) {
	return a.sqsService.Send(ctx, req)
}

func (a *UnicomActivities) NotifyWebhook(ctx context.Context, req model.ResponseChannelRequest) error {
	return nil
}

func (a *UnicomActivities) UpdateCommunicationStatus(ctx context.Context, workflowId string, status model.Status, externalId *string) error {
	return a.database.SetCommunicationStatus(ctx, workflowId, status, externalId)
}

func (a *UnicomActivities) SaveResponseChannelOutcome(ctx context.Context, id, externalId string, status model.Status) error {
	return a.database.SetResponseChannelStatus(ctx, id, externalId, status)
}
