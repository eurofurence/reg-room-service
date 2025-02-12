package middleware

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"net/http"
	"strings"

	"github.com/go-http-utils/headers"
	"github.com/golang-jwt/jwt/v4"

	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/authservice"
)

type openEndpoint struct {
	Method string
	Path   string
}

// openEndpoints configures the list of endpoints that can be called without authorization.
//
// If a token is provided, it is still parsed and processed.
var openEndpoints = []openEndpoint{
	{
		// countdown endpoint is allowed through (but AFTER auth processing)
		Method: http.MethodGet,
		Path:   "/api/rest/v1/countdown",
	},
}

// nolint
const (
	bearerPrefix = "Bearer "
	// TODO Remove after legacy system was replaced with 2FA
	// See reference https://github.com/eurofurence/reg-room-service/issues/57
	adminRequestHeader = "X-Admin-Request"
)

// --- getting the values from the request ---

func parseAuthCookie(r *http.Request, cookieName string) string {
	if cookieName == "" {
		// ok if not configured, don't accept cookies then
		return ""
	}

	authCookie, _ := r.Cookie(cookieName)
	if authCookie == nil {
		// missing cookie is not considered an error, either
		return ""
	}

	return authCookie.Value
}

func fromAuthHeader(r *http.Request) string {
	headerValue := r.Header.Get(headers.Authorization)

	if !strings.HasPrefix(headerValue, bearerPrefix) {
		return ""
	}

	return strings.TrimPrefix(headerValue, bearerPrefix)
}

func fromApiTokenHeader(r *http.Request) string {
	return r.Header.Get(common.ApiKeyHeader)
}

// --- validating the individual pieces ---

// important - if any of these return an error, you must abort processing via "return" and log the error message

func checkApiToken(ctx context.Context, conf *config.SecurityConfig, apiTokenValue string) (context.Context, bool, error) {
	if apiTokenValue != "" {
		// ignore jwt if set (may still need to pass it through to other service)
		if apiTokenValue == conf.Fixed.API {
			ctx = context.WithValue(ctx, common.CtxKeyAPIKey{}, apiTokenValue)
			return ctx, true, nil
		} else {
			return ctx, false, errors.New("token doesn't match the configured value")
		}
	}
	return ctx, false, nil
}

func checkAccessToken(ctx context.Context, conf *config.SecurityConfig, accessTokenValue string) (context.Context, bool, error) {
	if accessTokenValue != "" {
		if authservice.Get().IsEnabled() {
			authCtx := context.WithValue(ctx, common.CtxKeyAccessToken{}, accessTokenValue) // need this set for userinfo call

			userInfo, err := authservice.Get().UserInfo(authCtx)
			if err != nil {
				return ctx, false, fmt.Errorf("request failed access token check, denying: %s", err.Error())
			}

			if conf.Oidc.Audience != "" {
				if len(userInfo.Audiences) != 1 || userInfo.Audiences[0] != conf.Oidc.Audience {
					return ctx, false, errors.New("token audience does not match")
				}
			}

			overwriteClaims := common.AllClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:   conf.Oidc.Issuer,
					Subject:  userInfo.Subject,
					Audience: jwt.ClaimStrings{conf.Oidc.Audience},
				},
				CustomClaims: common.CustomClaims{
					EMail:         userInfo.Email,
					EMailVerified: userInfo.EmailVerified,
					Groups:        userInfo.Groups,
					Name:          userInfo.Name,
				},
			}

			ctx = context.WithValue(authCtx, common.CtxKeyClaims{}, &overwriteClaims)
			return ctx, true, nil
		} else {
			return ctx, false, errors.New("request failed access token check, denying: no userinfo endpoint configured")
		}
	}
	return ctx, false, nil
}

func keyFuncForKey(rsaPublicKey *rsa.PublicKey) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		return rsaPublicKey, nil
	}
}

func checkIdToken(
	ctx context.Context,
	conf *config.SecurityConfig,
	idTokenValue string,
) (retCtx context.Context, success bool, err error) {
	if idTokenValue != "" {
		tokenString := strings.TrimSpace(idTokenValue)

		errorMessage := ""
		for _, key := range parsedPEMs {
			claims := common.AllClaims{}
			token, err := jwt.ParseWithClaims(tokenString, &claims, keyFuncForKey(key), jwt.WithValidMethods([]string{"RS256", "RS512"}))
			if err == nil && token.Valid {
				parsedClaims, ok := token.Claims.(*common.AllClaims)
				if ok {
					if conf.Oidc.Audience != "" {
						if len(parsedClaims.Audience) != 1 || parsedClaims.Audience[0] != conf.Oidc.Audience {
							return ctx, false, errors.New("token audience does not match")
						}
					}

					if conf.Oidc.Issuer != "" {
						if parsedClaims.Issuer != conf.Oidc.Issuer {
							return ctx, false, errors.New("token issuer does not match")
						}
					}

					ctx = context.WithValue(ctx, common.CtxKeyIDToken{}, tokenString)
					ctx = context.WithValue(ctx, common.CtxKeyClaims{}, &claims)
					return ctx, true, nil
				}
				errorMessage = "empty claims substructure"
			} else if err != nil {
				errorMessage = err.Error()
			} else {
				errorMessage = "token parsed but invalid"
			}
		}
		return ctx, false, errors.New(errorMessage)
	}
	return ctx, false, nil
}

// --- top level ---.
func checkAllAuthentication(ctx context.Context, method string, urlPath string, conf *config.SecurityConfig, apiTokenHeaderValue string, authHeaderValue string, idTokenCookieValue string, accessTokenCookieValue string) (context.Context, string, error) {
	var success bool
	var err error

	// health check on / is allowed through
	if method == http.MethodGet && urlPath == "/" {
		return ctx, "", nil
	}

	// try api token first
	ctx, success, err = checkApiToken(ctx, conf, apiTokenHeaderValue)
	if err != nil {
		return ctx, "invalid api token", err
	}
	if success {
		return ctx, "", nil
	}

	// now try authorization header (gives only access token, so MUST use userinfo endpoint)
	ctx, success, err = checkAccessToken(ctx, conf, authHeaderValue)
	if err != nil {
		return ctx, "invalid bearer token", err
	}
	if success {
		return ctx, "", nil
	}

	// now try cookie pair
	ctx, success, err = checkIdToken(ctx, conf, idTokenCookieValue)
	if err != nil {
		return ctx, "invalid id token in cookie", err
	}
	if success {
		ctx, success, err = checkAccessToken(ctx, conf, accessTokenCookieValue)
		if err != nil {
			return ctx, "invalid or missing access token in cookie", err
		}
		if success {
			return ctx, "", nil
		}
	}

	// allow through (but still AFTER auth processing)
	for _, publicEndpoint := range openEndpoints {
		if method == publicEndpoint.Method && urlPath == publicEndpoint.Path {
			return ctx, "", nil
		}
	}

	return ctx, "you must be logged in for this operation", errors.New("no authorization presented")
}

// --- middleware validating the values and adding to context values ---

var parsedPEMs []*rsa.PublicKey

func CheckRequestAuthorization(conf *config.SecurityConfig) func(http.Handler) http.Handler {
	parsedPEMs = make([]*rsa.PublicKey, len(conf.Oidc.TokenPublicKeysPEM))

	for i, publicKey := range conf.Oidc.TokenPublicKeysPEM {
		rsaKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			panic("Couldn't parse configured pem " + publicKey)
		}

		parsedPEMs[i] = rsaKey
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			apiTokenHeaderValue := fromApiTokenHeader(r)
			authHeaderValue := fromAuthHeader(r)
			idTokenCookieValue := parseAuthCookie(r, conf.Oidc.IDTokenCookieName)
			accessTokenCookieValue := parseAuthCookie(r, conf.Oidc.AccessTokenCookieName)

			ctx, userFacingErrorMessage, err := checkAllAuthentication(ctx, r.Method, r.URL.Path, conf, apiTokenHeaderValue, authHeaderValue, idTokenCookieValue, accessTokenCookieValue)
			if err != nil {
				aulogging.WarnErrf(ctx, err, "authorization failed: %s: %s", userFacingErrorMessage, err.Error())
				web.SendUnauthorizedResponse(ctx, w, userFacingErrorMessage)
				return
			}

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
