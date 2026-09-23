package middelware_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jwtSer "github.com/IbrahimHYoussef/shared-utils-go/pkg/jwtutils"
	"github.com/IbrahimHYoussef/shared-utils-go/pkg/middelware"
	"github.com/golang-jwt/jwt/v5"
)

const (
	testKey    = "test-secret-0123456789abcdef0123"
	testIssuer = "test-issuer"
)

func sign(t *testing.T, method jwt.SigningMethod, key any, claims jwt.Claims) string {
	t.Helper()
	token, err := jwt.NewWithClaims(method, claims).SignedString(key)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return token
}

func claims(issuer string, expiresIn time.Duration) *jwtSer.Claims {
	return &jwtSer.Claims{
		Email: "a@b.c",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-1",
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		},
	}
}

// serve runs one request through middleware and returns the status, body and
// the subject the next handler saw ("" when it was not called).
func serve(mw func(http.Handler) http.Handler, header string) (int, string, string) {
	subject := ""
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Context().Value(middelware.UserClaimsKey).(*jwtSer.Claims)
		if c != nil {
			subject = c.Subject
		}
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if header != "" {
		req.Header.Set("Authorization", header)
	}
	rec := httptest.NewRecorder()
	mw(next).ServeHTTP(rec, req)
	body, _ := io.ReadAll(rec.Body)
	return rec.Code, string(body), subject
}

func TestAuthMiddleware(t *testing.T) {
	middelware.SetLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))

	valid := sign(t, jwt.SigningMethodHS256, []byte(testKey), claims(testIssuer, time.Hour))
	noExp := &jwtSer.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1", Issuer: testIssuer}}

	cases := []struct {
		name           string
		header         string
		fromServiceOK  bool
		plainFactoryOK bool
		wantExpired    bool
	}{
		{name: "valid HS256", header: "Bearer " + valid, fromServiceOK: true, plainFactoryOK: true},
		{name: "missing header", header: ""},
		{name: "short header does not panic", header: "Bear"},
		{name: "not a Bearer header", header: "Basic " + valid},
		{name: "empty Bearer token", header: "Bearer "},
		{name: "alg none", header: "Bearer " + sign(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, claims(testIssuer, time.Hour))},
		{name: "HS512 with the same key", header: "Bearer " + sign(t, jwt.SigningMethodHS512, []byte(testKey), claims(testIssuer, time.Hour))},
		{name: "wrong key", header: "Bearer " + sign(t, jwt.SigningMethodHS256, []byte("another-secret"), claims(testIssuer, time.Hour))},
		{name: "missing exp", header: "Bearer " + sign(t, jwt.SigningMethodHS256, []byte(testKey), noExp)},
		{name: "expired", header: "Bearer " + sign(t, jwt.SigningMethodHS256, []byte(testKey), claims(testIssuer, -time.Minute)), wantExpired: true},
		// Only the service-based factory knows the issuer.
		{name: "wrong issuer", header: "Bearer " + sign(t, jwt.SigningMethodHS256, []byte(testKey), claims("someone-else", time.Hour)), plainFactoryOK: true},
	}

	fromService := middelware.AuthMiddleWareFactoryFromService(jwtSer.NewJwtManager(testKey, 3600, 86400, testIssuer))
	plain := middelware.AuthMiddleWareFactory(testKey)

	for _, tc := range cases {
		for _, f := range []struct {
			name   string
			mw     func(http.Handler) http.Handler
			wantOK bool
		}{
			{"FromService", fromService, tc.fromServiceOK},
			{"Factory", plain, tc.plainFactoryOK},
		} {
			t.Run(f.name+"/"+tc.name, func(t *testing.T) {
				code, body, subject := serve(f.mw, tc.header)
				if f.wantOK {
					if code != http.StatusOK || subject != "user-1" {
						t.Fatalf("want 200 with claims in context, got %d subject=%q body=%s", code, subject, body)
					}
					return
				}
				if code != http.StatusUnauthorized || subject != "" {
					t.Fatalf("want 401 without reaching next, got %d subject=%q", code, subject)
				}
				if tc.wantExpired && !strings.Contains(body, "Token Expired") {
					t.Fatalf("want Token Expired message, got %s", body)
				}
			})
		}
	}
}
