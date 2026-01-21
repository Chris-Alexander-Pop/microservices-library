package secrets

import (
	"context"
)

// Config configures the Secret Manager.
type Config struct {
	// Provider specifies the secrets provider (memory, aws-secrets-manager, vault).
	Provider string `env:"SECURITY_SECRETS_PROVIDER" env-default:"memory"`
}

// SecretManager defines the interface for secrets management.
type SecretManager interface {
	Get(ctx context.Context, name string) (string, error)
	Set(ctx context.Context, name, value string) error
}
