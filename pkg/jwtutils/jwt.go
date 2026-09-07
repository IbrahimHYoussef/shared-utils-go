// Package jwtutils provides JWT claim types and helpers for generating and
// validating HMAC-signed tokens.
package jwtutils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenType represents the type of JWT token.
type TokenType string

const (
	// AccessTokenType identifies an access token.
	AccessTokenType TokenType = "access"
	// RefreshTokenType identifies a refresh token.
	RefreshTokenType TokenType = "refresh"
)

// IsValid reports whether t is one of the supported token types.
func (t TokenType) IsValid() bool {
	switch t {
	case AccessTokenType, RefreshTokenType:
		return true
	default:
		return false
	}
}

// String returns the string representation of the token type.
func (t TokenType) String() string {
	return string(t)
}

// Claims contains application-specific JWT claims plus the standard registered
// JWT claims.
type Claims struct {
	// UserName is the authenticated user's display or login name.
	UserName string `json:"user_name"`
	// Email is the authenticated user's email address.
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// JwtService generates and validates HMAC-signed JWTs.
type JwtService struct {
	secretKey []byte
	// AccessTokenExpiry is the duration added to access token expiry times.
	AccessTokenExpiry time.Duration
	// RefreshTokenExpiry is the configured refresh token lifetime.
	RefreshTokenExpiry time.Duration
	issuer             string
}

// NewJwtManager returns a JwtService configured with an HMAC secret, access and
// refresh token expiries in seconds, and an issuer value.
func NewJwtManager(secret string, accessTokenExpirySec int, refreshTokenExpriySec int, issuer string) *JwtService {
	return &JwtService{
		secretKey:          []byte(secret),
		AccessTokenExpiry:  time.Duration(accessTokenExpirySec) * time.Second,
		RefreshTokenExpiry: time.Duration(refreshTokenExpriySec) * time.Second,
		issuer:             issuer,
	}
}

// GenerateToken creates an HS256 access token for userID, email, and userName.
//
// The token subject is userID.String(), the issuer is the service issuer, and
// the expiration is based on AccessTokenExpiry.
func (tm *JwtService) GenerateToken(userID uuid.UUID, email string, userName string) (string, error) {

	userIDStr := userID.String()

	expiry := tm.AccessTokenExpiry

	claims := &Claims{
		UserName: userName,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    tm.issuer,
			Subject:   userIDStr,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.secretKey)
}

// GenerateRefreshSessionToken returns a cryptographically random URL-safe token.
//
// The length argument controls the number of random bytes before base64 URL
// encoding. When length is less than 1, 32 bytes are used.
func (tm *JwtService) GenerateRefreshSessionToken(length int) (string, error) {
	if length <= 0 {
		length = 32 // default length
	}
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// ParseWithClaims parses tokenString into claims using the service HMAC secret.
//
// ParseWithClaims rejects tokens signed with non-HMAC methods.
func (tm *JwtService) ParseWithClaims(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// ValidateToken parses and validates tokenString into Claims.
//
// ValidateToken returns the claims when the token is valid and signed with the
// service HMAC secret.
func (tm *JwtService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
