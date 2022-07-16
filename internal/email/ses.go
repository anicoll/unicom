package email

import (
	"bytes"
	"context"

	ses "github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type Service struct {
	sesClient *ses.Client
}

type Request struct {
	FromAddress      string
	Subject          string
	ReplyToAddresses []string
	CcAddresses      []string
	BccAddresses     []string
	ToAddresses      []string
	HtmlBody         string
	Attachments      []Attachment
}

type Attachment struct {
	Name string
	Data []byte
}

func NewService(client *ses.Client) *Service {
	return &Service{
		sesClient: client,
	}
}

func (es *Service) Send(ctx context.Context, args Request) (*string, error) {
	msg := NewMessage()
	msg.SetHeader("From", args.FromAddress)
	msg.SetHeader("To", args.ToAddresses...)
	msg.SetHeader("Cc", args.CcAddresses...)
	msg.SetHeader("Bcc", args.BccAddresses...)
	msg.SetHeader("Subject", args.Subject)
	msg.SetBody("text/html", args.HtmlBody)

	for _, attachment := range args.Attachments {
		msg.Attach(attachment.Name, attachment.Data)
	}

	var emailRaw bytes.Buffer
	_, err := msg.WriteTo(&emailRaw)
	if err != nil {
		return nil, err
	}

	output, err := es.sesClient.SendEmail(ctx, &ses.SendEmailInput{
		Content: &types.EmailContent{
			Raw: &types.RawMessage{
				Data: emailRaw.Bytes(),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return output.MessageId, nil
}
