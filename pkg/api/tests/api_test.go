package api_test

import (
	"github.com/chris-alexander-pop/system-design-library/pkg/api"
	"testing"

	"github.com/chris-alexander-pop/system-design-library/pkg/api/rest"
)

func TestNew(t *testing.T) {
	// Test REST factory
	cfg := api.Config{
		Protocol: api.ProtocolREST,
		Port:     "8081",
	}

	server, err := api.New(cfg)
	if err != nil {
		t.Fatalf("api.New(REST) failed: %v", err)
	}

	if _, ok := server.(*rest.Server); !ok {
		t.Errorf("Expected *rest.Server, got %T", server)
	}

	// Test Invalid api.Protocol
	_, err = api.New(api.Config{Protocol: "unknown"})
	if err == nil {
		t.Error("Expected error for unknown protocol")
	}
}
