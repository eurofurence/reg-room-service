package acceptance

import (
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
