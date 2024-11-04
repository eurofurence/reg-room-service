package acceptance

import (
	"encoding/json"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/mailservice"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/require"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type tstWebResponse struct {
	status      int
	body        string
	contentType string
	location    string
	header      http.Header
}

func tstWebResponseFromResponse(response *http.Response) tstWebResponse {
	status := response.StatusCode
	ct := ""
	if val, ok := response.Header[headers.ContentType]; ok {
		ct = val[0]
	}
	loc := ""
	if val, ok := response.Header[headers.Location]; ok {
		loc = val[0]
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponse{
		status:      status,
		body:        string(body),
		contentType: ct,
		location:    loc,
		header:      response.Header,
	}
}

func tstAddAuth(request *http.Request, token string) {
	if token == tstValidApiToken() || token == tstInvalidApiToken() {
		request.Header.Set(common.ApiKeyHeader, token)
	} else if strings.HasPrefix(token, "access") {
		request.Header.Set(headers.Authorization, "Bearer "+token)
	} else if token != "" {
		request.AddCookie(&http.Cookie{
			Name:     "JWT",
			Value:    token,
			Domain:   "localhost",
			Expires:  time.Now().Add(10 * time.Minute),
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
		request.AddCookie(&http.Cookie{
			Name:     "AUTH",
			Value:    "access" + token,
			Domain:   "localhost",
			Expires:  time.Now().Add(10 * time.Minute),
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
		if token == valid_JWT_is_admin_sub1234567890 {
			request.Header.Set("X-Admin-Request", "available")
		}
	}
}

func tstPerformGet(relativeUrlWithLeadingSlash string, token string) tstWebResponse {
	request, err := http.NewRequest(http.MethodGet, ts.URL+relativeUrlWithLeadingSlash, nil)
	if err != nil {
		log.Fatal(err)
	}
	tstAddAuth(request, token)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponseFromResponse(response)
}

func tstPerformPut(relativeUrlWithLeadingSlash string, requestBody string, token string) tstWebResponse {
	request, err := http.NewRequest(http.MethodPut, ts.URL+relativeUrlWithLeadingSlash, strings.NewReader(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	tstAddAuth(request, token)
	request.Header.Set(headers.ContentType, web.ContentTypeApplicationJSON)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponseFromResponse(response)
}

func tstPerformPost(relativeUrlWithLeadingSlash string, requestBody string, token string) tstWebResponse {
	request, err := http.NewRequest(http.MethodPost, ts.URL+relativeUrlWithLeadingSlash, strings.NewReader(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	tstAddAuth(request, token)
	request.Header.Set(headers.ContentType, web.ContentTypeApplicationJSON)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponseFromResponse(response)
}

func tstPerformPostNoBody(relativeUrlWithLeadingSlash string, token string) tstWebResponse {
	request, err := http.NewRequest(http.MethodPost, ts.URL+relativeUrlWithLeadingSlash, nil)
	if err != nil {
		log.Fatal(err)
	}
	tstAddAuth(request, token)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponseFromResponse(response)
}

func tstPerformDelete(relativeUrlWithLeadingSlash string, token string) tstWebResponse {
	request, err := http.NewRequest(http.MethodDelete, ts.URL+relativeUrlWithLeadingSlash, nil)
	if err != nil {
		log.Fatal(err)
	}
	tstAddAuth(request, token)
	request.Header.Set(headers.ContentType, web.ContentTypeApplicationJSON)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponseFromResponse(response)
}

func tstRenderJson(v interface{}) string {
	representationBytes, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	return string(representationBytes)
}

// tip: dto := &v1.Group{}
func tstParseJson(body string, dto interface{}) {
	err := json.Unmarshal([]byte(body), dto)
	if err != nil {
		log.Fatal(err)
	}
}

func p[T any](v T) *T {
	return &v
}

func tstReadGroup(t *testing.T, location string) modelsv1.Group {
	readAgainResponse := tstPerformGet(location, tstValidAdminToken(t))
	result := modelsv1.Group{}
	tstParseJson(readAgainResponse.body, &result)
	return result
}

func tstReadRoom(t *testing.T, location string) modelsv1.Room {
	readAgainResponse := tstPerformGet(location, tstValidAdminToken(t))
	result := modelsv1.Room{}
	tstParseJson(readAgainResponse.body, &result)
	return result
}

func tstRequireErrorResponse(t *testing.T, response tstWebResponse, expectedStatus int, expectedMessage string, expectedDetails interface{}) {
	require.Equal(t, expectedStatus, response.status, "unexpected http response status")
	errorDto := modelsv1.Error{}
	tstParseJson(response.body, &errorDto)
	require.Equal(t, expectedMessage, errorDto.Message, "unexpected error code")
	expectedDetailsStr, ok := expectedDetails.(string)
	if ok && expectedDetailsStr != "" {
		require.EqualValues(t, url.Values{"details": []string{expectedDetailsStr}}, errorDto.Details, "unexpected error details")
	}
	expectedDetailsUrlValues, ok := expectedDetails.(url.Values)
	if ok {
		require.EqualValues(t, expectedDetailsUrlValues, errorDto.Details, "unexpected error details")
	}
}

func tstRequireSuccessResponse(t *testing.T, response tstWebResponse, expectedStatus int, resultBodyPtr interface{}) {
	require.Equal(t, expectedStatus, response.status, "unexpected http response status")
	tstParseJson(response.body, resultBodyPtr)
}

func tstEqualResponseBodies(t *testing.T, expected interface{}, actual interface{}) {
	// render both values to yaml and then compare - this gives easiest to debug differences
	expectedYaml, err := yaml.Marshal(expected)
	if err != nil {
		t.Errorf("failed to marshal expected body to yaml: %s", err)
	}
	actualYaml, err := yaml.Marshal(actual)
	if err != nil {
		t.Errorf("failed to marshal actual body to yaml: %s", err)
	}
	require.Equal(t, string(expectedYaml), string(actualYaml))
}

func tstRequireMailRequests(t *testing.T, expectedMailRequests ...mailservice.MailSendDto) {
	require.Equal(t, len(expectedMailRequests), len(mailMock.Recording()))
	for i, expected := range expectedMailRequests {
		actual := mailMock.Recording()[i]
		require.Equal(t, len(expected.To), len(actual.To))
		for i := range expected.To {
			require.Contains(t, actual.To[i], expected.To[i])
		}
		actual.To = expected.To
		require.Equal(t, len(expected.Variables), len(actual.Variables))
		require.EqualValues(t, expected, actual)
	}

	mailMock.Reset()
}
