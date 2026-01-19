package trie

import (
	"testing"
)

func TestTrie(t *testing.T) {
	tr := New[int]()

	tr.Insert("apple", 1)
	tr.Insert("app", 2)

	if val, ok := tr.Get("apple"); !ok || val != 1 {
		t.Error("Expected apple to exist with value 1")
	}
	if val, ok := tr.Get("app"); !ok || val != 2 {
		t.Error("Expected app to exist with value 2")
	}
	if _, ok := tr.Get("appl"); ok {
		t.Error("Expected appl to be absent")
	}

	results := tr.PrefixSearch("app")
	if len(results) != 2 {
		t.Errorf("Expected 2 results for prefix 'app', got %d", len(results))
	}
}
