package reqid

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

var RequestIDHeader = "X-Request-Id"

type ctxKeyRequestID int
const RequestIDKey ctxKeyRequestID = 0

func createReqIdHandler(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		reqUuidStr := r.Header.Get(RequestIDHeader)
		if reqUuidStr == "" {
			reqUuid, err := uuid.NewRandom()
			if err == nil {
				reqUuidStr = reqUuid.String()[:8]
			} else {
				// this should not normally ever happen, but continue with this fixed requestId
				reqUuidStr ="ffffffff"
			}
		}
		ctx := r.Context()
		newCtx := context.WithValue(ctx, RequestIDKey, reqUuidStr)
		r = r.WithContext(newCtx)

		next.ServeHTTP(w, r)
	}
	return handlerFunc
}

// would not need this extra layer in the absence of parameters

func RequestIdMiddleware() func(http.Handler) http.Handler {
	middlewareCreator := func(next http.Handler) http.Handler {
		return http.HandlerFunc(createReqIdHandler(next))
	}
	return middlewareCreator
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return "00000000"
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return "ffffffff"
}
