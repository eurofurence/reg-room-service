package jwt

import (
	"context"
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/eurofurence/reg-room-service/internal/repository/logging"
	jwt "github.com/form3tech-oss/jwt-go"
)

const userProperty = "user"

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
	})

	return func(next http.Handler) http.Handler {
		return createHandlerFunction(jwtMiddleware, next)
	}
}

func getUserInformation(ctx context.Context) (*jwt.Token, error) {
	contextValue := ctx.Value(userProperty)

	if contextValue != nil {
		return contextValue.(*jwt.Token), nil
	}

	return nil, fmt.Errorf("no user in context")
}

func GetName(ctx context.Context) (string, error) {
	token, err := getUserInformation(ctx)

	if err != nil {
		return "", fmt.Errorf("failed to get name: %w", err)
	}

	return token.Claims.(jwt.MapClaims)["name"].(string), nil
}

func IsAdmin(ctx context.Context) (bool, error) {
	token, err := getUserInformation(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to get admin status: %w", err)
	}

	return token.Claims.(jwt.MapClaims)["admin"].(bool), nil
}
