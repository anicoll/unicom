package server

import (
	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"

	"github.com/anicoll/unicom/internal/email"
	"github.com/anicoll/unicom/internal/model"
	"github.com/anicoll/unicom/internal/workflows"
)

func mapEmailRequestIn(req *pb.EmailRequest) *email.Request {
	if req == nil {
		return nil
	}
	return &email.Request{
		FromAddress:      req.FromAddress,
		Subject:          req.Subject,
		ReplyToAddresses: []string{req.FromAddress},
		ToAddresses:      []string{req.ToAddress},
		HtmlBody:         req.Html,
		Attachments:      mapAttachmentsIn(req.Attachments),
	}

}

func mapAttachmentsIn(attachments []*pb.Attachment) []email.Attachment {
	resp := make([]email.Attachment, 0, len(attachments))

	for i, attachment := range attachments {
		resp[i] = email.Attachment{
			Name: attachment.GetName(),
			Data: attachment.GetData(),
		}
	}
	return resp
}

func mapWorkflowRequestToModel(workflowId string, req workflows.Request) *model.Communication {
	resp := &model.Communication{
		ID:               workflowId,
		Domain:           req.Domain,
		Type:             model.Email,
		ResponseChannels: make([]*model.ResponseChannel, len(req.ResponseRequests)),
	}
	for index, channel := range req.ResponseRequests {
		resp.ResponseChannels[index] = &model.ResponseChannel{
			Type:            channel.Type,
			Url:             channel.Url,
			ID:              channel.ID,
			CommunicationID: workflowId,
		}
	}
	return resp
}
