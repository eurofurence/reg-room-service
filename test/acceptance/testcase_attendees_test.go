package acceptance

import (
	"context"
	"fmt"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/mailservice"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

// squirrel has attending status, subject "101"
var squirrel = modelsv1.Member{
	ID:       42,
	Nickname: "Squirrel",
}

// snep has attending status, subject "202"
var snep = modelsv1.Member{
	ID:       43,
	Nickname: "Snep", // popular abbreviation for Snow Leopard
}

// panther has non-attending status by default (but some test cases may set her up differently), subject "1234567890"
var panther = modelsv1.Member{
	ID:       84,
	Nickname: "Panther",
}

func tstInfosBySubject(subject string) (string, string, string) {
	switch subject {
	case "101":
		return "42", "Squirrel", "squirrel@example.com"
	case "202":
		return "43", "Snep", "snep@example.com"
	default:
		return "84", "Panther", "panther@example.com"
	}
}

// --- test helper functions for use with these ---

func subjectUint(member modelsv1.Member) uint {
	switch member.ID {
	case 42:
		return 101
	case 43:
		return 202
	default:
		return 1234567890
	}
}

func subject(member modelsv1.Member) string {
	return fmt.Sprintf("%d", subjectUint(member))
}

func setupExistingGroup(t *testing.T, name string, public bool, subject string, additionalMemberSubjects ...string) string {
	flags := []string{}
	if public {
		flags = append(flags, "public")
	}
	badgeNo := registerSubject(subject)

	groupSent := modelsv1.GroupCreate{
		Name:     name,
		Flags:    flags,
		Comments: p("A nice comment for " + name),
		Owner:    badgeNo,
	}
	response := tstPerformPost("/api/rest/v1/groups", tstRenderJson(groupSent), tstValidAdminToken(t))
	require.Equal(t, http.StatusCreated, response.status, "unexpected http response status")
	require.Regexp(t, validGroupLocationRegex, response.location, "invalid location header in response")

	for _, addSubject := range additionalMemberSubjects {
		addBadgeNo := registerSubject(addSubject)
		addResponse := tstPerformPostNoBody(fmt.Sprintf("%s/members/%d?force=true", response.location, addBadgeNo), tstValidAdminToken(t))
		require.Equal(t, http.StatusNoContent, addResponse.status, "unexpected http response status")
	}
	mailMock.Reset()

	locs := strings.Split(response.location, "/")
	return locs[len(locs)-1]
}

func setupExistingRoom(t *testing.T, name string, final bool, occupants ...modelsv1.Member) string {
	roomSent := modelsv1.RoomCreate{
		Name:     name,
		Flags:    []string{},
		Comments: p("A nice comment for " + name),
		Size:     2,
	}
	if final {
		roomSent.Flags = []string{"final"}
	}
	response := tstPerformPost("/api/rest/v1/rooms", tstRenderJson(roomSent), tstValidAdminToken(t))

	require.Equal(t, http.StatusCreated, response.status, "unexpected http response status")
	require.Regexp(t, validRoomLocationRegex, response.location, "invalid location header in response")

	for _, addMember := range occupants {
		addBadgeNo := registerSubject(subject(addMember))
		require.Equal(t, addBadgeNo, addMember.ID) // ensure test case setup correctly
		addResponse := tstPerformPostNoBody(fmt.Sprintf("%s/occupants/%d", response.location, addBadgeNo), tstValidAdminToken(t))
		require.Equal(t, http.StatusNoContent, addResponse.status, "unexpected http response status")
	}

	return response.location
}

func registerSubject(subject string) int64 {
	switch subject {
	case "101":
		attMock.SetupRegistered("101", 42, attendeeservice.StatusApproved, "Squirrel", "squirrel@example.com")
		return 42

	case "202":
		attMock.SetupRegistered("202", 43, attendeeservice.StatusPaid, "Snep", "snep@example.com")
		return 43

	default:
		attMock.SetupRegistered("1234567890", 84, attendeeservice.StatusCancelled, "Panther", "panther@example.com")
		return 84
	}
}

func tstRoomLocationToRoomID(location string) string {
	locs := strings.Split(location, "/")
	return locs[len(locs)-1]
}

func tstSetupBan(t *testing.T, groupId string, subject uint) string {
	t.Helper()

	badgeNo := registerSubject(fmt.Sprintf("%d", subject))
	memberLocation := fmt.Sprintf("/api/rest/v1/groups/%s/members/%d", groupId, badgeNo)

	applyResponse := tstPerformPostNoBody(memberLocation, tstValidUserToken(t, subject))
	require.Equal(t, http.StatusNoContent, applyResponse.status, "setup ban step 1 failed")

	banResponse := tstPerformDelete(memberLocation+"?autodeny=true", tstValidAdminToken(t))
	require.Equal(t, http.StatusNoContent, banResponse.status, "setup ban step 2 failed")

	mailMock.Reset()

	return memberLocation
}

func tstRequireBanned(t *testing.T, groupId string, badgeNo int64, expected bool) {
	t.Helper()

	actual, err := db.HasGroupBan(context.TODO(), groupId, badgeNo)
	require.Nil(t, err)
	require.Equal(t, expected, actual)
}

func tstGroupState(t *testing.T, id string, location string, addMembers []modelsv1.Member, addInvites []modelsv1.Member) {
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
		Members:     []modelsv1.Member{squirrel},
		Invites:     nil,
	}
	expected.Members = append(expected.Members, addMembers...)
	expected.Invites = addInvites
	tstEqualResponseBodies(t, expected, actual)
}

func tstRoomState(t *testing.T, location string, occupants ...modelsv1.Member) {
	t.Helper()

	response := tstPerformGet(location, tstValidAdminToken(t))
	tstRoomGetResponse(t, location, response, occupants...)
}

func tstRoomGetResponse(t *testing.T, location string, response tstWebResponse, occupants ...modelsv1.Member) {
	t.Helper()

	locs := strings.Split(location, "/")
	require.Equal(t, 6, len(locs), "location for a room should be /api/rest/v1/rooms/{uuid}")
	id := locs[len(locs)-1]

	actual := modelsv1.Room{}
	tstRequireSuccessResponse(t, response, http.StatusOK, &actual)
	expected := modelsv1.Room{
		ID:       id,
		Name:     "31415",
		Flags:    []string{},
		Size:     2,
		Comments: p("A nice comment for 31415"),
	}
	expected.Occupants = append(expected.Occupants, occupants...)
	tstEqualResponseBodies(t, expected, actual)
}

func tstGroupMailToOwner(cid string, groupName string, target string, object string) mailservice.MailSendDto {
	_, targetNick, targetEmail := tstInfosBySubject(target)
	objectBadge, objectNick, _ := tstInfosBySubject(object)

	return mailservice.MailSendDto{
		CommonID: cid,
		Lang:     "en-US",
		To:       []string{targetEmail},
		Variables: map[string]string{
			"nickname":            targetNick,
			"groupname":           groupName,
			"object_badge_number": objectBadge,
			"object_nickname":     objectNick,
		},
	}
}

func tstGroupMailToMember(cid string, groupName string, target string, url string) mailservice.MailSendDto {
	_, targetNick, targetEmail := tstInfosBySubject(target)

	result := mailservice.MailSendDto{
		CommonID: cid,
		Lang:     "en-US",
		To:       []string{targetEmail},
		Variables: map[string]string{
			"nickname":  targetNick,
			"groupname": groupName,
		},
	}
	if url != "" {
		result.Variables["url"] = url
	}
	return result
}
