package acceptance

import (
	"fmt"
	"github.com/eurofurence/reg-room-service/docs"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRoomsUpdate_NotLoggedIn(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room")
	location := setupExistingRoom(t, "31415", false, squirrel)
	room := tstReadRoom(t, location)

	docs.Given("Given an anonymous user")
	token := tstNoToken()

	docs.When("When they try to update the room")
	room.Name = "27182"
	response := tstPerformPut(location, tstRenderJson(room), token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location, squirrel)
}

func TestRoomsUpdate_UserDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user, who is not an admin")
	token := tstValidUserToken(t, 101)

	docs.Given("Given a room")
	location := setupExistingRoom(t, "31415", false, squirrel, snep)
	room := tstReadRoom(t, location)

	docs.When("When they try to update the room")
	room.Name = "27182"
	response := tstPerformPut(location, tstRenderJson(room), token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "you are not authorized for this operation - the attempt has been logged")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location, squirrel, snep)
}

func TestRoomsUpdate_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room")
	location := setupExistingRoom(t, "27182", false, squirrel, snep)
	room := tstReadRoom(t, location)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they update the room")
	room.Name = "31415"                           // the name expected by tstRoomState
	room.Comments = p("A nice comment for 31415") // the comment expected by tstRoomState
	response := tstPerformPut(location, tstRenderJson(room), token)

	docs.Then("Then the request is successful")
	require.Equal(t, http.StatusNoContent, response.status)
	require.Equal(t, location, response.location)

	docs.Then("And the room has been updated")
	tstRoomState(t, location, squirrel, snep)
}

func TestRoomsUpdate_ApiTokenSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room")
	location := setupExistingRoom(t, "27182", false, squirrel, snep)
	room := tstReadRoom(t, location)

	docs.Given("Given a downstream service using a valid api token")
	token := tstValidApiToken()

	docs.When("When it updates the room")
	room.Name = "31415"                           // the name expected by tstRoomState
	room.Comments = p("A nice comment for 31415") // the comment expected by tstRoomState
	response := tstPerformPut(location, tstRenderJson(room), token)

	docs.Then("Then the request is successful")
	require.Equal(t, http.StatusNoContent, response.status)
	require.Equal(t, location, response.location)

	docs.Then("And the room has been updated")
	tstRoomState(t, location, squirrel, snep)
}

func TestRoomsUpdate_TooSmall(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room that is not empty")
	location := setupExistingRoom(t, "31415", false, squirrel, snep)
	room := tstReadRoom(t, location)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they try to change the room size to a value that is too small to hold the current occupants")
	room.Size = 1
	response := tstPerformPut(location, tstRenderJson(room), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusConflict, "room.size.too.small",
		"the room cannot be resized, too many occupants for new size")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location, squirrel, snep)
}

func TestRoomsUpdate_DuplicateName(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given two rooms")
	location1 := setupExistingRoom(t, "31415", false, squirrel)
	location2 := setupExistingRoom(t, "27182", false, snep)
	room1 := tstReadRoom(t, location1)
	room2 := tstReadRoom(t, location2)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they attempt to rename the first room, but use the name of the second room")
	room1.Name = room2.Name
	response := tstPerformPut(location1, tstRenderJson(room1), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusConflict, "room.data.duplicate", "another room with this name already exists")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location1, squirrel)
}

func TestRoomsUpdate_InvalidID(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room")
	location := setupExistingRoom(t, "31415", false)
	room := tstReadRoom(t, location)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they attempt to update a room, but supply an invalid id (not a uuid)")
	room.ID = "kittycats"
	room.Name = "27182"
	response := tstPerformPut("/api/rest/v1/rooms/kittycats", tstRenderJson(room), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "room.id.invalid", "you must specify a valid uuid")
}

func TestRoomsUpdate_InvalidBody(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room")
	location := setupExistingRoom(t, "31415", false)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they attempt to update the room, but supply syntactically invalid JSON in the body")
	response := tstPerformPut(location, `{"name":"invalid":"extra"`, token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "room.data.invalid", "invalid json provided")
}

func TestRoomsUpdate_NotFound(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	wrongId := "7ec0c20c-7dd4-491c-9b52-025be6950cdd"
	docs.When("When they attempt to update a room that does not exist")
	room := modelsv1.Room{
		ID:   wrongId,
		Name: "12345",
	}
	response := tstPerformPut(fmt.Sprintf("/api/rest/v1/rooms/%s", wrongId), tstRenderJson(room), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "room.id.notfound", "room does not exist")
}
