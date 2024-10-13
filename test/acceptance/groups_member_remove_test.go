package acceptance

import (
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eurofurence/reg-room-service/docs"
)

func TestGroupsRemoveMember_OwnerSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is owner of a group")
	docs.Given("Given another attendee who is a member of the group")
	id1 := setupExistingGroup(t, "kittens", false, "101", "202")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.When("When the group owner kicks them from the group")
	token := tstValidUserToken(t, 101)
	response := tstPerformDelete(groupLocation+"/members/43", token)

	docs.Then("Then the request is successful and they have been removed")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected http response status")
	tstGroupState(t, id1, groupLocation, nil, nil)

	docs.Then("And the expected mail has been sent to the removed attendee")
	tstRequireMailRequests(t,
		tstGroupMailToMember("group-member-kicked", "kittens", "202", ""))
}

func TestGroupsRemoveMember_AttendeeSelfSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a public group")
	docs.Given("Given another attendee who is a member of the group")
	id1 := setupExistingGroup(t, "kittens", false, "101", "202")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.When("When they leave the group")
	token := tstValidUserToken(t, 202)
	response := tstPerformDelete(groupLocation+"/members/43", token)

	docs.Then("Then the request is successful and they have been removed")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected http response status")
	tstGroupState(t, id1, groupLocation, nil, nil)

	docs.Then("And the expected mail is sent to the owner to inform them that someone has left their group")
	tstRequireMailRequests(t,
		tstGroupMailToOwner("group-member-left", "kittens", "101", "202"))
}

func TestGroupsRemoveMember_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	docs.Given("Given another attendee who is a member of the group")
	id1 := setupExistingGroup(t, "kittens", false, "101", "202")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given an unrelated user with admin rights")
	token := tstValidAdminToken(t)

	docs.When("When the admin removes the attendee from the group")
	response := tstPerformDelete(groupLocation+"/members/43", token)

	docs.Then("Then the request is successful and they have been removed")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected http response status")
	tstGroupState(t, id1, groupLocation, nil, nil)

	docs.Then("And the expected emails have been sent")
	tstRequireMailRequests(t,
		tstGroupMailToOwner("group-member-removed", "kittens", "101", "202"),
		tstGroupMailToMember("group-member-kicked", "kittens", "202", ""),
	)
}

// TODO invite leaves (success)

func TestGroupsRemoveMember_ThirdPartyDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	docs.Given("Given another attendee who is a member of the group")
	id1 := setupExistingGroup(t, "kittens", false, "101", "202")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given a third, different user, who also has an active registration and is not in any group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")
	token := tstValidUserToken(t, 1234567890)

	docs.When("When the user, who is neither group owner, nor the attendee, tries to remove the attendee from the group")
	response := tstPerformDelete(groupLocation+"/members/43", token)

	docs.Then("Then the request is denied with an appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "only the group owner or an admin can remove other people from a group")

	docs.Then("And the group is unchanged")
	tstGroupState(t, id1, groupLocation, []modelsv1.Member{
		{
			ID:       43,
			Nickname: "Snep",
		},
	}, nil)

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

func TestGroupsRemoveMember_AttendeeSelfWrongGroup(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a public group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who is owner of another group")
	_ = setupExistingGroup(t, "puppies", false, "202")

	docs.When("When they try to leave the first group, which they are not a member of")
	token := tstValidUserToken(t, 202)
	response := tstPerformDelete(groupLocation+"/members/43", token)

	docs.Then("Then the request is denied with an appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusConflict, "group.member.conflict", "this attendee is invited to a different group or in a different group")

	docs.Then("And the group is unchanged")
	tstGroupState(t, id1, groupLocation, nil, nil)

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

// TODO self member of other group

// TODO owner tries to leave

// TODO not in any group tries to leave

// TODO invited in other group tries to leave

func TestGroupsRemoveMember_NotLoggedIn(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	docs.Given("Given another attendee who is a member of the group")
	id1 := setupExistingGroup(t, "kittens", false, "101", "202")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.When("When an anonymous user tries to remove the group member")
	response := tstPerformDelete(groupLocation+"/members/43", tstNoToken())

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")

	docs.Then("And the group is unchanged")
	tstGroupState(t, id1, groupLocation, []modelsv1.Member{
		{
			ID:       43,
			Nickname: "Snep",
		},
	}, nil)

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

// TODO group not found, attendee not found

// TODO Bans - test after other member_remove tests done

func TestGroupsRemoveMember_BadgeNumberInvalid(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	docs.Given("Given another attendee who is a member of the group")
	id1 := setupExistingGroup(t, "kittens", false, "101", "202")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.When("When they attempt to leave the group, but supply an invalid badge number")
	token := tstValidUserToken(t, 202)
	response := tstPerformDelete(groupLocation+"/members/floof", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "request.parse.failed", "invalid badge number - must be positive integer")
}

func TestGroupsRemoveMember_BadgeNumberNegative(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	docs.Given("Given another attendee who is a member of the group")
	id1 := setupExistingGroup(t, "kittens", false, "101", "202")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.When("When they attempt to leave the group, but supply a negative badge number")
	token := tstValidUserToken(t, 202)
	response := tstPerformDelete(groupLocation+"/members/%2d144", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "request.parse.failed", "invalid badge number - must be positive integer")
}

// TODO more syntax failures

// TODO technical errors (downstream failures etc.)
