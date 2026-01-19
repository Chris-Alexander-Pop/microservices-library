package roundrobin

import (
	"context"
	"testing"
)

func TestRoundRobin(t *testing.T) {
	balancer := New("node1", "node2", "node3")
	ctx := context.Background()

	// 1
	n, err := balancer.Next(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != "node1" {
		t.Errorf("expected node1, got %s", n)
	}

	// 2
	n, err = balancer.Next(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != "node2" {
		t.Errorf("expected node2, got %s", n)
	}

	// 3
	n, err = balancer.Next(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != "node3" {
		t.Errorf("expected node3, got %s", n)
	}

	// 4 - Loop back
	n, err = balancer.Next(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != "node1" {
		t.Errorf("expected node1, got %s", n)
	}
}
