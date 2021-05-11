package acceptance

import (
	"net/http"
	"testing"

	"github.com/eurofurence/reg-room-service/api/v1/countdown"
	"github.com/eurofurence/reg-room-service/docs"
	"github.com/stretchr/testify/require"
)

// ------------------------------------------
// tokens for acceptance tests with various claims
// ------------------------------------------

const valid_JWT_is_admin = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.POstGetfAytaZS82wHcjoTyoqhMyxXiWdR7Nn7A29DNSl0EiXLdwJ6xC6AfgZWF1bOsS_TuYI3OG85AmiExREkrS6tDfTQ2B3WXlrr-wp5AokiRbz3_oB4OxG-W9KcEEbDRcZc0nH3L7LzYptiy1PtAylQGxHTWZXtGz4ht0bAecBgmpdgXMguEIcoqPJ1n3pIWk_dUZegpqx0Lka21H6XxUTxiy8OcaarA8zdnPUnV6AmNP3ecFawIFYdvJB_cm-GvpCSbr8G8y_Mllj8f4x9nBH8pQux89_6gUY618iYv7tuPWBFfEbLxtF2pZS6YC1aSfLQxeNe8djT9YjpvRZA`
const valid_JWT_is_not_admin = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOmZhbHNlLCJpYXQiOjE1MTYyMzkwMjJ9.ZWcaBvG4KlHKTEnKcEV3EB3h3L92SjSlJ7vCMcuJEUS3Ad7oWpOhK2aawdPshccD-JkUAh4lGmHNBy-MmxcBumO-5TbeUZaDY9BoCaHA_XH5uohK7d-WjLPOgHQ0pnyRXi90FfY4m1nQyx1dtAQk0rYYABKVN707OFIHegtIoEV_Ie5j1OmHFycCykfXkdx9qLPPCHaREgXtD0_5h9uVq83ODBy_5O0Lq8Ed0j6smgJPsUFuxHYB3oN61GUp4VkzdTY7VwATgzRAcCv4d5-CAz2s0czcUpSC_NEe0dQbYY9vNmJ90kjIXtDFJUTzG_jeZ2lvAWshNP5jUUxgrcL1oA`

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
	responseDto := countdown.CountdownResultDto{}
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
	responseDto := countdown.CountdownResultDto{}
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
	responseDto := countdown.CountdownResultDto{}
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
	response := tstPerformGet("/api/rest/v1/countdown", valid_JWT_is_admin)

	docs.Then("then a valid response is sent with countdown <= 0 that includes the secret")
	require.Equal(t, http.StatusOK, response.StatusCode, "unexpected http response status")
	responseDto := countdown.CountdownResultDto{}
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
	responseDto := countdown.CountdownResultDto{}
	tstParseJson(tstBodyToString(response), &responseDto)

	require.True(t, responseDto.CountdownSeconds <= 0, "unexpected countdown value is not negative")
	require.Equal(t, "3021-12-31T23:59:59+01:00", responseDto.TargetTimeIsoDateTime, "unexpected target time")
	require.Equal(t, "3022-12-31T23:59:59+01:00", responseDto.CurrentTimeIsoDateTime, "unexpected current time")
	require.Equal(t, "[demo-secret]", responseDto.Secret, "unexpected secret is not demo secret")
}

// security tests

func TestCountdownBeforeLaunch_DenyNonAdminToken(t *testing.T) {
	docs.Given("given a launch date in the future")
	tstSetup(tstDefaultConfigFileBeforeLaunch)
	defer tstShutdown()

	docs.When("when they request the countdown resource before the launch time has been reached, using a non-admin token")
	response := tstPerformGet("/api/rest/v1/countdown", valid_JWT_is_not_admin)

	docs.Then("then a valid response is sent with countdown > 0 that does not include the secret")
	require.Equal(t, http.StatusOK, response.StatusCode, "unexpected http response status")
	responseDto := countdown.CountdownResultDto{}
	tstParseJson(tstBodyToString(response), &responseDto)

	require.True(t, responseDto.CountdownSeconds > 0, "unexpected countdown value is not positive")
	require.Equal(t, "3021-12-31T23:59:59+01:00", responseDto.TargetTimeIsoDateTime, "unexpected target time")
	require.NotNil(t, responseDto.CurrentTimeIsoDateTime, "unexpected current time is nil")
	require.Equal(t, "", responseDto.Secret, "unexpected secret is not empty")
}
