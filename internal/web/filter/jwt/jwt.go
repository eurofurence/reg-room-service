package jwt

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eurofurence/reg-room-service/internal/repository/config"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"

	"github.com/eurofurence/reg-room-service/internal/repository/logging"
)

const userProperty = "user"

func fromCookie(r *http.Request) (string, error) {
	cookieName := config.JWTCookieName()
	if cookieName == "" {
		// ok if not configured, don't accept cookies then
		return "", nil
	}

	authCookie, _ := r.Cookie(cookieName)
	if authCookie == nil {
		// missing cookie is not considered an error, either
		return "", nil
	}

	return authCookie.Value, nil
}

func errorHandler(w http.ResponseWriter, r *http.Request, err string) {
	// avoid default error handler that may leak information from the error and prints an English message which confuses clients
}

func createHandlerFunction(jwtMiddleware *jwtmiddleware.JWTMiddleware, next http.Handler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := jwtMiddleware.CheckJWT(w, r)
		if err != nil {
			// missing JWT is not considered an error
			logging.Ctx(r.Context()).Warn("invalid JWT: ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func JwtMiddleware(publicKeyPEM string) func(http.Handler) http.Handler {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		logging.NoCtx().Fatal("invalid JWT public key ", err)
		// never reached
	}

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		},
		SigningMethod:       jwt.SigningMethodRS256,
		CredentialsOptional: true,
		UserProperty:        userProperty,
		Extractor:           jwtmiddleware.FromFirst(jwtmiddleware.FromAuthHeader, fromCookie),
		ErrorHandler:        errorHandler,
	})

	return func(next http.Handler) http.Handler {
		return createHandlerFunction(jwtMiddleware, next)
	}
}

func getToken(ctx context.Context) (*jwt.Token, error) {
	contextValue := ctx.Value(userProperty)

	if contextValue != nil {
		token, ok := contextValue.(*jwt.Token)
		if !ok {
			return nil, fmt.Errorf("token in context is of invalid data type (internal error, probably a re-used context key)")
		}
		return token, nil
	}

	return nil, fmt.Errorf("no token in context")
}

func getTokenClaims(ctx context.Context) (jwt.MapClaims, error) {
	token, err := getToken(ctx)
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}

func GetName(ctx context.Context) (string, error) {
	claims, err := getTokenClaims(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get name: failed to get token claims: %w", err)
	}

	global, ok := claims["global"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to get name: global section not found in claims")
	}

	name, ok := global["name"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get name: name not found in global section in claims")
	}

	return name, nil
}

func HasRole(ctx context.Context, roleName string) (bool, error) {
	claims, err := getTokenClaims(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check for '%s' role: failed to get token claims: %w", roleName, err)
	}

	global, ok := claims["global"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("failed to check for '%s' role: global section not found in claims", roleName)
	}

	roles, ok := global["roles"].([]interface{})
	if !ok {
		return false, fmt.Errorf("failed to check for '%s' role: roles list not found in global section in claims", roleName)
	}

	return contains(roles, roleName), nil
}

func contains(s []interface{}, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
