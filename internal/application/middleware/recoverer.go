package middleware

import (
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/application/web"
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
				aulogging.Error(ctx, "recovered from PANIC: "+stack)
				web.SendErrorWithStatusAndMessage(ctx, w, http.StatusInternalServerError, common.InternalErrorMessage, "")
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
