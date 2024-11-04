package acceptance

import (
	"fmt"
	"github.com/eurofurence/reg-room-service/docs"
	"net/http"
	"testing"
)

func TestRoomsGet_NotLoggedIn(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room")
	location := setupExistingRoom(t, "31415", false, squirrel, snep)

	docs.Given("Given an anonymous user")
	token := tstNoToken()

	docs.When("When they try to access the room information")
	response := tstPerformGet(location, token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")
}

func TestRoomsGet_UserDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user, who is not an admin")
	token := tstValidUserToken(t, 101)

	docs.Given("Given a room they are in")
	location := setupExistingRoom(t, "31415", false, squirrel, snep)

	docs.When("When they try to access the room information, but do not use the special 'find my room' endpoint")
	response := tstPerformGet(location, token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "you are not authorized for this operation - the attempt has been logged")
}

func TestRoomsGet_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an empty room")
	location := setupExistingRoom(t, "31415", false)

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they access the room information")
	response := tstPerformGet(location, token)

	docs.Then("Then the request is successful and the response is as expected")
	tstRoomGetResponse(t, location, response)
}

func TestRoomsGet_ApiTokenSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a room")
	location := setupExistingRoom(t, "31415", false, squirrel, snep)

	docs.Given("Given a downstream service using a valid api token")
	token := tstValidApiToken()

	docs.When("When it accesses the room information")
	response := tstPerformGet(location, token)

	docs.Then("Then the request is successful and the response is as expected")
	tstRoomGetResponse(t, location, response, squirrel, snep)
}

func TestRoomsGet_InvalidID(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they attempt to access room information, but supply an invalid id")
	response := tstPerformGet("/api/rest/v1/rooms/kittycats", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "room.id.invalid", "you must specify a valid uuid")
}

func TestRoomsGet_NotFound(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	wrongId := "7ec0c20c-7dd4-491c-9b52-025be6950cdd"
	docs.When("When they attempt to access room information, but supply a valid id for a room that does not exist")
	response := tstPerformGet(fmt.Sprintf("/api/rest/v1/rooms/%s", wrongId), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "room.id.notfound", "room does not exist")
}
