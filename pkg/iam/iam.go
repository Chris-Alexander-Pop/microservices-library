package iam

import "time"

// User represents an identity within the system.
type User struct {
	ID        string            `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Roles     []string          `json:"roles"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Token represents an issued authentication token.
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	IDToken      string    `json:"id_token,omitempty"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"` // Seconds
	IssuedAt     time.Time `json:"issued_at"`
}

// Credentials represents input for authentication.
type Credentials struct {
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	GrantType string `json:"grant_type,omitempty"` // password, refresh_token, etc.
}
