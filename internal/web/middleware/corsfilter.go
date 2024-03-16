package middleware

import (
	"github.com/go-http-utils/headers"

	"github.com/eurofurence/reg-room-service/internal/config"
	"github.com/eurofurence/reg-room-service/internal/logging"

	"log"
	"net/http"
)

func createCorsHeadersHandler(next http.Handler, config *config.SecurityConfig) func(w http.ResponseWriter, r *http.Request) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Example for cors middleware
		if config != nil && config.Cors.DisableCors {
			logging.LoggerFromContext(ctx).Warn("sending headers to disable CORS. This configuration is not intended for production use, only for local development!")
			w.Header().Set(headers.AccessControlAllowOrigin, config.Cors.AllowOrigin)
			w.Header().Set(headers.AccessControlAllowMethods, "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set(headers.AccessControlAllowHeaders, "content-type")
			w.Header().Set(headers.AccessControlAllowCredentials, "true")
			w.Header().Set(headers.AccessControlExposeHeaders, "Location, X-B3-TraceId")
		}

		if r.Method == http.MethodOptions {
			log.Println("INFO received OPTIONS request. Responding with OK.")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
	return handlerFunc
}

// would not need this extra layer in the absence of parameters

func CorsHeadersMiddleware(config *config.SecurityConfig) func(http.Handler) http.Handler {
	middlewareCreator := func(next http.Handler) http.Handler {
		return http.HandlerFunc(createCorsHeadersHandler(next, config))
	}
	return middlewareCreator
}
