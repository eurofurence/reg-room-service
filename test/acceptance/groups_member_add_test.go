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

func TestGroupsAddMember_OwnerFirstSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)
	token := tstValidUserToken(t, 101)

	docs.Given("Given another attendee with an active registration who is not in any group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.When("When the group owner requests an invite for the attendee, providing their nickname as proof")
	response := tstPerformPostNoBody(groupLocation+"/members/84?nickname=Panther", token)

	docs.Then("Then an invitation is successfully created")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected http response status")
	require.Contains(t, response.location, groupLocation+"/members/84?code=")

	docs.Then("And the expected mail has been sent to the invited attendee")
	tstRequireMailRequests(t,
		tstGroupMailToMember("group-invited", "kittens", "1234567890", response.location))

	docs.Then("And the personalized join link can be used by the invited attendee")
	joinResponse := tstPerformPostNoBody(response.location, tstValidUserToken(t, 1234567890))
	require.Equal(t, http.StatusNoContent, joinResponse.status, "unexpected http response status")

	docs.Then("And the expected mail is then sent to the owner if they join")
	tstRequireMailRequests(t,
		tstGroupMailToOwner("group-member-joined", "kittens", "101", "1234567890"))
}

func TestGroupsAddMember_AttendeeFirstSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a public group")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who is not in any group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")
	token := tstValidUserToken(t, 1234567890)

	docs.When("When they apply for the group")
	response := tstPerformPostNoBody(groupLocation+"/members/84", token)

	docs.Then("Then an invitation is successfully created")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected http response status")
	require.Contains(t, response.location, groupLocation+"/members/84")

	docs.Then("And the expected mail is sent to the owner to inform them about the application")
	tstRequireMailRequests(t,
		tstGroupMailToOwner("group-member-request", "kittens", "101", "1234567890"))

	docs.Then("And the owner can accept the application")
	acceptResponse := tstPerformPostNoBody(response.location, tstValidUserToken(t, 101))
	require.Equal(t, http.StatusNoContent, acceptResponse.status, "unexpected http response status")

	docs.Then("And the expected mail is then sent to the invited attendee")
	tstRequireMailRequests(t,
		tstGroupMailToMember("group-application-accepted", "kittens", "1234567890", ""))
}

func TestGroupsAddMember_AdminForceSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a private group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who is not in any group")
	attMock.SetupRegistered("202", 43, attendeeservice.StatusApproved, "Snep", "snep@example.com")

	docs.Given("Given an unrelated user with admin rights")
	token := tstValidAdminToken(t)

	docs.When("When the admin adds the attendee to the group, using force=true")
	response := tstPerformPostNoBody(groupLocation+"/members/43?force=true", token)

	docs.Then("Then the attendee is directly added")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected http response status")
	require.Contains(t, response.location, groupLocation+"/members/43")

	docs.Then("And the expected emails have been sent")
	tstRequireMailRequests(t,
		tstGroupMailToOwner("group-member-joined", "kittens", "101", "202"),
		// no mail to member
	)
}

func TestGroupsAddMember_AdminForceSuccessAfterInvite(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a private group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who is not in any group")
	attMock.SetupRegistered("202", 43, attendeeservice.StatusApproved, "Snep", "snep@example.com")

	docs.Given("Given the attendee has been invited to the group")
	inviteResponse := tstPerformPostNoBody(groupLocation+"/members/43?nickname=Snep", tstValidUserToken(t, 101))
	require.Equal(t, http.StatusNoContent, inviteResponse.status, "unexpected http response status")
	mailMock.Reset()

	docs.Given("Given an unrelated user with admin rights")
	token := tstValidAdminToken(t)

	docs.When("When the admin adds the attendee to the group, using force=true")
	response := tstPerformPostNoBody(groupLocation+"/members/43?force=true", token)

	docs.Then("Then the attendee is directly added")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected http response status")
	require.Contains(t, response.location, groupLocation+"/members/43")

	docs.Then("And the expected emails have been sent")
	tstRequireMailRequests(t,
		tstGroupMailToOwner("group-member-joined", "kittens", "101", "202"),
		// no mail to member
	)
}

func TestGroupsAddMember_NotLoggedIn(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who is not in any group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.When("When an anonymous user requests an invite for the attendee, even providing their nickname as proof")
	response := tstPerformPostNoBody(groupLocation+"/members/84?nickname=Panther", tstNoToken())

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")

	docs.Then("And the group is unchanged")
	tstGroupState(t, id1, groupLocation, nil, nil)

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

func TestGroupsAddMember_ThirdParty(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who is not in any group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.Given("Given a third, different user, who also has an active registration and is not in any group")
	attMock.SetupRegistered("202", 43, attendeeservice.StatusApproved, "Snep", "snep@example.com")
	token := tstValidUserToken(t, 202)

	docs.When("When the user, who is neither group owner, nor the attendee, tries to invite the attendee to the group")
	response := tstPerformPostNoBody(groupLocation+"/members/84?nickname=Panther", token)

	docs.Then("Then the request is denied with an appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "only the group owner or an admin can invite other people into a group")

	docs.Then("And the group is unchanged")
	tstGroupState(t, id1, groupLocation, nil, nil)

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

func TestGroupsAddMember_OwnerFirstNicknameMismatch(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)
	token := tstValidUserToken(t, 101)

	docs.Given("Given another attendee with an active registration who is not in any group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.When("When the group owner requests an invite for the attendee, but fails to provide their correct nickname")
	response := tstPerformPostNoBody(groupLocation+"/members/84?nickname=NoIdea", token)

	docs.Then("Then the request fails with an appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "group.invite.mismatch", "nickname did not match - you need to know the nickname to be able to invite this attendee")

	docs.Then("And the group is unchanged")
	tstGroupState(t, id1, groupLocation, nil, nil)

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

func TestGroupsAddMember_OwnerFirstAlreadyInADifferentGroup(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)
	token := tstValidUserToken(t, 101)

	docs.Given("Given another attendee with an active registration who is already in a group")
	_ = setupExistingGroup(t, "puppies", false, "202")

	docs.When("When the group owner requests an invite for the attendee")
	response := tstPerformPostNoBody(groupLocation+"/members/43?nickname=Snep", token)

	docs.Then("Then the request fails with an appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusConflict, "group.member.conflict", "this attendee is already invited to another group or in another group")

	docs.Then("And the group is unchanged")
	tstGroupState(t, id1, groupLocation, nil, nil)

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

func TestGroupsAddMember_OwnerFirstAlreadyInvitedElsewhere(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who has been invited into another group")
	id2 := setupExistingGroup(t, "puppies", false, "202")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")
	groupLocation2 := path.Join("/api/rest/v1/groups/", id2)
	inviteResponse := tstPerformPostNoBody(groupLocation2+"/members/84?nickname=Panther", tstValidUserToken(t, 202))
	require.Equal(t, http.StatusNoContent, inviteResponse.status, "setup of invite failed")
	mailMock.Reset()

	docs.When("When the group owner requests an invite for the attendee")
	token := tstValidUserToken(t, 101)
	response := tstPerformPostNoBody(groupLocation+"/members/84?nickname=Panther", token)

	docs.Then("Then the request fails with an appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusConflict, "group.member.conflict", "this attendee is already invited to another group or in another group")

	docs.Then("And the group is unchanged")
	tstGroupState(t, id1, groupLocation, nil, nil)

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

func TestGroupsAddMember_OwnerFirstAlreadyInTheGroup(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is owner of a group")
	docs.Given("Given another attendee with an active registration who is already in the group")
	id1 := setupExistingGroup(t, "kittens", false, "101", "202")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.When("When the group owner requests an invite for the attendee")
	token := tstValidUserToken(t, 101)
	response := tstPerformPostNoBody(groupLocation+"/members/43?nickname=Snep", token)

	docs.Then("Then the request fails with an appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusConflict, "group.member.duplicate", "this attendee is already a member of this group")

	docs.Then("And the group is unchanged")
	tstGroupState(t, id1, groupLocation, []modelsv1.Member{snep}, nil)

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

func TestGroupsAddMember_OwnerFirstConfirmThirdParty(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who has been invited into the group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")
	inviteResponse := tstPerformPostNoBody(groupLocation+"/members/84?nickname=Panther", tstValidUserToken(t, 101))
	require.Equal(t, http.StatusNoContent, inviteResponse.status, "setup invite failed")
	mailMock.Reset()

	docs.Given("Given a third, different user, who also has an active registration and is not in any group")
	attMock.SetupRegistered("202", 43, attendeeservice.StatusApproved, "Snep", "snep@example.com")
	token := tstValidUserToken(t, 202)

	docs.When("When the user, who is neither group owner, nor the attendee, tries to confirm the invitation, even knowing the code")
	response := tstPerformPostNoBody(inviteResponse.location, token)

	docs.Then("Then the request is denied with an appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "only the group owner or an admin can accept invitations from others into a group")

	docs.Then("And the group is unchanged")
	tstGroupState(t, id1, groupLocation, nil, []modelsv1.Member{panther})

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

func TestGroupsAddMember_OwnerFirstConfirmCodeMismatch(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who has been invited into the group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")
	inviteResponse := tstPerformPostNoBody(groupLocation+"/members/84?nickname=Panther", tstValidUserToken(t, 101))
	require.Equal(t, http.StatusNoContent, inviteResponse.status, "setup invite failed")
	mailMock.Reset()

	docs.When("When they try to accept the invitation, but supply a wrong invitation code")
	response := tstPerformPostNoBody(inviteResponse.location+"some_added_nonsense", tstValidUserToken(t, 1234567890))

	docs.Then("Then the request is denied with an appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "you must provide the invitation code you were sent in order to join")

	docs.Then("And the group is unchanged")
	tstGroupState(t, id1, groupLocation, nil, []modelsv1.Member{panther})

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

func TestGroupsAddMember_AttendeeFirstOtherOwner(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a public group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who is owner of another group")
	_ = setupExistingGroup(t, "puppies", false, "202")
	token := tstValidUserToken(t, 202)

	docs.When("When they try to apply for the group")
	response := tstPerformPostNoBody(groupLocation+"/members/43", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusConflict, "group.member.conflict", "this attendee is already invited to another group or in another group")

	docs.Then("And the group is unchanged")
	tstGroupState(t, id1, groupLocation, nil, nil)

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

func TestGroupsAddMember_AttendeeFirstOtherInvite(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a public group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation1 := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who is owner of another group")
	id2 := setupExistingGroup(t, "puppies", false, "202")
	groupLocation2 := path.Join("/api/rest/v1/groups/", id2)

	docs.Given("Given a third attendee with an active registration who has been invited to the second group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")
	inviteResponse := tstPerformPostNoBody(groupLocation2+"/members/84?nickname=Panther", tstValidUserToken(t, 202))
	require.Equal(t, http.StatusNoContent, inviteResponse.status, "setup invite failed")
	mailMock.Reset()

	docs.When("When they try to apply for the first group")
	token := tstValidUserToken(t, 1234567890)
	response := tstPerformPostNoBody(groupLocation1+"/members/84", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusConflict, "group.member.conflict", "this attendee is already invited to another group or in another group")

	docs.Then("And the first group is unchanged")
	tstGroupState(t, id1, groupLocation1, nil, nil)

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

// ban handling

func TestGroupsAddMember_AdminForceAddRemovesBan(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who has been banned from this group")
	memberLocation := tstSetupBan(t, id1, 202)

	docs.When("When an admin adds them to the group with the force parameter")
	response := tstPerformPostNoBody(memberLocation+"?force=true", tstValidAdminToken(t))

	docs.Then("Then the request is successful")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected status code")

	docs.Then("And the attendee has been added to the group")
	tstGroupState(t, id1, groupLocation, []modelsv1.Member{snep}, nil)

	docs.Then("And the ban has been removed")
	tstRequireBanned(t, id1, 43, false)
}

func TestGroupsAddMember_BannedUserReject(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who has been banned from this group")
	memberLocation := tstSetupBan(t, id1, 202)
	token := tstValidUserToken(t, 202)

	docs.When("When that user attempts to apply to the group again")
	response := tstPerformPostNoBody(memberLocation, token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "auth.forbidden", "you cannot join this group - please stop trying")

	docs.Then("And the attendee has not been added to the group")
	tstGroupState(t, id1, groupLocation, nil, nil)

	docs.Then("And the ban is still in place")
	tstRequireBanned(t, id1, 43, true)

	docs.Then("And no emails have been sent that would annoy the owner")
	tstRequireMailRequests(t)
}

func TestGroupsAddMember_OwnerInviteRemovesBan(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who has been banned from this group")
	memberLocation := tstSetupBan(t, id1, 202)

	docs.When("When the group owner invites them to the group")
	response := tstPerformPostNoBody(memberLocation+"?nickname=Snep", tstValidUserToken(t, 101))

	docs.Then("Then the request is successful")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected status code")

	docs.Then("And the attendee has been invited to the group")
	tstGroupState(t, id1, groupLocation, nil, []modelsv1.Member{snep})

	docs.Then("And the ban has been removed")
	tstRequireBanned(t, id1, 43, false)
}

func TestGroupsAddMember_GroupNotFound(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an attendee with an active registration who is not in any group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.When("When they attempt to apply to a group, but supply a group id that does not exist")
	token := tstValidUserToken(t, 1234567890)
	response := tstPerformPostNoBody("/api/rest/v1/groups/7a8d1116-d656-44eb-89dd-51eefef8a83b/members/84", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "group.id.notfound", "this group does not exist")
}

func TestGroupsAddMember_BadgeNumberInvalid(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who is not in any group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.When("When they attempt to apply to the group, but supply an invalid badge number")
	token := tstValidUserToken(t, 1234567890)
	response := tstPerformPostNoBody(groupLocation+"/members/floof", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "request.parse.failed", "invalid badge number - must be positive integer")
}

func TestGroupsAddMember_BadgeNumberNegative(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a user with an active registration who is owner of a group")
	id1 := setupExistingGroup(t, "kittens", false, "101")
	groupLocation := path.Join("/api/rest/v1/groups/", id1)

	docs.Given("Given another attendee with an active registration who is not in any group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.When("When they attempt to leave the group, but supply a negative badge number")
	token := tstValidUserToken(t, 1234567890)
	response := tstPerformPostNoBody(groupLocation+"/members/-144", token)

	docs.Then("Then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "group.data.invalid", "invalid badge number - must be positive integer")
}

// TODO technical errors (downstream failures etc.)
