package acceptance

import (
	"github.com/eurofurence/reg-room-service/docs"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"net/http"
	"testing"
)

func TestRoomsList_AdminSuccess(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given two registered attendees with an active registration who are in a room each")
	location1 := setupExistingRoom(t, "rodents", false, squirrel)
	location2 := setupExistingRoom(t, "cats", false, snep)

	docs.When("When an admin requests to list all rooms")
	token := tstValidAdminToken(t)
	response := tstPerformGet("/api/rest/v1/rooms", token)

	docs.Then("Then the request is successful and the response includes all room information")
	actual := modelsv1.RoomList{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	rm1 := modelsv1.Room{
		ID:        tstRoomLocationToRoomID(location1),
		Name:      "rodents",
		Flags:     []string{},
		Comments:  p("A nice comment for rodents"),
		Size:      2,
		Occupants: []modelsv1.Member{squirrel},
	}
	rm2 := modelsv1.Room{
		ID:        tstRoomLocationToRoomID(location2),
		Name:      "cats",
		Flags:     []string{},
		Comments:  p("A nice comment for cats"),
		Size:      2,
		Occupants: []modelsv1.Member{snep},
	}
	expected := modelsv1.RoomList{}
	if rm1.ID < rm2.ID {
		expected.Rooms = append(expected.Rooms, &rm1, &rm2)
	} else {
		expected.Rooms = append(expected.Rooms, &rm2, &rm1)
	}
	tstEqualResponseBodies(t, expected, actual)
}
