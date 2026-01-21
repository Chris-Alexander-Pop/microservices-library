// Package mfa provides Multi-Factor Authentication capabilities.
//
// This package supports various MFA methods including:
//   - TOTP (Time-based One-Time Password)
//   - SMS/Email OTP (via communication package)
//   - Recovery codes
//
// Usage:
//
//	mfaService := memory.New()
//	qrCode, secret, err := mfaService.Enroll(ctx, userID, mfa.TOTP)
//	err = mfaService.Verify(ctx, userID, code)
package mfa
