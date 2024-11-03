package acceptance

import (
	"github.com/eurofurence/reg-room-service/docs"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"net/http"
	"testing"
)

func TestRoomsMy_UserSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a registered attendee with an active registration who is in a finalized room")
	location1 := setupExistingRoom(t, "rodents", true, squirrel)

	docs.When("When the user requests their room")
	token := tstValidUserToken(t, subjectUint(squirrel))
	response := tstPerformGet("/api/rest/v1/rooms/my", token)

	docs.Then("Then the request is successful and the response is as expected")
	actual := modelsv1.Room{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.Room{
		ID:        tstRoomLocationToRoomID(location1),
		Name:      "rodents",
		Flags:     []string{"final"},
		Comments:  p("A nice comment for rodents"),
		Size:      2,
		Occupants: []modelsv1.Member{squirrel},
	}
	tstEqualResponseBodies(t, expected, actual)
}
