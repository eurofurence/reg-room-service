package middleware

import (
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/go-http-utils/headers"
	"net/http"
)

const TraceIdHeader = "X-Request-Id"

func CorsHandling(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if config.IsCorsDisabled() {
			w.Header().Set(headers.AccessControlAllowOrigin, "*")
			w.Header().Set(headers.AccessControlAllowMethods, "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set(headers.AccessControlAllowHeaders, "content-type")
			w.Header().Set(headers.AccessControlAllowCredentials, "true")
			w.Header().Set(headers.AccessControlExposeHeaders, "Location, "+TraceIdHeader)
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
