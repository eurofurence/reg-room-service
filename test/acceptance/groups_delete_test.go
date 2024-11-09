package acceptance

import (
	"github.com/eurofurence/reg-room-service/docs"
	"github.com/stretchr/testify/require"
	"net/http"
	"path"
	"testing"
)

func TestGroupsDelete_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an attendee with an active registration who is in a group with two members")
	id1 := setupExistingGroup(t, "kittens", true, "101", "202")

	docs.Given("Given an authorized admin (a different user)")
	token := tstValidAdminToken(t)

	docs.When("When they delete the group")
	response := tstPerformDelete(path.Join("/api/rest/v1/groups/", id1), token)

	docs.Then("Then the group is successfully deleted and all its members are now no longer in a group")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected http response status")

	user101group := tstPerformGet("/api/rest/v1/groups/my", tstValidUserToken(t, 101))
	require.Equal(t, http.StatusNotFound, user101group.status, "unexpected http response status")
	user202group := tstPerformGet("/api/rest/v1/groups/my", tstValidUserToken(t, 202))
	require.Equal(t, http.StatusNotFound, user202group.status, "unexpected http response status")
}

func TestGroupsDelete_AnonymousDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an attendee with an active registration who is in a group")
	id1 := setupExistingGroup(t, "kittens", true, "101")

	docs.Given("Given an unauthenticated user")
	token := tstNoToken()

	docs.When("When they try to delete the group")
	response := tstPerformDelete(path.Join("/api/rest/v1/groups/", id1), token)

	docs.Then("Then the request is denied with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")
}

func TestGroupsDelete_UserNotOwnerDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an attendee with an active registration who is in a group, but NOT its owner")
	id1 := setupExistingGroup(t, "kittens", false, "101", "202")
	token := tstValidUserToken(t, 202)

	docs.When("When they attempt to delete the group")
	response := tstPerformDelete(path.Join("/api/rest/v1/groups/", id1), token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "only the group owner or an admin can delete a group")
}

func TestGroupsDelete_UserNotMemberDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an attendee with an active registration who is in a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	_ = setupExistingGroup(t, "puppies", false, "202")
	token := tstValidUserToken(t, 202)

	docs.When("When they attempt to delete a different group they are not a member of")
	response := tstPerformDelete(path.Join("/api/rest/v1/groups/", id1), token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "only the group owner or an admin can delete a group")
}

func TestGroupsDelete_UserOwnerSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is the owner of a group")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	token := tstValidUserToken(t, 101)

	docs.When("When they delete the group")
	response := tstPerformDelete(path.Join("/api/rest/v1/groups/", id1), token)

	docs.Then("Then the request is successful and the group has been deleted")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected http response status")

	deletedResponse := tstPerformGet(path.Join("/api/rest/v1/groups/", id1), tstValidAdminToken(t))
	require.Equal(t, http.StatusNotFound, deletedResponse.status, "group was not correctly deleted")
}

func TestGroupsDelete_InvalidID(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration")
	_ = setupExistingGroup(t, "kittens", true, "101")
	token := tstValidUserToken(t, 101)

	docs.When("When they try to delete a group, but specify an invalid id")
	response := tstPerformDelete("/api/rest/v1/groups/kittycats", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "group.id.invalid", "'kittycats' is not a valid UUID")
}

func TestGroupsDelete_NotFound(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration")
	registerSubject("101")
	token := tstValidUserToken(t, 101)

	wrongId := "7ec0c20c-7dd4-491c-9b52-025be6950cdd"
	docs.When("When they try to delete a group, but specify a valid id for a group that does not exist")
	response := tstPerformDelete(path.Join("/api/rest/v1/groups/", wrongId), token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "group.id.notfound", "group does not exist")
}
