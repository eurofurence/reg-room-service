package mailservice

import (
	"context"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams"
	"net/http"
)

type Impl struct {
	client  aurestclientapi.Client
	baseUrl string
}

func New(mailServiceBaseUrl string, apiToken string) (MailService, error) {
	if mailServiceBaseUrl != "" {
		return newClient(mailServiceBaseUrl, apiToken)
	} else {
		aulogging.Logger.NoCtx().Warn().Printf("service.mail_service not configured. Using in-memory simulator for mail service (not useful for production!)")
		return NewMock(), nil
	}
}

func newClient(mailServiceBaseUrl string, apiToken string) (MailService, error) {
	apiTokenClient, err := downstreams.ClientWith(
		downstreams.ApiTokenRequestManipulator(apiToken),
		"mail-service-breaker",
	)
	if err != nil {
		return nil, err
	}

	return &Impl{
		client:  apiTokenClient,
		baseUrl: mailServiceBaseUrl,
	}, nil
}

func errByStatus(err error, status int) error {
	if err != nil {
		return err
	}
	if status >= 300 {
		return DownstreamError
	}
	return nil
}

func (i Impl) SendEmail(ctx context.Context, request MailSendDto) error {
	url := fmt.Sprintf("%s/api/v1/mail", i.baseUrl)
	response := aurestclientapi.ParsedResponse{}
	err := i.client.Perform(ctx, http.MethodPost, url, request, &response)
	return errByStatus(err, response.Status)
}
