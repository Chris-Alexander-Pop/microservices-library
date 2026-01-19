package network

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestTCPServer(t *testing.T) {
	// Find free port
	l, _ := net.Listen("tcp", "localhost:0")
	addr := l.Addr().String()
	l.Close()

	cfg := Config{
		Addr:        addr,
		ReadTimeout: 1 * time.Second,
	}

	done := make(chan bool)
	handler := func(conn net.Conn) {
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		conn.Write(buf[:n]) // Echo
		done <- true
	}

	server := NewTCPServer(cfg, handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server in background
	go func() {
		server.ListenAndServe(ctx)
	}()
	time.Sleep(100 * time.Millisecond) // Wait for startup

	// Client
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	msg := "hello"
	conn.Write([]byte(msg))

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for server handler")
	}
}
