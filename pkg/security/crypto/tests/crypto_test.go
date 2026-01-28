package tests

import (
	"encoding/base64"
	"testing"

	"github.com/chris-alexander-pop/system-design-library/pkg/security/crypto"
)

func TestAESEncryptor(t *testing.T) {
	key, err := crypto.GenerateAES256Key()
	if err != nil {
		t.Fatalf("GenerateAES256Key failed: %v", err)
	}

	enc, err := crypto.NewAESEncryptor(key)
	if err != nil {
		t.Fatalf("NewAESEncryptor failed: %v", err)
	}

	plaintext := "Hello, World!"
	ciphertext, err := enc.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if len(ciphertext) == 0 {
		t.Fatal("Ciphertext is empty")
	}

	decrypted, err := enc.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Errorf("Expected %q, got %q", plaintext, string(decrypted))
	}
}

func TestAESEncryptor_String(t *testing.T) {
	key, err := crypto.GenerateAES256Key()
	if err != nil {
		t.Fatalf("GenerateAES256Key failed: %v", err)
	}

	enc, err := crypto.NewAESEncryptor(key)
	if err != nil {
		t.Fatalf("NewAESEncryptor failed: %v", err)
	}

	plaintext := "Sensitive Data"
	encoded, err := enc.EncryptString(plaintext)
	if err != nil {
		t.Fatalf("EncryptString failed: %v", err)
	}

	if _, err := base64.StdEncoding.DecodeString(encoded); err != nil {
		t.Fatalf("EncryptString returned invalid base64: %v", err)
	}

	decrypted, err := enc.DecryptString(encoded)
	if err != nil {
		t.Fatalf("DecryptString failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Expected %q, got %q", plaintext, decrypted)
	}
}
