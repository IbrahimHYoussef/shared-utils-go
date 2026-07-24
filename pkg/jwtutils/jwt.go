package jwtutils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenType represents the type of JWT token
type TokenType string

const (
	AccessTokenType  TokenType = "access"
	RefreshTokenType TokenType = "refresh"
)

// IsValid checks if the token type is valid
func (t TokenType) IsValid() bool {
	switch t {
	case AccessTokenType, RefreshTokenType:
		return true
	default:
		return false
	}
}

// String returns the string representation of the token type
func (t TokenType) String() string {
	return string(t)
}

type Claims struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type JwtService struct {
	secretKey          []byte
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	issuer             string
}

func NewJwtManager(secret string, accessTokenExpirySec int, refreshTokenExpriySec int, issuer string) *JwtService {
	return &JwtService{
		secretKey:          []byte(secret),
		AccessTokenExpiry:  time.Duration(accessTokenExpirySec) * time.Second,
		RefreshTokenExpiry: time.Duration(refreshTokenExpriySec) * time.Second,
		issuer:             issuer,
	}
}

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
