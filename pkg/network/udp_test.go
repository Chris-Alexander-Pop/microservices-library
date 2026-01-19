package network

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestUDPServer(t *testing.T) {
	// Find free port
	conn, _ := net.ListenPacket("udp", "localhost:0")
	addr := conn.LocalAddr().String()
	conn.Close()

	cfg := Config{
		Addr: addr,
	}

	done := make(chan string)
	handler := func(addr net.Addr, data []byte) {
		done <- string(data)
	}

	server := NewUDPServer(cfg, handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		server.ListenAndServe(ctx)
	}()
	time.Sleep(100 * time.Millisecond)

	// Client
	c, err := net.Dial("udp", addr)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer c.Close()

	msg := "packet"
	c.Write([]byte(msg))

	select {
	case received := <-done:
		if received != msg {
			t.Errorf("Expected %s, got %s", msg, received)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for udp packet")
	}
}
