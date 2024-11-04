package acceptance

import (
	"github.com/eurofurence/reg-room-service/docs"
	v1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

const validRoomLocationRegex = "^\\/api\\/rest\\/v1\\/rooms\\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"

func TestRoomsCreate_NotLoggedIn(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an anonymous user")
	token := tstNoToken()

	docs.When("When they try to create a room")
	roomSent := v1.RoomCreate{
		Name:     "31415",
		Flags:    []string{},
		Comments: p("A pi comment"),
	}
	response := tstPerformPost("/api/rest/v1/rooms", tstRenderJson(roomSent), token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")
}

func TestRoomsCreate_UserDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user, who is not an admin")
	token := tstValidUserToken(t, 101)

	docs.When("When they try to create a room")
	roomSent := v1.RoomCreate{
		Name:     "31415",
		Flags:    []string{},
		Comments: p("A pi comment"),
	}
	response := tstPerformPost("/api/rest/v1/rooms", tstRenderJson(roomSent), token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "you are not authorized for this operation - the attempt has been logged")
}

func TestRoomsCreate_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they create a room with valid data")
	roomSent := v1.RoomCreate{
		Name:     "31415",
		Flags:    []string{},
		Comments: p("A pi comment"),
	}
	response := tstPerformPost("/api/rest/v1/rooms", tstRenderJson(roomSent), token)

	docs.Then("Then an empty room with the given name is successfully created")
	require.Equal(t, http.StatusCreated, response.status, "unexpected http response status")
	require.Regexp(t, validRoomLocationRegex, response.location, "invalid location header in response")
	roomReadAgain := tstReadRoom(t, response.location)
	require.Equal(t, roomSent.Name, roomReadAgain.Name)
	require.Equal(t, 0, len(roomReadAgain.Occupants))
}

func TestRoomsCreate_ApiTokenSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a downstream service using a valid api token")
	token := tstValidApiToken()

	docs.When("When they create a room with valid data")
	roomSent := v1.RoomCreate{
		Name:     "31415",
		Flags:    []string{},
		Comments: p("A pi comment"),
	}
	response := tstPerformPost("/api/rest/v1/rooms", tstRenderJson(roomSent), token)

	docs.Then("Then an empty room with the given name is successfully created")
	require.Equal(t, http.StatusCreated, response.status, "unexpected http response status")
	require.Regexp(t, validRoomLocationRegex, response.location, "invalid location header in response")
	roomReadAgain := tstReadRoom(t, response.location)
	require.Equal(t, roomSent.Name, roomReadAgain.Name)
	require.Equal(t, 0, len(roomReadAgain.Occupants))
}

func TestRoomsCreate_InvalidJSONSyntax(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they try to create a room, but supply syntactically invalid JSON")
	response := tstPerformPost("/api/rest/v1/rooms", `{"name":"invalid":"extra"`, token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "room.data.invalid", "invalid json provided")
}

func TestRoomsCreate_InvalidData(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a downstream service using a valid api token")
	token := tstValidApiToken()

	docs.When("When they try to create a room, but supply invalid information")
	roomSent := v1.RoomCreate{
		Name:  "",
		Flags: []string{"invalid"},
	}
	response := tstPerformPost("/api/rest/v1/rooms", tstRenderJson(roomSent), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "room.data.invalid", url.Values{"name": []string{"room name cannot be empty"}, "flags": []string{"no such flag 'invalid'"}})
}

func TestRoomsCreate_NameTooLong(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an admin")
	token := tstValidAdminToken(t)

	docs.When("When they try to create a room, but supply a name that is too long")
	roomSent := v1.RoomCreate{
		Name: "super long room name that is too long, longer than 50 characters, can you believe it?",
	}
	response := tstPerformPost("/api/rest/v1/rooms", tstRenderJson(roomSent), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "room.data.invalid", url.Values{"name": []string{"room name too long, max 50 characters"}})
}

// TODO duplicate name
