package linkedlist

import (
	"testing"
)

func TestLinkedList(t *testing.T) {
	l := New[int]()

	l.PushBack(1)
	l.PushBack(2)
	l.PushFront(0)

	if l.Len() != 3 {
		t.Errorf("Expected length 3, got %d", l.Len())
	}
}
