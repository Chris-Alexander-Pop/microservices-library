package heap

import (
	"testing"
)

func TestMinHeap(t *testing.T) {
	h := NewMinHeap[int]()

	h.PushItem(5, 5.0)
	h.PushItem(3, 3.0)
	h.PushItem(7, 7.0)

	val, score, ok := h.Peek()
	if !ok || val != 3 || score != 3.0 {
		t.Errorf("Expected min 3 with score 3.0, got %v with score %v", val, score)
	}

	val, score, ok = h.PopItem()
	if !ok || val != 3 {
		t.Errorf("Expected 3, got %v", val)
	}
	val, score, ok = h.PopItem()
	if !ok || val != 5 {
		t.Errorf("Expected 5, got %v", val)
	}
	val, score, ok = h.PopItem()
	if !ok || val != 7 {
		t.Errorf("Expected 7, got %v", val)
	}
}
