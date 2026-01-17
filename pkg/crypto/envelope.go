package crypto

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
)

// EnvelopeEncryption implements the envelope encryption pattern.
// Data is encrypted with a DEK (Data Encryption Key), and the DEK
// is encrypted with a KEK (Key Encryption Key) from a KMS.
//
// This pattern allows:
// - Fast local encryption with AES
// - Secure key management via KMS
// - Key rotation without re-encrypting all data
type EnvelopeEncryption struct {
	kms KeyProvider
}

// EnvelopePayload contains encrypted data and its encrypted DEK.
type EnvelopePayload struct {
	EncryptedData string `json:"encrypted_data"` // Base64-encoded encrypted data
	EncryptedDEK  string `json:"encrypted_dek"`  // Base64-encoded KMS-encrypted DEK
	KeyID         string `json:"key_id"`         // KMS key ID used
	Algorithm     string `json:"algorithm"`      // Encryption algorithm
}

// NewEnvelopeEncryption creates a new envelope encryptor.
func NewEnvelopeEncryption(kms KeyProvider) *EnvelopeEncryption {
	return &EnvelopeEncryption{kms: kms}
}

// Encrypt encrypts data using envelope encryption.
// 1. Generate a DEK from KMS
// 2. Encrypt data with DEK using AES-GCM
// 3. Return encrypted data + encrypted DEK
func (e *EnvelopeEncryption) Encrypt(ctx context.Context, plaintext []byte) (*EnvelopePayload, error) {
	// Generate a data encryption key from KMS
	dek, encryptedDEK, keyID, err := e.kms.GenerateDataKey(ctx)
	if err != nil {
		return nil, err
	}

	// Encrypt data with DEK
	encryptor, err := NewAESEncryptor(dek)
	if err != nil {
		return nil, err
	}

	ciphertext, err := encryptor.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}

	// Clear DEK from memory
	for i := range dek {
		dek[i] = 0
	}

	return &EnvelopePayload{
		EncryptedData: base64.StdEncoding.EncodeToString(ciphertext),
		EncryptedDEK:  base64.StdEncoding.EncodeToString(encryptedDEK),
		KeyID:         keyID,
		Algorithm:     "AES-256-GCM",
	}, nil
}

// Decrypt decrypts envelope-encrypted data.
// 1. Decrypt DEK using KMS
// 2. Decrypt data with DEK
func (e *EnvelopeEncryption) Decrypt(ctx context.Context, payload *EnvelopePayload) ([]byte, error) {
	// Decode encrypted DEK
	encryptedDEK, err := base64.StdEncoding.DecodeString(payload.EncryptedDEK)
	if err != nil {
		return nil, err
	}

	// Decrypt DEK using KMS
	dek, err := e.kms.DecryptDataKey(ctx, encryptedDEK, payload.KeyID)
	if err != nil {
		return nil, err
	}
	defer func() {
		for i := range dek {
			dek[i] = 0
		}
	}()

	// Decode ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(payload.EncryptedData)
	if err != nil {
		return nil, err
	}

	// Decrypt data
	encryptor, err := NewAESEncryptor(dek)
	if err != nil {
		return nil, err
	}

	return encryptor.Decrypt(ciphertext)
}

// EncryptToJSON encrypts and returns JSON-serialized envelope payload.
func (e *EnvelopeEncryption) EncryptToJSON(ctx context.Context, plaintext []byte) ([]byte, error) {
	payload, err := e.Encrypt(ctx, plaintext)
	if err != nil {
		return nil, err
	}
	return json.Marshal(payload)
}

// DecryptFromJSON decrypts JSON-serialized envelope payload.
func (e *EnvelopeEncryption) DecryptFromJSON(ctx context.Context, data []byte) ([]byte, error) {
	var payload EnvelopePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return e.Decrypt(ctx, &payload)
}

// =========================================================================
// In-Memory Key Provider (for testing/development)
// =========================================================================

// MemoryKeyProvider is an in-memory key provider for testing.
// DO NOT use in production - keys should come from a real KMS.
type MemoryKeyProvider struct {
	masterKey []byte
}

// NewMemoryKeyProvider creates a new in-memory key provider.
func NewMemoryKeyProvider(masterKey []byte) (*MemoryKeyProvider, error) {
	if len(masterKey) != 32 {
		return nil, errors.New("master key must be 32 bytes")
	}
	return &MemoryKeyProvider{masterKey: masterKey}, nil
}

func (m *MemoryKeyProvider) GetKey(ctx context.Context, keyID string) ([]byte, error) {
	return m.masterKey, nil
}

func (m *MemoryKeyProvider) GenerateDataKey(ctx context.Context) ([]byte, []byte, string, error) {
	// Generate random DEK
	dek, err := GenerateAES256Key()
	if err != nil {
		return nil, nil, "", err
	}

	// "Encrypt" DEK with master key (simplified - real KMS uses proper wrapping)
	encryptor, err := NewAESEncryptor(m.masterKey)
	if err != nil {
		return nil, nil, "", err
	}

	encryptedDEK, err := encryptor.Encrypt(dek)
	if err != nil {
		return nil, nil, "", err
	}

	return dek, encryptedDEK, "memory-key-1", nil
}

func (m *MemoryKeyProvider) DecryptDataKey(ctx context.Context, encryptedKey []byte, keyID string) ([]byte, error) {
	encryptor, err := NewAESEncryptor(m.masterKey)
	if err != nil {
		return nil, err
	}

	return encryptor.Decrypt(encryptedKey)
}
