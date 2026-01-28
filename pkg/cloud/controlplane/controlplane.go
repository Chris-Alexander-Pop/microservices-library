package controlplane

import (
	"context"

	"github.com/chris-alexander-pop/system-design-library/pkg/cloud"
)

// ControlPlane serves as the central brain of the private cloud.
// It tracks the state of all hosts and instances.
type ControlPlane interface {
	// RegisterHost adds a new host to the cluster.
	RegisterHost(ctx context.Context, host cloud.Host) error

	// DeregisterHost removes a host from the cluster.
	DeregisterHost(ctx context.Context, hostID string) error

	// UpdateHostStatus updates the status of a known host.
	UpdateHostStatus(ctx context.Context, hostID string, status cloud.HostStatus) error

	// GetHost retrieves a host by ID.
	GetHost(ctx context.Context, hostID string) (*cloud.Host, error)

	// ListHosts returns a list of all registered hosts.
	ListHosts(ctx context.Context) ([]cloud.Host, error)
}

// Config holds configuration for the Control Plane.
type Config struct {
	// Driver specifies the storage backend for state: "memory", "etcd", "postgres".
	Driver string `env:"CONTROLPLANE_DRIVER" env-default:"memory"`
}
