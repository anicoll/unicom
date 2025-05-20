package responsechannel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"
	"github.com/anicoll/unicom/internal/model"
)

type WebhookService struct {
	client *http.Client
}

func NewWebhookService(client *http.Client) *WebhookService {
	return &WebhookService{
		client: client,
	}
}

func (s *WebhookService) Send(ctx context.Context, req model.ResponseChannelRequest) (*string, error) {
	data, err := json.Marshal(pb.ResponseEvent{
		WorkflowId:   req.WorkflowId,
		Status:       req.Status,
		ErrorMessage: req.ErrorMessage,
	})
	if err != nil {
		return nil, err
	}

	reader := bytes.NewBuffer(data)
	httpRequest, err := http.NewRequest(http.MethodPost, req.Url, reader)
	if err != nil {
		return nil, err
	}

	response, err := s.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil, nil
	}

	body, err := io.ReadAll(response.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		log.Fatalln(err)
	}

	// TODO: parse error and return it
	fmt.Println("BODY")
	// Do I need to parse anything here?
	fmt.Println(string(body))

	return nil, nil
}
