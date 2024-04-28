package middleware

import (
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"net/http"
	"runtime/debug"
)

func PanicRecoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rvr := recover()
			if rvr != nil && rvr != http.ErrAbortHandler {
				ctx := r.Context()
				stack := string(debug.Stack())
				aulogging.Logger.Ctx(ctx).Error().Print("recovered from PANIC: " + stack)
				common.SendErrorWithStatusAndMessage(ctx, w, http.StatusInternalServerError, "internal.error", "")
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
