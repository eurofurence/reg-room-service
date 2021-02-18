package acceptance

import (
	"github.com/eurofurence/reg-room-service/api/v1/countdown"
	"github.com/eurofurence/reg-room-service/docs"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	//"time"
)

// ------------------------------------------
// acceptance tests for the countdown resource
// ------------------------------------------

func TestCountdownBeforeLaunch(t *testing.T) {
	tstSetup()
	defer tstShutdown()

	docs.Given("given a launch date in the future")
	//launchDate := time.Now().AddDate(0, 1, 0) // one month into the future

	docs.When("when they request the countdown resource before the launch time has been reached")
	response := tstPerformGet("/api/rest/v1/countdown")

	docs.Then("then a valid response is sent with countdown > 0")
	require.Equal(t, http.StatusOK, response.StatusCode, "unexpected http response status")
	responseDto := countdown.CountdownResultDto{}
	tstParseJson(tstBodyToString(response), &responseDto)

	require.True(t, responseDto.CountdownSeconds > 0, "unexpected countdown value is not positive")
	require.Equal(t, "2021-12-31T23:59:59+01:00", responseDto.TargetTimeIsoDateTime, "unexpected countdown value is not positive")
	require.NotNil(t, responseDto.CurrentTimeIsoDateTime, "unexpected countdown value is not positive")
	require.Equal(t, "", responseDto.Secret, "unexpected secret is not nil")
}

func TestCountdownAfterLaunch(t *testing.T) {
}
