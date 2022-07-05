package workflows

import (
	"context"

	"github.com/anicoll/unicom/internal/email"
	"github.com/anicoll/unicom/internal/model"
	"github.com/anicoll/unicom/internal/sqs"
)

type emailService interface {
	Send(ctx context.Context, args email.Request) (*string, error)
}

type sqsService interface {
	Send(ctx context.Context, args sqs.Request) (*string, error)
}

type postgres interface {
	SetCommunicationStatus(ctx context.Context, workflowId string, status model.Status) error
	CreateResponseChannel(ctx context.Context, channel model.ResponseChannel) error
	SetResponseChannelStatus(ctx context.Context, id, externalId string, status model.Status) error
}

type UnicomActivities struct {
	emailService emailService
	sqsService   sqsService
	database     postgres
}

func NewActivities(es emailService, sqs sqsService, db postgres) *UnicomActivities {
	return &UnicomActivities{
		emailService: es,
		sqsService:   sqs,
		database:     db,
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

func (a *UnicomActivities) MarkCommunicationAsFailed(ctx context.Context, workflowId string) error {
	return a.database.SetCommunicationStatus(ctx, workflowId, model.Failed)
}

func (a *UnicomActivities) MarkCommunicationAsSent(ctx context.Context, workflowId string) error {
	return a.database.SetCommunicationStatus(ctx, workflowId, model.Success)
}

func (a *UnicomActivities) SaveResponseChannelOutcome(ctx context.Context, id, externalId string, status model.Status) error {
	return a.database.SetResponseChannelStatus(ctx, id, externalId, status)
}
