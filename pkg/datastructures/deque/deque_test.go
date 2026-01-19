package deque

import (
	"testing"
)

func TestDeque(t *testing.T) {
	d := New[int](10)

	d.PushBack(1)
	d.PushFront(2)

	if d.Len() != 2 {
		t.Errorf("Expected length 2, got %d", d.Len())
	}

	if val, ok := d.PopFront(); !ok || val != 2 {
		t.Errorf("Expected 2, got %v", val)
	}
	if val, ok := d.PopBack(); !ok || val != 1 {
		t.Errorf("Expected 1, got %v", val)
	}
}
