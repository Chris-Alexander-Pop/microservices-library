package ahocorasick

import "container/list"

// Matcher implements the Aho-Corasick algorithm for multi-pattern string matching.
type Matcher struct {
	root *node
}

type node struct {
	children map[rune]*node
	fail     *node
	outputs  []string // patterns ending here
}

func New(patterns []string) *Matcher {
	m := &Matcher{
		root: &node{children: make(map[rune]*node)},
	}
	m.build(patterns)
	return m
}

func (m *Matcher) build(patterns []string) {
	// 1. Build Trie
	for _, p := range patterns {
		n := m.root
		for _, r := range p {
			if n.children[r] == nil {
				n.children[r] = &node{children: make(map[rune]*node)}
			}
			n = n.children[r]
		}
		n.outputs = append(n.outputs, p)
	}

	// 2. Build Failure Links (BFS)
	queue := list.New()

	// Init root children
	for _, child := range m.root.children {
		child.fail = m.root
		queue.PushBack(child)
	}

	for queue.Len() > 0 {
		curr := queue.Remove(queue.Front()).(*node)

		for r, child := range curr.children {
			// Find fail link for child
			fail := curr.fail
			for fail != nil {
				if next, ok := fail.children[r]; ok {
					child.fail = next
					break
				}
				fail = fail.fail
			}
			if child.fail == nil {
				child.fail = m.root
			}

			// Merge outputs
			child.outputs = append(child.outputs, child.fail.outputs...)

			queue.PushBack(child)
		}
	}
}

// Match returns all occurrences of patterns in text.
// Returns map[pattern][]index (start indices or end indices? let's do end indices for simplicity)
type Match struct {
	Pattern string
	Index   int
}

func (m *Matcher) FindAll(text string) []Match {
	var matches []Match
	curr := m.root

	for i, r := range text {
		for curr != nil {
			if next, ok := curr.children[r]; ok {
				curr = next
				break
			}
			curr = curr.fail
		}
		if curr == nil {
			curr = m.root
		}

		for _, pattern := range curr.outputs {
			// Index is the end position of the match
			matches = append(matches, Match{
				Pattern: pattern,
				Index:   i - len(pattern) + 1,
			})
		}
	}
	return matches
}
