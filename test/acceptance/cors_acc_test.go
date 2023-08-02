package acceptance

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eurofurence/reg-room-service/docs"
)

func TestCountdownNoCors(t *testing.T) {
	docs.Given("given a valid configuration for production")
	tstSetup(tstDefaultConfigFileBeforeLaunch)
	defer tstShutdown()

	docs.When("when they request the countdown resource")
	response := tstPerformGet("/api/rest/v1/countdown", "")

	docs.Then("then a valid response is sent that does not include the CORS headers")
	require.Equal(t, http.StatusOK, response.StatusCode, "unexpected http response status")
	require.Nil(t, response.Header["Access-Control-Allow-Origin"])
	require.Nil(t, response.Header["Access-Control-Allow-Methods"])
	require.Nil(t, response.Header["Access-Control-Allow-Headers"])
	require.Nil(t, response.Header["Access-Control-Expose-Headers"])
}
