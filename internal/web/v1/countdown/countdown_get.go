package countdown

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

// GetCountdownRequest is the request for the GetCountdown operation.
type GetCountdownRequest struct{}

// GetCountdown returns the countdown information.
// If the countdown has reached 0, also reveals the configured secret.
// You need to be logged in.
func (h *Handler) GetCountdown(ctx context.Context, req *GetCountdownRequest, w http.ResponseWriter) (*modelsv1.Countdown, error) {
	return nil, nil
}

// GetCountdownRequest validates and creates the request for the GetCountdown operation.
func (h *Handler) GetCountdownRequest(r *http.Request, w http.ResponseWriter) (*GetCountdownRequest, error) {
	return nil, nil
}

// GetCountdownResponse writes out the response for the GetCountdown operation.
func (h *Handler) GetCountdownResponse(ctx context.Context, res *modelsv1.Countdown, w http.ResponseWriter) error {
	return nil
}
