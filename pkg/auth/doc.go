/*
Package auth provides authentication and authorization primitives.

This package defines the common `Verifier` interface and `Claims` structure used across different authentication strategies (JWT, OIDC, etc.).

Usage:

	import "github.com/chris-alexander-pop/system-design-library/pkg/auth"
	import "github.com/chris-alexander-pop/system-design-library/pkg/auth/adapters/jwt"

	verifier := jwt.New(jwt.Config{Secret: "secret"})
	claims, err := verifier.Verify(ctx, tokenString)
*/
package auth
