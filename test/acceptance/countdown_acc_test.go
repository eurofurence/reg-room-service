package acceptance

import (
	"github.com/eurofurence/reg-room-service/internal/api/v1"
	"net/http"
	"testing"

	"github.com/eurofurence/reg-room-service/docs"
	"github.com/stretchr/testify/require"
)

// ------------------------------------------
// tokens for acceptance tests with various claims
// ------------------------------------------

const valid_JWT_is_staff = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZ2xvYmFsIjp7Im5hbWUiOiJKb2huIERvZSIsInJvbGVzIjpbInN0YWZmIl19LCJpYXQiOjE1MTYyMzkwMjJ9.XKf8Eqrs7JmGSNhVUS-5RLIMnSuQHh65VMiUJPHaE5AEFZ55EY7MxsD2Sqdc6QV9cX0zA5weGXX2cGOAR0CNcjOGGsQSVogAcoEuwjve4WXLVvHPb41p95Jkbe9Md5bSPrk9oJwopJCVDI5DU1rLg0FIbt2yWORinZQiGvxZlPSZyNQuFoAXJXQPv4TNfTaBZcKhzUeO0u6_AQzIKGrF8VmbE4cMHq0fEAflnzroDmo9oJ-8dKJc2BNEyFQYHi9Jp3h3C85BvxEsdRzL3e9Qjw2SpFS0A8pPr4HEQikIn2nOEXav2RAcZMGN3YmdUeUBHwnfQ9ubY-0KilK9zNfGBw`
const valid_JWT_is_not_staff = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZ2xvYmFsIjp7Im5hbWUiOiJKb2huIERvZSIsInJvbGVzIjpbXX0sImlhdCI6MTUxNjIzOTAyMn0.IH3Q46k85RZsvgWD3wC9kNCtCRujTEOpzzCw6rqrKF4QoDcmn6Pd-Y2qQ8IZydrtzGrCu7yUiVziL634gxDlRvVliyHU6KkIMMsXDtnJWOGrKkpJgr_PZCA2LIlYD0GsXYzzQBuOg3eeXgidkGD7WVjHuKcuJe5By9nc6cTHlBHV-XeRIeCCy9jq10pbqyNv1kfjhdKuUQpFogV2JIKlTi3cR5pZalahYLe4o2iArcQHz3_VRsYd7frWN2kkF4ARwQl3UlOHH6jOSzT5h6PtnOJ1pDpIGME5NqG3TDvQnom5TAKW-XiZckk5lJAp3I51qGvDjve1AyZCRPfDHMsKAA`

// ------------------------------------------
// acceptance tests for the countdown resource
// ------------------------------------------

func TestCountdownBeforeLaunch(t *testing.T) {
	docs.Given("given a launch date in the future")
	tstSetup(tstDefaultConfigFileBeforeLaunch)
	defer tstShutdown()

	docs.When("when they request the countdown resource before the launch time has been reached")
	response := tstPerformGet("/api/rest/v1/countdown", "")

	docs.Then("then a valid response is sent with countdown > 0 that does not include the secret")
	require.Equal(t, http.StatusOK, response.StatusCode, "unexpected http response status")
	responseDto := v1.CountdownResultDto{}
	tstParseJson(tstBodyToString(response), &responseDto)

	require.True(t, responseDto.CountdownSeconds > 0, "unexpected countdown value is not positive")
	require.Equal(t, "3021-12-31T23:59:59+01:00", responseDto.TargetTimeIsoDateTime, "unexpected target time")
	require.NotNil(t, responseDto.CurrentTimeIsoDateTime, "unexpected current time is nil")
	require.Equal(t, "", responseDto.Secret, "unexpected secret is not empty")
}

func TestCountdownAfterPublicLaunch(t *testing.T) {
	docs.Given("given a public launch date in the past")
	tstSetup(tstDefaultConfigFileAfterPublicLaunch)
	defer tstShutdown()

	docs.When("when they request the countdown resource after the public launch time has been reached")
	response := tstPerformGet("/api/rest/v1/countdown", "")

	docs.Then("then a valid response is sent with countdown <= 0 that includes the secret")
	require.Equal(t, http.StatusOK, response.StatusCode, "unexpected http response status")
	responseDto := v1.CountdownResultDto{}
	tstParseJson(tstBodyToString(response), &responseDto)

	require.True(t, responseDto.CountdownSeconds <= 0, "unexpected countdown value is not negative")
	require.Equal(t, "2020-12-31T23:59:59+01:00", responseDto.TargetTimeIsoDateTime, "unexpected target time")
	require.NotNil(t, responseDto.CurrentTimeIsoDateTime, "unexpected current time is nil")
	require.Equal(t, "Kaiser-Wilhelm-Koog", responseDto.Secret, "unexpected secret")
}

func TestCountdownAfterStaffLaunchWithoutStaffClaim(t *testing.T) {
	docs.Given("given a staff launch date in the past")
	tstSetup(tstDefaultConfigFileAfterStaffLaunch)
	defer tstShutdown()

	docs.When("when they request the countdown resource after the staff launch time has been reached")
	response := tstPerformGet("/api/rest/v1/countdown", "")

	docs.Then("then a valid response is sent with countdown <= 0 that includes the secret")
	require.Equal(t, http.StatusOK, response.StatusCode, "unexpected http response status")
	responseDto := v1.CountdownResultDto{}
	tstParseJson(tstBodyToString(response), &responseDto)

	require.True(t, responseDto.CountdownSeconds > 0, "unexpected countdown value is not positive")
	require.Equal(t, "3021-12-31T23:59:59+01:00", responseDto.TargetTimeIsoDateTime, "unexpected target time")
	require.NotNil(t, responseDto.CurrentTimeIsoDateTime, "unexpected current time is nil")
	require.Equal(t, "", responseDto.Secret, "unexpected secret is not empty")
}

func TestCountdownAfterStaffLaunchWithStaffClaim(t *testing.T) {
	docs.Given("given a staff launch date in the past")
	tstSetup(tstDefaultConfigFileAfterStaffLaunch)
	defer tstShutdown()

	docs.When("when they request the countdown resource after the staff launch time has been reached")
	response := tstPerformGet("/api/rest/v1/countdown", valid_JWT_is_staff)

	docs.Then("then a valid response is sent with countdown <= 0 that includes the secret")
	require.Equal(t, http.StatusOK, response.StatusCode, "unexpected http response status")
	responseDto := v1.CountdownResultDto{}
	tstParseJson(tstBodyToString(response), &responseDto)

	require.True(t, responseDto.CountdownSeconds <= 0, "unexpected countdown value is not negative")
	require.Equal(t, "2020-12-31T23:59:59+01:00", responseDto.TargetTimeIsoDateTime, "unexpected target time")
	require.NotNil(t, responseDto.CurrentTimeIsoDateTime, "unexpected current time is nil")
	require.Equal(t, "Dithmarschen", responseDto.Secret, "unexpected secret")
}

func TestCountdownBeforeLaunchWithMockTime(t *testing.T) {
	docs.Given("given a launch date in the future")
	tstSetup(tstDefaultConfigFileBeforeLaunch)
	defer tstShutdown()

	docs.When("when they request the countdown resource before the launch time has been reached using a mock time in the future")
	response := tstPerformGet("/api/rest/v1/countdown?currentTimeIso=3022-12-31T23:59:59%2B01:00", "")

	docs.Then("then a valid response is sent with countdown <= 0 that does not include the real secret")
	require.Equal(t, http.StatusOK, response.StatusCode, "unexpected http response status")
	responseDto := v1.CountdownResultDto{}
	tstParseJson(tstBodyToString(response), &responseDto)

	require.True(t, responseDto.CountdownSeconds <= 0, "unexpected countdown value is not negative")
	require.Equal(t, "3021-12-31T23:59:59+01:00", responseDto.TargetTimeIsoDateTime, "unexpected target time")
	require.Equal(t, "3022-12-31T23:59:59+01:00", responseDto.CurrentTimeIsoDateTime, "unexpected current time")
	require.Equal(t, "[demo-secret]", responseDto.Secret, "unexpected secret is not demo secret")
}

// security tests

func TestCountdownBeforeLaunch_DenyNonStaffToken(t *testing.T) {
	docs.Given("given a launch date in the future")
	tstSetup(tstDefaultConfigFileBeforeLaunch)
	defer tstShutdown()

	docs.When("when they request the countdown resource before the launch time has been reached, using a non-staff token")
	response := tstPerformGet("/api/rest/v1/countdown", valid_JWT_is_not_staff)

	docs.Then("then a valid response is sent with countdown > 0 that does not include the secret")
	require.Equal(t, http.StatusOK, response.StatusCode, "unexpected http response status")
	responseDto := v1.CountdownResultDto{}
	tstParseJson(tstBodyToString(response), &responseDto)

	require.True(t, responseDto.CountdownSeconds > 0, "unexpected countdown value is not positive")
	require.Equal(t, "3021-12-31T23:59:59+01:00", responseDto.TargetTimeIsoDateTime, "unexpected target time")
	require.NotNil(t, responseDto.CurrentTimeIsoDateTime, "unexpected current time is nil")
	require.Equal(t, "", responseDto.Secret, "unexpected secret is not empty")
}
