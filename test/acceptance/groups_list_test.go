package acceptance

import (
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"net/http"
	"testing"

	"github.com/eurofurence/reg-room-service/docs"
)

// list groups

func TestGroupsList_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given two registered attendees with an active registration who are in a group each")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	id2 := setupExistingGroup(t, "puppies", false, "202")

	docs.When("When an admin requests to list all groups")
	token := tstValidAdminToken(t)
	response := tstPerformGet("/api/rest/v1/groups", token)

	docs.Then("Then the request is successful and the response includes all group information")
	actual := modelsv1.GroupList{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.GroupList{
		Groups: []*modelsv1.Group{
			{
				ID:          id1,
				Name:        "kittens",
				Flags:       []string{"public"},
				Comments:    p("A nice comment for kittens"),
				MaximumSize: p(int32(6)),
				Owner:       42,
				Members: []modelsv1.Member{
					{
						ID:       42,
						Nickname: "",
					},
				},
				Invites: nil,
			},
			{
				ID:          id2,
				Name:        "puppies",
				Flags:       []string{},
				Comments:    p("A nice comment for puppies"),
				MaximumSize: p(int32(6)),
				Owner:       43,
				Members: []modelsv1.Member{
					{
						ID:       43,
						Nickname: "",
					},
				},
				Invites: nil,
			},
		},
	}
	tstEqualResponseBodies(t, expected, actual)
}
