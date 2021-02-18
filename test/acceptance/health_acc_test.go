package acceptance

import (
	"github.com/eurofurence/reg-room-service/docs"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

// ----------------------------------------
// acceptance tests for the health endpoint
// ----------------------------------------

func TestHealthEndpoint(t *testing.T) {
	docs.Given("given an unauthenticated user")
	tstSetup(tstDefaultConfigFileBeforeLaunch)
	defer tstShutdown()

	docs.When("when the user accesses the health endpoint")
	response := tstPerformGet("/")

	docs.Then("then the operation is successful")
	require.Equal(t, http.StatusOK, response.StatusCode, "unexpected http status")
}
