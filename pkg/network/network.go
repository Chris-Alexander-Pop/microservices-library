package network

import (
	"time"
)

// Config holds generic network server configuration.
type Config struct {
	// Addr is the address to listen on (e.g., ":9000", "localhost:8080").
	Addr string `env:"NETWORK_ADDR" env-default:":9000"`

	// ReadTimeout is the timeout for reading from connections.
	ReadTimeout time.Duration `env:"NETWORK_READ_TIMEOUT" env-default:"30s"`

	// WriteTimeout is the timeout for writing to connections.
	WriteTimeout time.Duration `env:"NETWORK_WRITE_TIMEOUT" env-default:"30s"`
}
