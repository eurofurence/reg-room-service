package jwt

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/eurofurence/reg-room-service/internal/repository/logging"
	jwt "github.com/form3tech-oss/jwt-go"
	"net/http"
)

func createHandlerFunction(jwtMiddleware *jwtmiddleware.JWTMiddleware, next http.Handler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := jwtMiddleware.CheckJWT(w, r)
		if err != nil {
			// missing JWT is not considered an error
			logging.Ctx(r.Context()).Warn("invalid JWT: ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// XXX TODO: store authorization header in ctx

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
	})

	return func(next http.Handler) http.Handler {
		return createHandlerFunction(jwtMiddleware, next)
	}
}

// XXX TODO:  add accessor methods for the various JWT fields
