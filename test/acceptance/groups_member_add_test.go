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
	location := path.Join("/api/rest/v1/groups/", id1)
	token := tstValidUserToken(t, 101)

	docs.Given("Given another attendee with an active registration who is not in any group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved, "Panther", "panther@example.com")

	docs.When("When the group owner requests an invite for the attendee")
	response := tstPerformPostNoBody(location+"/members/84", token)

	docs.Then("Then an invitation is successfully created")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected http response status")
	// TODO validate body

	docs.Then("And its code can be used by the invited attendee")
	// TODO
}
