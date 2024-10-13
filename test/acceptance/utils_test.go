package acceptance

import (
	"encoding/json"
	"fmt"
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
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

func setupExistingGroup(t *testing.T, name string, public bool, subject string, additionalMemberSubjects ...string) string {
	flags := []string{}
	if public {
		flags = append(flags, "public")
	}
	badgeNo := registerSubject(subject)

	groupSent := modelsv1.GroupCreate{
		Name:     name,
		Flags:    flags,
		Comments: p("A nice comment for " + name),
		Owner:    badgeNo,
	}
	response := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), tstValidAdminToken(t))
	require.Equal(t, http.StatusCreated, response.status, "unexpected http response status")
	require.Regexp(t, validGroupLocationRegex, response.location, "invalid location header in response")

	for _, addSubject := range additionalMemberSubjects {
		addBadgeNo := registerSubject(addSubject)
		addResponse := tstPerformPostNoBody(fmt.Sprintf("%s/members/%d?force=true", response.location, addBadgeNo), tstValidAdminToken(t))
		require.Equal(t, http.StatusNoContent, addResponse.status, "unexpected http response status")
	}
	mailMock.Reset()

	locs := strings.Split(response.location, "/")
	return locs[len(locs)-1]
}

func registerSubject(subject string) int64 {
	switch subject {
	case "101":
		attMock.SetupRegistered("101", 42, attendeeservice.StatusApproved, "Squirrel", "squirrel@example.com")
		return 42

	case "202":
		attMock.SetupRegistered("202", 43, attendeeservice.StatusPaid, "Snep", "snep@example.com")
		return 43

	default:
		attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusCancelled, "Panther", "panther@example.com")
		return 84
	}
}

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
	if strings.HasPrefix(token, "access") {
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

func tstGroupState(t *testing.T, id string, location string, addMembers []modelsv1.Member, addInvites []modelsv1.Member) {
	t.Helper()

	response := tstPerformGet(location, tstValidAdminToken(t))
	actual := modelsv1.Group{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.Group{
		ID:          id,
		Name:        "kittens",
		Flags:       []string{},
		Comments:    p("A nice comment for kittens"),
		MaximumSize: 6,
		Owner:       42,
		Members: []modelsv1.Member{
			{
				ID:       42,
				Nickname: "Squirrel",
			},
		},
		Invites: nil,
	}
	expected.Members = append(expected.Members, addMembers...)
	expected.Invites = addInvites
	tstEqualResponseBodies(t, expected, actual)
}

func tstRequireErrorResponse(t *testing.T, response tstWebResponse, expectedStatus int, expectedMessage string, expectedDetails interface{}) {
	require.Equal(t, expectedStatus, response.status, "unexpected http response status")
	errorDto := modelsv1.Error{}
	tstParseJson(response.body, &errorDto)
	require.Equal(t, expectedMessage, string(errorDto.Message), "unexpected error code")
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

func tstGroupMailToOwner(cid string, groupName string, target string, object string) mailservice.MailSendDto {
	_, targetNick, targetEmail := tstInfosBySubject(target)
	objectBadge, objectNick, _ := tstInfosBySubject(object)

	return mailservice.MailSendDto{
		CommonID: cid,
		Lang:     "en-US",
		To:       []string{targetEmail},
		Variables: map[string]string{
			"nickname":            targetNick,
			"groupname":           groupName,
			"object_badge_number": objectBadge,
			"object_nickname":     objectNick,
		},
	}
}

func tstGroupMailToMember(cid string, groupName string, target string, url string) mailservice.MailSendDto {
	_, targetNick, targetEmail := tstInfosBySubject(target)

	result := mailservice.MailSendDto{
		CommonID: cid,
		Lang:     "en-US",
		To:       []string{targetEmail},
		Variables: map[string]string{
			"nickname":  targetNick,
			"groupname": groupName,
		},
	}
	if url != "" {
		result.Variables["url"] = url
	}
	return result
}

func tstInfosBySubject(subject string) (string, string, string) {
	switch subject {
	case "101":
		return "42", "Squirrel", "squirrel@example.com"
	case "202":
		return "43", "Snep", "snep@example.com"
	default:
		return "84", "Panther", "panther@example.com"
	}
}
