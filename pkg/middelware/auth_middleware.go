package middelware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/errorsutil"
	jwtSer "github.com/IbrahimHYoussef/shared-utils-go/pkg/jwtutils"
	"github.com/IbrahimHYoussef/shared-utils-go/pkg/responses"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

// UserClaimsKey is the context key used to store parsed JWT claims.
const UserClaimsKey contextKey = "claims"

// NOT_ALLOWED is the default unauthorized error returned by authentication
// middleware when access is denied.
var NOT_ALLOWED = errorsutil.NewUnauthorizedError("Not Allowed To Access This Endpoint")

// TOKEN_EXPIRE is the default unauthorized error returned when a JWT is expired.
var TOKEN_EXPIRE = errorsutil.NewUnauthorizedError("Token Expired")

const bearerPrefix = "Bearer "

// parseFunc parses a raw token into claims and reports whether it is valid.
type parseFunc func(tokenString string, claims *jwtSer.Claims) (*jwt.Token, error)

// AuthMiddleWareFactoryFromService returns middleware that authenticates requests
// with jwtService.
//
// The middleware expects an Authorization header containing a Bearer token. It
// accepts only HS256 tokens with an exp claim whose iss claim matches the
// service issuer, and stores the *jwtutils.Claims in the request context under
// UserClaimsKey. Prefer this factory: it is the only one that checks the issuer.
func AuthMiddleWareFactoryFromService(jwtService *jwtSer.JwtService) func(http.Handler) http.Handler {
	return authMiddleware(func(tokenString string, claims *jwtSer.Claims) (*jwt.Token, error) {
		return jwtService.ParseWithClaims(tokenString, claims)
	})
}

// AuthMiddleWareFactory returns middleware that authenticates requests with a
// JWT HMAC secret.
//
// The middleware requires an Authorization header in the form "Bearer <token>".
// It accepts only HS256 tokens with an exp claim. It cannot check the issuer;
// use AuthMiddleWareFactoryFromService for that. On success, it stores the
// *jwtutils.Claims in the request context under UserClaimsKey.
func AuthMiddleWareFactory(JWTKey string) func(http.Handler) http.Handler {
	key := []byte(JWTKey)
	return authMiddleware(func(tokenString string, claims *jwtSer.Claims) (*jwt.Token, error) {
		return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		}, jwtSer.ParseOptions("")...)
	})
}

func authMiddleware(parse parseFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Error("Missing Authorization header")
				responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
				return
			}
			if !strings.HasPrefix(authHeader, bearerPrefix) || len(authHeader) <= len(bearerPrefix) {
				// The header is not logged: it may contain a credential.
				logger.Error("Authorization header is not a Bearer token")
				responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
				return
			}

			claims := &jwtSer.Claims{}
			token, err := parse(authHeader[len(bearerPrefix):], claims)
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				logger.Error("Expired token", "error", err)
				responses.RespondWithError(w, http.StatusUnauthorized, TOKEN_EXPIRE)
				return
			case err != nil:
				logger.Error("Rejected token", "error", err)
				responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
				return
			case !token.Valid:
				logger.Error("Invalid token")
				responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
