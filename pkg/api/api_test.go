package api

import (
	"testing"

	"github.com/chris-alexander-pop/system-design-library/pkg/api/rest"
)

func TestNew(t *testing.T) {
	// Test REST factory
	cfg := Config{
		Protocol: ProtocolREST,
		Port:     "8081",
	}

	server, err := New(cfg)
	if err != nil {
		t.Fatalf("New(REST) failed: %v", err)
	}

	if _, ok := server.(*rest.Server); !ok {
		t.Errorf("Expected *rest.Server, got %T", server)
	}

	// Test Invalid Protocol
	_, err = New(Config{Protocol: "unknown"})
	if err == nil {
		t.Error("Expected error for unknown protocol")
	}
}
