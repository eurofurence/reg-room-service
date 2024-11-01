package acceptance

import (
	"fmt"
	"github.com/eurofurence/reg-room-service/docs"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRoomsRemoveOccupant_NotLoggedIn(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with an occupied bed")
	location := setupExistingRoom(t, "31415", squirrel)
	occupantLoc := fmt.Sprintf("%s/occupants/%d", location, squirrel.ID)

	docs.Given("Given an anonymous user")
	token := tstNoToken()

	docs.When("When they try to remove the attendee from the room")
	response := tstPerformDelete(occupantLoc, token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location, squirrel)
}

func TestRoomsRemoveOccupant_UserDenyOtherRemove(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with an occupied bed")
	location := setupExistingRoom(t, "31415", snep)
	occupantLoc := fmt.Sprintf("%s/occupants/%d", location, snep.ID)

	docs.Given("Given another user, who is not an admin")
	token := tstValidUserToken(t, 101)

	docs.When("When they try to remove the other attendee from the room")
	response := tstPerformDelete(occupantLoc, token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "you are not authorized for this operation - the attempt has been logged")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location, snep)
}

func TestRoomsRemoveOccupant_UserDenySelfRemove(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with an occupied bed")
	location := setupExistingRoom(t, "31415", squirrel)
	occupantLoc := fmt.Sprintf("%s/occupants/%d", location, squirrel.ID)

	docs.Given("Given the attendee occupying the bed")
	token := tstValidUserToken(t, 101)

	docs.When("When they try to remove themselves from the room")
	response := tstPerformDelete(occupantLoc, token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "you are not authorized for this operation - the attempt has been logged")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location, squirrel)
}

func TestRoomsRemoveOccupant_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with occupied beds")
	location := setupExistingRoom(t, "31415", squirrel, snep)
	squirrelLoc := fmt.Sprintf("%s/occupants/%d", location, squirrel.ID)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When the admin removes the attendee from the room")
	response := tstPerformDelete(squirrelLoc, token)

	docs.Then("Then the request is successful")
	require.Equal(t, http.StatusNoContent, response.status)

	docs.Then("And the attendee has been removed from the room")
	tstRoomState(t, location, snep)
}

func TestRoomsRemoveOccupant_ApiTokenSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with occupied beds")
	location := setupExistingRoom(t, "31415", squirrel, snep)
	snepLoc := fmt.Sprintf("%s/occupants/%d", location, snep.ID)

	docs.Given("Given a downstream service using a valid api token")
	token := tstValidApiToken()

	docs.When("When it removes an attendee from the room")
	response := tstPerformDelete(snepLoc, token)

	docs.Then("Then the request is successful")
	require.Equal(t, http.StatusNoContent, response.status)

	docs.Then("And the attendee has been removed from the room")
	tstRoomState(t, location, squirrel)
}

func TestRoomsRemoveOccupant_NotInRoom(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with occupied beds")
	location := setupExistingRoom(t, "31415", squirrel)

	docs.Given("Given an attendee with an active registration who is not in any room")
	registerSubject(subject(snep))
	snepLoc := fmt.Sprintf("%s/occupants/%d", location, snep.ID) // not actually in the room

	docs.When("When an admin tries to remove the attendee from the room (which they are not actually in)")
	token := tstValidAdminToken(t)
	response := tstPerformDelete(snepLoc, token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "room.occupant.notfound", "this attendee is not in any room")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location, squirrel)
}

func TestRoomsRemoveOccupant_InAnotherRoom(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with occupied beds")
	location := setupExistingRoom(t, "31415", squirrel)

	docs.Given("Given an attendee with an active registration who is in another room")
	_ = setupExistingRoom(t, "27182", snep)

	docs.When("When an admin tries to remove the attendee from the room (which they are not actually in)")
	token := tstValidAdminToken(t)
	wrongSnepLoc := fmt.Sprintf("%s/occupants/%d", location, snep.ID) // not actually in this room
	response := tstPerformDelete(wrongSnepLoc, token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusConflict, "room.occupant.conflict", "this attendee is in a different room")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location, squirrel)
}

func TestRoomsRemoveOccupant_RoomNotFound(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they try to remove an attendee from a room, but specify a room that does not exist")
	response := tstPerformDelete("/api/rest/v1/rooms/7a8d1116-d656-44eb-89dd-51eefef8a83b/occupants/84", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "room.id.notfound", "room does not exist")
}

func TestRoomsRemoveOccupant_AttendeeNotFound(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with occupied beds")
	location := setupExistingRoom(t, "31415", snep)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they try to remove an attendee to the room, but specify a badge number that does not exist")
	response := tstPerformDelete(location+"/occupants/4711", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "attendee.notfound", "no such attendee")
}

func TestRoomsRemoveOccupant_CancelledRemoveSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with a bed occupied by a cancelled attendee")
	location := setupExistingRoom(t, "31415", snep)
	snepLoc := fmt.Sprintf("%s/occupants/%d", location, snep.ID)
	// now cancel snep - has to be done after adding to room
	attMock.SetupRegistered("202", 43, attendeeservice.StatusCancelled, "Snep", "snep@example.com")

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When the admin removes the attendee from the room")
	response := tstPerformDelete(snepLoc, token)

	docs.Then("Then the request is successful, even though their registration is not in attending status")
	require.Equal(t, http.StatusNoContent, response.status)

	docs.Then("And the attendee has been removed from the room")
	tstRoomState(t, location)
}

func TestRoomsRemoveOccupant_InvalidRoomID(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they attempt to remove an attendee from a room, but specify an invalid room id")
	response := tstPerformDelete("/api/rest/v1/rooms/kittycats/occupants/84", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "room.id.invalid", "you must specify a valid uuid")
}

func TestRoomsRemoveOccupant_BadgeNumberInvalid(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room")
	location := setupExistingRoom(t, "31415", snep)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they attempt to remove an occupant from the room, but supply an invalid badge number")
	response := tstPerformDelete(location+"/occupants/floof", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "request.parse.failed", "invalid badge number - must be positive integer")
}

func TestRoomsRemoveOccupant_BadgeNumberNotPositive(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with free beds")
	location := setupExistingRoom(t, "31415", snep)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they attempt to remove an occupant from the room, but supply a zero badge number")
	response := tstPerformDelete(location+"/occupants/0", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "room.data.invalid", "invalid badge number - must be positive integer")
}

// TODO technical errors (downstream failures etc.)
