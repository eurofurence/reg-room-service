package acceptance

import (
	"fmt"
	"github.com/eurofurence/reg-room-service/docs"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRoomsDelete_NotLoggedIn(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an empty room")
	location := setupExistingRoom(t, "31415")

	docs.Given("Given an anonymous user")
	token := tstNoToken()

	docs.When("When they try to delete the room")
	response := tstPerformDelete(location, token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location)
}

func TestRoomsDelete_UserDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user, who is not an admin")
	token := tstValidUserToken(t, 101)

	docs.Given("Given an empty room")
	location := setupExistingRoom(t, "31415")

	docs.When("When they try to delete the room")
	response := tstPerformDelete(location, token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "you are not authorized for this operation - the attempt has been logged")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location)
}

func TestRoomsDelete_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an empty room")
	location := setupExistingRoom(t, "31415")

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they request deletion of the room")
	response := tstPerformDelete(location, token)

	docs.Then("Then the request is successful")
	require.Equal(t, http.StatusNoContent, response.status)

	docs.Then("And the room has been deleted")
	readAgainResponse := tstPerformGet(location, tstValidAdminToken(t))
	require.Equal(t, http.StatusNotFound, readAgainResponse.status)
}

func TestRoomsDelete_ApiTokenSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an empty room")
	location := setupExistingRoom(t, "31415")

	docs.Given("Given a downstream service using a valid api token")
	token := tstValidApiToken()

	docs.When("When it requests deletion of the room")
	response := tstPerformDelete(location, token)

	docs.Then("Then the request is successful")
	require.Equal(t, http.StatusNoContent, response.status)

	docs.Then("And the room has been deleted")
	readAgainResponse := tstPerformGet(location, tstValidAdminToken(t))
	require.Equal(t, http.StatusNotFound, readAgainResponse.status)
}

func TestRoomsDelete_NotEmpty(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room that is not empty")
	location := setupExistingRoom(t, "31415", squirrel)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they attempt to delete the room")
	response := tstPerformDelete(location, token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusConflict, "room.not.empty",
		"room is not empty and room deletion is a dangerous operation - please remove all occupants first "+
			"to ensure you really mean this (also prevents possible problems with concurrent updates)")

	docs.Then("And the room is unchanged")
	tstRoomState(t, location, squirrel)
}

func TestRoomsDelete_InvalidID(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they attempt to delete a room, but supply an invalid id")
	response := tstPerformDelete("/api/rest/v1/rooms/kittycats", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "room.id.invalid", "you must specify a valid uuid")
}

func TestRoomsDelete_NotFound(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	wrongId := "7ec0c20c-7dd4-491c-9b52-025be6950cdd"
	docs.When("When they attempt to delete a room, but supply a valid id for a room that does not exist")
	response := tstPerformDelete(fmt.Sprintf("/api/rest/v1/rooms/%s", wrongId), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "room.id.notfound", "room does not exist")
}
