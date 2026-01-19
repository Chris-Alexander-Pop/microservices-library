package vectorclock

import (
	"testing"
)

func TestVectorClock(t *testing.T) {
	vc1 := New()
	vc2 := New()

	vc1.Increment("A")
	vc2.Increment("B")

	if vc1.Compare(vc2) != Concurrent {
		t.Error("Expected A:1, B:0 and A:0, B:1 to be concurrent")
	}

	vc1.Merge(vc2) // vc1 now A:1, B:1

	if vc1.Compare(vc2) != After {
		t.Error("Expected A:1, B:1 to be after A:0, B:1")
	}

	ts := NewThreadSafe()
	ts.Increment("C")
	snap := ts.Snapshot()
	if snap["C"] != 1 {
		t.Errorf("Expected C:1, got %v", snap["C"])
	}
}
