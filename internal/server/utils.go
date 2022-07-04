package server

import (
	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"

	"github.com/anicoll/unicom/internal/email"
)

func mapEmailRequestIn(request interface{}) *email.Request {
	switch req := request.(type) {
	case *pb.SendAsyncRequest:
		return &email.Request{
			FromAddress:      req.FromAddress,
			Subject:          req.Subject,
			ReplyToAddresses: []string{req.FromAddress},
			ToAddresses:      []string{req.ToAddress},
			HtmlBody:         req.Html,
			Attachments:      mapAttachmentsIn(req.Attachments),
		}
	case *pb.SendSyncRequest:
		return &email.Request{
			FromAddress:      req.FromAddress,
			Subject:          req.Subject,
			ReplyToAddresses: []string{req.FromAddress},
			ToAddresses:      []string{req.ToAddress},
			HtmlBody:         req.Html,
			Attachments:      mapAttachmentsIn(req.Attachments),
		}
	}
	return nil
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
