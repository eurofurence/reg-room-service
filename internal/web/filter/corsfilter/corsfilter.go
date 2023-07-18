package corsfilter

import (
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/logging"
	"github.com/go-http-utils/headers"
	"net/http"
)

func createCorsHeadersHandler(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if config.IsCorsDisabled() {
			logging.Ctx(ctx).Warn("sending headers to disable CORS. This configuration is not intended for production use, only for local development!")
			w.Header().Set(headers.AccessControlAllowOrigin, "*")
			w.Header().Set(headers.AccessControlAllowMethods, "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set(headers.AccessControlAllowHeaders, "content-type")
			w.Header().Set(headers.AccessControlExposeHeaders, "Location, X-B3-TraceId")
		}

		if r.Method == http.MethodOptions {
			logging.Ctx(ctx).Info("received OPTIONS request. Responding with OK.")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
	return handlerFunc
}

// would not need this extra layer in the absence of parameters

func CorsHeadersMiddleware() func(http.Handler) http.Handler {
	middlewareCreator := func(next http.Handler) http.Handler {
		return http.HandlerFunc(createCorsHeadersHandler(next))
	}
	return middlewareCreator
}
