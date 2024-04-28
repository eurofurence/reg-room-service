package acceptance

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eurofurence/reg-room-service/docs"
	v1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

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
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	responseDto := v1.Countdown{}
	tstParseJson(response.body, &responseDto)

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
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	responseDto := v1.Countdown{}
	tstParseJson(response.body, &responseDto)

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
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	responseDto := v1.Countdown{}
	tstParseJson(response.body, &responseDto)

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
	response := tstPerformGet("/api/rest/v1/countdown", valid_JWT_is_staff_sub202)

	docs.Then("then a valid response is sent with countdown <= 0 that includes the secret")
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	responseDto := v1.Countdown{}
	tstParseJson(response.body, &responseDto)

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
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	responseDto := v1.Countdown{}
	tstParseJson(response.body, &responseDto)

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
	response := tstPerformGet("/api/rest/v1/countdown", valid_JWT_is_registered_sub101)

	docs.Then("then a valid response is sent with countdown > 0 that does not include the secret")
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	responseDto := v1.Countdown{}
	tstParseJson(response.body, &responseDto)

	require.True(t, responseDto.CountdownSeconds > 0, "unexpected countdown value is not positive")
	require.Equal(t, "3021-12-31T23:59:59+01:00", responseDto.TargetTimeIsoDateTime, "unexpected target time")
	require.NotNil(t, responseDto.CurrentTimeIsoDateTime, "unexpected current time is nil")
	require.Equal(t, "", responseDto.Secret, "unexpected secret is not empty")
}
