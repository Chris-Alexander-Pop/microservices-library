package network

import (
	"context"
	"net"

	"github.com/chris-alexander-pop/system-design-library/pkg/logger"
)

type UDPHandler func(addr net.Addr, data []byte)

type UDPServer struct {
	Addr       string
	Handler    UDPHandler
	BufferSize int
}

func NewUDPServer(addr string, handler UDPHandler) *UDPServer {
	return &UDPServer{
		Addr:       addr,
		Handler:    handler,
		BufferSize: 4096, // Tunable
	}
}

func (s *UDPServer) ListenAndServe(ctx context.Context) error {
	pc, err := net.ListenPacket("udp", s.Addr)
	if err != nil {
		return err
	}
	defer pc.Close()

	logger.L().InfoContext(ctx, "started udp server", "addr", s.Addr)

	go func() {
		<-ctx.Done()
		pc.Close()
	}()

	buf := make([]byte, s.BufferSize)
	for {
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			logger.L().ErrorContext(ctx, "udp read error", "error", err)
			continue
		}

		// Cloning buffer for async handling or block?
		// Block is safer for simple UDP loop unless we allocate per packet.
		// For "Overengineered", we allocate valid data slice.
		data := make([]byte, n)
		copy(data, buf[:n])

		go s.Handler(addr, data)
	}
}
