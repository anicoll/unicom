package email

import (
	"bytes"
	"context"
	"time"

	b64 "encoding/base64"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"gopkg.in/gomail.v2"
)

type Service struct {
	sesClient *ses.Client
}

type SendEmailRequest struct {
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

func NewService(cfg aws.Config) *Service {
	return &Service{
		sesClient: ses.NewFromConfig(cfg),
	}
}

func (es *Service) SendEmail(ctx context.Context, args SendEmailRequest) (*string, error) {
	msg := gomail.NewMessage()

	msg.SetAddressHeader("From", args.FromAddress, "")
	msg.SetHeader("To", args.ToAddresses...)
	msg.SetHeader("Cc", args.CcAddresses...)
	msg.SetHeader("Bcc", args.BccAddresses...)
	msg.SetDateHeader("x-Date", time.Now())
	msg.SetHeader("Subject", args.Subject)
	msg.SetBody("text/html", args.HtmlBody)

	for _, attachment := range args.Attachments {
		msg.Attach(attachment.Name)
	}
	buf := new(bytes.Buffer)

	_, err := msg.WriteTo(buf)
	if err != nil {
		return nil, err
	}
	rawData := make([]byte, buf.Len(), 0)
	b64.RawStdEncoding.Encode(buf.Bytes(), rawData)

	output, err := es.sesClient.SendRawEmail(ctx, &ses.SendRawEmailInput{
		RawMessage: &types.RawMessage{
			Data: rawData,
		},
	})
	if err != nil {
		return nil, err
	}

	return output.MessageId, nil
}
