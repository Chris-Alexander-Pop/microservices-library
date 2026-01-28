package firewall

import (
	"context"
)

// FirewallManager manages network security rules and groups.
type FirewallManager interface {
	// CreateSecurityGroup creates a new group of firewall rules.
	CreateSecurityGroup(ctx context.Context, spec SecurityGroupSpec) (string, error)

	// DeleteSecurityGroup removes a security group.
	DeleteSecurityGroup(ctx context.Context, groupID string) error

	// AddRule adds a rule to a security group.
	AddRule(ctx context.Context, groupID string, rule Rule) error

	// RemoveRule removes a rule from a security group.
	RemoveRule(ctx context.Context, groupID string, ruleID string) error
}

// SecurityGroupSpec defines a new security group.
type SecurityGroupSpec struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	VPCID       string            `json:"vpc_id"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// SecurityGroup represents a collection of firewall rules.
type SecurityGroup struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	VPCID    string            `json:"vpc_id"`
	Inbound  []Rule            `json:"inbound"`
	Outbound []Rule            `json:"outbound"`
	Tags     map[string]string `json:"tags,omitempty"`
}

// Rule defines a single firewall permission.
type Rule struct {
	ID          string `json:"id"`
	Direction   string `json:"direction"` // "inbound" or "outbound"
	Protocol    string `json:"protocol"`  // "tcp", "udp", "icmp", "any"
	PortStart   int    `json:"port_start"`
	PortEnd     int    `json:"port_end"`
	CIDR        string `json:"cidr,omitempty"`
	SourceGroup string `json:"source_group,omitempty"` // For referencing other SGs
}

// Config holds configuration for the Firewall service.
type Config struct {
	// Driver specifies the enforcement backend: "memory", "iptables", "nftables".
	Driver string `env:"FIREWALL_DRIVER" env-default:"memory"`
}
