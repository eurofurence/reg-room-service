package acceptance

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eurofurence/reg-room-service/docs"
)

// ----------------------------------------
// acceptance tests for the health endpoint
// ----------------------------------------

func TestHealthEndpoint(t *testing.T) {
	docs.Given("given an unauthenticated user")
	tstSetup(tstDefaultConfigFileBeforeLaunch)
	defer tstShutdown()

	docs.When("when the user accesses the health endpoint")
	response := tstPerformGet("/", "")

	docs.Then("then the operation is successful")
	require.Equal(t, http.StatusOK, response.StatusCode, "unexpected http status")
}
