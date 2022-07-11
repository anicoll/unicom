package model

import (
	"time"
)

type Model struct {
	CreatedAt time.Time
	SentAt    time.Time
}

type Status string

const (
	Pending Status = "PENDING"
	Success Status = "SUCCESS"
	Failed  Status = "FAILED"
)

type NotificationType string

const (
	Email NotificationType = "EMAIL"
	Sms   NotificationType = "SMS"
	Push  NotificationType = "PUSH"
)

type Communication struct {
	Model
	ID               string
	ExternalId       *string
	Domain           string
	Status           Status
	Type             NotificationType
	ResponseChannels []*ResponseChannel
}

type ResponseChannelType string

const (
	Sqs         ResponseChannelType = "SQS"
	Webhook     ResponseChannelType = "WEBHOOK"
	EventBridge ResponseChannelType = "EVENT_BRIDGE"
)

type ResponseChannel struct {
	Model
	ID              string
	CommunicationID string
	Type            ResponseChannelType
	Status          Status
	Url             string
}
