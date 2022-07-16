package server

import (
	"io/ioutil"
	"net/http"

	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"

	"github.com/anicoll/unicom/internal/email"
	"github.com/anicoll/unicom/internal/model"
	"github.com/anicoll/unicom/internal/push"
	"github.com/anicoll/unicom/internal/workflows"
)

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
			Arabic:  req.GetHeading().GetArabic(),
			English: req.GetHeading().GetEnglish(),
		}
	}

	return notification
}

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

func mapWorkflowRequestToModel(workflowId string, req workflows.Request) *model.Communication {
	resp := &model.Communication{
		ID:               workflowId,
		Type:             model.Email,
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

// DownloadFile will download the file and return the data.
func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
