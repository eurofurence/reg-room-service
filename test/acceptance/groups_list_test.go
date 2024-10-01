package acceptance

import (
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"net/http"
	"net/url"
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
	grp1 := modelsv1.Group{
		ID:          id1,
		Name:        "kittens",
		Flags:       []string{"public"},
		Comments:    p("A nice comment for kittens"),
		MaximumSize: 6,
		Owner:       42,
		Members: []modelsv1.Member{
			{
				ID:       42,
				Nickname: "",
			},
		},
		Invites: nil,
	}
	grp2 := modelsv1.Group{
		ID:          id2,
		Name:        "puppies",
		Flags:       []string{},
		Comments:    p("A nice comment for puppies"),
		MaximumSize: 6,
		Owner:       43,
		Members: []modelsv1.Member{
			{
				ID:       43,
				Nickname: "",
			},
		},
		Invites: nil,
	}
	expected := modelsv1.GroupList{}
	if id1 < id2 {
		expected.Groups = append(expected.Groups, &grp1, &grp2)
	} else {
		expected.Groups = append(expected.Groups, &grp2, &grp1)
	}
	tstEqualResponseBodies(t, expected, actual)
}

func TestGroupsList_AdminSuccess_Filtered(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given two registered attendees with an active registration who are in a group each")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	_ = setupExistingGroup(t, "puppies", false, "202")

	docs.When("When an admin requests to list groups containing a certain attendee")
	token := tstValidAdminToken(t)
	response := tstPerformGet("/api/rest/v1/groups?member_ids=42", token)

	docs.Then("Then the request is successful and the response includes the requested group information")
	actual := modelsv1.GroupList{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.GroupList{
		Groups: []*modelsv1.Group{
			{
				ID:          id1,
				Name:        "kittens",
				Flags:       []string{"public"},
				Comments:    p("A nice comment for kittens"),
				MaximumSize: 6,
				Owner:       42,
				Members: []modelsv1.Member{
					{
						ID:       42,
						Nickname: "",
					},
				},
				Invites: nil,
			},
		},
	}
	tstEqualResponseBodies(t, expected, actual)
}

func TestGroupsList_UserSuccess_Public(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a public group and a registered attendee who is not in the group")
	id1 := setupExistingGroup(t, "kittens", true, "101")
	_ = registerSubject("202")

	docs.When("When the attendee not in the group requests to list groups")
	token := tstValidUserToken(t, 202)
	response := tstPerformGet("/api/rest/v1/groups", token)

	docs.Then("Then the request is successful and the response includes only public information about public groups")
	actual := modelsv1.GroupList{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.GroupList{
		Groups: []*modelsv1.Group{
			// Owner and badge numbers are omitted!
			{
				ID:          id1,
				Name:        "kittens",
				Flags:       []string{"public"},
				MaximumSize: 6,
				Members: []modelsv1.Member{
					{
						Nickname: "",
					},
				},
			},
		},
	}
	tstEqualResponseBodies(t, expected, actual)
}

func TestGroupsList_UserSuccess_NonPublic(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given a private group and a registered attendee who is not in the group")
	_ = setupExistingGroup(t, "puppies", false, "101")
	_ = registerSubject("202")

	docs.When("When the attendee not in the group requests to list groups")
	token := tstValidUserToken(t, 202)
	response := tstPerformGet("/api/rest/v1/groups", token)

	docs.Then("Then the request is successful but the response does not include the group")
	actual := modelsv1.GroupList{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.GroupList{
		Groups: []*modelsv1.Group{},
	}
	tstEqualResponseBodies(t, expected, actual)
}

func TestGroupsList_AnonymousDeny(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an unauthenticated user")
	token := tstNoToken()

	docs.When("When they attempt to list groups")
	response := tstPerformGet("/api/rest/v1/groups", token)

	docs.Then("Then the request is denied")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")
}

func TestGroupsList_UserNoReg(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with NO registration")
	token := tstValidUserToken(t, 101)

	docs.When("When they try to list groups")
	response := tstPerformGet("/api/rest/v1/groups", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "attendee.notfound", "you do not have a valid registration")
}

func TestGroupsList_UserNonAttendingReg(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with a registration in non-attending status")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusNew)
	token := tstValidUserToken(t, 101)

	docs.When("When they try to list groups")
	response := tstPerformGet("/api/rest/v1/groups", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusForbidden, "attendee.status.not.attending", "registration is not in attending status")
}

func TestGroupsCreate_InvalidQueryParams(t *testing.T) {
	tstSetup(tstDefaultConfigFileRoomGroups)
	defer tstShutdown()

	docs.Given("Given an authorized user with a registration in attending status")
	attMock.SetupRegistered("101", 42, attendeeservice.StatusApproved)
	token := tstValidUserToken(t, 101)

	docs.When("When they try to list groups, but supply invalid parameters")
	response := tstPerformGet("/api/rest/v1/groups?member_ids=kittycat,-999", token)

	docs.Then("Then the request fails with the expected error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "request.parse.failed", url.Values{"details": []string{"member ids must be numeric and valid. Invalid member id: kittycat"}})
}
