package middelware

import (
	"context"
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

// AuthMiddleWareFactoryFromService returns middleware that authenticates requests
// with jwtService.
//
// The middleware expects an Authorization header containing a Bearer token,
// validates the token into jwtutils.Claims, and stores those claims in the
// request context under UserClaimsKey.
func AuthMiddleWareFactoryFromService(jwtService *jwtSer.JwtService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth_header := r.Header.Get("Authorization")
			if len(auth_header) == 0 {
				logger.Error("Missing Authorization header")
				responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
				return
			}

			tokenString := auth_header[7:]

			claims := &jwtSer.Claims{}

			token, err := jwtService.ParseWithClaims(tokenString, claims)
			if err != nil {
				if err == jwt.ErrTokenSignatureInvalid {
					responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
					logger.Error("Invalid token signature", "error", err)
					return
				}
				if err == jwt.ErrTokenMalformed {
					responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
					logger.Error("Malformed token", "error", err)
					return
				}
				if err == jwt.ErrTokenExpired {
					responses.RespondWithError(w, http.StatusUnauthorized, TOKEN_EXPIRE)
					logger.Error("Expired token", "error", err)
					return
				}

				logger.Error("Error parsing token", "error", err)
				responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
				return
			}
			if !token.Valid {
				responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
				logger.Error("Invalid token", "error", err)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthMiddleWareFactory returns middleware that authenticates requests with a
// JWT HMAC secret.
//
// The middleware requires an Authorization header in the form "Bearer <token>".
// On success, it stores the parsed token claims in the request context under
// UserClaimsKey.
func AuthMiddleWareFactory(JWTKey string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth_header := r.Header.Get("Authorization")
			if len(auth_header) == 0 {
				logger.Error("Missing Authorization header")
				responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
				return
			}
			if !strings.HasPrefix(auth_header, "Bearer ") || len(auth_header) <= 7 {
				logger.Error("Auth Header malformed", "header", auth_header)
				responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
				return
			}
			tokenString := auth_header[7:]

			// Initialize the generic type properly
			// var claims T
			// Create a new instance of the type
			// claimsPtr := new(T)
			// claims = *claimsPtr
			claims := &jwtSer.Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(JWTKey), nil
			})
			if err != nil {
				if err == jwt.ErrTokenSignatureInvalid {
					responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
					logger.Error("Invalid token signature", "error", err)
					return
				}
				if err == jwt.ErrTokenMalformed {
					responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
					logger.Error("Malformed token", "error", err)
					return
				}
				if err == jwt.ErrTokenExpired {
					responses.RespondWithError(w, http.StatusUnauthorized, TOKEN_EXPIRE)
					logger.Error("Expired token", "error", err)
					return
				}
				logger.Error("Error parsing token", "error", err)
				responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
				return
			}
			if !token.Valid {
				responses.RespondWithError(w, http.StatusUnauthorized, NOT_ALLOWED)
				logger.Error("Invalid token", "error", err)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, token.Claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
