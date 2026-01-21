package memory

import (
	"context"

	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
)

// Verifier implements captcha.Verifier using simple memory checks.
type Verifier struct {
	validToken string
}

// New creates a new memory captcha verifier.
// It accepts a magic token that is considered valid. All others are invalid.
// Defaults to "valid-token" if empty.
func New(magicToken string) *Verifier {
	if magicToken == "" {
		magicToken = "valid-token"
	}
	return &Verifier{
		validToken: magicToken,
	}
}

func (v *Verifier) Verify(ctx context.Context, token string) error {
	if token == "" {
		return errors.InvalidArgument("captcha token missing", nil)
	}
	if token != v.validToken {
		return errors.Forbidden("invalid captcha token", nil)
	}
	return nil
}
