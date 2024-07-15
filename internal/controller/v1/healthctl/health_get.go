package healthctl

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

// GetHealthRequest is the request for the GetCountdown operation.
type GetHealthRequest struct{}

// GetHealth is used as a simple health check.
func (*Handler) GetHealth(ctx context.Context, req *GetHealthRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

// GetHealthRequest validates and creates the request for the GetCountdown operation.
func (*Handler) GetHealthRequest(r *http.Request, w http.ResponseWriter) (*GetHealthRequest, error) {
	return nil, nil
}

// GetHealthResponse writes out the response for the GetCountdown operation.
func (*Handler) GetHealthResponse(ctx context.Context, res *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}
