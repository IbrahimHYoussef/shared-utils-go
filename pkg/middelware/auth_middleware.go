package middelware

import (
	"context"
	"net/http"
	"strings"

	responses "github.com/IbrahimHYoussef/project-management/shared-utils-go/Responses"
	"github.com/IbrahimHYoussef/project-management/shared-utils-go/errorsutil"
	jwtSer "github.com/IbrahimHYoussef/project-management/shared-utils-go/jwtutils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserClaimsKey contextKey = "claims"

var NOT_ALLOWED = errorsutil.NewUnauthorizedError("Not Allowed To Access This Endpoint")
var TOKEN_EXPIRE = errorsutil.NewUnauthorizedError("Token Expired")

// this is the middleware that will authenticate the user using the JWT token
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
