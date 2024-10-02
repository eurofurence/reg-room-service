package acceptance

import (
	"github.com/eurofurence/reg-room-service/docs"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"path"
	"testing"
)

func TestGroupsUpdate_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an attendee with an active registration who is in a group")
	id1 := setupExistingGroup(t, "kittens", true, "101")

	docs.Given("Given an authorized admin (a different user)")
	token := tstValidAdminToken(t)

	docs.When("When they retrieve the group and update its name")
	getGroup := tstReadGroup(t, path.Join("/api/rest/v1/groups/", id1))
	savedID := getGroup.ID
	getGroup.Name = "dogs"

	response := tstPerformPut(path.Join("/api/rest/v1/groups/", getGroup.ID), tstRenderJson(getGroup), token)

	docs.Then("Then the group is successfully updated")
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	require.Regexp(t, validGroupLocationRegex, response.location, "invalid location header in response")

	getGroup = tstReadGroup(t, response.location)
	require.Equal(t, savedID, getGroup.ID)
	require.Equal(t, "dogs", getGroup.Name)
	require.Len(t, getGroup.Flags, 1)
}

func TestGroupsUpdate_AnonymousDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an attendee with an active registration who is in a group")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	location := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given an unauthenticated user")
	token := tstNoToken()

	docs.When("When they try to update the group")
	getGroup := tstReadGroup(t, location)
	getGroup.Name = "dogs"
	response := tstPerformPut(path.Join("/api/rest/v1/groups/", getGroup.ID), tstRenderJson(getGroup), token)

	docs.Then("Then the request should be denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")

	docs.Then("And the group is unchanged")
	getGroupAgain := tstReadGroup(t, location)
	require.Equal(t, "kittens", getGroupAgain.Name)
}

func TestGroupsUpdate_UserDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an attendee with an active registration who is in a group, but NOT its owner")
	id1 := setupExistingGroup(t, "kittens", true, "101", "202")
	location := path.Join("/api/rest/v1/groups/", id1)
	token := tstValidUserToken(t, 202)

	docs.When("When they attempt to update the name of the group")
	getGroup := tstReadGroup(t, location)
	getGroup.Name = "dogs"

	response := tstPerformPut(path.Join("/api/rest/v1/groups/", getGroup.ID), tstRenderJson(getGroup), token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "only the group owner or an admin can change a group")

	docs.Then("And the group is unchanged")
	getGroupAgain := tstReadGroup(t, location)
	require.Equal(t, "kittens", getGroupAgain.Name)
}

func TestGroupsUpdate_UserOwnerSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is the owner of a group")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	token := tstValidUserToken(t, 101)

	docs.When("When they retrieve the group and update its name")
	getGroup := tstReadGroup(t, path.Join("/api/rest/v1/groups/", id1))
	savedID := getGroup.ID
	getGroup.Name = "dogs"

	response := tstPerformPut(path.Join("/api/rest/v1/groups/", getGroup.ID), tstRenderJson(getGroup), token)

	docs.Then("Then the group is successfully updated")
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	require.Regexp(t, validGroupLocationRegex, response.location, "invalid location header in response")

	getGroup = tstReadGroup(t, response.location)
	require.Equal(t, savedID, getGroup.ID)
	require.Equal(t, "dogs", getGroup.Name)
	require.Len(t, getGroup.Flags, 1)
}

func TestGroupsUpdate_InvalidJSONSyntax(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is the owner of a group")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	token := tstValidUserToken(t, 101)

	docs.When("When they try to update the room group, but supply syntactically invalid JSON")
	response := tstPerformPut(path.Join("/api/rest/v1/groups/", id1), `{"name":"invalid":"extra"`, token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "group.data.invalid", "invalid json provided")
}

func TestGroupsUpdate_InvalidData(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is the owner of a group")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	token := tstValidUserToken(t, 101)

	docs.When("When they try to update the group but supply invalid information")
	getGroup := tstReadGroup(t, path.Join("/api/rest/v1/groups/", id1))
	getGroup.Flags = []string{"invalid"}
	getGroup.Name = ""
	response := tstPerformPut(path.Join("/api/rest/v1/groups/", getGroup.ID), tstRenderJson(getGroup), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "group.data.invalid", url.Values{"name": []string{"group name cannot be empty"}, "flags": []string{"no such flag 'invalid'"}})
}

func TestGroupsUpdate_InvalidID(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	token := tstValidUserToken(t, 101)

	docs.When("When they try to update the group but supply invalid information")
	getGroup := tstReadGroup(t, path.Join("/api/rest/v1/groups/", id1))
	response := tstPerformPut("/api/rest/v1/groups/kittycats", tstRenderJson(getGroup), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "group.id.invalid", "'kittycats' is not a valid UUID")
}

func TestGroupsUpdate_NotFound(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	token := tstValidUserToken(t, 101)

	wrongId := "7ec0c20c-7dd4-491c-9b52-025be6950cdd"
	if wrongId == id1 {
		wrongId = "7ec0c20c-7dd4-491c-9b52-025be6950cef"
	}
	docs.When("When they try to update a group, but specify a valid id for a group that does not exist")
	getGroup := tstReadGroup(t, path.Join("/api/rest/v1/groups/", id1))
	response := tstPerformPut(path.Join("/api/rest/v1/groups/", wrongId), tstRenderJson(getGroup), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "group.id.notfound", "group does not exist")
}
