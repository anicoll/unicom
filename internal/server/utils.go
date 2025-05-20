package server

import (
	"io"
	"net/http"

	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"

	"github.com/anicoll/unicom/internal/email"
	"github.com/anicoll/unicom/internal/model"
	"github.com/anicoll/unicom/internal/push"
	"github.com/anicoll/unicom/internal/workflows"
)

// mapEmailRequestIn maps a protobuf EmailRequest to an internal email.Request structure,
// including attachments. Returns an error if attachment mapping fails.
func mapEmailRequestIn(req *pb.EmailRequest) (*email.Request, error) {
	if req == nil {
		return nil, nil
	}

	attachments, err := mapAttachmentsIn(req.Attachments)
	if err != nil {
		return nil, err
	}
	return &email.Request{
		FromAddress:      req.FromAddress,
		Subject:          req.Subject,
		ReplyToAddresses: []string{req.FromAddress},
		ToAddresses:      []string{req.ToAddress},
		HtmlBody:         req.Html,
		Attachments:      attachments,
	}, nil
}

// mapPushNotificationIn maps a protobuf PushRequest to an internal push.Notification structure,
// including language-specific content and optional subtitle.
func mapPushNotificationIn(req *pb.PushRequest) *push.Notification {
	notification := &push.Notification{
		IdempotencyKey:     req.GetIdempotencyKey(),
		ExternalCustomerId: req.GetExternalCustomerId(),
		Content: push.LanguageContent{
			Arabic:  req.GetContent().GetArabic(),
			English: req.GetContent().GetEnglish(),
		},
		Heading: push.LanguageContent{
			Arabic:  req.GetHeading().GetArabic(),
			English: req.GetHeading().GetEnglish(),
		},
	}

	if req.GetSubTitle() != nil {
		notification.SubTitle = &push.LanguageContent{
			Arabic:  req.GetSubTitle().GetArabic(),
			English: req.GetSubTitle().GetEnglish(),
		}
	}

	return notification
}

// mapAttachmentsIn converts a slice of protobuf Attachment objects to internal email.Attachment objects,
// downloading file data if a URL is provided. Returns an error if any download fails.
func mapAttachmentsIn(attachments []*pb.Attachment) ([]email.Attachment, error) {
	resp := make([]email.Attachment, len(attachments))

	for i, attachment := range attachments {
		if attachment.Url != nil {
			data, err := downloadFile(*attachment.Url)
			if err != nil {
				return nil, err
			}
			resp[i] = email.Attachment{
				Name: attachment.GetName(),
				Data: data,
			}
			continue
		}
		resp[i] = email.Attachment{
			Name: attachment.GetName(),
			Data: attachment.GetData(),
		}
	}
	return resp, nil
}

// mapWorkflowRequestToModel maps a workflow request and its ID to a model.Communication object,
// setting the communication type and response channels.
func mapWorkflowRequestToModel(workflowId string, req workflows.Request) *model.Communication {
	communicationType := model.Email
	if req.EmailRequest == nil {
		communicationType = model.Push
	}
	resp := &model.Communication{
		ID:               workflowId,
		Type:             communicationType,
		Domain:           req.Domain,
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

// downloadFile downloads the file from the specified URL and returns its contents as a byte slice.
// Returns an error if the download fails.
func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return io.ReadAll(resp.Body)
}
