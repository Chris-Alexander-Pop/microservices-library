package dhcp

import (
	"context"
	"time"
)

// IPAM (IP Address Management) manages the allocation of IP addresses within pools.
type IPAM interface {
	// AllocateIP allocates an available IP from the given subnet/pool.
	AllocateIP(ctx context.Context, poolID string) (*IPAllocation, error)

	// ReserveIP reserves a specific IP address.
	ReserveIP(ctx context.Context, poolID string, ip string) error

	// ReleaseIP releases an allocated IP back to the pool.
	ReleaseIP(ctx context.Context, allocationID string) error
}

// IPAllocation represents an allocated IP address.
type IPAllocation struct {
	ID          string    `json:"id"`
	PoolID      string    `json:"pool_id"`
	IP          string    `json:"ip"`
	AllocatedAt time.Time `json:"allocated_at"`
	OwnerID     string    `json:"owner_id,omitempty"` // e.g., VM ID or Interface ID
}

// Config holds configuration for the DHCP/IPAM service.
type Config struct {
	// Driver specifies the IPAM backend: "memory", "phpipam", "netbox".
	Driver string `env:"IPAM_DRIVER" env-default:"memory"`
}
