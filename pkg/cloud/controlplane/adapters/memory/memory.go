package memory

import (
	"context"

	"github.com/chris-alexander-pop/system-design-library/pkg/cloud"
	"github.com/chris-alexander-pop/system-design-library/pkg/cloud/controlplane"
	"github.com/chris-alexander-pop/system-design-library/pkg/concurrency"
)

// MemoryControlPlane is an in-memory implementation of the ControlPlane interface.
type MemoryControlPlane struct {
	hosts map[string]cloud.Host
	mu    *concurrency.SmartRWMutex
}

// New creates a new MemoryControlPlane.
func New() *MemoryControlPlane {
	return &MemoryControlPlane{
		hosts: make(map[string]cloud.Host),
		mu: concurrency.NewSmartRWMutex(concurrency.MutexConfig{
			Name: "memory-controlplane",
		}),
	}
}

func (c *MemoryControlPlane) RegisterHost(ctx context.Context, host cloud.Host) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.hosts[host.ID]; ok {
		return controlplane.ErrHostAlreadyRegistered
	}

	c.hosts[host.ID] = host
	return nil
}

func (c *MemoryControlPlane) DeregisterHost(ctx context.Context, hostID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.hosts[hostID]; !ok {
		return controlplane.ErrHostNotFound
	}

	delete(c.hosts, hostID)
	return nil
}

func (c *MemoryControlPlane) UpdateHostStatus(ctx context.Context, hostID string, status cloud.HostStatus) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	host, ok := c.hosts[hostID]
	if !ok {
		return controlplane.ErrHostNotFound
	}

	host.Status = status
	c.hosts[hostID] = host
	return nil
}

func (c *MemoryControlPlane) GetHost(ctx context.Context, hostID string) (*cloud.Host, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	host, ok := c.hosts[hostID]
	if !ok {
		return nil, controlplane.ErrHostNotFound
	}

	return &host, nil
}

func (c *MemoryControlPlane) ListHosts(ctx context.Context) ([]cloud.Host, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hosts := make([]cloud.Host, 0, len(c.hosts))
	for _, h := range c.hosts {
		hosts = append(hosts, h)
	}
	return hosts, nil
}
