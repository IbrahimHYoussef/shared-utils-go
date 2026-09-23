package jwtutils_test

import (
	"testing"
	"time"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/jwtutils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	testKey    = "test-secret-0123456789abcdef0123"
	testIssuer = "test-issuer"
)

func signed(t *testing.T, method jwt.SigningMethod, issuer string, exp *jwt.NumericDate) string {
	t.Helper()
	token, err := jwt.NewWithClaims(method, &jwtutils.Claims{RegisteredClaims: jwt.RegisteredClaims{
		Subject:   "user-1",
		Issuer:    issuer,
		ExpiresAt: exp,
	}}).SignedString([]byte(testKey))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return token
}

func TestValidateToken(t *testing.T) {
	svc := jwtutils.NewJwtManager(testKey, 3600, 86400, testIssuer)
	inAnHour := jwt.NewNumericDate(time.Now().Add(time.Hour))

	generated, err := svc.GenerateToken(uuid.New(), "a@b.c", "a")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	if _, err := svc.ValidateToken(generated); err != nil {
		t.Fatalf("token from GenerateToken rejected: %v", err)
	}

	rejected := map[string]string{
		"HS512":        signed(t, jwt.SigningMethodHS512, testIssuer, inAnHour),
		"wrong issuer": signed(t, jwt.SigningMethodHS256, "someone-else", inAnHour),
		"no issuer":    signed(t, jwt.SigningMethodHS256, "", inAnHour),
		"missing exp":  signed(t, jwt.SigningMethodHS256, testIssuer, nil),
		"expired":      signed(t, jwt.SigningMethodHS256, testIssuer, jwt.NewNumericDate(time.Now().Add(-time.Minute))),
	}
	for name, token := range rejected {
		if _, err := svc.ValidateToken(token); err == nil {
			t.Errorf("%s: ValidateToken accepted the token", name)
		}
		if _, err := svc.ParseWithClaims(token, &jwtutils.Claims{}); err == nil {
			t.Errorf("%s: ParseWithClaims accepted the token", name)
		}
	}
}

func TestEmptyIssuerSkipsIssuerCheck(t *testing.T) {
	svc := jwtutils.NewJwtManager(testKey, 3600, 86400, "")
	token := signed(t, jwt.SigningMethodHS256, "anyone", jwt.NewNumericDate(time.Now().Add(time.Hour)))
	if _, err := svc.ValidateToken(token); err != nil {
		t.Fatalf("service without issuer should not check iss: %v", err)
	}
}
