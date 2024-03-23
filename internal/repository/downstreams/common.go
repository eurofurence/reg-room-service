package downstreams

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/eurofurence/reg-room-service/internal/config"

	aurestlogging "github.com/StephanHCB/go-autumn-restclient/implementation/requestlogging"
	"github.com/go-chi/chi/v5/middleware"

	aurestbreaker "github.com/StephanHCB/go-autumn-restclient-circuitbreaker/implementation/breaker"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	auresthttpclient "github.com/StephanHCB/go-autumn-restclient/implementation/httpclient"
	"github.com/go-http-utils/headers"

	"github.com/eurofurence/reg-room-service/internal/logging"
	"github.com/eurofurence/reg-room-service/internal/web/common"
)

//nolint
const apiKeyHeader = "X-Api-Key"

var (
	ErrDownStreamUnavailable = errors.New("downstream unavailable - see log for details")
)

func requestIDFromContext(ctx context.Context) string {
	if reqID, ok := ctx.Value(logging.RequestIDKey).(string); ok {
		return reqID
	}

	return "ffffffff"
}

func ApiTokenRequestManipulator(fixedApiToken string) aurestclientapi.RequestManipulatorCallback {
	return func(ctx context.Context, r *http.Request) {
		r.Header.Add(apiKeyHeader, fixedApiToken)
		r.Header.Add(middleware.RequestIDHeader, requestIDFromContext(ctx))
	}
}

func AccessTokenForwardingRequestManipulator() aurestclientapi.RequestManipulatorCallback {
	return func(ctx context.Context, r *http.Request) {
		accessToken, ok := ctx.Value(common.CtxKeyAccessToken{}).(string)
		if ok {
			r.Header.Add(headers.Authorization, "Bearer "+accessToken)
		}
		r.Header.Add(middleware.RequestIDHeader, requestIDFromContext(ctx))
	}
}

func CookiesOrAuthHeaderForwardingRequestManipulator(conf config.SecurityConfig) aurestclientapi.RequestManipulatorCallback {
	return func(ctx context.Context, r *http.Request) {
		r.Header.Add(middleware.RequestIDHeader, requestIDFromContext(ctx))

		idToken, ok2 := ctx.Value(common.CtxKeyIDToken{}).(string)
		accessToken, ok3 := ctx.Value(common.CtxKeyAccessToken{}).(string)

		if ok2 && ok3 {
			r.AddCookie(&http.Cookie{
				Name:     conf.Oidc.IDTokenCookieName,
				Value:    idToken,
				Domain:   "localhost",
				Expires:  time.Now().Add(10 * time.Minute),
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			})
			r.AddCookie(&http.Cookie{
				Name:     conf.Oidc.AccessTokenCookieName,
				Value:    accessToken,
				Domain:   "localhost",
				Expires:  time.Now().Add(10 * time.Minute),
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			})
		} else {
			// downstream service may need to contact idp to get identity info, but better than nothing
			r.Header.Add(headers.Authorization, "Bearer "+accessToken)
		}
	}
}

func ClientWith(requestManipulator aurestclientapi.RequestManipulatorCallback, circuitBreakerName string) (aurestclientapi.Client, error) {
	httpClient, err := auresthttpclient.New(0, nil, requestManipulator)
	if err != nil {
		return nil, err
	}

	requestLoggingClient := aurestlogging.New(httpClient)

	circuitBreakerClient := aurestbreaker.New(requestLoggingClient,
		circuitBreakerName,
		10,
		2*time.Minute,
		30*time.Second,
		15*time.Second,
	)

	return circuitBreakerClient, nil
}

func ErrByStatus(err error, status int) error {
	if err != nil {
		return err
	}
	if status >= 300 {
		return ErrDownStreamUnavailable
	}
	return nil
}
