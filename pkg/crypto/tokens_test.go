package crypto_test

import (
	"encoding/base64"
	"testing"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/crypto"
)

func TestGenerateSecretToken(t *testing.T) {
	token, err := crypto.GenerateSecretToken(32)
	if err != nil {
		t.Fatalf("GenerateSecretToken: %v", err)
	}
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(raw) != 32 {
		t.Fatalf("want 32 bytes of unpadded base64url, got %q (%d bytes, err %v)", token, len(raw), err)
	}

	other, _ := crypto.GenerateSecretToken(32)
	if other == token {
		t.Fatal("two tokens are equal")
	}

	def, _ := crypto.GenerateSecretToken(0)
	if raw, _ := base64.RawURLEncoding.DecodeString(def); len(raw) != crypto.DefaultSecretTokenBytes {
		t.Fatalf("numBytes 0 should use %d bytes, got %d", crypto.DefaultSecretTokenBytes, len(raw))
	}
}

func TestHashToken(t *testing.T) {
	// SHA-256("abc") test vector.
	want := "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if got := crypto.HashToken("abc"); got != want {
		t.Fatalf("HashToken(abc) = %s, want %s", got, want)
	}
}

func TestCheckTokenHash(t *testing.T) {
	hash := crypto.HashToken("123456")
	if !crypto.CheckTokenHash("123456", hash) {
		t.Fatal("matching token rejected")
	}
	for _, wrong := range []string{"123457", "", "1234567"} {
		if crypto.CheckTokenHash(wrong, hash) {
			t.Fatalf("%q accepted", wrong)
		}
	}
	if crypto.CheckTokenHash("123456", "") {
		t.Fatal("empty hash accepted")
	}
}
