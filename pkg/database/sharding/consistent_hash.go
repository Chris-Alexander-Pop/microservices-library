package sharding

import (
	"hash/crc32"
	"sort"
	"strconv"

	"github.com/chris-alexander-pop/system-design-library/pkg/concurrency"
)

// Strategy defines the interface for sharding strategies
type Strategy interface {
	GetShard(key string) string
	AddShard(shardID string)
	RemoveShard(shardID string) // Optional: for dynamic rebalancing (not fully impl in v1)
}

// ConsistentHash implements a consistent hashing ring
type ConsistentHash struct {
	replicas int            // Number of virtual nodes per shard
	keys     []int          // Sorted list of hash keys
	hashMap  map[int]string // Map of hash key to shard ID
	mu       *concurrency.SmartRWMutex
}

// NewConsistentHash creates a new ConsistentHash strategy
func NewConsistentHash(replicas int, shards []string) *ConsistentHash {
	m := &ConsistentHash{
		replicas: replicas,
		hashMap:  make(map[int]string),
		mu:       concurrency.NewSmartRWMutex(concurrency.MutexConfig{Name: "ConsistentHash"}),
	}
	for _, shard := range shards {
		m.AddShard(shard)
	}
	return m
}

func (m *ConsistentHash) AddShard(shardID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := 0; i < m.replicas; i++ {
		hash := int(crc32.ChecksumIEEE([]byte(strconv.Itoa(i) + shardID)))
		m.keys = append(m.keys, hash)
		m.hashMap[hash] = shardID
	}
	sort.Ints(m.keys)
}

func (m *ConsistentHash) RemoveShard(shardID string) {
	// Not implemented for v1 simplicity
}

func (m *ConsistentHash) GetShard(key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.keys) == 0 {
		return ""
	}

	hash := int(crc32.ChecksumIEEE([]byte(key)))

	// Binary search for appropriate replica
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// Wrap around to 0 if we went past the end
	if idx == len(m.keys) {
		idx = 0
	}

	return m.hashMap[m.keys[idx]]
}
