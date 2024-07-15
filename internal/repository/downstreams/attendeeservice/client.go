package attendeeservice

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/eurofurence/reg-room-service/internal/repository/config"

	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"

	"github.com/eurofurence/reg-room-service/internal/repository/downstreams"
)

type Impl struct {
	myTokenClient aurestclientapi.Client
	baseUrl       string
}

func New(attendeeServiceBaseUrl string) (AttendeeService, error) {
	if attendeeServiceBaseUrl == "" {
		return nil, errors.New("service.attendee_service not configured. This service cannot function without the attendee service, though you can run it in inmemory database mode for development.")
	}
	conf, err := config.GetApplicationConfig()
	if err != nil {
		return nil, err
	}

	myTokenClient, err := downstreams.ClientWith(
		downstreams.CookiesOrAuthHeaderForwardingRequestManipulator(conf.Security),
		"attendee-service-breaker",
	)
	if err != nil {
		return nil, err
	}

	return &Impl{
		myTokenClient: myTokenClient,
		baseUrl:       attendeeServiceBaseUrl,
	}, nil
}

type AttendeeIdList struct {
	Ids []int64 `json:"ids"`
}

func (i *Impl) ListMyRegistrationIds(ctx context.Context) ([]int64, error) {
	url := fmt.Sprintf("%s/api/rest/v1/attendees", i.baseUrl)
	bodyDto := AttendeeIdList{
		Ids: make([]int64, 0),
	}
	response := aurestclientapi.ParsedResponse{
		Body: &bodyDto,
	}
	err := i.myTokenClient.Perform(ctx, http.MethodGet, url, nil, &response)
	if response.Status == http.StatusNotFound {
		// not really an error - this user has no registrations
		return make([]int64, 0), nil
	}
	return bodyDto.Ids, downstreams.ErrByStatus(err, response.Status)
}

type StatusDto struct {
	Status Status `json:"status"`
}

func (i *Impl) GetStatus(ctx context.Context, id int64) (Status, error) {
	url := fmt.Sprintf("%s/api/rest/v1/attendees/%d/status", i.baseUrl, id)
	bodyDto := StatusDto{
		Status: StatusDeleted,
	}
	response := aurestclientapi.ParsedResponse{
		Body: &bodyDto,
	}
	err := i.myTokenClient.Perform(ctx, http.MethodGet, url, nil, &response)
	if response.Status == http.StatusNotFound {
		// not really an error - this user has no registrations - treat as deleted
		return StatusDeleted, nil
	}
	return bodyDto.Status, downstreams.ErrByStatus(err, response.Status)
}
