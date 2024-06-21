package acceptance

import (
	"github.com/eurofurence/reg-room-service/docs"
	v1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/stretchr/testify/require"
	"net/http"
	"path"
	"testing"
)

func TestGroupsUpdate_UserSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is the owner of a group")
	token := tstValidUserToken(t, 101)
	attMock.SetupRegistered("101", 42, attendeeservice.StatusApproved)

	groupSent := v1.GroupCreate{
		Name:     "kittens",
		Flags:    []string{"public"},
		Comments: p("A nice comment"),
		Owner:    0, // myself
	}

	response := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), token)

	docs.When("When they retrieve the group and update the name of a group")
	getGroup := tstReadGroup(t, response.location)
	savedID := getGroup.ID
	getGroup.Name = "dogs"

	response = tstPerformPut(path.Join("/api/rest/v1/groups/", getGroup.ID), tstRenderJson(getGroup), token)
	docs.Then("Then the group should be successfully updated")

	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	require.Regexp(t, validGroupLocationRegex, response.location, "invalid location header in response")

	getGroup = tstReadGroup(t, response.location)
	require.Equal(t, savedID, getGroup.ID)
	require.Equal(t, "dogs", getGroup.Name)
	require.Len(t, getGroup.Flags, 1)
}

func TestGroupsUpdate_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized admin")
	token := tstValidAdminToken(t)

	docs.Given("Given an attendee with an active registration")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusApproved)

	docs.When("When the admin creates a room group with that attendee as owner")
	groupSent := v1.GroupCreate{
		Name:     "kittens",
		Flags:    []string{"public"},
		Comments: p("A nice comment"),
		Owner:    42,
	}

	response := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), token)

	docs.When("When they retrieve the group and update the name of a group")
	getGroup := tstReadGroup(t, response.location)
	savedID := getGroup.ID
	getGroup.Name = "dogs"

	response = tstPerformPut(path.Join("/api/rest/v1/groups/", getGroup.ID), tstRenderJson(getGroup), token)
	docs.Then("Then the group should be successfully updated")

	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	require.Regexp(t, validGroupLocationRegex, response.location, "invalid location header in response")

	getGroup = tstReadGroup(t, response.location)
	require.Equal(t, savedID, getGroup.ID)
	require.Equal(t, "dogs", getGroup.Name)
	require.Len(t, getGroup.Flags, 1)
}

func TestGroupsUpdate_AnonymousDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an existing group that was created by an authenticated user")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusApproved)
	authenticatedToken := tstValidUserToken(t, 101)

	groupSent := v1.GroupCreate{
		Name:     "kittens",
		Flags:    []string{"public"},
		Comments: p("A nice comment"),
		Owner:    42,
	}

	response := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), authenticatedToken)

	docs.Given("Given an unauthenticated user")
	token := tstNoToken()

	docs.When("When they try to update the group")
	getGroup := tstReadGroup(t, response.location)

	getGroup.Name = "dogs"
	response = tstPerformPut(path.Join("/api/rest/v1/groups/", getGroup.ID), tstRenderJson(getGroup), token)

	docs.Then("Then the request should be denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")
}
