package acceptance

import (
	"net/http"
	"testing"

	"github.com/eurofurence/reg-room-service/docs"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------
// acceptance tests for the health endpoint
// ----------------------------------------

func TestHealthEndpoint(t *testing.T) {
	t.Skip("Skipping until implementation can be done properly")
	docs.Given("given an unauthenticated user")
	tstSetup(tstDefaultConfigFileBeforeLaunch)
	defer tstShutdown()

	docs.When("when the user accesses the health endpoint")
	response := tstPerformGet("/", "")

	docs.Then("then the operation is successful")
	require.Equal(t, http.StatusOK, response.status, "unexpected http status")
}
