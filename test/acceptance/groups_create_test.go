package acceptance

import (
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eurofurence/reg-room-service/docs"
	v1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

const validGroupLocationRegex = "^\\/api\\/rest\\/v1\\/groups\\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"

func TestGroupsCreate_UserSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is not in any group")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusApproved)
	token := tstValidUserToken(t, 101)

	docs.When("When they create a room group with valid data")
	groupSent := v1.GroupCreate{
		Name:     "kittens",
		Flags:    []string{"public"},
		Comments: p("A nice comment"),
		Owner:    0, // myself
	}
	response := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), token)

	docs.Then("Then the group is successfully created")
	require.Equal(t, http.StatusCreated, response.status, "unexpected http response status")
	require.Regexp(t, validGroupLocationRegex, response.location, "invalid location header in response")

	docs.Then("And it can be read again by an admin")
	groupReadAgain := tstReadGroup(t, response.location)
	require.Equal(t, groupSent.Name, groupReadAgain.Name)

	docs.Then("And it contains exactly the user as owner and no invites")
	require.Equal(t, int32(42), groupReadAgain.Owner)
	require.Equal(t, 1, len(groupReadAgain.Members))
	require.Equal(t, int32(42), groupReadAgain.Members[0].ID)
	require.Equal(t, 0, len(groupReadAgain.Invites))
}

func TestGroupsCreate_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized admin")
	token := tstValidAdminToken(t)

	docs.Given("And a registered attendee with an active registration who is not in any group")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusApproved)

	docs.When("When the admin creates a room group with that attendee as owner")
	groupSent := v1.GroupCreate{
		Name:     "kittens",
		Flags:    []string{"public"},
		Comments: p("A nice comment"),
		Owner:    42,
	}
	response := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), token)

	docs.Then("Then the group is successfully created")
	require.Equal(t, http.StatusCreated, response.status, "unexpected http response status")
	require.Regexp(t, validGroupLocationRegex, response.location, "invalid location header in response")

	docs.Then("And it can be read again")
	groupReadAgain := tstReadGroup(t, response.location)
	require.Equal(t, groupSent.Name, groupReadAgain.Name)

	docs.Then("And it contains exactly the attendee as owner and no invites")
	require.Equal(t, int32(42), groupReadAgain.Owner)
	require.Equal(t, 1, len(groupReadAgain.Members))
	require.Equal(t, int32(42), groupReadAgain.Members[0].ID)
	require.Equal(t, 0, len(groupReadAgain.Invites))
}

func TestGroupsCreate_AnonymousDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an unauthenticated user")
	token := tstNoToken()

	docs.When("When they attempt create a room group")
	groupSent := v1.GroupCreate{
		Name:     "kittens",
		Flags:    []string{"public"},
		Comments: p("A nice comment"),
	}
	response := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")
}

func TestGroupsCreate_CrossUserDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized non-admin user with an active registration")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusApproved)
	token := tstValidUserToken(t, 101)

	docs.Given("Given another user with an active registration who is not in any group")
	attMock.SetupRegistered("1234567890", 43, attendeeservice.StatusApproved)

	docs.When("When the non-admin user tries to create a room group with a different owner than themselves")
	groupSent := v1.GroupCreate{
		Name:     "kittens",
		Flags:    []string{"public"},
		Comments: p("A nice comment"),
		Owner:    43, // not myself
	}
	response := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), token)

	docs.Then("Then the group is successfully created, but with the non-admin user as owner")
	require.Equal(t, http.StatusCreated, response.status, "unexpected http response status")
	require.Regexp(t, validGroupLocationRegex, response.location, "invalid location header in response")

	docs.Then("And it can be read again")
	groupReadAgain := tstReadGroup(t, response.location)
	require.Equal(t, groupSent.Name, groupReadAgain.Name)

	docs.Then("And it contains exactly the non-admin attendee as owner and no invites")
	require.Equal(t, int32(42), groupReadAgain.Owner)
	require.Equal(t, 1, len(groupReadAgain.Members))
	require.Equal(t, int32(42), groupReadAgain.Members[0].ID)
	require.Equal(t, 0, len(groupReadAgain.Invites))
}

func TestGroupsCreate_UserNoReg(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with NO registration")
	token := tstValidUserToken(t, 101)

	docs.When("When they try to create a room group")
	groupSent := v1.GroupCreate{
		Name:     "kittens",
		Flags:    []string{"public"},
		Comments: p("A nice comment"),
		Owner:    0, // myself
	}
	response := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "attendee.notfound", "you do not have a valid registration")
}

func TestGroupsCreate_UserNonAttendingReg(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with a registration in non-attending status")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusCancelled)
	token := tstValidUserToken(t, 101)

	docs.When("When they try to create a room group")
	groupSent := v1.GroupCreate{
		Name:     "kittens",
		Flags:    []string{"public"},
		Comments: p("A nice comment"),
		Owner:    0, // myself
	}
	response := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "attendee.status.not.attending", "registration is not in attending status")
}

func TestGroupsCreate_InvalidJSONSyntax(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with a registration in non-attending status")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusCancelled)
	token := tstValidUserToken(t, 101)

	docs.When("When they try to create a room group, but supply syntactically invalid JSON")
	response := tstPerformPost("/api/rest/v1/groups", `{"name":"invalid":"extra"`, token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "group.data.invalid", "please check if your provided JSON is valid")
}
