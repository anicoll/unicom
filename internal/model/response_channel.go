package model

type ResponseChannelRequest struct {
	Url          string
	WorkflowId   string
	Status       string
	ErrorMessage *string
}
