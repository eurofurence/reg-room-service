package acceptance

import (
	"net/http"
	"testing"

	"github.com/eurofurence/reg-room-service/docs"
	"github.com/stretchr/testify/require"
)

func TestCountdownNoCors(t *testing.T) {
	t.Skip("Skipping until implementation can be done properly")
	docs.Given("given a valid configuration for production")
	tstSetup(tstDefaultConfigFileBeforeLaunch)
	defer tstShutdown()

	docs.When("when they request the countdown resource")
	response := tstPerformGet("/api/rest/v1/countdown", "")

	docs.Then("then a valid response is sent that does not include the CORS headers")
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	require.Nil(t, response.header["Access-Control-Allow-Origin"])
	require.Nil(t, response.header["Access-Control-Allow-Methods"])
	require.Nil(t, response.header["Access-Control-Allow-Headers"])
	require.Nil(t, response.header["Access-Control-Expose-Headers"])
}
