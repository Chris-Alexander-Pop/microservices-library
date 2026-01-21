package captcha

import (
	"context"
)

// Config configures the captcha system.
type Config struct {
	// Provider specifies which captcha provider to use (memory, recaptcha, turnstile, etc.).
	Provider string `env:"SECURITY_CAPTCHA_PROVIDER" env-default:"memory"`
	// SecretKey is the server-side secret for verification.
	SecretKey string `env:"SECURITY_CAPTCHA_SECRET"`
}

// Verifier defines the interface for captcha verification.
type Verifier interface {
	Verify(ctx context.Context, token string) error
}
