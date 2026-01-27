// Package memory provides an in-memory implementation of cdn.CDNManager.
package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
	"github.com/chris-alexander-pop/system-design-library/pkg/network/cdn"
	"github.com/google/uuid"
)

// Manager implements an in-memory CDN manager for testing.
type Manager struct {
	mu            sync.RWMutex
	distributions map[string]*cdn.Distribution
	invalidations map[string][]*cdn.Invalidation // distID -> invalidations
}

// New creates a new in-memory CDN manager.
func New() *Manager {
	return &Manager{
		distributions: make(map[string]*cdn.Distribution),
		invalidations: make(map[string][]*cdn.Invalidation),
	}
}

func (m *Manager) CreateDistribution(ctx context.Context, opts cdn.CreateDistributionOptions) (*cdn.Distribution, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.NewString()[:8]
	dist := &cdn.Distribution{
		ID:         id,
		DomainName: fmt.Sprintf("d%s.cloudfront.net", id),
		Status:     cdn.StatusDeployed,
		Origins: []cdn.Origin{{
			ID:         "origin-1",
			DomainName: opts.OriginDomain,
			Protocol:   "https",
			HTTPPort:   80,
			HTTPSPort:  443,
		}},
		CacheBehaviors:    []cdn.CacheBehavior{},
		Enabled:           opts.Enabled,
		SSLCertificateARN: opts.SSLCertificateARN,
		PriceClass:        opts.PriceClass,
		CreatedAt:         time.Now(),
		LastModified:      time.Now(),
	}

	if dist.PriceClass == "" {
		dist.PriceClass = "PriceClass_All"
	}

	m.distributions[id] = dist
	return dist, nil
}

func (m *Manager) GetDistribution(ctx context.Context, id string) (*cdn.Distribution, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	dist, ok := m.distributions[id]
	if !ok {
		return nil, errors.NotFound("distribution not found", nil)
	}
	return dist, nil
}

func (m *Manager) ListDistributions(ctx context.Context) ([]*cdn.Distribution, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*cdn.Distribution, 0, len(m.distributions))
	for _, dist := range m.distributions {
		result = append(result, dist)
	}
	return result, nil
}

func (m *Manager) UpdateDistribution(ctx context.Context, id string, opts cdn.CreateDistributionOptions) (*cdn.Distribution, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	dist, ok := m.distributions[id]
	if !ok {
		return nil, errors.NotFound("distribution not found", nil)
	}

	if opts.OriginDomain != "" {
		dist.Origins[0].DomainName = opts.OriginDomain
	}
	if opts.SSLCertificateARN != "" {
		dist.SSLCertificateARN = opts.SSLCertificateARN
	}
	dist.Enabled = opts.Enabled
	dist.LastModified = time.Now()

	return dist, nil
}

func (m *Manager) DeleteDistribution(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.distributions[id]; !ok {
		return errors.NotFound("distribution not found", nil)
	}
	delete(m.distributions, id)
	delete(m.invalidations, id)
	return nil
}

func (m *Manager) DisableDistribution(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	dist, ok := m.distributions[id]
	if !ok {
		return errors.NotFound("distribution not found", nil)
	}
	dist.Enabled = false
	dist.Status = cdn.StatusDisabled
	return nil
}

func (m *Manager) EnableDistribution(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	dist, ok := m.distributions[id]
	if !ok {
		return errors.NotFound("distribution not found", nil)
	}
	dist.Enabled = true
	dist.Status = cdn.StatusDeployed
	return nil
}

func (m *Manager) Invalidate(ctx context.Context, distributionID string, paths []string) (*cdn.Invalidation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.distributions[distributionID]; !ok {
		return nil, errors.NotFound("distribution not found", nil)
	}

	inv := &cdn.Invalidation{
		ID:             uuid.NewString()[:8],
		DistributionID: distributionID,
		Paths:          paths,
		Status:         "Completed",
		CreatedAt:      time.Now(),
		CompletedAt:    time.Now(),
	}

	m.invalidations[distributionID] = append(m.invalidations[distributionID], inv)
	return inv, nil
}

func (m *Manager) GetInvalidation(ctx context.Context, distributionID, invalidationID string) (*cdn.Invalidation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	invs, ok := m.invalidations[distributionID]
	if !ok {
		return nil, errors.NotFound("invalidation not found", nil)
	}

	for _, inv := range invs {
		if inv.ID == invalidationID {
			return inv, nil
		}
	}
	return nil, errors.NotFound("invalidation not found", nil)
}
