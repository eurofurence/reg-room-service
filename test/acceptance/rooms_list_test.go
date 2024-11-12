package acceptance

import (
	"github.com/eurofurence/reg-room-service/docs"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"net/http"
	"net/url"
	"testing"
)

func TestRoomsList_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given two registered attendees with an active registration who are in a room each")
	location1 := setupExistingRoom(t, "rodents", false, squirrel)
	location2 := setupExistingRoom(t, "cats", false, snep)

	docs.When("When an admin requests to list all rooms")
	token := tstValidAdminToken(t)
	response := tstPerformGet("/api/rest/v1/rooms", token)

	docs.Then("Then the request is successful and the response includes all room information")
	actual := modelsv1.RoomList{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	rm2 := modelsv1.Room{
		ID:        tstRoomLocationToRoomID(location2),
		Name:      "cats",
		Flags:     []string{},
		Comments:  p("A nice comment for cats"),
		Size:      2,
		Occupants: []modelsv1.Member{snep},
	}
	rm1 := modelsv1.Room{
		ID:        tstRoomLocationToRoomID(location1),
		Name:      "rodents",
		Flags:     []string{},
		Comments:  p("A nice comment for rodents"),
		Size:      2,
		Occupants: []modelsv1.Member{squirrel},
	}
	expected := modelsv1.RoomList{}
	expected.Rooms = append(expected.Rooms, &rm2, &rm1) // sorted by name
	tstEqualResponseBodies(t, expected, actual)
}

func TestRoomsList_AdminSuccess_Filtered(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given two registered attendees with an active registration who are in a room each")
	location1 := setupExistingRoom(t, "rodents", true, squirrel)
	_ = setupExistingRoom(t, "cats", true, snep)

	docs.When("When an admin requests to list rooms containing a certain attendee")
	token := tstValidAdminToken(t)
	response := tstPerformGet("/api/rest/v1/rooms?occupant_ids=42", token)

	docs.Then("Then the request is successful and the response includes the requested room information")
	actual := modelsv1.RoomList{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	rm1 := modelsv1.Room{
		ID:        tstRoomLocationToRoomID(location1),
		Name:      "rodents",
		Flags:     []string{"final"},
		Comments:  p("A nice comment for rodents"),
		Size:      2,
		Occupants: []modelsv1.Member{squirrel},
	}
	expected := modelsv1.RoomList{Rooms: []*modelsv1.Room{&rm1}}
	tstEqualResponseBodies(t, expected, actual)
}

func TestRoomsList_ApiTokenSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given two registered attendees with an active registration who are in a room each")
	location1 := setupExistingRoom(t, "rodents", false, squirrel)
	location2 := setupExistingRoom(t, "cats", false, snep)

	docs.Given("Given a downstream service using a valid api token")
	token := tstValidApiToken()

	docs.When("When it requests to list all rooms")
	response := tstPerformGet("/api/rest/v1/rooms", token)

	docs.Then("Then the request is successful and the response includes all room information")
	actual := modelsv1.RoomList{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	rm1 := modelsv1.Room{
		ID:        tstRoomLocationToRoomID(location1),
		Name:      "rodents",
		Flags:     []string{},
		Comments:  p("A nice comment for rodents"),
		Size:      2,
		Occupants: []modelsv1.Member{squirrel},
	}
	rm2 := modelsv1.Room{
		ID:        tstRoomLocationToRoomID(location2),
		Name:      "cats",
		Flags:     []string{},
		Comments:  p("A nice comment for cats"),
		Size:      2,
		Occupants: []modelsv1.Member{snep},
	}
	expected := modelsv1.RoomList{}
	if rm1.ID < rm2.ID {
		expected.Rooms = append(expected.Rooms, &rm1, &rm2)
	} else {
		expected.Rooms = append(expected.Rooms, &rm2, &rm1)
	}
	tstEqualResponseBodies(t, expected, actual)
}

func TestRoomsList_AnonymousDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an unauthenticated user")
	token := tstNoToken()

	docs.When("When they attempt to list rooms")
	response := tstPerformGet("/api/rest/v1/rooms", token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")
}

func TestRoomsList_NonAdminDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a registered attendee with an active registration who is in a room")
	_ = setupExistingRoom(t, "rodents", false, squirrel)

	docs.When("When they attempt to list rooms using the find rooms endpoint")
	token := tstValidUserToken(t, subjectUint(squirrel))
	response := tstPerformGet("/api/rest/v1/rooms", token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "you are not authorized for this operation - the attempt has been logged")
}

func TestRoomsList_InvalidQueryParams(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they try to list rooms, but supply invalid parameters")
	response := tstPerformGet("/api/rest/v1/rooms?occupant_ids=kittycat,-999", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "request.parse.failed", url.Values{"details": []string{"member ids must be numeric and valid. Invalid member id: kittycat"}})
}
