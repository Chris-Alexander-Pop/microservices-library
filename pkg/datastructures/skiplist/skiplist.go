package skiplist

import (
	"math/rand"
	"sync"

	"golang.org/x/exp/constraints"
)

const maxLevel = 16
const p = 0.5

type SkipList[K constraints.Ordered, V any] struct {
	head  *node[K, V]
	level int
	mu    sync.RWMutex
}

type node[K constraints.Ordered, V any] struct {
	key     K
	value   V
	forward []*node[K, V]
}

func New[K constraints.Ordered, V any]() *SkipList[K, V] {
	return &SkipList[K, V]{
		head:  &node[K, V]{forward: make([]*node[K, V], maxLevel)},
		level: 0,
	}
}

func (s *SkipList[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	update := make([]*node[K, V], maxLevel)
	current := s.head

	for i := s.level; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].key < key {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]

	if current != nil && current.key == key {
		current.value = value
		return
	}

	lvl := randomLevel()
	if lvl > s.level {
		for i := s.level + 1; i <= lvl; i++ {
			update[i] = s.head
		}
		s.level = lvl
	}

	newNode := &node[K, V]{
		key:     key,
		value:   value,
		forward: make([]*node[K, V], lvl+1),
	}

	for i := 0; i <= lvl; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
}

func (s *SkipList[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	current := s.head
	for i := s.level; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].key < key {
			current = current.forward[i]
		}
	}

	current = current.forward[0]
	if current != nil && current.key == key {
		return current.value, true
	}

	var zero V
	return zero, false
}

func (s *SkipList[K, V]) Delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	update := make([]*node[K, V], maxLevel)
	current := s.head

	for i := s.level; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].key < key {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]

	if current != nil && current.key == key {
		for i := 0; i <= s.level; i++ {
			if update[i].forward[i] != current {
				break
			}
			update[i].forward[i] = current.forward[i]
		}

		for s.level > 0 && s.head.forward[s.level] == nil {
			s.level--
		}
	}
}

func randomLevel() int {
	lvl := 0
	for rand.Float32() < p && lvl < maxLevel-1 {
		lvl++
	}
	return lvl
}
