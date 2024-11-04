package acceptance

import (
	"github.com/eurofurence/reg-room-service/docs"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"net/http"
	"testing"
)

func TestRoomsMy_UserSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a registered attendee with an active registration who is in a finalized room")
	location1 := setupExistingRoom(t, "rodents", true, squirrel)

	docs.When("When the user requests their room")
	token := tstValidUserToken(t, subjectUint(squirrel))
	response := tstPerformGet("/api/rest/v1/rooms/my", token)

	docs.Then("Then the request is successful and the response is as expected")
	actual := modelsv1.Room{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.Room{
		ID:        tstRoomLocationToRoomID(location1),
		Name:      "rodents",
		Flags:     []string{"final"},
		Comments:  p("A nice comment for rodents"),
		Size:      2,
		Occupants: []modelsv1.Member{squirrel},
	}
	tstEqualResponseBodies(t, expected, actual)
}

// TODO admin success

// TODO api token produces correct error

func TestRoomsMy_AnonymousDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an unauthenticated user")
	token := tstNoToken()

	docs.When("When they request their room")
	response := tstPerformGet("/api/rest/v1/rooms/my", token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")
}

func TestRoomsMy_UserNoReg(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with NO registration")
	token := tstValidUserToken(t, 101)

	docs.When("When they request their room")
	response := tstPerformGet("/api/rest/v1/rooms/my", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "attendee.notfound", "you do not have a valid registration")
}

func TestRoomsMy_UserNonAttendingReg(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with a registration in non-attending status")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusNew, "Squirrel", "squirrel@example.com")
	token := tstValidUserToken(t, 101)

	docs.When("When they request their room")
	response := tstPerformGet("/api/rest/v1/rooms/my", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "attendee.status.not.attending", "registration is not in attending status")
}

func TestRoomsMy_UserNoRoom(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with a registration in attending status")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusPartiallyPaid, "Squirrel", "squirrel@example.com")
	token := tstValidUserToken(t, 101)

	docs.Given("Given they are not in any room")

	docs.When("When they request their room")
	response := tstPerformGet("/api/rest/v1/rooms/my", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "room.occupant.notfound", "not in a room, or final flag not set on room")
}

// TODO not finalized = 404
