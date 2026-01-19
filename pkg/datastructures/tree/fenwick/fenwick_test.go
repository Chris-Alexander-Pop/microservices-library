package fenwick

import (
	"testing"
)

func TestFenwickTree(t *testing.T) {
	ft := New(10)

	ft.Add(1, 5)
	ft.Add(3, 2)

	if sum := ft.Query(5); sum != 7 {
		t.Errorf("Expected sum 7, got %d", sum)
	}

	if sum := ft.Query(2); sum != 5 {
		t.Errorf("Expected sum 5, got %d", sum)
	}
}
