package crypto

import (
	"encoding/hex"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := hex.EncodeToString([]byte("0123456789abcdef0123456789abcdef")) // 32 bytes → 64 hex chars
	plain := "my-secret-password"

	enc, err := Encrypt(key, plain)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if len(enc) == 0 {
		t.Fatal("encrypted output is empty")
	}

	dec, err := Decrypt(key, enc)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if dec != plain {
		t.Fatalf("got %q, want %q", dec, plain)
	}
}

func TestEncryptEmpty(t *testing.T) {
	key := hex.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	enc, err := Encrypt(key, "")
	if err != nil || enc != nil {
		t.Fatalf("empty plaintext should return nil, nil; got %v, %v", enc, err)
	}
	dec, err := Decrypt(key, nil)
	if err != nil || dec != "" {
		t.Fatalf("nil ciphertext should return empty; got %q, %v", dec, err)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1 := hex.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	key2 := hex.EncodeToString([]byte("fedcba9876543210fedcba9876543210"))

	enc, _ := Encrypt(key1, "secret")
	// Wrong key → fallback to plaintext (returns raw bytes as string)
	dec, err := Decrypt(key2, enc)
	if err != nil {
		t.Fatalf("expected fallback, got error: %v", err)
	}
	// Fallback returns the raw ciphertext as string (not the original plaintext)
	if dec == "secret" {
		t.Fatal("should not decrypt correctly with wrong key")
	}
}

func TestDecryptLegacyPlaintext(t *testing.T) {
	key := hex.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	// Simulate legacy data: stored as plain bytes without encryption
	legacy := []byte("my-smtp-password")
	dec, err := Decrypt(key, legacy)
	if err != nil {
		t.Fatalf("legacy decrypt error: %v", err)
	}
	if dec != "my-smtp-password" {
		t.Fatalf("got %q, want %q", dec, "my-smtp-password")
	}
}
