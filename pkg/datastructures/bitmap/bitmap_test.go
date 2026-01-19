package bitmap

import (
	"testing"
)

func TestBitmap(t *testing.T) {
	bm := New(100)

	bm.Set(10)
	bm.Set(50)

	if !bm.Get(10) {
		t.Error("Expected bit 10 to be set")
	}
	if !bm.Get(50) {
		t.Error("Expected bit 50 to be set")
	}
	if bm.Get(20) {
		t.Error("Expected bit 20 to be unset")
	}

	bm.Clear(10)
	if bm.Get(10) {
		t.Error("Expected bit 10 to be cleared")
	}
}
