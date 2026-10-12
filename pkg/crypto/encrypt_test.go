package crypto_test

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/crypto"
)

func TestEncryptRoundTrip(t *testing.T) {
	key, err := crypto.ParseEncryptionKey(base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef")))
	if err != nil {
		t.Fatal(err)
	}
	sealed, err := crypto.Encrypt(key, "EAAG...token")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(sealed, "v1:") || strings.Contains(sealed, "token") {
		t.Fatalf("ciphertext: %q", sealed)
	}
	opened, err := crypto.Decrypt(key, sealed)
	if err != nil || opened != "EAAG...token" {
		t.Fatalf("round trip: %q %v", opened, err)
	}
	again, _ := crypto.Encrypt(key, "EAAG...token")
	if again == sealed {
		t.Error("two encryptions of one text are equal: the nonce is not random")
	}

	other, _ := crypto.ParseEncryptionKey(base64.StdEncoding.EncodeToString([]byte("fedcba9876543210fedcba9876543210")))
	if _, err := crypto.Decrypt(other, sealed); err == nil {
		t.Error("another key opened the ciphertext")
	}
	tampered := sealed[:len(sealed)-2] + "AA"
	if _, err := crypto.Decrypt(key, tampered); err == nil {
		t.Error("a changed ciphertext opened")
	}
	if _, err := crypto.ParseEncryptionKey("dG9vIHNob3J0"); err == nil {
		t.Error("a short key was accepted")
	}
}
