package pqc

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/curve25519"
)

// X25519KEM implements classical ECDH using Curve25519.
type X25519KEM struct{}

func (k *X25519KEM) KeyGen() (publicKey, privateKey []byte, err error) {
	privateKey = make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, privateKey); err != nil {
		return nil, nil, err
	}

	publicKey, err = curve25519.X25519(privateKey, curve25519.Basepoint)
	if err != nil {
		return nil, nil, err
	}

	return publicKey, privateKey, nil
}

func (k *X25519KEM) Encapsulate(publicKey []byte) (sharedSecret, ciphertext []byte, err error) {
	// Generate ephemeral keypair
	ephemeralPriv := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, ephemeralPriv); err != nil {
		return nil, nil, err
	}

	ephemeralPub, err := curve25519.X25519(ephemeralPriv, curve25519.Basepoint)
	if err != nil {
		return nil, nil, err
	}

	// Compute shared secret
	sharedSecret, err = curve25519.X25519(ephemeralPriv, publicKey)
	if err != nil {
		return nil, nil, err
	}

	return sharedSecret, ephemeralPub, nil
}

func (k *X25519KEM) Decapsulate(privateKey, ciphertext []byte) (sharedSecret []byte, err error) {
	return curve25519.X25519(privateKey, ciphertext)
}

func (k *X25519KEM) PublicKeySize() int    { return 32 }
func (k *X25519KEM) CiphertextSize() int   { return 32 }
func (k *X25519KEM) SharedSecretSize() int { return 32 }
