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
		tstGroupMailToMember("group-application-accepted", "kittens", "1234567890", response.location))
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
	tstGroupUnchanged(t, id1, groupLocation, nil, nil)

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
	tstGroupUnchanged(t, id1, groupLocation, nil, nil)

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
	tstGroupUnchanged(t, id1, groupLocation, nil, nil)

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
	tstGroupUnchanged(t, id1, groupLocation, nil, nil)

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
	tstGroupUnchanged(t, id1, groupLocation, nil, nil)

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
	tstGroupUnchanged(t, id1, groupLocation, []modelsv1.Member{
		{
			ID:       43,
			Nickname: "Snep",
		},
	}, nil)

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
	tstGroupUnchanged(t, id1, groupLocation, nil, []modelsv1.Member{
		{
			ID:       84,
			Nickname: "Panther",
		},
	})

	docs.Then("And no emails have been sent")
	tstRequireMailRequests(t)
}

// TODO attendee goes first error cases

// TODO Bans

// TODO syntax failures

// --- helpers ---

func tstGroupUnchanged(t *testing.T, id string, location string, addMembers []modelsv1.Member, addInvites []modelsv1.Member) {
	t.Helper()

	response := tstPerformGet(location, tstValidAdminToken(t))
	actual := modelsv1.Group{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.Group{
		ID:          id,
		Name:        "kittens",
		Flags:       []string{},
		Comments:    p("A nice comment for kittens"),
		MaximumSize: 6,
		Owner:       42,
		Members: []modelsv1.Member{
			{
				ID:       42,
				Nickname: "Squirrel",
			},
		},
		Invites: nil,
	}
	expected.Members = append(expected.Members, addMembers...)
	expected.Invites = addInvites
	tstEqualResponseBodies(t, expected, actual)
}
