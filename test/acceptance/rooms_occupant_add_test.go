package acceptance

import (
	"fmt"
	"github.com/eurofurence/reg-room-service/docs"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRoomsAddOccupant_NotLoggedIn(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with free beds")
	location := setupExistingRoom(t, "31415")

	docs.Given("Given an attendee with an active registration who is not in any room")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.Given("Given an anonymous user")
	token := tstNoToken()

	docs.When("When they try to add the attendee to the room")
	response := tstPerformPostNoBody(location+"/occupants/84", token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location)
}

func TestRoomsAddOccupant_UserDenyOtherAdd(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with free beds")
	location := setupExistingRoom(t, "31415")

	docs.Given("Given an attendee with an active registration who is not in any room")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.Given("Given another user, who is not an admin")
	token := tstValidUserToken(t, 101)

	docs.When("When they try to add the other attendee to the room")
	response := tstPerformPostNoBody(location+"/occupants/84", token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "you are not authorized for this operation - the attempt has been logged")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location)
}

func TestRoomsAddOccupant_UserDenySelfAdd(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with free beds")
	location := setupExistingRoom(t, "31415")

	docs.Given("Given an attendee with an active registration who is not in any room")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.When("When they try to add themselves to the room")
	token := tstValidUserToken(t, 84)
	response := tstPerformPostNoBody(location+"/occupants/84", token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "you are not authorized for this operation - the attempt has been logged")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location)
}

func TestRoomsAddOccupant_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with free beds")
	location := setupExistingRoom(t, "31415")

	docs.Given("Given an attendee with an active registration who is not in any room")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When the admin adds the attendee to the room")
	response := tstPerformPostNoBody(location+"/occupants/84", token)

	docs.Then("Then the request is successful")
	require.Equal(t, http.StatusNoContent, response.status)

	docs.Then("And the attendee has been added to the room")
	tstRoomState(t, location, panther)
}

func TestRoomsAddOccupant_ApiTokenSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with free beds")
	location := setupExistingRoom(t, "31415", snep)

	docs.Given("Given an attendee with an active registration who is not in any room")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.Given("Given a downstream service using a valid api token")
	token := tstValidApiToken()

	docs.When("When it adds the attendee to the room")
	response := tstPerformPostNoBody(location+"/occupants/84", token)

	docs.Then("Then the request is successful")
	require.Equal(t, http.StatusNoContent, response.status)

	docs.Then("And the attendee has been added to the room")
	tstRoomState(t, location, snep, panther)
}

func TestRoomsAddOccupant_Full(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a full room with no free beds")
	location := setupExistingRoom(t, "31415", squirrel, snep)

	docs.Given("Given an attendee with an active registration who is not in any room")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they try to add the attendee to the room")
	response := tstPerformPostNoBody(location+"/occupants/84", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusConflict, "room.size.full", "this room is full")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location, squirrel, snep)
}

func TestRoomsAddOccupant_Duplicate(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with an occupied bed")
	location := setupExistingRoom(t, "31415", squirrel)
	occupantLoc := fmt.Sprintf("%s/occupants/%d", location, squirrel.ID)

	docs.Given("Given an attendee with an active registration who is already in the room")
	registerSubject(subject(squirrel))

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When the admin tries to add the attendee to the room again")
	response := tstPerformPostNoBody(occupantLoc, token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusConflict, "room.occupant.duplicate", "this attendee is already in this room")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location, squirrel)
}

func TestRoomsAddOccupant_InAnotherRoom(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with free beds")
	location := setupExistingRoom(t, "31415")
	occupantLoc := fmt.Sprintf("%s/occupants/%d", location, squirrel.ID)

	docs.Given("Given an attendee who is already in another room")
	_ = setupExistingRoom(t, "27182", squirrel)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When the admin tries to also add the attendee to the first room")
	response := tstPerformPostNoBody(occupantLoc, token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusConflict, "room.occupant.conflict", "this attendee is already in another room")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location)
}

func TestRoomsAddOccupant_RoomNotFound(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an attendee with an active registration who is not in any room")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they try to add the attendee to a room, but specify a room that does not exist")
	response := tstPerformPostNoBody("/api/rest/v1/rooms/7a8d1116-d656-44eb-89dd-51eefef8a83b/occupants/84", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "room.id.notfound", "room does not exist")
}

func TestRoomsAddOccupant_AttendeeNotFound(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with free beds")
	location := setupExistingRoom(t, "31415", snep)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they try to add an attendee to the room, but specify a badge number that does not exist")
	response := tstPerformPostNoBody(location+"/occupants/4711", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "attendee.notfound", "no such attendee")
}

func TestRoomsAddOccupant_AttendeeNotAttending(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with free beds")
	location := setupExistingRoom(t, "31415", snep)

	docs.Given("Given a cancelled attendee")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusCancelled, "Panther", "panther@example.com")

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they try to add the cancelled attendee to the room")
	response := tstPerformPostNoBody(location+"/occupants/84", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusConflict, "attendee.status.not.attending", "registration is not in attending status")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location, snep)
}

func TestRoomsAddOccupant_InvalidRoomID(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.Given("Given an attendee with an active registration who is not in any room")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.When("When they attempt to add the attendee to a room, but specify an invalid room id")
	response := tstPerformPostNoBody("/api/rest/v1/rooms/kittycats/occupants/84", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "room.id.invalid", "you must specify a valid uuid")
}

func TestRoomsAddOccupant_BadgeNumberInvalid(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with free beds")
	location := setupExistingRoom(t, "31415", snep)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they attempt to add an occupant to the room, but supply an invalid badge number")
	response := tstPerformPostNoBody(location+"/occupants/floof", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "request.parse.failed", "invalid badge number - must be positive integer")
}

func TestRoomsAddOccupant_BadgeNumberNegative(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room with free beds")
	location := setupExistingRoom(t, "31415", snep)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they attempt to add an occupant to the room, but supply a negative badge number")
	response := tstPerformPostNoBody(location+"/occupants/-144", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "room.data.invalid", "invalid badge number - must be positive integer")
}

// TODO technical errors (downstream failures etc.)
