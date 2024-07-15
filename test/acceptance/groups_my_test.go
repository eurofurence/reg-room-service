package acceptance

import (
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eurofurence/reg-room-service/docs"
)

// find my group

func TestGroupsMy_UserSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a registered attendee with an active registration who is in a group")
	id1 := setupExistingGroup(t, "kittens", true, "101")

	docs.When("When the user requests their group")
	token := tstValidUserToken(t, 101)
	response := tstPerformGet("/api/rest/v1/groups/my", token)

	docs.Then("Then the request is successful and the response is as expected")
	actual := modelsv1.Group{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.Group{
		ID:          id1,
		Name:        "kittens",
		Flags:       []string{"public"},
		Comments:    p("A nice comment for kittens"),
		MaximumSize: p(int32(6)),
		Owner:       42,
		Members: []modelsv1.Member{
			{
				ID:       42,
				Nickname: "",
			},
		},
		Invites: nil,
	}
	require.EqualValues(t, expected, actual, "unexpected differences in response body")
}

func TestGroupsMy_AnonymousDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an unauthenticated user")
	token := tstNoToken()

	docs.When("When they request their group")
	response := tstPerformGet("/api/rest/v1/groups/my", token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")
}

func TestGroupsMy_UserNoReg(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with NO registration")
	token := tstValidUserToken(t, 101)

	docs.When("When they request their group")
	response := tstPerformGet("/api/rest/v1/groups/my", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "attendee.notfound", "you do not have a valid registration")
}

func TestGroupsMy_UserNonAttendingReg(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with a registration in non-attending status")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusNew)
	token := tstValidUserToken(t, 101)

	docs.When("When they request their group")
	response := tstPerformGet("/api/rest/v1/groups/my", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "attendee.status.not.attending", "registration is not in attending status")
}

func TestGroupsMy_UserNoGroup(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with a registration in attending status")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusPartiallyPaid)
	token := tstValidUserToken(t, 101)

	docs.Given("Given they are not in any group")

	docs.When("When they request their group")
	response := tstPerformGet("/api/rest/v1/groups/my", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "group.member.notfound", "not in a group")
}
