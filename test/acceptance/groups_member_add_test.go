package acceptance

import (
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eurofurence/reg-room-service/docs"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

func TestGroupsAddMember_OwnerFirstSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is owner of a group")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusApproved)
	token := tstValidUserToken(t, 101)
	groupSent := modelsv1.GroupCreate{
		Name:     "kittens",
		Flags:    []string{"public"},
		Comments: p("A nice comment"),
		Owner:    0, // myself
	}
	group := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), token)

	docs.Given("Given another attendee with an active registration who is not in any group")
	attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusApproved)

	docs.When("When the group owner requests an invite for the attendee")
	response := tstPerformPostNoBody(group.location+"/members/84", token)

	docs.Then("Then an invitation is successfully created")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected http response status")
	// TODO validate body

	docs.Then("And its code can be used by the invited attendee")
	// TODO
}
