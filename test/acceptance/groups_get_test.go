package acceptance

import (
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"net/http"
	"testing"

	"github.com/eurofurence/reg-room-service/docs"
)

func TestGroupsGet_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized admin")
	token := tstValidAdminToken(t)

	docs.Given("And a registered attendee with an active registration who is in a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")

	docs.When("When the admin requests the group information")
	response := tstPerformGet("/api/rest/v1/groups/"+id1, token)

	docs.Then("Then the response is as expected and includes all information")
	actual := modelsv1.Group{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.Group{
		ID:          id1,
		Name:        "kittens",
		Flags:       []string{},
		Comments:    p("A nice comment for kittens"),
		MaximumSize: 6,
		Owner:       42,
		Members: []modelsv1.Member{
			{
				ID:       42,
				Nickname: "",
			},
		},
		Invites: nil,
	}
	tstEqualResponseBodies(t, expected, actual)
}

func TestGroupsGet_UserSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is in a group")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	token := tstValidUserToken(t, 101)

	docs.When("When they request the group information")
	response := tstPerformGet("/api/rest/v1/groups/"+id1, token)

	docs.Then("Then the response is as expected and includes all user visible fields")
	actual := modelsv1.Group{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.Group{
		ID:          id1,
		Name:        "kittens",
		Flags:       []string{"public"},
		Comments:    p("A nice comment for kittens"),
		MaximumSize: 6,
		Owner:       42,
		Members: []modelsv1.Member{
			{
				ID:       42,
				Nickname: "",
			},
		},
	}
	tstEqualResponseBodies(t, expected, actual)
}

func TestGroupsGet_AnonymousDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a group exists")
	id1 := setupExistingGroup(t, "kittens", true, "101")

	docs.Given("Given an unauthenticated user")
	token := tstNoToken()

	docs.When("When they attempt to access the group information")
	response := tstPerformGet("/api/rest/v1/groups/"+id1, token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")
}

func TestGroupsGet_UserNotMemberAllow(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given someone with an active registration who is in a non-public group")
	id1 := setupExistingGroup(t, "kittens", false, "101")

	docs.Given("Given another user with an active registration who is not in the group but knows the secret group id")
	attMock.SetupRegistered("1234567890", 43, attendeeservice.StatusApproved)
	token := tstValidUserToken(t, 1234567890)

	docs.When("When they attempt to access the group information")
	response := tstPerformGet("/api/rest/v1/groups/"+id1, token)

	docs.Then("Then the request is successful and the user visible information is returned")
	actual := modelsv1.Group{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.Group{
		ID:          id1,
		Name:        "kittens",
		Flags:       []string{},
		Comments:    p("A nice comment for kittens"),
		MaximumSize: 6,
		Owner:       42,
		Members: []modelsv1.Member{
			{
				ID:       42,
				Nickname: "",
			},
		},
	}
	tstEqualResponseBodies(t, expected, actual)
}

func TestGroupsGet_UserNoReg(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given someone with an active registration who is in a group")
	id1 := setupExistingGroup(t, "kittens", true, "101")

	docs.Given("Given another authorized non-admin user with NO registration")
	token := tstValidUserToken(t, 1234567890)

	docs.When("When they attempt to access the group information")
	response := tstPerformGet("/api/rest/v1/groups/"+id1, token)

	docs.Then("Then the request is denied with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "attendee.notfound", "you do not have a valid registration")
}

func TestGroupsGet_UserNonAttendingReg(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given someone with an active registration who is in a group")
	id1 := setupExistingGroup(t, "kittens", true, "101")

	docs.Given("Given another authorized user with a registration in non-attending status")
	attMock.SetupRegistered("202", 43, attendeeservice.StatusWaiting)
	token := tstValidUserToken(t, 202)

	docs.When("When they attempt to access the group information")
	response := tstPerformGet("/api/rest/v1/groups/"+id1, token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "attendee.status.not.attending", "registration is not in attending status")
}

func TestGroupsGet_InvalidID(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with a registration in attending status")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusApproved)
	token := tstValidUserToken(t, 101)

	docs.When("When they attempt to access group information, but supply an invalid id")
	response := tstPerformGet("/api/rest/v1/groups/kittycats", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "group.id.invalid", "you must specify a valid uuid")
}
