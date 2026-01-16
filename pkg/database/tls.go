package database

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
)

// LoadTLSConfig helper to create a tls.Config with custom CA if provided.
// Returns nil, nil if no SSL/TLS is configured.
func LoadTLSConfig(cfg Config) (*tls.Config, error) {
	if cfg.SSLMode == "disable" || cfg.SSLMode == "" {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// 1. Insecure Skip Verify (Use with caution)
	if cfg.SSLMode == "skip-verify" {
		tlsConfig.InsecureSkipVerify = true
	}

	// 2. Load Custom CA if provided
	if cfg.SSLRootCert != "" {
		caCert, err := os.ReadFile(cfg.SSLRootCert)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read ssl root cert")
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return nil, errors.New(errors.CodeInternal, "failed to parse ssl root cert", nil)
		}
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}
