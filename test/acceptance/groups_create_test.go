package acceptance

import (
	"github.com/eurofurence/reg-room-service/docs"
	v1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const validGroupLocationRegex = "^\\/api\\/rest\\/v1\\/groups\\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"

func TestGroupsCreate_Success(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with an active registration who is not in any group")
	// TODO - set up mock for badge number 42 and status approved
	token := tstValidUserToken(t, 101)

	docs.When("When they create a room group with valid data")
	groupSent := v1.GroupCreate{
		Name:     "kittens",
		Flags:    []string{"public"},
		Comments: p("A nice comment"),
		Owner:    0, // myself
	}
	response := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), token)

	docs.Then("Then the group is successfully created")
	require.Equal(t, http.StatusCreated, response.status, "unexpected http response status")
	require.Regexp(t, validGroupLocationRegex, response.location, "invalid location header in response")

	docs.Then("And it can be read again by an admin")
	groupReadAgain := tstReadGroup(t, response.location)
	require.Equal(t, groupSent.Name, groupReadAgain.Name)

	docs.Then("And it contains exactly the user as owner and no invites")
	require.Equal(t, groupReadAgain.Owner, int32(42))
	require.Equal(t, len(groupReadAgain.Members), 1)
	require.Equal(t, groupReadAgain.Members[0].ID, int32(42))
	require.Equal(t, len(groupReadAgain.Invites), 0)
}
