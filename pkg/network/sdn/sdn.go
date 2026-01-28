package sdn

import (
	"context"
)

// NetworkManager defines the interface for managing Software Defined Networks (VPCs).
type NetworkManager interface {
	// CreateNetwork creates a new isolated virtual network (VPC).
	CreateNetwork(ctx context.Context, spec NetworkSpec) (string, error)

	// DeleteNetwork removes a virtual network.
	DeleteNetwork(ctx context.Context, networkID string) error

	// CreateSubnet adds a subnet to a network.
	CreateSubnet(ctx context.Context, networkID string, spec SubnetSpec) (string, error)

	// DeleteSubnet removes a subnet.
	DeleteSubnet(ctx context.Context, subnetID string) error

	// GetNetwork retrieves network details.
	GetNetwork(ctx context.Context, networkID string) (*Network, error)
}

// NetworkSpec defines the properties of a virtual network.
type NetworkSpec struct {
	Name string            `json:"name"`
	CIDR string            `json:"cidr"` // e.g. "10.0.0.0/16"
	Tags map[string]string `json:"tags,omitempty"`
}

// Network represents a virtual network (VPC).
type Network struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	CIDR    string            `json:"cidr"`
	Subnets []Subnet          `json:"subnets"`
	Tags    map[string]string `json:"tags,omitempty"`
}

// SubnetSpec defines the properties of a subnet.
type SubnetSpec struct {
	Name string            `json:"name"`
	CIDR string            `json:"cidr"` // e.g. "10.0.1.0/24"
	Zone string            `json:"zone"`
	Tags map[string]string `json:"tags,omitempty"`
}

// Subnet represents a segment of a virtual network.
type Subnet struct {
	ID        string            `json:"id"`
	NetworkID string            `json:"network_id"`
	Name      string            `json:"name"`
	CIDR      string            `json:"cidr"`
	Zone      string            `json:"zone"`
	Tags      map[string]string `json:"tags,omitempty"`
}

// Config holds configuration for the SDN provider.
type Config struct {
	// Driver specifies the SDN backend: "memory", "vxlan", "calico".
	Driver string `env:"SDN_DRIVER" env-default:"memory"`
}
