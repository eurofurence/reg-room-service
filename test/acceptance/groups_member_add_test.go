package acceptance

import (
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
	// TODO - set up mock for badge number 42 and status approved
	token := tstValidUserToken(t, 101)
	groupSent := modelsv1.GroupCreate{
		Name:     "kittens",
		Flags:    []string{"public"},
		Comments: p("A nice comment"),
		Owner:    0, // myself
	}
	group := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), token)

	docs.Given("Given another attendee with an active registration who is not in any group")
	// TODO - set up mock for badge number 84 and status approved

	docs.When("When the group owner requests an invite for the attendee")
	response := tstPerformPostNoBody(group.location+"/members/84", token)

	docs.Then("Then an invitation is successfully created")
	require.Equal(t, http.StatusNoContent, response.status, "unexpected http response status")
	// TODO validate body

	docs.Then("And its code can be used by the invited attendee")
	// TODO
}
