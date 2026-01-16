package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

// Load reads configuration from .env file or environment variables and validates it.
func Load[T any](cfg *T) error {
	// 1. Load from .env if it exists
	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		// If .env doesn't exist or we just want to rely on env vars,
		// we fallback to ReadEnv to pick up environment variables processing.
		// cleanenv.ReadConfig already does ReadEnv if file fails?
		// Actually cleanenv.ReadConfig returns error if file not found.
		// So we fallback to ReadEnv.
		if err := cleanenv.ReadEnv(cfg); err != nil {
			return fmt.Errorf("failed to read env config: %w", err)
		}
	}

	// 2. Validate the struct
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}
